package rem

import (
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/emilia/reze"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/rei"
)

// GalleryVendored returns vendored local path of the gallery item.
func GalleryVendored(item GalleryItem) yunyun.RelativePathFile {
	return yunyun.JoinRelativePaths(emilia.Config.Project.DarknessVendorDirectory, galleryItemHash(item))
}

// galleryVendorItem vendors given item and returns the full path of the file.
//
// Only call this function on remote images, it's up to the user to make the
// .IsExternal check before calling this. SLOW function because of network calls.
//
// If the vendoring fails at any point, fallback to the remote image path.
func galleryVendorItem(item GalleryItem) yunyun.FullPathFile {
	// Process only one vendor request at a time.
	vendorLock.Lock()
	// Unlock so the next vendor request can get processed.
	defer vendorLock.Unlock()

	vendoredImagePath := GalleryVendored(item)
	localVendoredPath := filepath.Join(emilia.Config.WorkDir, string(vendoredImagePath))

	// Create the two types of return.
	fallbackReturn := yunyun.FullPathFile(item.Item)
	expectedReturn := emilia.JoinPath(vendoredImagePath)

	// Check if the image was already vendored, if it was, return it immediately.
	if exists, err := rei.FileExists(localVendoredPath); exists {
		return expectedReturn
	} else if err != nil {
		emilia.Logger.Errorf("checking for vendored path existence %s: %v", localVendoredPath, err)
		return fallbackReturn
	}

	img, err := reze.DownloadImage(string(item.Item), "vendor", "", string(galleryItemHash(item)))
	if err != nil {
		emilia.Logger.Errorf("vendoring %s: %v", item.Item, err)
		return fallbackReturn
	}

	// Open the file writer and encode the image there.
	imgFile, err := os.Create(filepath.Clean(localVendoredPath))
	if err != nil {
		emilia.Logger.Errorf("creating vendored file %s: %v", localVendoredPath, err)
		return fallbackReturn
	}
	defer func() {
		if err := imgFile.Close(); err != nil {
			emilia.Logger.Errorf("closing vendored file %s: %v", localVendoredPath, err)
		}
	}()

	// Decode the image into the file.
	if err := imaging.Encode(imgFile, img, imaging.JPEG); err != nil {
		emilia.Logger.Errorf("encoding vendored file %s: %v", vendoredImagePath, err)
		return fallbackReturn
	}

	// Finally.
	return expectedReturn
}
