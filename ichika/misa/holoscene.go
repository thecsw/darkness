package misa

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/narumi"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/ichika/hizuru"
)

const (
	holosceneTitlesTempDir = "temp-holoscene"
)

// UpdateHoloceneTitles adds holoscene titles to the output files.
func UpdateHoloceneTitles(conf *alpha.DarknessConfig, dryRun bool) {
	if dryRun {
		if err := os.Mkdir(holosceneTitlesTempDir, 0o750); err != nil {
			puck.Logger.Fatalf("creating temporary directory %s: %v", holosceneTitlesTempDir, err)
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
			puck.Logger.Errorf("opening file %s: %v", v, err)
			continue
		}
		actuallyFound = append(actuallyFound, file)
	}

	// Add holoscene titles to all the output files.
	fmt.Printf("Adding holoscene titles to %d output files\n", len(actuallyFound))

	// Add holoscene titles to all the output files.
	for _, foundOutput := range actuallyFound {
		// Read the output file.
		filename := foundOutput.Name()
		output, err := io.ReadAll(foundOutput)
		if err := foundOutput.Close(); err != nil {
			puck.Logger.Errorf("closing file %s: %v", filename, err)
		}
		if err != nil {
			puck.Logger.Errorf("reading file %s: %v", filename, err)
			continue
		}

		// Add holoscene titles to the output.
		newOutput := narumi.AddHolosceneTitles(string(output), -1)
		var file *os.File
		if dryRun {
			file, err = os.CreateTemp(holosceneTitlesTempDir,
				filepath.Base(filename))
		} else {
			file, err = os.Create(filepath.Clean(filename))
		}
		if err != nil {
			puck.Logger.Errorf("overwriting %s: %v", filename, err)
			continue
		}

		// Write the new output to the file.
		written, err := io.Copy(file, strings.NewReader(newOutput))
		if err := file.Close(); err != nil {
			puck.Logger.Errorf("closing file %s: %v", file.Name(), err)
		}
		if err != nil {
			puck.Logger.Errorf("writing file %s: %v", file.Name(), err)
			continue
		}

		puck.Logger.Printf("Wrote %d bytes to %s", written, file.Name())
		if dryRun {
			fmt.Printf(": %s", strings.TrimPrefix(filename, string(conf.Runtime.WorkDir)))
		}
	}

	if dryRun {
		if err := os.RemoveAll(holosceneTitlesTempDir); err != nil {
			puck.Logger.Errorf("clearing temporary directory %s: %v", holosceneTitlesTempDir, err)
		}
	}
}
