package main

import (
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

// oneFileCommandFunc builds a single file.
func oneFileCommandFunc() {
	fileCmd := darknessFlagset(oneFileCommand)
	fileCmd.StringVar(&filename, "input", "index.org", "file on input")
	emilia.InitDarkness(getEmiliaOptions(fileCmd))
	fmt.Println(emilia.InputToOutput(filename))
}

// build builds the entire directory.
func buildCommandFunc() {
	emilia.InitDarkness(getEmiliaOptions(darknessFlagset(buildCommand)))
	build()
	// Check that we actually processed some files before reporting.
	if emilia.NumFoundFiles < 0 {
		fmt.Println("no files found")
		return
	}

	fmt.Println("farewell")
}

// build uses set flags and emilia data to build the local directory.
func build() {
	start := time.Now()

	// Create the channel to feed read files.
	inputFilenames := make(chan string, customChannelCapacity)

	// Create the worker that will read files and push tuples.
	inputFiles := gana.GenericWorkers(inputFilenames, func(v string) gana.Tuple[string, string] {
		data, err := ioutil.ReadFile(filepath.Clean(v))
		if err != nil {
			fmt.Printf("Failed to open %s: %s\n", v, err.Error())
		}
		return gana.NewTuple(v, string(data))
	}, 1, customChannelCapacity)

	// Create the workers for parsing and converting orgmode to Page.
	pages := gana.GenericWorkers(inputFiles, func(v gana.Tuple[string, string]) *yunyun.Page {
		return emilia.ParserBuilder.BuildParser(fdb(v.UnpackRef())).Parse()
	}, customNumWorkers, customChannelCapacity)

	// Create the workers for building Page's into html documents.
	results := gana.GenericWorkers(pages, func(v *yunyun.Page) gana.Tuple[string, string] {
		return gana.NewTuple(emilia.InputFilenameToOutput(v.File), emilia.EnrichAndExportPage(v))
	}, customNumWorkers, customChannelCapacity)

	// This will block darkness from exiting until all the files are done.
	wg := &sync.WaitGroup{}

	// Add a block here so the file explorer has a bit of time to spin
	// up and start filling up its channel.
	wg.Add(1)

	// Run a discovery for files and feed to the reader worker.
	go emilia.FindFilesByExt(inputFilenames, workDir, emilia.Config.Project.Input, wg)

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
	fmt.Printf("Processed %d files in %d ms\n", emilia.NumFoundFiles, time.Since(start).Milliseconds())
}
