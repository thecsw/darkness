package ichika

import (
	"fmt"

	"github.com/thecsw/darkness/emilia"
)

// MisaCommandFunc will support many different tools that darkness can support,
// such as creating gallery previews, etc. WIP.
func MisaCommandFunc() {
	misaCmd := darknessFlagset(misaCommand)

	buildGalleryPreviews := misaCmd.Bool("gallery-previews", false, "build gallery previews")
	removeGalleryPreviews := misaCmd.Bool("no-gallery-previews", false, "delete gallery previews")
	addHolosceneTitles := misaCmd.Bool("holoscene-titles", false, "add holoscene titles")
	rss := misaCmd.Bool("rss", false, "generate rss.xml")
	dryRun := misaCmd.Bool("dry-run", false, "skip writing files (but do the reading)")

	options := getEmiliaOptions(misaCmd)
	if !*rss {
		options.Dev = true
	}
	emilia.InitDarkness(options)

	if *buildGalleryPreviews {
		buildGalleryFiles(*dryRun)
	}
	if *removeGalleryPreviews {
		removeGalleryFiles(*dryRun)
	}
	if *addHolosceneTitles {
		updateHolosceneTitles(*dryRun)
	}
	if *rss {
		rssf(*dryRun)
	}

	if misaCmd.NFlag() == 0 {
		fmt.Println("I don't know what you want me to do, see -help")
	}
}
