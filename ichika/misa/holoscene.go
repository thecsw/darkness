package misa

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/narumi"
	"github.com/thecsw/darkness/v3/ichika/hizuru"
	"github.com/thecsw/darkness/v3/ichika/kuroko"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/rei"
)

const (
	holosceneTitlesTempDir = "temp-holoscene"
)

// UpdateHoloceneTitles adds holoscene titles to the output files.
func UpdateHoloceneTitles(conf *alpha.DarknessConfig, dryRun bool) {
	initLog()
	if dryRun {
		if err := rei.Mkdir(holosceneTitlesTempDir); err != nil {
			logger.Fatalf("creating temporary directory %s: %v", holosceneTitlesTempDir, err)
		}
	}

	// Find all the files that need to be updated.
	inputFilenames := hizuru.FindFilesByExtSimple(conf)

	// Convert the input filenames to output filenames.
	outputs := make([]string, len(inputFilenames))
	for i, inputFilename := range inputFilenames {
		outputs[i] = conf.Project.InputFilenameToOutput(inputFilename)
	}

	// Open all the output files.
	actuallyFound := make([]*os.File, 0, len(outputs))
	for _, v := range outputs {
		v := filepath.Clean(v)
		file, err := os.Open(v)
		if err != nil {
			logger.Errorf("opening file %s: %v", v, err)
			continue
		}
		actuallyFound = append(actuallyFound, file)
	}

	// Add holoscene titles to all the output files.
	fmt.Printf("Adding holoscene titles to %d output files\n", len(actuallyFound))

	skipped := atomic.Int32{}

	// Add holoscene titles to all the output files.
	for _, foundOutput := range actuallyFound {
		start := time.Now()

		// Read the output file.
		filename := foundOutput.Name()
		output, err := io.ReadAll(foundOutput)
		if err := foundOutput.Close(); err != nil {
			logger.Errorf("closing file %s: %v", filename, err)
		}
		if err != nil {
			logger.Errorf("reading file %s: %v", filename, err)
			continue
		}

		// Add holoscene titles to the output.
		newOutput := narumi.AddHolosceneTitles(string(output), -1)

		// Skip if the same, unless forced.
		if !kuroko.Force {
			if len(output) == len(newOutput) {
				skipped.Add(1)
				continue
			}
		}

		var file *os.File
		if dryRun {
			file, err = os.CreateTemp(holosceneTitlesTempDir,
				filepath.Base(filename))
		} else {
			file, err = os.Create(filepath.Clean(filename))
		}
		if err != nil {
			logger.Errorf("overwriting %s: %v", filename, err)
			continue
		}

		// Write the new output to the file.
		_, err = io.Copy(file, strings.NewReader(newOutput))
		if err := file.Close(); err != nil {
			logger.Errorf("closing file %s: %v", file.Name(), err)
		}
		if err != nil {
			logger.Errorf("writing file %s: %v", file.Name(), err)
			continue
		}

		logger.Info("Added holoscene titles",
			"loc", conf.Runtime.WorkDir.Rel(yunyun.FullPathFile(file.Name())),
			"elapsed", time.Since(start),
		)
		if dryRun {
			fmt.Printf(": %s", strings.TrimPrefix(filename, string(conf.Runtime.WorkDir)))
		}
	}

	if dryRun {
		if err := os.RemoveAll(holosceneTitlesTempDir); err != nil {
			logger.Errorf("clearing temporary directory %s: %v", holosceneTitlesTempDir, err)
		}
	}

	// Write a notice if we skipped any preview generations.
	if numSkipped := skipped.Load(); numSkipped > 0 {
		logger.Warn("Some outputs weren't affected, use -force to overwrite", "skipped", numSkipped)
	}
}
