package emilia

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// Pack cleans the filename from absolute workspace prefix.
func Pack(filename yunyun.FullPathFile, data string) (yunyun.RelativePathFile, string) {
	return FullPathToWorkDirRel(filename), data
}

// PackRef cleans the filename from absolute workspace prefix by refs.
func PackRef(filename *yunyun.FullPathFile, data *string) (yunyun.RelativePathFile, string) {
	return FullPathToWorkDirRel(*filename), *data
}

// relPathToWorkdir returns path trimmed by the workspace
func FullPathToWorkDirRel(filename yunyun.FullPathFile) yunyun.RelativePathFile {
	return yunyun.RelativePathFile(strings.TrimPrefix(string(filename), Config.WorkDir+`/`))
}
