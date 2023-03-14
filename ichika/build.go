package ichika

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
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

	filesPool := gana.NewPool(openFile, customNumWorkers, 0, "files")
	parserPool := gana.NewPool(parsePage, customNumWorkers, 0, "parser")
	exporterPool := gana.NewPool(exportPage, customNumWorkers, 0, "exporter")

	filesPool.Connect(parserPool)
	parserPool.Connect(exporterPool)

	// Run a discovery for files and feed to the reader worker.
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go emilia.FindFilesByExt(filesPool, emilia.Config.Project.Input, wg)

	// Build a wait group to ensure we always read and write the same
	// number of files, such that after the file has been read, parsed,
	// enriched, and exported -- this goroutine would pick them up and
	// save it at the right spot, marking itself Done and leaving.
	go func() {
		for result := range exporterPool.Outputs() {
			if _, err := writeFile(result.First, result.Second); err != nil {
				fmt.Println("writing file:", err.Error())
			}
		}
	}()

	wg.Wait()
	// time.Sleep(time.Second)

	// filesPool.Wait()
	// parserPool.Wait()
	// exporterPool.Wait()

	// filesPool.Close()
	// parserPool.Close()
	// exporterPool.Close()

	//exporterPool.Close()

	// Clear the download progress bar if present by wiping out the line.
	fmt.Print("\r\033[2K")

	fmt.Printf("Processed %d files in %d ms\n", emilia.NumFoundFiles, time.Since(start).Milliseconds())
}

//go:inline
func openPage(v yunyun.FullPathFile) gana.Tuple[yunyun.FullPathFile, *os.File] {
	return openFile(v)
}

//go:inline
func parsePage(v gana.Tuple[yunyun.FullPathFile, *os.File]) *yunyun.Page {
	return emilia.ParserBuilder.BuildParserReader(emilia.FullPathToWorkDirRel(v.First), v.Second).Parse()
}

//go:inline
func exportPage(v *yunyun.Page) gana.Tuple[string, *bufio.Reader] {
	return gana.NewTuple(emilia.InputFilenameToOutput(emilia.JoinWorkdir(v.File)), emilia.EnrichExportPageAsBufio(v))
}
