package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/html"
	"github.com/thecsw/darkness/internals"
	"github.com/thecsw/darkness/orgmode"
)

var (
	// workDir is the directory to look for files
	workDir = "."
	// darknessToml is the location of darkness.toml
	darknessToml = "darkness.toml"
	// filename is the file to build
	filename = "index.org"
)

// oneFile builds a single file
func oneFile() {
	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	fileCmd.StringVar(&filename, "i", "index.org", "file on input")
	fileCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	fileCmd.Parse(os.Args[2:])
	emilia.InitDarkness(darknessToml)
	fmt.Println(orgToHTML(filename))
}

// build builds the entire directory
func build() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	buildCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	disableParallel := buildCmd.Bool("disable-parallel", false, "disable parallel build (only use one worker)")
	customNumWorkers := buildCmd.Int("workers", defaultNumOfWorkers, "number of workers to spin up")
	customChannelCapacity := buildCmd.Int("capacity", defaultNumOfWorkers, "worker channels' capacity")
	buildCmd.Parse(os.Args[2:])

	// Read the config and initialize emilia settings.
	emilia.InitDarkness(darknessToml)

	// Initialize some of the custom exporter settings.
	html.InitializeExporter()

	// Find all the appropriate orgmode files and save the list.
	start := time.Now()

	// Set the channel capacity to user input.
	channelCapacity = *customChannelCapacity

	// If parallel processing is disabled, only provision one workers
	// per each processing stage.
	if *disableParallel {
		*customNumWorkers = 1
	}

	// Create the channel to feed read files.
	orgfiles := make(chan string, channelCapacity)

	// Create the worker that will read files and push bundles.
	orgmodes := genericWorkers(orgfiles, func(v string) *bundle {
		data, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Printf("Failed to open %s: %s\n", v, err.Error())
		}
		return &bundle{
			File: v,
			Data: string(data),
		}
	}, 1)

	// Create the workers for parsing and converting orgmode to Page.
	pages := genericWorkers(orgmodes, func(v *bundle) *internals.Page {
		return orgmode.Parse(v.Data, v.File)
	}, *customNumWorkers)

	// Create the workers for building Page's into html documents.
	results := genericWorkers(pages, func(v *internals.Page) *bundle {
		return &bundle{
			File: getTarget(v.File),
			Data: exportAndEnrich(v),
		}
	}, *customNumWorkers)

	// This will block darkness from exiting until all the files are done.
	wg := &sync.WaitGroup{}

	// Add a block here so the file explorer has a bit of time to spin
	// up and start filling up its channel.
	wg.Add(1)

	// Run a discovery for files and feed to the reader worker.
	go findFilesByExt(workDir, emilia.Config.Project.Input, orgfiles, wg)

	// Build a wait group to ensure we always read and write the same
	// number of files, such that after the file has been read, parsed,
	// enriched, and exported -- this goroutine would pick them up and
	// save it at the right spot, marking itself Done and leaving.
	go func(wg *sync.WaitGroup) {
		for result := range results {
			os.WriteFile(result.File, []byte(result.Data), savePerms)
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
