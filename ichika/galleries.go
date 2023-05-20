package ichika

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/gana"
	"github.com/thecsw/rei"
)

const (
	// galleryPreviewImageSize is the size of the preview image.
	galleryPreviewImageSize = 250
	// galleryPreviewImageBlur is the blur radius of the preview image.
	galleryPreviewImageBlur = 20
)

// buildGalleryFiles finds all the gallery entries and build a resized
// preview version of it.
func buildGalleryFiles(dryRun bool) {
	// Make sure the preview directory exists
	previewDirectory := filepath.Join(emilia.Config.WorkDir, string(emilia.Config.Project.DarknessPreviewDirectory))
	if err := rei.Mkdir(previewDirectory); err != nil {
		puck.Logger.Fatalf("creating preview directory %s: %v", previewDirectory, err)
	}

	// Get all the gallery files.
	galleryFiles := getGalleryFiles()

	// Filter out all the files that already exist.
	missingFiles := gana.Filter(func(item *emilia.GalleryItem) bool {
		return !rei.FileMustExist(string(emilia.GalleryPreview(item)))
	}, galleryFiles)

	// Build all the missing files.
	for i, galleryFile := range missingFiles {
		newFile := emilia.GalleryPreview(galleryFile)

		// Retrieve image contents reader:
		// - For local files, it's a reader of the file.
		// - For remote files, it's a reader of the response body,
		//   unless it's vendored, then it's a read of the vendored file.
		prefix := fmt.Sprintf("[%d/%d] ", i+1, len(missingFiles))
		sourceImage, err := emilia.GalleryItemToImage(galleryFile, "preview", prefix)
		if err != nil {
			puck.Logger.Errorf("parsing a gallery item: %v", err)
			continue
		}

		// Encode preview image into a buffer.
		previewImage := resizeAndBlur(sourceImage)

		// Don't save the file if it's in dry run mode.
		if !dryRun {
			file, err := os.Create(string(newFile))

			// Create a progress bar.
			bar := emilia.ProgressBar(-1, "misa", prefix, "Resizing", string(emilia.FullPathToWorkDirRel(newFile)))
			if err != nil {
				puck.Logger.Errorf("creating file %s: %v", newFile, err)
				continue
			}

			// Write the final preview image file.
			if err := imaging.Encode(io.MultiWriter(file, bar), previewImage, imaging.JPEG); err != nil {
				puck.Logger.Errorf("encoding image: %v", err)
				continue
			}

			// Close the file.
			if err := file.Close(); err != nil {
				puck.Logger.Errorf("closing image preview file %s: %v", newFile, err)
			}
		}
	}
	fmt.Print("\r\033[2K")
}

// resizeAndBlur takes an image object and modifies it to preview standards.
func resizeAndBlur(img image.Image) *image.NRGBA {
	// Resize the image to save up on storage.
	img = imaging.Resize(img, galleryPreviewImageSize, 0, imaging.Lanczos)
	// Blur the image to make it look better.
	blurred := imaging.Blur(img, galleryPreviewImageBlur)
	return blurred
}

func dryRemove(val string) error {
	return nil
}

// removeGalleryFiles removes all generate gallery previews.
func removeGalleryFiles(dryRun bool) {
	removeFunc := os.Remove
	if dryRun {
		removeFunc = dryRemove
	}
	for _, galleryFile := range getGalleryFiles() {
		newFile := emilia.GalleryPreview(galleryFile)
		if err := removeFunc(string(newFile)); err != nil && !os.IsNotExist(err) {
			puck.Logger.Errorf("deleting %s: %v", newFile, err)
		}
	}
}

// getGalleryFiles returns a slice of all gallery images represented as `emilia.GalleryItem`.
func getGalleryFiles() []*emilia.GalleryItem {
	galleryFiles := make([]*emilia.GalleryItem, 0, 32)
	for _, page := range buildPagesSimple(nil) {
		for _, gc := range page.Contents.Galleries() {
			for _, item := range gc.List {
				galleryFiles = append(galleryFiles, emilia.NewGalleryItem(page, gc, item))
			}
		}
	}

	return galleryFiles
}
