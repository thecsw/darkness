package html

import (
	"fmt"
	"html"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// metaTopTag is the top tag for all meta tags
func (e state) metaTags() []string {
	// Find the first paragraph for description
	description := ""
	for _, content := range e.page.Contents {
		// We are only looking for paragraphs
		if !content.IsParagraph() {
			continue
		}
		// Skip holoscene times
		paragraph := strings.TrimSpace(content.Paragraph)
		if paragraph == "" || puck.HEregex.MatchString(paragraph) {
			continue
		}
		description = flattenFormatting(
			paragraph[:gana.Min(len(paragraph), e.conf.Website.DescriptionLength)]) + "..."
		break
	}
	basic := addBasic(e.conf, e.page, description)
	openGraph := addOpenGraph(e.conf, e.page, description)
	twitter := addTwitterMeta(e.conf, e.page, description)

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
func addBasic(conf *alpha.DarknessConfig, page *yunyun.Page, description string) []string {
	return append(metaTopTag, gana.Map(metaTag, []meta{
		{"viewport", "viewport", "width=device-width, initial-scale=1.0"},
		{"generator", "generator", "Darkness"},
		{"author", "author", conf.Author.Name},
		{"date", "date", page.Date},
		{"theme-color", "theme-color", conf.Website.Color},
		{"description", "description", html.EscapeString(description)},
	})...)
}

// addOpenGraph adds the opengraph preview meta tags
func addOpenGraph(conf *alpha.DarknessConfig, page *yunyun.Page, description string) []string {
	return gana.Map(metaTag, []meta{
		{"og:title", "og:title", html.EscapeString(flattenFormatting(page.Title))},
		{"og:site_name", "og:site_name", html.EscapeString(conf.Title)},
		{"og:url", "og:url", string(conf.Runtime.Join(yunyun.RelativePathFile(page.Location)))},
		{"og:locale", "og:locale", conf.Website.Locale},
		{"og:type", "og:type", "website"},
		{"og:image", "og:image", string(conf.Runtime.Join(yunyun.JoinRelativePaths(page.Location, yunyun.RelativePathFile(page.Accoutrement.Preview))))},
		{"og:image:alt", "og:image:alt", "Preview"},
		{"og:image:type", "og:image:type", "image/" + strings.TrimLeft(filepath.Ext(page.Accoutrement.Preview), ".")},
		{"og:image:width", "og:image:width", page.Accoutrement.PreviewWidth},    // default: "1200"
		{"og:image:height", "og:image:height", page.Accoutrement.PreviewHeight}, // default: "700"
		{"og:description", "og:description", html.EscapeString(description)},
	})
}

// addTwitterMeta adds the twitter preview meta tags
func addTwitterMeta(conf *alpha.DarknessConfig, page *yunyun.Page, description string) []string {
	return gana.Map(metaTag, []meta{
		{"twitter:card", "twitter:card", "summary_large_image"},
		{"twitter:site", "twitter:site", html.EscapeString(conf.Title)},
		{"twitter:creator", "twitter:creator", conf.Website.Twitter},
		{"twitter:image:src", "twitter:image:src", string(conf.Runtime.Join(yunyun.JoinRelativePaths(page.Location, yunyun.RelativePathFile(page.Accoutrement.Preview))))},
		{"twitter:url", "twitter:url", string(conf.Runtime.Join(yunyun.RelativePathFile(page.Location)))},
		{"twitter:title", "twitter:title", html.EscapeString(flattenFormatting(page.Title))},
		{"twitter:description", "twitter:description", html.EscapeString(description)},
	})
}
