package ichika

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/export"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
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
func build(conf alpha.DarknessConfig) {

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
	parserPool := komi.NewWithSettings(komi.Work(func(v gana.Tuple[yunyun.FullPathFile, string]) *yunyun.Page {
		return parser.Do(conf.Runtime.WorkDir.Rel(v.First), v.Second)
	}), &komi.Settings{
		Name:     "Komi Parsing 🧹 ",
		Laborers: customNumWorkers,
		Debug:    debugEnabled,
	})

	// Create a pool that that takes yunyun pages and exports them into request format.
	exporterPool := komi.NewWithSettings(komi.Work(func(v *yunyun.Page) gana.Tuple[string, io.Reader] {
		return gana.NewTuple(conf.InputFilenameToOutput(conf.Runtime.WorkDir.Join(v.File)), exporter.Do(EnrichPage(conf, v)))
	}), &komi.Settings{
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

	<-emilia.FindFilesByExt(conf, filesPool)

	writerPool.Close()

	finish := time.Now()

	// Clear the download progress bar if present by wiping out the line.
	fmt.Print("\r\033[2K")

	fmt.Printf("Processed %d files in %d ms\n", exporterPool.JobsSucceeded(), finish.Sub(start).Milliseconds())
}

//go:inline
func openPage(v yunyun.FullPathFile) (gana.Tuple[yunyun.FullPathFile, string], error) {
	file, err := os.ReadFile(filepath.Clean(string(v)))
	if err != nil {
		return gana.NewTuple[yunyun.FullPathFile, string]("", ""), err
	}
	return gana.NewTuple(v, string(file)), nil
}

//go:inline
func writePage(v gana.Tuple[string, io.Reader]) error {
	_, err := writeFile(v.First, v.Second)
	if err != nil {
		return fmt.Errorf("writing page %s: %v", v.First, err)
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
