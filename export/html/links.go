package html

import (
	"fmt"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// rel is a struct for holding the rel and href of a link
type rel struct {
	Rel  string
	Href yunyun.FullPathFile
	Type string
}

// linkTag returns a string of the form <link rel="..." href="..." />
func linkTag(val rel) string {
	return fmt.Sprintf(`<link rel="%s" href="%s" type="%s"/>`, val.Rel, val.Href, val.Type)
}

// linkTags returns a string of the form <link rel="..." href="..." /> for an entire page
func (e *state) linkTags() []string {
	return gana.Map(linkTag, []rel{
		{"canonical", e.conf.Runtime.Join(yunyun.RelativePathFile(e.page.Location)), ""},
		{"shortcut icon", e.conf.Runtime.Join("assets/favicon.ico"), "image/x-icon"},
		{"apple-touch-icon", e.conf.Runtime.Join("assets/apple-touch-icon.png"), "image/png"},
		{"image_src", e.conf.Runtime.Join("assets/android-chrome-512x512.png"), "image/png"},
		{"icon", e.conf.Runtime.Join("assets/favicon.ico"), ""},
	})
}
