package akane

import (
	"fmt"

	"github.com/thecsw/darkness/emilia/alpha"
)

// Do starts going through the requests and processes them.
func Do(conf *alpha.DarknessConfig) {
	fmt.Println()
	logger.Info("Starting to process requests...")

	if len(pagePreviewsToGenerate) > 0 {
		// Do page previews generation.
		logger.Info("Generating page previews...",
			"page_previews", len(pagePreviewsToGenerate))
		doPagePreviews(conf)
	}

	if conf.Runtime.VendorGalleries {
		// Do the gallery vendoring.
		logger.Info("Generating gallery vendors...",
			"gallery_vendors", len(galleryVendorsToDownload))
		doGalleryVendors(conf)
	}
}
