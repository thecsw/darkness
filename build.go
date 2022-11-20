package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

const (
	// savePerms tells us what permissions to use for the
	// final export files.
	savePerms = fs.FileMode(0644)
)

var (
	// workDir is the directory to look for files.
	workDir = "."
	// darknessToml is the location of `darkness.toml`.
	darknessToml = "darkness.toml"
	// filename is the file to build.
	filename = "index.org"
	// defaultNumOfWorkers gives us the number of workers to
	// spin up in each stage: parsing and processing.
	defaultNumOfWorkers = 14
)

// oneFile builds a single file.
func oneFile() {
	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	fileCmd.StringVar(&filename, "i", "index.org", "file on input")
	fileCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	fileCmd.Parse(os.Args[2:])
	emilia.InitDarkness(&emilia.EmiliaOptions{DarknessConfig: darknessToml})
	fmt.Println(inputToOutput(filename))
}

// build builds the entire directory.
func build() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	buildCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	disableParallel := buildCmd.Bool("disable-parallel", false, "disable parallel build (only use one worker)")
	customNumWorkers := buildCmd.Int("workers", defaultNumOfWorkers, "number of workers to spin up")
	customChannelCapacity := buildCmd.Int("capacity", defaultNumOfWorkers, "worker channels' capacity")
	useCurrentDirectory := buildCmd.Bool("dev", false, "use local path for urls (development)")
	buildCmd.Parse(os.Args[2:])

	// Read the config and initialize emilia settings.
	emilia.InitDarkness(&emilia.EmiliaOptions{
		DarknessConfig: darknessToml,
		Dev:            *useCurrentDirectory,
	})

	var err error
	workDir, err = filepath.Abs(workDir)
	if err != nil {
		fmt.Println("Couldn't determine absolute path of", workDir)
		os.Exit(1)
	}

	// If parallel processing is disabled, only provision one workers
	// per each processing stage.
	if *disableParallel {
		*customNumWorkers = 1
	}

	// Create the channel to feed read files.
	inputFilenames := make(chan string, *customChannelCapacity)

	// Create the worker that will read files and push tuples.
	inputFiles := gana.GenericWorkers(inputFilenames, func(v string) gana.Tuple[string, string] {
		data, err := ioutil.ReadFile(filepath.Clean(v))
		if err != nil {
			fmt.Printf("Failed to open %s: %s\n", v, err.Error())
		}
		return gana.NewTuple(v, string(data))
	}, 1, *customChannelCapacity)

	// Create the workers for parsing and converting orgmode to Page.
	pages := gana.GenericWorkers(inputFiles, func(v gana.Tuple[string, string]) *yunyun.Page {
		//return getParser().Parse()
		return emilia.ParserBuilder.BuildParser(v.Unpack()).Parse()
	}, *customNumWorkers, *customChannelCapacity)

	// Create the workers for building Page's into html documents.
	results := gana.GenericWorkers(pages, func(v *yunyun.Page) gana.Tuple[string, string] {
		return gana.NewTuple(getTarget(v.File), exportAndEnrich(v))
	}, *customNumWorkers, *customChannelCapacity)

	// This will block darkness from exiting until all the files are done.
	wg := &sync.WaitGroup{}

	// Add a block here so the file explorer has a bit of time to spin
	// up and start filling up its channel.
	wg.Add(1)

	// Find all the appropriate orgmode files and save the list.
	start := time.Now()

	// Run a discovery for files and feed to the reader worker.
	go findFilesByExt(inputFilenames, wg)

	// Build a wait group to ensure we always read and write the same
	// number of files, such that after the file has been read, parsed,
	// enriched, and exported -- this goroutine would pick them up and
	// save it at the right spot, marking itself Done and leaving.
	go func(wg *sync.WaitGroup) {
		for result := range results {
			if err := os.WriteFile(result.First, []byte(result.Second), savePerms); err != nil {
				fmt.Printf("failed to write %s: %s", result.First, err.Error())
			}
			wg.Done()
		}
		// Remove the artificial block we made before discovery.
		wg.Done()
	}(wg)

	// Wait for all the files to get saved and then leave.
	wg.Wait()

	// Report back on some of the results
	fmt.Printf("Processed in %d ms\n", time.Since(start).Milliseconds())
	fmt.Println("farewell")
}
