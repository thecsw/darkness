package misaka

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/yunyun"
)

const (
	// reportDirectory is the directory where reports are stored.
	reportDirectory = ".darkness/reports/"
)

// WriteReport writes a report of the current run.
func WriteReport(conf *alpha.DarknessConfig) {
	// See if the report directory exists.
	reportDir := conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(reportDirectory))
	if _, err := os.Stat(string(reportDir)); os.IsNotExist(err) {
		// Create the directory.
		err = os.Mkdir(string(reportDir), 0755)
		if err != nil {
			puck.Logger.Error("failed to create report directory", "err", err)
			return
		}
	}

	reportOutputFilename := conf.Runtime.WorkDir.Join(
		reportDirectory + yunyun.RelativePathFile(time.Now().Format(time.RFC3339)) + ".csv")

	reportOutputHandler, err := os.Create(filepath.Clean(string(reportOutputFilename)))
	if err != nil {
		puck.Logger.Error("failed to create report file", "err", err)
		return
	}
	defer func(reportOutputHandler *os.File) {
		err := reportOutputHandler.Close()
		if err != nil {
			puck.Logger.Error("failed to close report file", "err", err)
		}
	}(reportOutputHandler)

	written, err := io.Copy(reportOutputHandler, buildCSVReport(conf))
	if err != nil {
		puck.Logger.Error("failed to write report file", "err", err)
		return
	}
	if written == 0 {
		puck.Logger.Warn("no report written")
		return
	}
	puck.Logger.Warn("Build report produced", "file", conf.Runtime.WorkDir.Rel(yunyun.FullPathFile(reportOutputFilename)))
}

// buildCSVReport builds a CSV report of the current run.
func buildCSVReport(conf *alpha.DarknessConfig) *bytes.Buffer {
	fullReport := GetFullReport()
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	writer.Write([]string{"#", "Input", "Output", "Read Tme", "Parse Time", "Export Time", "Write Time", "Total Time"})
	num := 1
	for inputFile, report := range fullReport {
		readTime := int(report[readIndex])
		parseTime := int(report[parseIndex])
		exportTime := int(report[exportIndex])
		writeTime := int(report[writeIndex])
		totalTime := readTime + parseTime + exportTime + writeTime
		writer.Write([]string{
			strconv.Itoa(num),
			string(conf.Runtime.WorkDir.Rel(inputFile)),
			string(conf.Runtime.WorkDir.Rel(yunyun.FullPathFile(conf.Project.InputFilenameToOutput(inputFile)))),
			strconv.Itoa(readTime) + "μs",
			strconv.Itoa(parseTime) + "μs",
			strconv.Itoa(exportTime) + "μs",
			strconv.Itoa(writeTime) + "μs",
			strconv.Itoa(totalTime) + "μs",
		})
		num++
	}
	writer.Flush()
	return buf
}
