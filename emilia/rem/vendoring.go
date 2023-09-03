package rem

import (
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/reze"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/rei"
)

// GalleryVendored returns vendored local path of the gallery item.
func GalleryVendored(vendorDir yunyun.RelativePathDir, item GalleryItem) yunyun.RelativePathFile {
	return yunyun.JoinRelativePaths(vendorDir, galleryItemHash(item))
}

// galleryVendorItemFilename returns the path to the vendored item.
func galleryVendorItemFilename(conf *alpha.DarknessConfig, item GalleryItem) yunyun.FullPathFile {
	vendoredImagePath := GalleryVendored(conf.Project.DarknessVendorDirectory, item)
	expectedReturn := conf.Runtime.Join(vendoredImagePath)
	return expectedReturn
}

// galleryVendorItemFilenameLocalPath is the actual local filesystem path where to save vendored files.
func galleryVendorItemFilenameLocalPath(conf *alpha.DarknessConfig, item GalleryItem) string {
	// This is the vendored path that will be put into the outputs.
	vendoredImagePath := GalleryVendored(conf.Project.DarknessVendorDirectory, item)
	// This is the actual physical location of the vendored image.
	localVendoredPath := string(conf.Runtime.WorkDir.Join(vendoredImagePath))
	return localVendoredPath
}

// GalleryVendorItem vendors given item and returns the full path of the file.
//
// Only call this function on remote images, it's up to the user to make the
// .IsExternal check before calling this. SLOW function because of network calls.
//
// If the vendoring fails at any point, fallback to the remote image path.
func GalleryVendorItem(conf *alpha.DarknessConfig, item GalleryItem) (yunyun.FullPathFile, bool) {
	// Create the two types of return.
	fallbackReturn := yunyun.FullPathFile(item.Item)
	localVendoredPath := galleryVendorItemFilenameLocalPath(conf, item)
	expectedReturn := galleryVendorItemFilename(conf, item)

	// Check if the image was already vendored, if it was, return it immediately.
	if exists, err := rei.FileExists(localVendoredPath); exists {
		return expectedReturn, false
	} else if err != nil {
		logger.Error("checking for vendored path existence", "path", localVendoredPath, "err", err)
		return fallbackReturn, false
	}

	img, err := reze.DownloadImage(string(item.Item), "vendor", "", string(galleryItemHash(item)))
	if err != nil {
		logger.Error("downloading vendored image", "item", item.Item, "err", err)
		return fallbackReturn, false
	}

	// Open the file writer and encode the image there.
	imgFile, err := os.Create(filepath.Clean(localVendoredPath))
	if err != nil {
		logger.Error("creating vendored file", "path", localVendoredPath, "err", err)
		return fallbackReturn, false
	}
	defer func() {
		if err := imgFile.Close(); err != nil {
			logger.Error("closing vendored file", "file", localVendoredPath, "err", err)
		}
	}()

	// Decode the image into the file.
	if err := imaging.Encode(imgFile, img, imaging.JPEG); err != nil {
		logger.Error("encoding vendored file", "file", localVendoredPath, "err", err)
		return fallbackReturn, false
	}

	// Finally.
	return expectedReturn, true
}
