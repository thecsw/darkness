package html

import (
	"fmt"
	"html"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// metaTopTag is the top tag for all meta tags
func (e ExporterHTML) metaTags() []string {
	// Find the first paragraph for description
	description := ""
	for _, content := range e.page.Contents {
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
			paragraph[:gana.Min(len(paragraph), emilia.Config.Website.DescriptionLength)]) + "..."
		break
	}
	basic := addBasic(e.page, description)
	openGraph := addOpenGraph(e.page, description)
	twitter := addTwitterMeta(e.page, description)

	metas := make([]string, 0, len(basic)+len(openGraph)+len(twitter))
	metas = append(metas, basic...)
	metas = append(metas, openGraph...)
	metas = append(metas, twitter...)

	return metas
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
var metaTopTag = []string{
	`<meta charset="UTF-8">`,
	`<meta http-equiv="X-UA-Compatible" content="IE=edge">`,
}

// addBasic adds the basic meta tags
func addBasic(page *yunyun.Page, description string) []string {
	return append(metaTopTag, gana.Map(metaTag, []meta{
		{"viewport", "viewport", "width=device-width, initial-scale=1.0"},
		{"generator", "generator", "Darkness"},
		{"author", "author", emilia.Config.Author.Name},
		{"date", "date", page.Date},
		{"theme-color", "theme-color", emilia.Config.Website.Color},
		{"description", "description", html.EscapeString(description)},
	})...)
}

// addOpenGraph adds the opengraph preview meta tags
func addOpenGraph(page *yunyun.Page, description string) []string {
	return gana.Map(metaTag, []meta{
		{"og:title", "og:title", html.EscapeString(flattenFormatting(page.Title))},
		{"og:site_name", "og:site_name", html.EscapeString(emilia.Config.Title)},
		{"og:url", "og:url", string(emilia.JoinPathGeneric[yunyun.RelativePathDir, yunyun.FullPathDir](page.Location))},
		{"og:locale", "og:locale", emilia.Config.Website.Locale},
		{"og:type", "og:type", "website"},
		{"og:image", "og:image", string(emilia.JoinPath(
			yunyun.JoinRelativePaths(page.Location, yunyun.RelativePathFile(page.Accoutrement.Preview))))},
		{"og:image:alt", "og:image:alt", "Preview"},
		{"og:image:type", "og:image:type", "image/" + strings.TrimLeft(filepath.Ext(page.Accoutrement.Preview), ".")},
		{"og:image:width", "og:image:width", "1280"},
		{"og:image:height", "og:image:height", "640"},
		{"og:description", "og:description", html.EscapeString(description)},
	})
}

// addTwitterMeta adds the twitter preview meta tags
func addTwitterMeta(page *yunyun.Page, description string) []string {
	return gana.Map(metaTag, []meta{
		{"twitter:card", "twitter:card", "summary_large_image"},
		{"twitter:site", "twitter:site", html.EscapeString(emilia.Config.Title)},
		{"twitter:creator", "twitter:creator", emilia.Config.Website.Twitter},
		{"twitter:image:src", "twitter:image:src", string(emilia.JoinPath(
			yunyun.JoinRelativePaths(page.Location, yunyun.RelativePathFile(page.Accoutrement.Preview))))},
		{"twitter:url", "twitter:url", string(emilia.JoinPathGeneric[yunyun.RelativePathDir, yunyun.FullPathDir](page.Location))},
		{"twitter:title", "twitter:title", html.EscapeString(flattenFormatting(page.Title))},
		{"twitter:description", "twitter:description", html.EscapeString(description)},
	})
}
