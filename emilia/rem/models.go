package rem

import "github.com/thecsw/darkness/yunyun"

// GalleryItem is a struct that holds the gallery item path
// and a flag whether it is external (Url regexp matches).
type GalleryItem struct {
	// Item is the link that was provided.
	Item yunyun.RelativePathFile
	// Path is the path of the local gallery source file.
	Path yunyun.RelativePathDir
	// Text found through the link regexp.
	Text string
	// Description found through the link regexp.
	Description string
	// OriginalLine is the original line that include org options.
	OriginalLine string
	// Link is an optional parameter that the gallery item should
	// also link to something.
	Link string
	// IsExternal runs a Url regexp check.
	IsExternal bool
}
