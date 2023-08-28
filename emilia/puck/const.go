package puck

import "github.com/thecsw/darkness/yunyun"

const (
	ExtensionOrgmode  = ".org"
	ExtensionMarkdown = ".md"
	ExtensionHtml     = ".html"

	DefaultPreviewFile = "preview.png"

	// DefaultVendorDirectory is the name of the dir where vendor images are stored.
	DefaultVendorDirectory yunyun.RelativePathDir = "darkness_vendor"
	// DefaultPreviewDirectory is the name of the dir where all gallery previews are stored.
	DefaultPreviewDirectory yunyun.RelativePathDir = "darkness_gallery_previews"
)
