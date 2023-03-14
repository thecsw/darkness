package ichika

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// writeFile takes a filename and a bufio reader and writes it.
func writeFile(filename string, reader *bufio.Reader) (int64, error) {
	target, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return -1, errors.Wrap(err, "failed to create "+filename)
	}
	written, err := io.Copy(target, reader)
	if err != nil {
		return -1, errors.Wrap(err, "failed to copy to "+filename)
	}
	if target.Close() != nil {
		return -1, errors.Wrap(err, "failed to close "+filename)
	}
	return written, nil
}

// openFile attemps to open the full path and return tuple, empty tuple otherwise.
func openFile(v yunyun.FullPathFile) gana.Tuple[yunyun.FullPathFile, *os.File] {
	file, err := os.Open(filepath.Clean(string(v)))
	if err != nil {
		log.Printf("failed to open %s: %s\n", v, err)
		return gana.NewTuple[yunyun.FullPathFile, *os.File]("", nil)
	}
	return gana.NewTuple(v, file)
}
