package html

import (
	"darkness/internals"
	"fmt"
)

type rel struct {
	Rel  string
	Href string
	Type string
}

func linkTag(val rel) string {
	return fmt.Sprintf(
		`<link rel="%s" href="%s" type="%s"/>`+"\n",
		val.Rel, val.Href, val.Type,
	)
}

func linkTags(page *internals.Page) string {
	toAdd := []rel{
		{
			Rel:  "canonical",
			Href: page.URL,
			Type: "",
		},
		{
			Rel:  "shortcut icon",
			Href: "/favicon.ico",
			Type: "image/x-icon",
		},
		{
			Rel:  "icon",
			Href: "/favicon.ico",
			Type: "",
		},
	}
	content := ""
	for _, add := range toAdd {
		content += linkTag(add)
	}
	return content
}
