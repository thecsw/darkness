package ichika

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/export"
	"github.com/thecsw/darkness/v3/ichika/akane"
	"github.com/thecsw/darkness/v3/ichika/hizuru"
	"github.com/thecsw/darkness/v3/ichika/kuroko"
	"github.com/thecsw/darkness/v3/ichika/makima"
	"github.com/thecsw/darkness/v3/ichika/misaka"
	"github.com/thecsw/darkness/v3/parse"
	"github.com/thecsw/darkness/v3/parse/orgmode"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/komi"
	"github.com/thecsw/rei"
)

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

	// Let's complete the akane requests when done building.
	if !kuroko.Akaneless {
		defer akane.Do(conf)
	}

	// Before we kick off the entire parsing loop, let's see if we have global macros defined.
	recordGlobalMacros(conf)

	// Create the pool that reads files and returns their handles.
	filesPool := komi.NewWithSettings(komi.WorkWithErrors(makima.Woof.Read), &komi.Settings{
		Name:     "Komi Reading 📚 ",
		Laborers: runtime.NumCPU(),
		Debug:    kuroko.DebugEnabled,
	})
	go logErrors("reading", rei.Must(filesPool.Errors()))

	// Create a pool that take a files handle and parses it out into yunyun pages.
	parserPool := komi.NewWithSettings(komi.Work(makima.Woof.Parse), &komi.Settings{
		Name:     "Komi Parsing 🧹 ",
		Laborers: kuroko.CustomNumWorkers,
		Debug:    kuroko.DebugEnabled,
	})

	// Create a pool that that takes yunyun pages and exports them into request format.
	exporterPool := komi.NewWithSettings(komi.Work(makima.Woof.Export), &komi.Settings{
		Name:     "Komi Exporting 🥂 ",
		Laborers: kuroko.CustomNumWorkers,
		Debug:    kuroko.DebugEnabled,
	})

	// Create a pool that reads the exported data and writes them to target files.
	writerPool := komi.NewWithSettings(komi.WorkSimpleWithErrors(makima.Woof.Write), &komi.Settings{
		Name:     "Komi Writing 🎸",
		Laborers: runtime.NumCPU(),
		Debug:    kuroko.DebugEnabled,
	})
	go logErrors("writer", rei.Must(writerPool.Errors()))

	// Connect all the pools between each other, so the relationship is as follows,
	//
	//           Reading 📚                      Parsing 🧹
	//   path  ┌───────────┐   file handler   ┌─────────────┐
	// ──────> │ filesPool │ ───────────────> │  parserPool │
	//         └───────────┘                  └─────────────┘
	//          log errors                          │
	//                                              │    parsed files
	//                                              │  aka yunyun pages
	//                                              │
	//   file  ┌────────────┐  exported data  ┌──────────────┐
	//  <───── │ writerPool │ <────────────── │ exporterPool │
	//         └────────────┘              	  └──────────────┘
	//           Writing 🎸                     Exporting 🥂
	//
	rei.Try(filesPool.Connect(parserPool))
	rei.Try(parserPool.Connect(exporterPool))
	rei.Try(exporterPool.Connect(writerPool))

	// Find all the files that need to be parsed.
	inputFilenames := make(chan yunyun.FullPathFile, 8)
	go hizuru.FindFilesByExt(conf, inputFilenames)

	// Record the start time.
	start := time.Now()

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

	// Let's process the misaka report if user wants to see it.
	if kuroko.BuildReport {
		misaka.WriteReport(conf)
	}

	// Let's write the report time to a special file, last_built.txt
	nowUtc := time.Now().UTC().Format(time.RFC3339)
	if err := os.WriteFile(puck.LastBuildTimestampFile, []byte(nowUtc), 0o600); err != nil {
		conf.Runtime.Logger.Warnf("couldn't write the last_built.txt: %v", err)
	}
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

// recordGlobalMacros will read the global macros and inject them into pages.
func recordGlobalMacros(conf *alpha.DarknessConfig) {
	// Recall that this is primarily an orgmode feature, so we will lock it to that input ext.
	if conf.Project.Input != puck.ExtensionOrgmode {
		return
	}
	globalMacrosFile := yunyun.RelativePathFile(globalMacrosFileBasename + conf.Project.Input)
	globalMacrosFileFull := string(conf.Runtime.WorkDir.Join(globalMacrosFile))
	if exists, err := rei.FileExists(globalMacrosFileFull); exists {
		file, err := os.ReadFile(filepath.Clean(globalMacrosFileFull))
		if err != nil {
			conf.Runtime.Logger.Warn("Failed reading global macros file", "file", globalMacrosFileFull, "err", err)
		}
		if orgmode.CollectGlobalMacros(conf, globalMacrosFile, string(file)) {
			conf.Runtime.Logger.Info("Loaded global macros")
		}

	} else if err != nil {
		conf.Runtime.Logger.Warn("Failed checking global macros file", "err", err)
	}
}
