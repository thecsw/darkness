package misaka

import (
	"bytes"
	"encoding/csv"
	"github.com/thecsw/rei"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/yunyun"
)

const (
	// reportDirectory is the directory where reports are stored.
	reportDirectory = ".darkness/reports/"
)

// WriteReport writes a report of the current run.
func WriteReport(conf *alpha.DarknessConfig) {
	start := time.Now()

	// See if the report directory exists.
	reportDir := conf.Runtime.WorkDir.Join(reportDirectory)
	if _, err := os.Stat(string(reportDir)); os.IsNotExist(err) {
		// Create the directory.
		err = os.Mkdir(string(reportDir), 0755)
		if err != nil {
			logger.Error("failed to create report directory", "err", err)
			return
		}
	}

	reportOutputFilename := conf.Runtime.WorkDir.Join(
		reportDirectory + yunyun.RelativePathFile(time.Now().Format(time.RFC3339)) + ".csv")

	reportOutputHandler, err := os.Create(filepath.Clean(string(reportOutputFilename)))
	if err != nil {
		logger.Error("failed to create report file", "err", err)
		return
	}
	defer func(reportOutputHandler *os.File) {
		err := reportOutputHandler.Close()
		if err != nil {
			logger.Error("failed to close report file", "err", err)
		}
	}(reportOutputHandler)

	written, err := io.Copy(reportOutputHandler, buildCSVReport(conf))
	if err != nil {
		logger.Error("failed to write report file", "err", err)
		return
	}
	if written == 0 {
		logger.Warn("no report written")
		return
	}

	logger.Warn(
		"Build report produced",
		"loc", conf.Runtime.WorkDir.Rel(reportOutputFilename),
		"elapsed", time.Since(start),
	)
}

// buildCSVReport builds a CSV report of the current run.
func buildCSVReport(conf *alpha.DarknessConfig) *bytes.Buffer {
	fullReport := GetFullReport()
	buf := &bytes.Buffer{}
	unit := ", μs"
	writer := csv.NewWriter(buf)
	rei.Try(writer.Write([]string{
		"№",
		"Input",
		"Output",
		"Read Tme" + unit,
		"Parse Time" + unit,
		"Export Time" + unit,
		"Write Time" + unit,
		"Total Time" + unit,
	}))
	num := 1
	for inputFile, report := range fullReport {
		readTime := int64(report[readIndex])
		parseTime := int64(report[parseIndex])
		exportTime := int64(report[exportIndex])
		writeTime := int64(report[writeIndex])
		totalTime := readTime + parseTime + exportTime + writeTime
		fullpath := yunyun.FullPathFile(conf.Project.InputFilenameToOutput(inputFile))
		rei.Try(writer.Write([]string{
			strconv.Itoa(num),
			string(conf.Runtime.WorkDir.Rel(inputFile)),
			string(conf.Runtime.WorkDir.Rel(fullpath)),
			humanize.Comma(readTime),
			humanize.Comma(parseTime),
			humanize.Comma(exportTime),
			humanize.Comma(writeTime),
			humanize.Comma(totalTime),
		}))
		num++
	}
	writer.Flush()
	return buf
}
