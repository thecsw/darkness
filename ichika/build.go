package ichika

import (
	"fmt"
	"runtime"
	"time"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/ichika/akane"
	"github.com/thecsw/darkness/ichika/hizuru"
	"github.com/thecsw/darkness/ichika/makima"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/komi"
	"github.com/thecsw/rei"
)

var akaneless = false

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

	if !akaneless {
		// Let's complete the akane requests when done building.
		defer akane.Do(conf)
	}

	// Create the pool that reads files and returns their handles.
	filesPool := komi.NewWithSettings(komi.WorkWithErrors(makima.Woof.Read), &komi.Settings{
		Name:     "Komi Reading 📚 ",
		Laborers: runtime.NumCPU(),
		Debug:    debugEnabled,
	})
	filesError := rei.Must(filesPool.Errors())
	go logErrors("reading", filesError)

	// Create a pool that take a files handle and parses it out into yunyun pages.
	parserPool := komi.NewWithSettings(komi.Work(makima.Woof.Parse), &komi.Settings{
		Name:     "Komi Parsing 🧹 ",
		Laborers: customNumWorkers,
		Debug:    debugEnabled,
	})

	// Create a pool that that takes yunyun pages and exports them into request format.
	exporterPool := komi.NewWithSettings(komi.Work(makima.Woof.Export), &komi.Settings{
		Name:     "Komi Exporting 🥂 ",
		Laborers: customNumWorkers,
		Debug:    debugEnabled,
	})

	// Create a pool that reads the exported data and writes them to target files.
	writerPool := komi.NewWithSettings(komi.WorkSimpleWithErrors(makima.Woof.Write), &komi.Settings{
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

	// Record the start time.
	start := time.Now()

	// Find all the files that need to be parsed.
	inputFilenames := make(chan yunyun.FullPathFile, 8)
	go hizuru.FindFilesByExt(conf, inputFilenames)

	// Submit all the files to the pool.
	for inputFilename := range inputFilenames {
		// Submit the job to the pool.
		rei.Try(filesPool.Submit(&makima.Control{
			Conf:          conf,
			Parser:        parser,
			Exporter:      exporter,
			InputFilename: inputFilename,
		}))
	}

	// Wait for all the pools to finish.
	writerPool.Close()

	// Record the time it took to finish.
	finish := time.Now()

	// Clear the download progress bar if present by wiping out the line.
	fmt.Print("\r\033[2K")

	fmt.Printf("Processed %d files in %d ms\n", exporterPool.JobsSucceeded(), finish.Sub(start).Milliseconds())
}

// logErrors is a helper function that logs errors from a pool. It is meant to be
// used as a goroutine.
func logErrors[T any](name string, vv chan komi.PoolError[T]) {
	for v := range vv {
		if v.Error != nil {
			puck.Logger.Error("job failed", "err", v.Error, "pool", name)
		}
	}
}
