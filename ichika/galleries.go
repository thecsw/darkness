package ichika

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
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
	previewDirectory := filepath.Join(emilia.Config.WorkDir, string(emilia.GalleryPreviewDirectory))
	fmt.Println("GALLERY:", previewDirectory)
	if err := emilia.Mkdir(previewDirectory); err != nil {
		fmt.Println("fatal: couldn't create preview directory:", err.Error())
		os.Exit(1)
	}
	galleryFiles := getGalleryFiles()
	for i, galleryFile := range galleryFiles {
		fmt.Printf("[%d/%d] ", i+1, len(galleryFiles))
		newFile := emilia.GalleryPreview(galleryFile)
		if emilia.FileExists(string(newFile)) {
			fmt.Printf("%s already exists\n", emilia.RelPathToWorkdir(newFile))
			continue
		}
		fmt.Printf("%s... ", emilia.RelPathToWorkdir(newFile))

		// Mark the processing start time.
		start := time.Now()

		// Retrieve image contents reader:
		// - For local files, it's a reader of the file.
		// - For remote files, it's a reader of the response body.
		imgReader, err := galleryItemToReader(galleryFile)
		if err != nil {
			fmt.Println("gallery item to reader:", err.Error())
			continue
		}

		// Encode preview image into a buffer.
		img, err := readAndBlur(imgReader)
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
			if err := imaging.Encode(file, img, imaging.JPEG); err != nil {
				fmt.Printf("failed to encode image: %s", err.Error())
				continue
			}
			file.Close()
		}

		fmt.Printf("%d ms\n", time.Since(start).Milliseconds())
	}
}

// galleryItemToReader takes in a gallery item and returns an `io.ReadCloser`
// for the image's contents.
func galleryItemToReader(item *emilia.GalleryItem) (io.ReadCloser, error) {
	// If it's a local file, simply open the os file.
	if !item.IsExternal {
		file := emilia.JoinWorkdir(yunyun.JoinRelativePaths(item.Path, item.Item))
		return os.Open(string(file))
	}
	// Check if the item has been vendored by any chance?
	vendorPath := filepath.Join(emilia.Config.WorkDir, string(emilia.GalleryVendored(item)))
	if emilia.FileExists(vendorPath) {
		fmt.Printf(" (vendored) ")
		return os.Open(vendorPath)
	}
	// If it's a remote file, then run a get request and return
	// the body reader.
	resp, err := http.Get(string(item.Item))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't retrieve external gallery item ("+string(item.Item)+")")
	}
	return resp.Body, nil
}

// readAndBlur decodes image from `source`, makes a preview out of it,
// and finally encodes it into `target`.
func readAndBlur(source io.ReadCloser) (*image.NRGBA, error) {
	// Respect EXIF flags with AutoOrientation turned on.
	img, err := imaging.Decode(source, imaging.AutoOrientation(true))
	defer source.Close()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read the image")
	}
	// Resize the image to save up on storage.
	img = imaging.Resize(img, galleryPreviewImageSize, 0, imaging.Lanczos)
	blurred := imaging.Blur(img, galleryPreviewImageBlur)
	return blurred, nil
}

// removeGalleryFiles removes all generate gallery previews.
func removeGalleryFiles() {
	for _, galleryFile := range getGalleryFiles() {
		newFile := emilia.GalleryPreview(galleryFile)
		if err := os.Remove(string(newFile)); err != nil && !os.IsNotExist(err) {
			fmt.Println("Couldn't delete", newFile, "| reason:", err.Error())
		}
	}
}

// getGalleryFiles returns a slice of all gallery images represented as `emilia.GalleryItem`.
func getGalleryFiles() []*emilia.GalleryItem {
	inputFilenames := make(chan yunyun.FullPathFile, customChannelCapacity)
	pages := gana.GenericWorkers(gana.GenericWorkers(inputFilenames,
		func(v yunyun.FullPathFile) gana.Tuple[yunyun.FullPathFile, string] {
			data, err := ioutil.ReadFile(filepath.Clean(string(v)))
			if err != nil {
				fmt.Printf("Failed to open %s: %s\n", v, err.Error())
			}
			return gana.NewTuple(v, string(data))
		}, 1, customChannelCapacity), func(v gana.Tuple[yunyun.FullPathFile, string]) *yunyun.Page {
		return emilia.ParserBuilder.BuildParser(emilia.PackRef(v.UnpackRef())).Parse()
	}, customNumWorkers, customChannelCapacity)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go emilia.FindFilesByExt(inputFilenames, emilia.Config.Project.Input, wg)

	// Launch a second discovery for gallery files.
	galleryFiles := make([]*emilia.GalleryItem, 0, 32)
	go func(wg *sync.WaitGroup) {
		for page := range pages {
			for _, gc := range page.Contents.Galleries() {
				for _, item := range gc.List {
					galleryFiles = append(galleryFiles, emilia.NewGalleryItem(page, gc, item))
				}
			}
			wg.Done()
		}
		wg.Done()
	}(wg)

	wg.Wait()
	return galleryFiles
}
