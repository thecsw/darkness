package emilia

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// Pack cleans the filename from absolute workspace prefix.
func Pack(filename yunyun.FullPathFile, data string) (yunyun.RelativePathFile, string) {
	return RelPathToWorkdir(filename), data
}

// PackRef cleans the filename from absolute workspace prefix by refs.
func PackRef(filename *yunyun.FullPathFile, data *string) (yunyun.RelativePathFile, string) {
	return RelPathToWorkdir(*filename), *data
}

// relPathToWorkdir returns path trimmed by the workspace
func RelPathToWorkdir(filename yunyun.FullPathFile) yunyun.RelativePathFile {
	return yunyun.RelativePathFile(strings.TrimPrefix(string(filename), Config.WorkDir+`/`))
}

// sha256String hashes given string to sha256.
func sha256String(what string) string {
	ans := sha256.Sum256([]byte(what))
	return hex.EncodeToString(ans[:])
}
