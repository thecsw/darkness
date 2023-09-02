package hizuru

import (
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
	"github.com/thecsw/gana/prelude"
	"os"
	"path/filepath"
)

// openFile attemps to open the full path and return tuple, empty tuple otherwise.
func openFile(v yunyun.FullPathFile) prelude.Option[gana.Tuple[yunyun.FullPathFile, *os.File]] {
	file, err := os.Open(filepath.Clean(string(v)))
	if err != nil {
		//log.Printf("failed to open %s: %s\n", v, err)
		return prelude.None[gana.Tuple[yunyun.FullPathFile, *os.File]]()
	}
	return prelude.Some(gana.NewTuple(v, file))
}
