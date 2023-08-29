package ichika

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/ichika/makima"

	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/komi"
	"github.com/thecsw/rei"
)

const (
	// savePerms tells us what permissions to use for the
	// final export files.
	savePerms = fs.FileMode(0o644)
)

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
var vendorGalleryImages bool

// OneFileCommandFunc builds a single file.
func OneFileCommandFunc() {
	panic("todo")
	// fileCmd := darknessFlagset(oneFileCommand)
	// fileCmd.StringVar(&filename, "input", "index.org", "file on input")
	// conf := alpha.BuildConfig(getAlphaOptions(fileCmd))
	// fmt.Println(conf.Runtime.InputToOutput(conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(filename))))
	//fmt.Println(emilia.InputToOutput(emilia.Join(yunyun.RelativePathFile(filename))))
}

// BuildCommandFunc builds the entire directory.
func BuildCommandFunc() {
	cmd := darknessFlagset(buildCommand)
	conf := alpha.BuildConfig(getAlphaOptions(cmd))
	build(conf)
	fmt.Println("farewell")
}

// build uses set flags and emilia data to build the local directory.
func build(conf *alpha.DarknessConfig) {

	parser := parse.BuildParser(conf)
	exporter := export.BuildExporter(conf)

	// Create the pool that reads files and returns their handles.
	filesPool := komi.NewWithSettings(komi.WorkWithErrors(openPage), &komi.Settings{
		Name:     "Komi Reading 📚 ",
		Laborers: runtime.NumCPU(),
		Debug:    debugEnabled,
	})
	filesError := rei.Must(filesPool.Errors())
	go logErrors("reading", filesError)

	// Create a pool that take a files handle and parses it out into yunyun pages.
	parserPool := komi.NewWithSettings(komi.Work(parsePage), &komi.Settings{
		Name:     "Komi Parsing 🧹 ",
		Laborers: customNumWorkers,
		Debug:    debugEnabled,
	})

	// Create a pool that that takes yunyun pages and exports them into request format.
	exporterPool := komi.NewWithSettings(komi.Work(exportPage), &komi.Settings{
		Name:     "Komi Exporting 🥂 ",
		Laborers: customNumWorkers,
		Debug:    debugEnabled,
	})

	// Create a pool that reads the exported data and writes them to target files.
	writerPool := komi.NewWithSettings(komi.WorkSimpleWithErrors(writePage), &komi.Settings{
		Name:     "Komi Writing 🎸",
		Laborers: runtime.NumCPU(),
		Debug:    debugEnabled,
	})
	writersErrors := rei.Must(writerPool.Errors())
	go logErrors("writer", writersErrors)

	// Connect all the pools between each other, so the relationship is as follows,
	//
	//           Reading 📚                      Parsing 🧹
	//   path  ┌───────────┐   file handler   ┌─────────────┐
	// ──────> │ filesPool │ ───────────────> │  parserPool │
	//	   └───────────┘	          └─────────────┘
	//	    log errors				│
	//						│    parsed files
	//					       	│  aka yunyun pages
	//                                              │
	//   file  ┌────────────┐  exported data  ┌──────────────┐
	//  <───── │ writerPool │ <────────────── │ exporterPool │
	// 	   └────────────┘              	  └──────────────┘
	//           Writing 🎸                     Exporting 🥂
	//
	rei.Try(filesPool.Connect(parserPool))
	rei.Try(parserPool.Connect(exporterPool))
	rei.Try(exporterPool.Connect(writerPool))

	start := time.Now()

	freshContext := makima.Control{
		Conf:     conf,
		Parser:   parser,
		Exporter: exporter,
	}

	<-FindFilesByExt(conf, filesPool, freshContext)

	writerPool.Close()

	finish := time.Now()

	// Clear the download progress bar if present by wiping out the line.
	fmt.Print("\r\033[2K")

	fmt.Printf("Processed %d files in %d ms\n", exporterPool.JobsSucceeded(), finish.Sub(start).Milliseconds())
}

//go:inline
func openPage(c *makima.Control) (*makima.Control, error) {
	file, err := os.ReadFile(filepath.Clean(string(c.InputFilename)))
	if err != nil {
		return nil, err
	}
	c.Input = string(file)
	return c, nil
}

//go:inline
func parsePage(c *makima.Control) *makima.Control {
	return c.Parse()
}

//go:inline
func exportPage(c *makima.Control) *makima.Control {
	return c.Export()
}

//go:inline
func writePage(c *makima.Control) error {
	_, err := writeFile(c.OutputFilename, c.Output)
	if err != nil {
		return fmt.Errorf("writing page %s: %v", c.OutputFilename, err)
	}
	return nil
}

// logErrors is a helper function that logs errors from a pool. It is meant to be
// used as a goroutine.
func logErrors[T any](name string, vv chan komi.PoolError[T]) {
	for v := range vv {
		if v.Error != nil {
			puck.Logger.Errorf("pool %s encountered an error: %v", name, v.Error)
		}
	}
}
