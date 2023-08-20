package rem

import (
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/rei"
)

// galleryItemHash returns a hashed name of a gallery item link.
func galleryItemHash(item GalleryItem) yunyun.RelativePathFile {
	return yunyun.RelativePathFile(rei.Sha256([]byte(item.Item))[:7] + ".jpg")
}

// galleryPreviewRelative takes gallery item and returns relative path to it.
func galleryPreviewRelative(item GalleryItem) yunyun.RelativePathFile {
	if item.IsExternal {
		return galleryItemHash(item)
	}
	filename := filepath.Base(string(item.Item))
	ext := filepath.Ext(filename)
	return yunyun.RelativePathFile(strings.TrimSuffix(filename, ext) + "_small.jpg")
}
