package misaka

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/thecsw/darkness/v3/ichika/kuroko"

	"github.com/thecsw/darkness/v3/yunyun"
)

var (
	recordedFiles        = sync.Map{}
	recordedFilesCounter = atomic.Int32{}

	readTimes   = sync.Map{}
	parseTimes  = sync.Map{}
	exportTimes = sync.Map{}
	writeTimes  = sync.Map{}
)

const (
	readIndex = iota
	parseIndex
	exportIndex
	writeIndex
)

// RecordReadTime records the time it took to read a file.
//
//go:inline
func RecordReadTime(inputFile yunyun.FullPathFile, duration time.Duration) {
	recordTime(inputFile, duration, &readTimes)
}

// RecordParseTime records the time it took to parse a file.
//
//go:inline
func RecordParseTime(inputFile yunyun.FullPathFile, duration time.Duration) {
	recordTime(inputFile, duration, &parseTimes)
}

// RecordExportTime records the time it took to export a file.
//
//go:inline
func RecordExportTime(inputFile yunyun.FullPathFile, duration time.Duration) {
	recordTime(inputFile, duration, &exportTimes)
}

// RecordWriteTime records the time it took to write a file.
//
//go:inline
func RecordWriteTime(inputFile yunyun.FullPathFile, duration time.Duration) {
	recordTime(inputFile, duration, &writeTimes)
}

// recordTime records the time it took to do something in its respective sync.Map.
//
//go:inline
func recordTime(inputFile yunyun.FullPathFile, duration time.Duration, times *sync.Map) {
	// Only record the file if we are recording build reports.
	if kuroko.BuildReport {
		times.Store(inputFile, duration.Microseconds())
		recordedFiles.Store(inputFile, true)
		recordedFilesCounter.Add(1)
	}
}

// GetNumberReports returns the number of files that have been recorded.
func GetNumberReports() int {
	return int(recordedFilesCounter.Load())
}

// GetFullReport returns a map of all the files and their times.
func GetFullReport() map[yunyun.FullPathFile][]int64 {
	fullReport := make(map[yunyun.FullPathFile][]int64)
	recordedFiles.Range(func(key, value any) bool {
		inputFile := key.(yunyun.FullPathFile)
		fullReport[inputFile] = make([]int64, 4)
		loadIntoFullReport(inputFile, fullReport, &readTimes, readIndex)
		loadIntoFullReport(inputFile, fullReport, &parseTimes, parseIndex)
		loadIntoFullReport(inputFile, fullReport, &exportTimes, exportIndex)
		loadIntoFullReport(inputFile, fullReport, &writeTimes, writeIndex)
		// Signal to continue.
		return true
	})
	return fullReport
}

// loadIntoFullReport loads the time into the full report by given index.
func loadIntoFullReport(
	inputFile yunyun.FullPathFile,
	fullReport map[yunyun.FullPathFile][]int64,
	times *sync.Map,
	index int,
) {
	whateverTime, ok := times.Load(inputFile)
	if ok {
		fullReport[inputFile][index] = whateverTime.(int64)
	}
}
