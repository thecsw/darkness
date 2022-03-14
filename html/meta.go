package html

import (
	"fmt"
	"html"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
)

const (
	DescriptionLength = 100
)

func metaTags(page *internals.Page) string {
	// Find the first paragraph for description
	description := ""
	for _, content := range page.Contents {
		// We are only looking for paragraphs
		if !content.IsParagraph() {
			continue
		}
		// Skip holoscene times
		paragraph := strings.TrimSpace(content.Paragraph)
		if paragraph == "" || emilia.HEregex.MatchString(paragraph) {
			continue
		}
		description = paragraph[:min(len(paragraph), DescriptionLength)] + "..."
		break
	}

	content := addBasic(page, description)
	content += addOpenGraph(page, description)
	content += addTwitterMeta(page, description)

	return content
}

type meta struct {
	Name     string
	Property string
	Content  string
}

func metaTag(val meta) string {
	return fmt.Sprintf(
		`<meta name="%s" property="%s" content="%s">`+"\n",
		val.Name, val.Property, html.EscapeString(val.Content),
	)
}

func addBasic(page *internals.Page, description string) string {
	toAdd := []meta{
		{
			Name:     "viewport",
			Property: "viewport",
			Content:  "width=device-width, initial-scale=1.0",
		},
		{
			Name:     "generator",
			Property: "generator",
			Content:  "Darkness",
		},
		{
			Name:     "author",
			Property: "author",
			Content:  emilia.Config.Author.Name,
		},
		{
			Name:     "theme-color",
			Property: "theme-color",
			Content:  emilia.Config.Website.Color,
		},
		{
			Name:     "description",
			Property: "description",
			Content:  html.EscapeString(description),
		},
	}
	content := `<meta charset="UTF-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
`
	for _, add := range toAdd {
		content += metaTag(add)
	}
	return content
}

func addOpenGraph(page *internals.Page, description string) string {
	toAdd := []meta{
		{
			Name:     "og:title",
			Property: "og:title",
			Content:  html.EscapeString(page.Title),
		},
		{
			Name:     "og:site_name",
			Property: "og:site_name",
			Content:  html.EscapeString(emilia.Config.Title),
		},
		{
			Name:     "og:url",
			Property: "og:url",
			Content:  page.URL,
		},
		{
			Name:     "og:locale",
			Property: "og:locale",
			Content:  emilia.Config.Website.Locale,
		},
		{
			Name:     "og:type",
			Property: "og:type",
			Content:  "website",
		},
		{
			Name:     "og:image",
			Property: "og:image",
			Content:  page.URL + "/preview.png",
		},
		{
			Name:     "og:image:alt",
			Property: "og:image:alt",
			Content:  "Preview",
		},
		{
			Name:     "og:image:type",
			Property: "og:image:type",
			Content:  "image/png",
		},
		{
			Name:     "og:image:width",
			Property: "og:image:width",
			Content:  "1280",
		},
		{
			Name:     "og:image:height",
			Property: "og:image:height",
			Content:  "640",
		},
		{
			Name:     "og:description",
			Property: "og:description",
			Content:  html.EscapeString(description),
		},
	}
	content := ""
	for _, add := range toAdd {
		content += metaTag(add)
	}
	return content
}

func addTwitterMeta(page *internals.Page, description string) string {
	toAdd := []meta{
		{
			Name:     "twitter:card",
			Property: "twitter:card",
			Content:  "summary_large_image",
		},
		{
			Name:     "twitter:site",
			Property: "twitter:site",
			Content:  html.EscapeString(emilia.Config.Title),
		},
		{
			Name:     "twitter:creator",
			Property: "twitter:creator",
			Content:  emilia.Config.Website.Twitter,
		},
		{
			Name:     "twitter:image:src",
			Property: "twitter:image:src",
			Content:  page.URL + "/preview.png",
		},
		{
			Name:     "twitter:url",
			Property: "twitter:url",
			Content:  page.URL,
		},
		{
			Name:     "twitter:title",
			Property: "twitter:title",
			Content:  html.EscapeString(page.Title),
		},
		{
			Name:     "twitter:description",
			Property: "twitter:description",
			Content:  html.EscapeString(description),
		},
	}
	content := ""
	for _, add := range toAdd {
		content += metaTag(add)
	}
	return content
}

type Number interface {
	int | float64
}

func min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}
