package ichika

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/narumi"
	"github.com/thecsw/darkness/emilia/puck"
)

const (
	holosceneTitlesTempDir = "temp-holoscene"
)

func updateHolosceneTitles(conf *alpha.DarknessConfig, dryRun bool) {
	if dryRun {
		if err := os.Mkdir(holosceneTitlesTempDir, 0o750); err != nil {
			puck.Logger.Fatalf("creating temporary directory %s: %v", holosceneTitlesTempDir, err)
		}
	}

	inputFilenames := FindFilesByExtSimple(conf)
	outputs := make([]string, len(inputFilenames))
	for i, inputFilename := range inputFilenames {
		outputs[i] = conf.Project.InputFilenameToOutput(inputFilename)
	}

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

	fmt.Printf("Adding holoscene titles to %d output files\n", len(actuallyFound))

	for _, foundOutput := range actuallyFound {
		filename := foundOutput.Name()
		output, err := io.ReadAll(foundOutput)
		if err := foundOutput.Close(); err != nil {
			puck.Logger.Errorf("closing file %s: %v", filename, err)
		}
		if err != nil {
			puck.Logger.Errorf("reading file %s: %v", filename, err)
			continue
		}

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
