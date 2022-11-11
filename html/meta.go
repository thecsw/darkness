package html

import (
	"fmt"
	"html"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/echidna"
)

// metaTopTag is the top tag for all meta tags
func (e *ExporterHTML) metaTags() string {
	// Find the first paragraph for description
	description := ""
	for _, content := range e.Page.Contents {
		// We are only looking for paragraphs
		if !content.IsParagraph() {
			continue
		}
		// Skip holoscene times
		paragraph := strings.TrimSpace(content.Paragraph)
		if paragraph == "" || emilia.HEregex.MatchString(paragraph) {
			continue
		}
		description = flattenFormatting(
			paragraph[:echidna.Min(len(paragraph), emilia.Config.Website.DescriptionLength)]) + "..."
		break
	}
	tags := make([]string, 3)
	tags[0] = addBasic(e.Page, description)
	tags[1] = addOpenGraph(e.Page, description)
	tags[2] = addTwitterMeta(e.Page, description)
	return strings.Join(tags, "")
}

// meta is a struct for meta tags
type meta struct {
	Name     string
	Property string
	Content  string
}

// metaTag returns a string of the form <meta name="..." content="..." />
func metaTag(val meta) string {
	return fmt.Sprintf(
		`<meta name="%s" property="%s" content="%s">`,
		val.Name, val.Property, html.EscapeString(val.Content),
	)
}

// metaTopTag is the top tag for all meta tags
const metaTopTag = `<meta charset="UTF-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">`

// addBasic adds the basic meta tags
func addBasic(page *yunyun.Page, description string) string {
	return metaTopTag + strings.Join(echidna.Map(metaTag, []meta{
		{"viewport", "viewport", "width=device-width, initial-scale=1.0"},
		{"generator", "generator", "Darkness"},
		{"author", "author", emilia.Config.Author.Name},
		{"date", "date", page.Date},
		{"theme-color", "theme-color", emilia.Config.Website.Color},
		{"description", "description", html.EscapeString(description)}}), "\n")
}

// addOpenGraph adds the opengraph preview meta tags
func addOpenGraph(page *yunyun.Page, description string) string {
	return strings.Join(echidna.Map(metaTag, []meta{
		{"og:title", "og:title", html.EscapeString(flattenFormatting(page.Title))},
		{"og:site_name", "og:site_name", html.EscapeString(emilia.Config.Title)},
		{"og:url", "og:url", page.URL},
		{"og:locale", "og:locale", emilia.Config.Website.Locale},
		{"og:type", "og:type", "website"},
		{"og:image", "og:image", strings.TrimRight(page.URL, "/") + "/" + emilia.Config.Website.Preview},
		{"og:image:alt", "og:image:alt", "Preview"},
		{"og:image:type", "og:image:type", "image/" + strings.TrimLeft(filepath.Ext(emilia.Config.Website.Preview), ".")},
		{"og:image:width", "og:image:width", "1280"},
		{"og:image:height", "og:image:height", "640"},
		{"og:description", "og:description", html.EscapeString(description)}}), "\n")
}

// addTwitterMeta adds the twitter preview meta tags
func addTwitterMeta(page *yunyun.Page, description string) string {
	return strings.Join(echidna.Map(metaTag, []meta{
		{"twitter:card", "twitter:card", "summary_large_image"},
		{"twitter:site", "twitter:site", html.EscapeString(emilia.Config.Title)},
		{"twitter:creator", "twitter:creator", emilia.Config.Website.Twitter},
		{"twitter:image:src", "twitter:image:src", strings.TrimRight(page.URL, "/") + "/" + emilia.Config.Website.Preview},
		{"twitter:url", "twitter:url", page.URL},
		{"twitter:title", "twitter:title", html.EscapeString(flattenFormatting(page.Title))},
		{"twitter:description", "twitter:description", html.EscapeString(description)}}), "\n")
}
