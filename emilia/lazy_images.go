package emilia

import (
	"fmt"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

const (
	JSPlaceholder = `<script src="%s" async=""></script>`
	lazysizesJS   = `scripts/lazysizes.min.js`
)

// WithLazyGalleries adds the
func WithLazyGalleries() yunyun.PageOption {
	return func(page *yunyun.Page) {
		if gana.Anyf(func(v *yunyun.Content) bool { return v.IsGallery() }, page.Contents) {
			page.Scripts = append(page.Scripts,
				fmt.Sprintf(JSPlaceholder, JoinPath(lazysizesJS)))
		}
	}
}
