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
	dryRun := misaCmd.Bool("dry-run", false, "skip writing files (but do the reading)")

	options := getEmiliaOptions(misaCmd)
	options.Dev = true
	emilia.InitDarkness(options)

	if *buildGalleryPreviews {
		buildGalleryFiles(*dryRun)
		return
	}
	if *removeGalleryPreviews {
		removeGalleryFiles()
		return
	}

	fmt.Println("I don't know what you want me to do, see -help")
}
