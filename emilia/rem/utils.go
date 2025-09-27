package rem

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/rei"
)

// galleryItemHash returns a hashed name of a gallery item link.
func galleryItemHash(item GalleryItem) yunyun.RelativePathFile {
	return yunyun.RelativePathFile(rei.Sha256([]byte(item.Item))[:7] + ".jpg")
}

// galleryPreviewRelative takes gallery item and returns relative path to it.
func galleryPreviewRelative(item GalleryItem) yunyun.RelativePathFile {
	prefix := rei.Sha256([]byte(item.Path))[:7]
	if item.IsExternal {
		return galleryItemHash(item)
	}
	filename := filepath.Base(string(item.Item))
	ext := filepath.Ext(filename)

	// Multiple directories can have filenames named the same, so we differentiate them
	// by hashing the directory they're coming from. The preview is always going to be jpg.
	final_base := fmt.Sprintf("%s-%s.jpg", prefix, strings.TrimSuffix(filename, ext))
	return yunyun.RelativePathFile(final_base)
}
