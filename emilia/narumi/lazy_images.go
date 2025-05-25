package narumi

import (
	"fmt"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
)

const (
	JSPlaceholder                         = `<script src="%s" async=""></script>`
	lazysizesJS   yunyun.RelativePathFile = `scripts/lazysizes.min.js`
)

// WithLazyGalleries adds the lazy image loading scripts
// (thanks to https://afarkas.github.io/lazysizes/index.html) if any
// gallery blocks are found.
func WithLazyGalleries(conf *alpha.DarknessConfig) yunyun.PageOption {
	return func(page *yunyun.Page) {
		if page == nil || page.Contents == nil || conf == nil {
			return
		}
		if gana.Anyf(func(v *yunyun.Content) bool { return v.IsGallery() }, page.Contents) {
			page.Scripts = append(page.Scripts,
				fmt.Sprintf(JSPlaceholder, conf.Runtime.Join(lazysizesJS)))
		}
	}
}
