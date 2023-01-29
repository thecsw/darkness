package ichika

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/thecsw/darkness/emilia"
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
		fmt.Println("fatal: couldn't create preview directory:", err.Error())
		os.Exit(1)
	}
	galleryFiles := getGalleryFiles()
	for i, galleryFile := range galleryFiles {
		fmt.Printf("[%d/%d] ", i+1, len(galleryFiles))
		newFile := emilia.GalleryPreview(galleryFile)
		if emilia.FileExists(string(newFile)) {
			fmt.Printf("%s already exists\n", emilia.FullPathToWorkDirRel(newFile))
			continue
		}
		fmt.Printf("%s... ", emilia.FullPathToWorkDirRel(newFile))

		// Mark the processing start time.
		start := time.Now()

		// Retrieve image contents reader:
		// - For local files, it's a reader of the file.
		// - For remote files, it's a reader of the response body,
		//   unless it's vendored, then it's a read of the vendored file.
		sourceImage, err := emilia.GalleryItemToImage(galleryFile)
		if err != nil {
			fmt.Println("gallery item to reader:", err.Error())
			continue
		}

		// Encode preview image into a buffer.
		previewImage, err := resizeAndBlur(sourceImage)
		if err != nil {
			fmt.Println("gallery reader to writer:", err.Error())
			continue
		}

		// Don't save the file if it's in dry run mode.
		if !dryRun {
			file, err := os.Create(string(newFile))
			if err != nil {
				fmt.Printf("failed to create file %s: %s\n", newFile, err.Error())
				continue
			}
			// Write the final preview image file.
			if err := imaging.Encode(file, previewImage, imaging.JPEG); err != nil {
				fmt.Printf("failed to encode image: %s", err.Error())
				continue
			}
			file.Close()
		}

		fmt.Printf("%d ms\n", time.Since(start).Milliseconds())
	}
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
			fmt.Println("Couldn't delete", newFile, "| reason:", err.Error())
		}
	}
}

// getGalleryFiles returns a slice of all gallery images represented as `emilia.GalleryItem`.
func getGalleryFiles() []*emilia.GalleryItem {
	galleryFiles := make([]*emilia.GalleryItem, 0, 32)
	for _, page := range buildPagesSimple() {
		for _, gc := range page.Contents.Galleries() {
			for _, item := range gc.List {
				galleryFiles = append(galleryFiles, emilia.NewGalleryItem(page, gc, item))
			}
		}
	}

	return galleryFiles
}
