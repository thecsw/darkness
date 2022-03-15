package html

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
)

type rel struct {
	Rel  string
	Href string
	Type string
}

func linkTag(val rel) string {
	return fmt.Sprintf(`<link rel="%s" href="%s" type="%s"/>`, val.Rel, val.Href, val.Type)
}

var (
	// used in `linkTags`
	links = make([]string, 3)
)

func linkTags(page *internals.Page) string {
	toAdd := []rel{
		{"canonical", page.URL, ""},
		{"shortcut icon", emilia.JoinPath("favicon.ico"), "image/x-icon"},
		{"icon", emilia.JoinPath("favicon.ico"), ""},
	}
	for i, add := range toAdd {
		links[i] = linkTag(add)
	}
	return strings.Join(links, "\n")
}
