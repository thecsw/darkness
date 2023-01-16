package emilia

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
)

// OpenImage opens local path image and returns decoded image.
func OpenImage(path string) (image.Image, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "OpenImage: opening file "+path)
	}
	// Respect the EXIF orientation flags.
	return imaging.Decode(file, imaging.AutoOrientation(true))
}

// DownloadImage attempts to download an image and returns it
// with any fatal errors (if occured).
func DownloadImage(link string) (image.Image, error) {
	// Build the request.
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadImage: create request")
	}

	// Attempt to make the request.
	resp, err := vendorClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "DownloadImage: Do request")
	}

	// If we got not found or server issue, bail.
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, errors.Wrap(err,
			fmt.Sprintf("DownloadImage: Bad status: %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	// Attempt to decode.
	img, err := imaging.Decode(resp.Body, imaging.AutoOrientation(true))
	if err != nil {
		return nil, errors.Wrap(err, "DownloadImage: failed to decode")
	}

	return img, nil
}
