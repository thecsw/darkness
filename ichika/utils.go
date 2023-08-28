package ichika

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
	"github.com/thecsw/gana/prelude"
)

// writeFile takes a filename and a bufio reader and writes it.
func writeFile(filename string, reader io.Reader) (int64, error) {
	filename = filepath.Clean(filename)
	target, err := os.Create(filename)
	if err != nil {
		return -1, fmt.Errorf("creating file %s: %v", filename, err)
	}
	written, err := io.Copy(target, reader)
	if err != nil {
		return -1, fmt.Errorf("copying to file %s: %v", filename, err)
	}
	if target.Close() != nil {
		return -1, fmt.Errorf("closing file %s: %v", filename, err)
	}
	return written, nil
}

// openFile attemps to open the full path and return tuple, empty tuple otherwise.
func openFile(v yunyun.FullPathFile) prelude.Option[gana.Tuple[yunyun.FullPathFile, *os.File]] {
	file, err := os.Open(filepath.Clean(string(v)))
	if err != nil {
		log.Printf("failed to open %s: %s\n", v, err)
		return prelude.None[gana.Tuple[yunyun.FullPathFile, *os.File]]()
	}
	return prelude.Some(gana.NewTuple(v, file))
}
