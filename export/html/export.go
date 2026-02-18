package html

import (
	_ "embed"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/ichika/akane"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
)

const (
	tombEnding = " ◼"
)

var (
	//go:embed banner.txt
	darknessBannerSource string
	// darknessBanner wrapes `darknessBannerSource` in a comment block.
	darknessBanner = "<!--\n" + darknessBannerSource + "\n-->\n"
)

func (e ExporterHtml) Do(page *yunyun.Page) io.Reader {
	s := &state{conf: e.Config, page: page}

	// All of these indices in this array MUST match the defined types in
	// yunyun/flags.go. Since all of the types are defined as incremental iota
	// elements, these are set up to "just work" that way. The terminal type
	// is an exception, as it should never be used anywhere.
	s.contentFunctions = []func(*yunyun.Content) string{
		s.heading,
		s.paragraph,
		s.list,
		s.listNumbered,
		s.link,
		s.sourceCode,
		s.rawHtml,
		s.horizontalLine,
		s.attentionBlock,
		s.table,
		s.details,
		s.toc,
	}
	return s.export()
}

// Export runs the process of exporting
func (e *state) export() io.Reader {
	// Initialize the html mapping after yunyun built regexes.
	markupHtmlMappingSetOnce.Do(func() {
		markupHtmlMapping = map[*regexp.Regexp]string{
			yunyun.BoldItalicText:    `$l<strong><em>$text</em></strong>$r`,
			yunyun.ItalicBoldText:    `$l<em><strong>$text</strong></em>$r`,
			yunyun.ItalicText:        `$l<em>$text</em>$r`,
			yunyun.BoldText:          `$l<strong>$text</strong>$r`,
			yunyun.VerbatimText:      `$l<code>$text</code>$r`,
			yunyun.StrikethroughText: `$l<s>$text</s>$r`,
			yunyun.UnderlineText:     `$l<u>$text</u>$r`,
			yunyun.SuperscriptText:   `$l<sup>$text</sup>$r`,
			yunyun.SubscriptText:     `$l<sub>$text</sub>$r`,
		}
	})

	// Add the red tomb to the last paragraph on given directories.
	// Only trigger if the tombs were manually flipped.
	if e.page.Accoutrement.Tomb.IsEnabled() {
		e.addTomb()
	}
	// If the page hasn't set a custom preview, default to emilia.
	if len(e.page.Accoutrement.Preview) < 1 {
		e.page.Accoutrement.Preview = string(e.conf.Website.Preview)
	}

	// If the user sets the toc option to true, then we put it as the very first
	// element on the whole page. For finer control, it is recommended for users
	// to use the #+toc control instead, where they can control where toc goes
	// (even multiple times!)
	if e.page.Accoutrement.Toc.IsEnabled() {
		e.page.Contents = append(e.tocAsContent(), e.page.Contents...)
	}

	// If the page requests a preview, generate it.
	if e.page.Accoutrement.PreviewGenerate.IsEnabled() {
		// Set the default preview width and height that the generator will use.
		e.page.Accoutrement.PreviewWidth = puck.PagePreviewWidthString
		e.page.Accoutrement.PreviewHeight = puck.PagePreviewHeightString

		// Send the page to the preview generator.
		akane.RequestPagePreview(e.page.Location, e.page.Title, e.page.Date,
			e.page.Accoutrement.PreviewGenerateBg, e.page.Accoutrement.PreviewGenerateFg)
	}

	// Build the HTML (string) representation of each content
	content := make([]string, 0, len(e.page.Contents))
	for i, v := range e.page.Contents {
		e.currentContentIndex = i
		e.currentContent = v
		content = append(content, e.buildContent())
	}

	output := fmt.Sprintf(`%s<!DOCTYPE html>
<html lang="en">
<head>
%s
<title>%s</title>
</head>
<body class="article">
%s
%s
%s
</body>
</html>`,
		darknessBanner,
		e.combineAndFilterHtmlHead(),
		processTitle(flattenFormatting(e.page.Title)),
		e.authorHeader(),
		strings.Join(content, ""),
		e.addFootnotes(),
	)

	return strings.NewReader(output)
}

// buildContent builds the HTML representation of a content.
func (e *state) buildContent() string {
	// Build the HTML (string) representation of each content.
	built := e.contentFunctions[e.currentContent.Type](e.currentContent)

	// Set the content flags, like whether it's in writing mode or not.
	e.setContentFlags(e.currentContent)

	// If the content is in writing mode, wrap it in a writing div.
	// otherwise, wrap it in other divs, depending on the content type.
	return e.resolveDivTags(built)
}

func (e *state) combineAndFilterHtmlHead() string {
	// Build the array of all head elements (except page's specific head options).
	allHead := [][]string{e.linkTags(), e.metaTags(), e.styleTags(), e.scriptTags()}

	// Go through all the head elements and filter them out depending on page's specific exclusion rules.
	var finalHead strings.Builder
	for _, head := range allHead {
		finalHead.WriteString(strings.Join(gana.Filter(e.page.Accoutrement.ExcludeHtmlHeadContains.ShouldKeep, head), "\n"))
	}

	// User can provide multiple inserts with the same property key, however, in almost all browsers, only the
	// first one is respected, rest are ignored. This could create a jarring experience if say some file is using
	// an imported org manifest, or even nested, and trying to override some setting wouldn't work. So, let's only
	// allow the latest named property.
	e.conf.Website.ExtraHead = filterByLatestMetaName(e.conf.Website.ExtraHead)
	e.page.HtmlHead = filterByLatestMetaName(e.page.HtmlHead)

	// Compile it all together.
	extraHeads := strings.Join(e.conf.Website.ExtraHead, "\n") + "\n" + strings.Join(e.page.HtmlHead, "\n")

	// Then collect it all together.
	return finalHead.String() + extraHeads + "\n"
}

