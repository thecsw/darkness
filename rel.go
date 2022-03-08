package main

import "fmt"

type rel struct {
	Rel  string
	Href string
	Type string
}

func relTag(val rel) string {
	return fmt.Sprintf(
		`<rel rel="%s" href="%s" type="%s">`+"\n",
		val.Rel, val.Href, val.Type,
	)
}

func addRelTags(page *Page) string {
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
		content += relTag(add)
	}
	return content
}
