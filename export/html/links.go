package html

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/gana"
)

// rel is a struct for holding the rel and href of a link
type rel struct {
	Rel  string
	Href string
	Type string
}

// linkTag returns a string of the form <link rel="..." href="..." />
func linkTag(val rel) string {
	return fmt.Sprintf(`<link rel="%s" href="%s" type="%s"/>`, val.Rel, val.Href, val.Type)
}

// linkTags returns a string of the form <link rel="..." href="..." /> for an entire page
func (e ExporterHTML) linkTags() string {
	return strings.Join(gana.Map(linkTag, []rel{
		{"canonical", e.page.URL, ""},
		{"shortcut icon", emilia.JoinPath("assets/favicon.ico"), "image/x-icon"},
		{"apple-touch-icon", emilia.JoinPath("assets/apple-touch-icon.png"), "image/png"},
		{"image_src", emilia.JoinPath("assets/android-chrome-512x512.png"), "image/png"},
		{"icon", emilia.JoinPath("assets/favicon.ico"), ""},
	}), "\n")
}
