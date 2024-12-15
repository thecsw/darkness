package misa

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/emilia/rem"
	"github.com/thecsw/darkness/v3/emilia/reze"
	"github.com/thecsw/darkness/v3/ichika/hizuru"
	"github.com/thecsw/darkness/v3/ichika/kuroko"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
	"github.com/thecsw/rei"
)

const (
	// galleryPreviewImageSize is the size of the preview image.
	galleryPreviewImageSize = 250
	// galleryPreviewImageBlur is the blur radius of the preview image.
	galleryPreviewImageBlur = 37
	galleryJPEGQuality      = 90
)

// BuildGalleryFiles finds all the gallery entries and build a resized blurred
// preview version of it.
func BuildGalleryFiles(conf *alpha.DarknessConfig, dryRun bool) {
	// Make sure the preview directory exists
	previewDirectory := string(conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(conf.Project.DarknessPreviewDirectory)))
	if err := rei.Mkdir(previewDirectory); err != nil {
		puck.Logger.Fatalf("creating preview directory %s: %v", previewDirectory, err)
	}

	// Get all the gallery files.
	galleryFiles := getGalleryFiles(conf)

	// Filter out all the files that already exist.
	missingFiles := gana.Filter(func(item rem.GalleryItem) bool {
		return !rei.FileMustExist(string(rem.GalleryPreview(conf, item))) || kuroko.Force
	}, galleryFiles)

	// Build all the missing files.
	for i, galleryFile := range missingFiles {
		newFile := rem.GalleryPreview(conf, galleryFile)

		// Retrieve image contents reader:
		// - For local files, it's a reader of the file.
		// - For remote files, it's a reader of the response body,
		//   unless it's vendored, then it's a read of the vendored file.
		prefix := fmt.Sprintf("[%d/%d] ", i+1, len(missingFiles))
		sourceImage, err := rem.GalleryItemToImage(conf, galleryFile, "preview", prefix)
		if err != nil {
			puck.Logger.Warnf("\nopening a gallery image (%s): %v",
				filepath.Join(string(galleryFile.Path), string(galleryFile.Item)), err)
			continue
		}

		// Encode preview image into a buffer.
		previewImage := resizeAndBlur(sourceImage)

		// Don't save the file if it's in dry run mode.
		if !dryRun {
			file, err := os.Create(string(newFile))

			// Create a progress bar.
			bar := reze.ProgressBar(-1, "misa", prefix, "Resizing", string(conf.Runtime.WorkDir.Rel(newFile)))
			if err != nil {
				puck.Logger.Errorf("creating file %s: %v", newFile, err)
				continue
			}

			// Write the final preview image file.
			if err := imgio.JPEGEncoder(galleryJPEGQuality)(io.MultiWriter(file, bar), previewImage); err != nil {
				puck.Logger.Errorf("encoding image: %v", err)
				continue
			}

			// Close the file.
			if err := file.Close(); err != nil {
				puck.Logger.Errorf("closing image preview file %s: %v", newFile, err)
			}

			// Clear the progressbar.
			fmt.Print("\r\033[2K")
			// Log the thing.
			logger.Info("Resized item",
				"path", conf.Runtime.WorkDir.Rel(newFile),
				"dir", galleryFile.Path)
			rei.Try(bar.Close())
		}
	}
}

// resizeAndBlur takes an image object and modifies it to preview standards.
func resizeAndBlur(img image.Image) *image.RGBA {
	// Resize the image to save up on storage.
	newWidth := galleryPreviewImageSize
	newHeight := reze.PreserveImageHeightRatio(img, newWidth)
	resized := transform.Resize(img, newWidth, newHeight, transform.Lanczos)
	// Blur the image to make it look better.
	blurred := blur.Gaussian(resized, galleryPreviewImageBlur)
	return blurred

}

func dryRemove(val string) error {
	return nil
}

// RemoveGalleryFiles removes all generate gallery previews.
func RemoveGalleryFiles(conf *alpha.DarknessConfig, dryRun bool) {
	removeFunc := os.Remove
	if dryRun {
		removeFunc = dryRemove
	}
	for _, galleryFile := range getGalleryFiles(conf) {
		newFile := rem.GalleryPreview(conf, galleryFile)
		if err := removeFunc(string(newFile)); err != nil && !os.IsNotExist(err) {
			puck.Logger.Errorf("deleting %s: %v", newFile, err)
		}
	}
}

// getGalleryFiles returns a slice of all gallery images represented as `rem.GalleryItem`.
func getGalleryFiles(conf *alpha.DarknessConfig) []rem.GalleryItem {
	galleryFiles := make([]rem.GalleryItem, 0, 32)
	for _, page := range hizuru.BuildPagesSimple(conf, nil) {
		for _, gc := range page.Contents.Galleries() {
			for _, item := range gc.List {
				if strings.Contains(item.Text, ":no-preview") {
					continue
				}
				galleryFiles = append(galleryFiles, rem.NewGalleryItem(page, gc, item.Text))
			}
		}
	}
	return galleryFiles
}
