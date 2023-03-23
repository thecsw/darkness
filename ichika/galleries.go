package ichika

import (
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/gana"
)

const (
	galleryPreviewImageSize = 250
	galleryPreviewImageBlur = 20
)

// buildGalleryFiles finds all the gallery entries and build a resized
// preview version of it.
func buildGalleryFiles(dryRun bool) {
	// Make sure the preview directory exists
	previewDirectory := filepath.Join(emilia.Config.WorkDir, string(emilia.Config.Project.DarknessPreviewDirectory))
	if err := emilia.Mkdir(previewDirectory); err != nil {
		fmt.Println("fatal: couldn't create preview directory:", err)
		os.Exit(1)
	}
	galleryFiles := getGalleryFiles()

	missingFiles := gana.Filter(func(item *emilia.GalleryItem) bool {
		return !emilia.FileExists(string(emilia.GalleryPreview(item)))
	}, galleryFiles)

	for i, galleryFile := range missingFiles {
		newFile := emilia.GalleryPreview(galleryFile)

		// Retrieve image contents reader:
		// - For local files, it's a reader of the file.
		// - For remote files, it's a reader of the response body,
		//   unless it's vendored, then it's a read of the vendored file.
		prefix := fmt.Sprintf("[%d/%d] ", i+1, len(missingFiles))
		sourceImage, err := emilia.GalleryItemToImage(galleryFile, "preview", prefix)
		if err != nil {
			fmt.Println("gallery item to reader:", err)
			continue
		}

		// Encode preview image into a buffer.
		previewImage, err := resizeAndBlur(sourceImage)
		if err != nil {
			fmt.Println("gallery reader to writer:", err)
			continue
		}

		// Don't save the file if it's in dry run mode.
		if !dryRun {
			file, err := os.Create(string(newFile))
			bar := emilia.ProgressBar(-1, "misa", prefix, "Resizing", string(emilia.FullPathToWorkDirRel(newFile)))
			if err != nil {
				fmt.Printf("failed to create file %s: %s\n", newFile, err)
				continue
			}
			// Write the final preview image file.
			if err := imaging.Encode(io.MultiWriter(file, bar), previewImage, imaging.JPEG); err != nil {
				fmt.Printf("failed to encode image: %s", err)
				continue
			}
			if err := file.Close(); err != nil {
				fmt.Printf("failed to close image preview file %s: %s\n", newFile, err)
			}
		}
	}
	fmt.Print("\r\033[2K")
}

// resizeAndBlur takes an image object and modifies it to preview standards.
func resizeAndBlur(img image.Image) (*image.NRGBA, error) {
	// Resize the image to save up on storage.
	img = imaging.Resize(img, galleryPreviewImageSize, 0, imaging.Lanczos)
	blurred := imaging.Blur(img, galleryPreviewImageBlur)
	return blurred, nil
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
			fmt.Println("Couldn't delete", newFile, "| reason:", err)
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
