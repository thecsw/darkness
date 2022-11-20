package emilia

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func Untar(reader io.Reader, dirName string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		path, err := SanitizeArchivePath(dirName, header.Name)
		if err != nil {
			return errors.Wrap(err, "untarring")
		}
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(
			filepath.Clean(path),
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			info.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Printf("failed to close file in untar: %s\n", err.Error())
			}
		}()
		for {
			_, err = io.CopyN(file, tarReader, 1024)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}
	}
	return nil
}

// Sanitize archive file pathing from "G305: Zip Slip vulnerability".
// Found in https://github.com/securego/gosec/issues/324
func SanitizeArchivePath(d, t string) (v string, err error) {
	v = filepath.Join(d, t)
	if strings.HasPrefix(v, filepath.Clean(d)) {
		return v, nil
	}
	return "", fmt.Errorf("%s: %s", "content filepath is tainted", t)
}
