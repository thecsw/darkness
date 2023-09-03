package akane

import (
	"github.com/thecsw/darkness/emilia/alpha"
)

func Do(conf *alpha.DarknessConfig) {
	logger.Info("Starting to process requests...")

	// Do page previews generation.
	logger.Info("Generating page previews...", "page_previews", len(pagePreviewsToGenerate))
	doPagePreviews(conf)
}
