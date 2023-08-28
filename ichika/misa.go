package ichika

import (
	"fmt"
	"os"
	"strings"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
)

// MisaCommandFunc will support many different tools that darkness can support,
// such as creating gallery previews, etc. WIP.
func MisaCommandFunc() {
	misaCmd := darknessFlagset(misaCommand)

	buildGalleryPreviews := misaCmd.Bool("gallery-previews", false, "build gallery previews")
	removeGalleryPreviews := misaCmd.Bool("no-gallery-previews", false, "delete gallery previews")
	addHolosceneTitles := misaCmd.Bool("holoscene-titles", false, "add holoscene titles")
	rss := misaCmd.String("rss", "", "generate an rss file")
	rssDirectories := misaCmd.String("rss-dirs", "", "look up specific dirs")
	dryRun := misaCmd.Bool("dry-run", false, "skip writing files (but do the reading)")

	options := getAlphaOptions(misaCmd)
	options.Dev = true

	puck.Logger.SetPrefix("Misa 🍎 ")

	if len(*rss) > 0 {
		options.Dev = false
	}
	conf := alpha.BuildConfig(options)

	if *buildGalleryPreviews {
		buildGalleryFiles(conf, *dryRun)
		os.Exit(0)
	}
	if *removeGalleryPreviews {
		removeGalleryFiles(conf, *dryRun)
		os.Exit(0)
	}
	if *addHolosceneTitles {
		updateHolosceneTitles(conf, *dryRun)
		os.Exit(0)
	}
	if len(*rss) > 0 {
		rssf(conf, *rss, strings.Split(*rssDirectories, ","), *dryRun)
		os.Exit(0)
	}

	if misaCmd.NFlag() == 0 {
		fmt.Println("I don't know what you want me to do, see -help")
	}
}
