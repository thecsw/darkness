package misaka

import (
	"sync"
	"time"

	"github.com/thecsw/darkness/yunyun"
)

var (
	recordedFiles = sync.Map{}

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
func RecordReadTime(inputFile yunyun.FullPathFile, duration time.Duration) {
	recordTime(inputFile, duration, &readTimes)
}

// RecordParseTime records the time it took to parse a file.
func RecordParseTime(inputFile yunyun.FullPathFile, duration time.Duration) {
	recordTime(inputFile, duration, &parseTimes)
}

// RecordExportTime records the time it took to export a file.
func RecordExportTime(inputFile yunyun.FullPathFile, duration time.Duration) {
	recordTime(inputFile, duration, &exportTimes)
}

// RecordWriteTime records the time it took to write a file.
func RecordWriteTime(inputFile yunyun.FullPathFile, duration time.Duration) {
	recordTime(inputFile, duration, &writeTimes)
}

// recordTime records the time it took to do something in its respective sync.Map.
func recordTime(inputFile yunyun.FullPathFile, duration time.Duration, times *sync.Map) {
	times.Store(inputFile, duration.Microseconds())
	recordedFiles.Store(inputFile, true)
}

// GetNumberReports returns the number of files that have been recorded.
func GetNumberReports() int {
	total := 0
	recordedFiles.Range(func(key, value interface{}) bool {
		total++
		return true
	})
	return total
}

// GetFullReport returns a map of all the files and their times.
func GetFullReport() map[yunyun.FullPathFile][]int64 {
	fullReport := make(map[yunyun.FullPathFile][]int64)
	recordedFiles.Range(func(key, value interface{}) bool {
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
