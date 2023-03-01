package emilia

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/k0kubun/go-ansi"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"github.com/thecsw/haruhi"
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
func DownloadImage(link string, authority, prefix, name string) (image.Image, error) {
	resp, cancel, err := haruhi.URL(link).Response()
	defer cancel()
	if err != nil {
		return nil, err
	}
	// If we got not found or server issue, bail.
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, errors.Wrap(err,
			fmt.Sprintf("DownloadImage: Bad status: %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	bar := ProgressBar(resp.ContentLength, authority, prefix, "Downloading", name)
	if _, err := io.Copy(io.MultiWriter(buf, bar), resp.Body); err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "failed to read the image data")
	}

	// Attempt to decode.
	img, err := imaging.Decode(buf, imaging.AutoOrientation(true))
	if err != nil {
		return nil, errors.Wrap(err, "DownloadImage: failed to decode")
	}

	return img, nil
}

// ProgressBar will return darkness style progress bar.
func ProgressBar(size int64, authority, prefix, action, name string) *progressbar.ProgressBar {
	return progressbar.NewOptions64(size,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetDescription(
			fmt.Sprintf("[cyan][%s][reset] %s%s %s", authority, prefix, action, name)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[magenta]=[reset]",
			SaucerHead:    "[yellow]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}
