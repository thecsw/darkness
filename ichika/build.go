package ichika

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
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
	// vendorGalleryImages is a flag that dictates whether we should
	// store a local copy of all remote gallery images and stub them
	// in the gallery links instead of the remote links.
	//
	// Turning this option on would result in a VERY slow build the
	// first time, as it would need to retrieve however many images
	// from remote services.
	//
	// All images will be put in "darkness_vendor" directory, which
	// will be skipped in discovery process AND should be put it
	// .gitignore by user, so they don't pollute their git objects.
	vendorGalleryImages bool
)

// OneFileCommandFunc builds a single file.
func OneFileCommandFunc() {
	fileCmd := darknessFlagset(oneFileCommand)
	fileCmd.StringVar(&filename, "input", "index.org", "file on input")
	emilia.InitDarkness(getEmiliaOptions(fileCmd))
	fmt.Println(emilia.InputToOutput(emilia.JoinWorkdir(yunyun.RelativePathFile(filename))))
}

// build builds the entire directory.
func BuildCommandFunc() {
	cmd := darknessFlagset(buildCommand)
	emilia.InitDarkness(getEmiliaOptions(cmd))
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
	inputFilenames := make(chan yunyun.FullPathFile, customChannelCapacity)

	// Create the worker that will read files and push tuples.
	inputFiles := gana.GenericWorkers(inputFilenames, func(v yunyun.FullPathFile) gana.Tuple[yunyun.FullPathFile, *os.File] {
		return openFile(v)
	}, 1, customChannelCapacity)

	// Create the workers for parsing and converting orgmode to Page.
	pages := gana.GenericWorkers(inputFiles, func(v gana.Tuple[yunyun.FullPathFile, *os.File]) *yunyun.Page {
		return emilia.ParserBuilder.BuildParserReader(emilia.FullPathToWorkDirRel(v.First), v.Second).Parse()
	}, customNumWorkers, customChannelCapacity)

	// Create the workers for building Page's into html documents.
	results := gana.GenericWorkers(pages, func(v *yunyun.Page) gana.Tuple[string, *bufio.Reader] {
		return gana.NewTuple(emilia.InputFilenameToOutput(emilia.JoinWorkdir(v.File)), emilia.EnrichExportPageAsBufio(v))
	}, customNumWorkers, customChannelCapacity)

	// This will block darkness from exiting until all the files are done.
	wg := &sync.WaitGroup{}

	// Add a block here so the file explorer has a bit of time to spin
	// up and start filling up its channel.
	wg.Add(1)

	// Run a discovery for files and feed to the reader worker.
	go emilia.FindFilesByExt(inputFilenames, emilia.Config.Project.Input, wg)

	// Build a wait group to ensure we always read and write the same
	// number of files, such that after the file has been read, parsed,
	// enriched, and exported -- this goroutine would pick them up and
	// save it at the right spot, marking itself Done and leaving.
	go func(wg *sync.WaitGroup) {
		for result := range results {
			if _, err := writeFile(result.First, result.Second); err != nil {
				fmt.Println("writing file:", err.Error())
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

// writeFile takes a filename and a bufio reader and writes it.
func writeFile(filename string, reader *bufio.Reader) (int64, error) {
	target, err := os.Create(filename)
	if err != nil {
		return -1, errors.Wrap(err, "failed to create "+filename)
	}
	written, err := io.Copy(target, reader)
	if err != nil {
		return -1, errors.Wrap(err, "failed to copy to "+filename)
	}
	if target.Close() != nil {
		return -1, errors.Wrap(err, "failed to close "+filename)
	}
	return written, nil
}

// openFile attemps to open the full path and return tuple, empty tuple otherwise.
func openFile(v yunyun.FullPathFile) gana.Tuple[yunyun.FullPathFile, *os.File] {
	file, err := os.Open(filepath.Clean(string(v)))
	if err != nil {
		log.Printf("failed to open %s: %s\n", v, err)
		return gana.NewTuple[yunyun.FullPathFile, *os.File]("", nil)
	}
	return gana.NewTuple(v, file)
}
