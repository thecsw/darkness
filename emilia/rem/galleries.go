package rem

import (
	"fmt"
	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/reze"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/rei"
	"image"
	"strings"
)

// NewGalleryItem creates a new helper `GalleryItem` and
// decides whether the passed item is an external link or not.
func NewGalleryItem(page *yunyun.Page, content *yunyun.Content, wholeLine string) GalleryItem {
	extractedLinks := yunyun.ExtractLinks(wholeLine)
	// If line wasn't found, then the whole line should be counted as the line path.
	line := wholeLine
	text := ""
	description := ""
	if len(extractedLinks) > 0 {
		line = extractedLinks[0].Link
		text = extractedLinks[0].Text
		description = extractedLinks[0].Description
	}
	optionalLink := ""
	if len(extractedLinks) > 1 {
		optionalLink = extractedLinks[1].Link
	}
	return GalleryItem{
		Item:         yunyun.RelativePathFile(line),
		Path:         yunyun.JoinPaths(page.Location, content.GalleryPath),
		IsExternal:   strings.HasPrefix(line, "http"),
		Text:         text,
		Description:  description,
		OriginalLine: wholeLine,
		Link:         optionalLink,
	}
}

// GalleryImage takes a gallery item and returns its full path depending
// on the option, so whether it's an image that needs to be vendored (downloaded).
func GalleryImage(conf *alpha.DarknessConfig, item GalleryItem) (yunyun.FullPathFile, bool) {
	if item.IsExternal {
		// If it's vendored, then retrieve a local copy (if doesn't already
		// exist) and stub it in as the full path
		if conf.Runtime.VendorGalleries {
			// Return the path to the vendored image.
			return galleryVendorItemFilename(conf, item), true
		}
		return yunyun.FullPathFile(item.Item), false
	}
	return conf.Runtime.Join(yunyun.JoinRelativePaths(item.Path, item.Item)), false
}

// GalleryPreview takes an original image's path and returns
// the preview path of it. Previews are always .jpg
func GalleryPreview(conf *alpha.DarknessConfig, item GalleryItem) yunyun.FullPathFile {
	return conf.Runtime.Join(yunyun.JoinRelativePaths(conf.Project.DarknessPreviewDirectory, galleryPreviewRelative(item)))
}

// GalleryItemToImage takes in a gallery item and returns an image object.
func GalleryItemToImage(conf *alpha.DarknessConfig, item GalleryItem, authority, prefix string) (image.Image, error) {
	// If it's a local file, simply open the os file.
	if !item.IsExternal {
		file := conf.Runtime.WorkDir.Join(yunyun.JoinRelativePaths(item.Path, item.Item))
		return reze.OpenImage(string(file))
	}

	// Check if the item has been vendored by any chance?
	vendorPath := string(conf.Runtime.WorkDir.Join(GalleryVendored(conf.Project.DarknessVendorDirectory, item)))
	if exists, err := rei.FileExists(vendorPath); exists {
		return reze.OpenImage(vendorPath)
	} else if err != nil {
		return nil, fmt.Errorf("checking for vendored file existence %s: %v", vendorPath, err)
	}

	// If it's a remote file, then ask Emilia to try and fetch it.
	return reze.DownloadImage(string(item.Item), authority, prefix, string(galleryItemHash(item)))
}
