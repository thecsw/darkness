package alpha

import (
	"strings"

	"github.com/thecsw/darkness/v3/yunyun"
)

const (
	// We flush pages in a json format.
	debugStructExtension = ".json"
)

// InputFilenameToOutput converts input filename to the filename to write.
func (p ProjectConfig) InputFilenameToOutput(file yunyun.FullPathFile) string {
	return strings.Replace(string(file), p.Input, p.Output, 1)
}

// InputFilenameToDebugStruct flushes pages as json for debugging.
func (p ProjectConfig) InputFilenameToDebugStruct(file yunyun.FullPathFile) string {
	return strings.Replace(string(file), p.Input, debugStructExtension, 1)
}
