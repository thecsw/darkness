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
		description = paragraph[:internals.Min(len(paragraph), DescriptionLength)] + "..."
		break
	}
	tags := make([]string, 3)
	tags[0] = addBasic(page, description)
	tags[1] = addOpenGraph(page, description)
	tags[2] = addTwitterMeta(page, description)
	return strings.Join(tags, "")
}

type meta struct {
	Name     string
	Property string
	Content  string
}

func metaTag(val meta) string {
	return fmt.Sprintf(
		`<meta name="%s" property="%s" content="%s">`,
		val.Name, val.Property, html.EscapeString(val.Content),
	)
}

const metaTopTag = `<meta charset="UTF-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">`

// Pre-allocate the tags slices, because they're same sized
var (
	// used in `addBasic`
	basicTags = make([]string, 5)
	// used in `addOpenGraph`
	opengraphTags = make([]string, 11)
	// used in `addTwitterMeta`
	twitterTags = make([]string, 7)
)

func addBasic(page *internals.Page, description string) string {
	toAdd := []meta{
		{"viewport", "viewport", "width=device-width, initial-scale=1.0"},
		{"generator", "generator", "Darkness"},
		{"author", "author", emilia.Config.Author.Name},
		{"theme-color", "theme-color", emilia.Config.Website.Color},
		{"description", "description", html.EscapeString(description)},
	}
	for i, add := range toAdd {
		basicTags[i] = metaTag(add)
	}
	return metaTopTag + strings.Join(basicTags, "\n")
}

func addOpenGraph(page *internals.Page, description string) string {
	toAdd := []meta{
		{"og:title", "og:title", html.EscapeString(page.Title)},
		{"og:site_name", "og:site_name", html.EscapeString(emilia.Config.Title)},
		{"og:url", "og:url", page.URL},
		{"og:locale", "og:locale", emilia.Config.Website.Locale},
		{"og:type", "og:type", "website"},
		{"og:image", "og:image", page.URL + "/preview.png"},
		{"og:image:alt", "og:image:alt", "Preview"},
		{"og:image:type", "og:image:type", "image/png"},
		{"og:image:width", "og:image:width", "1280"},
		{"og:image:height", "og:image:height", "640"},
		{"og:description", "og:description", html.EscapeString(description)}}
	for i, add := range toAdd {
		opengraphTags[i] = metaTag(add)
	}
	return strings.Join(opengraphTags, "\n")
}

func addTwitterMeta(page *internals.Page, description string) string {
	toAdd := []meta{
		{"twitter:card", "twitter:card", "summary_large_image"},
		{"twitter:site", "twitter:site", html.EscapeString(emilia.Config.Title)},
		{"twitter:creator", "twitter:creator", emilia.Config.Website.Twitter},
		{"twitter:image:src", "twitter:image:src", page.URL + "/preview.png"},
		{"twitter:url", "twitter:url", page.URL},
		{"twitter:title", "twitter:title", html.EscapeString(page.Title)},
		{"twitter:description", "twitter:description", html.EscapeString(description)},
	}
	for i, add := range toAdd {
		twitterTags[i] = metaTag(add)
	}
	return strings.Join(twitterTags, "\n")
}
