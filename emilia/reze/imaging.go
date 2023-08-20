package reze

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/thecsw/haruhi"
)

// OpenImage opens local path image and returns decoded image.
func OpenImage(path string) (image.Image, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("opening image file %s: %v", path, err)
	}
	// Respect the EXIF orientation flags.
	return imaging.Decode(file, imaging.AutoOrientation(true))
}

// DownloadImage attempts to download an image and returns it
// with any fatal errors (if occured).
func DownloadImage(link string, authority, prefix, name string) (image.Image, error) {
	resp, cancel, err := haruhi.URL(link).Response()
	defer cancel()
	if err != nil {
		return nil, fmt.Errorf("downloading image: %v", err)
	}
	// If we got not found or server issue, bail.
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("downloading image got bad status %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	bar := ProgressBar(resp.ContentLength, authority, prefix, "Downloading", name)
	if _, err := io.Copy(io.MultiWriter(buf, bar), resp.Body); err != nil && err != io.EOF {
		return nil, fmt.Errorf("reading downloaded image data: %v", err)
	}

	// Attempt to decode.
	img, err := imaging.Decode(buf, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("decoding downloaded image: %v", err)
	}

	return img, nil
}
