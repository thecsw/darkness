package akane

import (
	"fmt"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/rem"
)

// galleryVendorRequest is a request to download a gallery vendor.
type galleryVendorRequest struct {
	Item rem.GalleryItem
}

// galleryVendorsToDownload is a list of vendors to download.
var galleryVendorsToDownload = make([]galleryVendorRequest, 0, 16)

// RequestGalleryVendor adds a gallery vendor to the list of vendors to download.
func RequestGalleryVendor(item rem.GalleryItem) {
	galleryVendorsToDownload = append(galleryVendorsToDownload, galleryVendorRequest{
		Item: item,
	})
}

// Go through gallery requests and download the images.
func doGalleryVendors(conf *alpha.DarknessConfig) {
	// Clear the request list when we're done.
	defer func() {
		galleryVendorsToDownload = galleryVendorsToDownload[:0]
	}()

	// Go through each gallery vendor request.
	for _, galleryVendorRequestItem := range galleryVendorsToDownload {
		item := galleryVendorRequestItem.Item
		path, downloaded := rem.GalleryVendorItem(conf, item)
		if downloaded {
			// Clear the progressbar.
			fmt.Print("\r\033[2K")
			// Log the thing.
			logger.Info("Vendored item", "path", conf.Runtime.Rel(path), "dir", item.Path)
		}
	}
}
