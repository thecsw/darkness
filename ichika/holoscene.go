package ichika

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia"
)

const (
	holosceneTitlesTempDir = "temp-holoscene"
)

func updateHolosceneTitles(dryRun bool) {
	if dryRun {
		if err := os.Mkdir(holosceneTitlesTempDir, 0755); err != nil {
			fmt.Printf("Failed to create temp dir: %s", err)
			os.Exit(1)
		}
	}

	inputs := emilia.FindFilesByExtSimple(emilia.Config.Project.Input)
	outputs := make([]string, len(inputs))
	for i, v := range inputs {
		outputs[i] = emilia.InputFilenameToOutput(v)
	}

	actuallyFound := make([]*os.File, 0, len(outputs))
	for _, v := range outputs {
		file, err := os.Open(v)
		if err != nil {
			fmt.Printf("Couldn't open %s: %s\n", v, err)
			continue
		}
		actuallyFound = append(actuallyFound, file)
	}

	fmt.Printf("Adding holoscene titles to %d output files\n",
		len(actuallyFound))

	for _, foundOutput := range actuallyFound {
		filename := foundOutput.Name()
		output, err := io.ReadAll(foundOutput)
		if err := foundOutput.Close(); err != nil {
			fmt.Printf("Failed to close %s: %s\n",
				filename, err)
		}
		if err != nil {
			fmt.Printf("Couldn't read %s: %s\n",
				filename, err)
			continue
		}

		newOutput := emilia.AddHolosceneTitles(string(output), -1)
		var file *os.File
		if dryRun {
			file, err = os.CreateTemp(holosceneTitlesTempDir,
				filepath.Base(filename))
		} else {
			file, err = os.Create(filename)
		}
		if err != nil {
			fmt.Printf("Failed to overwrite %s: %s\n",
				filename, err)
			continue
		}

		written, err := io.Copy(file, strings.NewReader(newOutput))
		if err := file.Close(); err != nil {
			fmt.Printf("Failed to close (2) %s: %s\n",
				file.Name(), err)
		}
		if err != nil {
			fmt.Printf("Failed to write %s: %s\n", file.Name(), err)
			continue
		}

		fmt.Printf("Wrote %d bytes to %s", written, file.Name())
		if dryRun {
			fmt.Printf(": %s",
				strings.TrimPrefix(filename, emilia.Config.WorkDir))
		}
		fmt.Println()
	}

	if dryRun {
		if err := os.RemoveAll(holosceneTitlesTempDir); err != nil {
			fmt.Printf("Failed to clear temp dir: %s\n", err)
		}
	}
}
