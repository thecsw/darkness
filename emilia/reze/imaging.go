package reze

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thecsw/haruhi"
	"github.com/thecsw/rei"
)

// OpenImage opens local path image and returns decoded image.
func OpenImage(path string) (image.Image, error) {
	if strings.HasSuffix(path, ".gif") {
		return nil, fmt.Errorf("GIF images are not supported")
	}
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("opening image file %s: %v", path, err)
	}
	// Respect the EXIF orientation flags.
	img, _, err := image.Decode(file)
	return img, err
}

var (
	// vendorClient is a client that is used to download images from the internet.
	vendorClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			MaxConnsPerHost:     1,
		},
		Timeout: 10 * time.Second,
	}
)

// DownloadImage attempts to download an image and returns it
// with any fatal errors (if occured).
func DownloadImage(link string, authority, prefix, name string) (image.Image, error) {
	resp, cancel, err := haruhi.URL(link).Client(vendorClient).Response()
	defer cancel()
	if err != nil {
		return nil, fmt.Errorf("downloading image: %v", err)
	}
	// If we got not found or server issue, bail.
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("downloading image got bad status %d", resp.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Warn("closing response body", "name", name, "err", err)
		}
	}(resp.Body)

	buf := new(bytes.Buffer)
	bar := ProgressBar(resp.ContentLength, authority, prefix, "Downloading", name)
	if _, err := io.Copy(io.MultiWriter(buf, bar), resp.Body); err != nil && err != io.EOF {
		return nil, fmt.Errorf("reading downloaded image data: %v", err)
	}

	// Attempt to decode.
	img, _, err := image.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("decoding downloaded image: %v", err)
	}

	fmt.Print("\r\033[2K")
	rei.Try(bar.Close())
	return img, nil
}

// PreserveImageHeightRatio calculates the new height of an image,
// given the new width, while preserving the original height ratio.
func PreserveImageHeightRatio(img image.Image, newWidth int) (height int) {
	if img == nil {
		return 0
	}
	// Calculate the new width and height.
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()
	tmpH := float64(newWidth) * float64(originalHeight) / float64(originalWidth)
	return int(math.Max(1.0, math.Floor(tmpH+0.5)))
}
