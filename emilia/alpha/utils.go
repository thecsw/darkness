package alpha

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// InputFilenameToOutput converts input filename to the filename to write.
func (p ProjectConfig) InputFilenameToOutput(file yunyun.FullPathFile) string {
	return strings.Replace(string(file), p.Input, p.Output, 1)
}