// styleTags is the processed style tags.
func (e *state) styleTags() []string {
	content := make([]string, len(e.conf.Website.Styles)+len(e.page.Stylesheets))
	for i, style := range e.conf.Website.Styles {
		stylePath := yunyun.FullPathFile(style)
		if !strings.HasPrefix(string(style), "http") {
			stylePath = e.conf.Runtime.Join(style)
		}
		content[i] = fmt.Sprintf(
			`<link rel="stylesheet" type="text/css" href="%s">`+"\n", stylePath,
		)
	}
	return append(content, e.page.Stylesheets...)
}

// defaultScripts are the default scripts.
var defaultScripts = []string{
	`<script type="module">document.documentElement.classList.remove("no-js");document.documentElement.classList.add("js");</script>`,
	`<script async src="https://sandyuraz.com/scripts/time.js"></script>`,
}

// scriptTags returns the script tags.
func (e *state) scriptTags() []string {
	return append(defaultScripts, e.page.Scripts...)
}

func (e *state) rssLink() string {
	if !e.conf.RSS.Enable {
		return ""
	}
	return `<span><a href="/feed.xml" class="rss-link"><img src="/assets/rss.svg" class="rss-icon"></a></span><br>` + "\n"
}

func (e *state) authorEmail() string {
	if !e.conf.Author.EmailEnable {
		return ""
	}
	return `<span id="email" class="email">` + e.conf.Author.Email + `</span><br>` + "\n"
}

// authorHeader returns the author header.
func (e *state) authorHeader() string {
	content := fmt.Sprintf(`
<div class="header">
<h1 class="section-1">%s%s</h1>
<div class="menu">
%s%s`,
		e.authorImage(), processTitle(e.page.Title),
		e.rssLink(), e.authorEmail(),
	)
	content += `<span id="revdate">` + "\n"

	// Build the navigation links.
	navLinks := make([]string, 0, len(e.conf.Navigation))

	// Go through elements.
	for i := 1; i <= len(e.conf.Navigation); i++ {
		// Get the navigation element read from Darkness' toml.
		v := e.conf.Navigation[strconv.FormatInt(int64(i), 10)]

		// How about adding a "Go Up ☝️" button?
		fullLoc := filepath.Join("/", string(e.page.Location))

		// If the nav element wants to hide in this location, then skip it.
		if yunyun.RelativePathDir(fullLoc) == v.Hide {
			continue
		}

		whatToJoin := v.Link
		// Relative path.
		if !strings.HasPrefix(string(v.Link), "/") {
			whatToJoin = yunyun.RelativePathDir(filepath.Join(string(e.page.Location), string(v.Link)))
		}

		// If it matched hideif, then exist.
		if len(v.HideIf) > 0 && yunyun.RelativePathDir(filepath.Join("/", string(whatToJoin))) == v.HideIf {
			continue
		}

		// Otherwise, join against the relative path of this page.
		navLinks = append(navLinks,
			fmt.Sprintf(`<a href="%s">%s</a>`,
				e.conf.Runtime.Join(yunyun.RelativePathFile(whatToJoin)),
				v.Title,
			))

	}

	// Close the navigation links span.
	content += strings.Join(navLinks, " | ") + `</span>`

	// Add the Holoscene time element.
	content += `
</div>
<div id="hetime" class="menu"></div>
</div>`
	// Return the website header.
	return content
}

// authorHeader returns img element if author header image is given.
func (e *state) authorImage() string {
	// Return nothing if it's not provided.
	if e.conf.Author.Image == "" || e.page.Accoutrement.AuthorImage.IsDisabled() {
		return ""
	}
	return fmt.Sprintf(`<img id="myface" src="%s" alt="avatar">`, e.conf.Author.ImagePreComputed)
}

// addTomb adds the tomb to the last paragraph.
func (e *state) addTomb() {
	// Empty???
	if len(e.page.Contents) < 1 {
		return
	}
	// Find the last paragraph and attached the tomb.
	for i := len(e.page.Contents) - 1; i >= 0; i-- {
		// Skip if it's not a paragraph.
		if !e.page.Contents[i].IsParagraph() {
			continue
		}
		// Add the tomb and break out.
		e.page.Contents[i].Paragraph += tombEnding
		break
	}
}

// toc returns the table of contents.
func (e *state) tocAsContent() []*yunyun.Content {
	return []*yunyun.Content{
		// First, add the table of contents header.
		{
			Type:                 yunyun.TypeHeading,
			Heading:              "table of Contents",
			HeadingLevel:         3,
			HeadingLevelAdjusted: 1,
		},
		// Then, add the table of contents.
		{
			Type: yunyun.TypeList,
			// overload the summary field to indicate
			// that this is the table of contents.
			Summary: "toc",
			List:    GenerateTableOfContents(e.page),
		},
	}
}
