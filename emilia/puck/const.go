package puck

import (
	"strconv"

	"github.com/thecsw/darkness/v3/yunyun"
)

const (
	// ExtensionOrgmode is the extension of orgmode files.
	ExtensionOrgmode = ".org"
	// ExtensionMarkdown is the extension of markdown files.
	ExtensionMarkdown = ".md"
	// ExtensionHtml is the extension of html files.
	ExtensionHtml = ".html"

	// DefaultPreviewFile is the name of the file where the preview of the gallery is stored.
	DefaultPreviewFile = "preview.png"

	// DefaultVendorDirectory is the name of the dir where vendor images are stored.
	DefaultVendorDirectory yunyun.RelativePathDir = "darkness_vendor"
	// DefaultPreviewDirectory is the name of the dir where all gallery previews are stored.
	DefaultPreviewDirectory yunyun.RelativePathDir = "darkness_gallery_previews"

	// DefaultPreviewWidth is the default width of the gallery preview.
	PagePreviewWidth = 1200

	// DefaultPreviewHeight is the default height of the gallery preview.
	PagePreviewHeight = 700
)

var (
	// PagePreviewWidthString is the string representation of PagePreviewWidth.
	PagePreviewWidthString = strconv.Itoa(PagePreviewWidth)

	// PagePreviewHeightString is the string representation of PagePreviewHeight.
	PagePreviewHeightString = strconv.Itoa(PagePreviewHeight)
)
