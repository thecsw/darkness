package html

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
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
func linkTags(page *internals.Page) string {
	toAdd := []rel{
		{"canonical", page.URL, ""},
		{"shortcut icon", emilia.JoinPath("favicon.ico"), "image/x-icon"},
		{"icon", emilia.JoinPath("favicon.ico"), ""},
	}
	links := make([]string, 3)
	for i, add := range toAdd {
		links[i] = linkTag(add)
	}
	return strings.Join(links, "\n")
}
