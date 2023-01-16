package emilia

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/thecsw/darkness/yunyun"
)

// GalleryItem is a struct that holds the gallery item path
// and a flag whether it is external (URL regexp matches).
type GalleryItem struct {
	// Item is the link that was provided.
	Item yunyun.RelativePathFile
	// Path is the path of the local gallery source file.
	Path yunyun.RelativePathDir
	// IsExternal runs a URL regexp check.
	IsExternal bool
	// Text found through the link regexp.
	Text string
	// Description found through the link regexp.
	Description string
	// OriginalLine is the original line that include org options.
	OriginalLine string
	// Link is an optional parameter that the gallery item should
	// also link to something.
	Link string
}

// NewGalleryItem creates a new helper `GalleryItem` and
// decides whether the passed item is an external link or not.
func NewGalleryItem(page *yunyun.Page, content *yunyun.Content, wholeLine string) *GalleryItem {
	extractedLinks := yunyun.ExtractLinks(wholeLine)
	// If image wasn't found, then the whole line should be counted as the image path.
	image := wholeLine
	text := ""
	description := ""
	if len(extractedLinks) > 0 {
		image = extractedLinks[0].Link
		text = extractedLinks[0].Text
		description = extractedLinks[0].Description
	}
	optionalLink := ""
	if len(extractedLinks) > 1 {
		optionalLink = extractedLinks[1].Link
	}
	return &GalleryItem{
		Item:         yunyun.RelativePathFile(image),
		Path:         yunyun.JoinPaths(page.Location, content.GalleryPath),
		IsExternal:   yunyun.URLRegexp.MatchString(image),
		Text:         text,
		Description:  description,
		OriginalLine: wholeLine,
		Link:         optionalLink,
	}
}

func GalleryImage(item *GalleryItem) yunyun.FullPathFile {
	if item.IsExternal {
		// If it's vendored, then retrieve a local copy (if doesn't already
		// exist) and stub it in as the full path
		if Config.VendorGalleries {
			return galleryVendorItem(item)
		}
		return yunyun.FullPathFile(item.Item)
	}
	return JoinPath(yunyun.JoinRelativePaths(item.Path, item.Item))
}

// galleryPreviewRelative takes gallery item and returns relative path to it.
func galleryPreviewRelative(item *GalleryItem) yunyun.RelativePathFile {
	if item.IsExternal {
		return galleryItemHash(item)
	}
	filename := filepath.Base(string(item.Item))
	ext := filepath.Ext(filename)
	return yunyun.RelativePathFile(strings.TrimSuffix(filename, ext) + "_small.jpg")
}

// GalleryPreview takes an original image's path and returns
// the preview path of it. Previews are always .jpg
func GalleryPreview(item *GalleryItem) yunyun.FullPathFile {
	return JoinPath(yunyun.JoinRelativePaths(Config.Project.DarknessPreviewDirectory, galleryPreviewRelative(item)))
}

const (
	// defaultVendorDirectory is the name of the dir where vendor images are stored.
	defaultVendorDirectory yunyun.RelativePathDir = "darkness_vendor"
	// defaultPreviewDirectory is the name of the dir where all gallery previews are stored.
	defaultPreviewDirectory yunyun.RelativePathDir = "darkness_gallery_previews"
)

var (
	vendorClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			MaxConnsPerHost:     1,
		},
		Timeout: 10 * time.Second,
	}
	vendorLock = &sync.Mutex{}
)

// GalleryVendored returns vendored local path of the gallery item.
func GalleryVendored(item *GalleryItem) yunyun.RelativePathFile {
	return yunyun.JoinRelativePaths(Config.Project.DarknessVendorDirectory, galleryItemHash(item))
}

// galleryVendorItem vendors given item and returns the full path of the file.
//
// Only call this function on remote images, it's up to the user to make the
// .IsExternal check before calling this. SLOW function because of network calls.
//
// If the vendoring fails at any point, fallback to the remote image path.
func galleryVendorItem(item *GalleryItem) yunyun.FullPathFile {
	// Process only one vendor request at a time.
	vendorLock.Lock()
	// Unlock so the next vendor request can get processed.
	defer vendorLock.Unlock()

	vendoredImagePath := GalleryVendored(item)
	localVendoredPath := filepath.Join(Config.WorkDir, string(vendoredImagePath))

	// Create the two types of return.
	fallbackReturn := yunyun.FullPathFile(item.Item)
	expectedReturn := JoinPath(vendoredImagePath)

	// Check if the image was already vendored, if it was, return it immediately.
	if FileExists(localVendoredPath) {
		return expectedReturn
	}

	start := time.Now()
	fmt.Printf("Vendoring %s... ", vendoredImagePath)

	img, err := DownloadImage(string(item.Item))
	if err != nil {
		fmt.Printf("failed to vendor: %s\n", err.Error())
		return fallbackReturn
	}

	// Open the file writer and encode the image there.
	imgFile, err := os.Create(filepath.Clean(localVendoredPath))
	if err != nil {
		fmt.Printf("Failed to create file %s: %s\n", localVendoredPath, err.Error())
		return fallbackReturn
	}
	defer imgFile.Close()

	// Decode the image into the file.
	if err := imaging.Encode(imgFile, img, imaging.JPEG); err != nil {
		fmt.Printf("Failde to encode %s: %s\n", vendoredImagePath, err.Error())
		return fallbackReturn
	}

	finish := time.Now()
	fmt.Printf("done in %d ms\n", finish.Sub(start).Milliseconds())

	// Finally.
	return expectedReturn
}

// GalleryItemToImage takes in a gallery item and returns an image object.
func GalleryItemToImage(item *GalleryItem) (image.Image, error) {
	// If it's a local file, simply open the os file.
	if !item.IsExternal {
		file := JoinWorkdir(yunyun.JoinRelativePaths(item.Path, item.Item))
		return OpenImage(string(file))
	}

	// Check if the item has been vendored by any chance?
	vendorPath := filepath.Join(Config.WorkDir, string(GalleryVendored(item)))
	if FileExists(vendorPath) {
		fmt.Printf(" (vendored) ")
		return OpenImage(vendorPath)
	}

	// If it's a remote file, then ask Emilia to try and fetch it.
	return DownloadImage(string(item.Item))
}

// galleryItemHash returns a hashed name of a gallery item link.
func galleryItemHash(item *GalleryItem) yunyun.RelativePathFile {
	return yunyun.RelativePathFile(sha256String(string(item.Item))[:7] + ".jpg")
}
