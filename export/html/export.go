package html

import (
	_ "embed"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/ichika/akane"
	"github.com/thecsw/darkness/yunyun"
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
	}
	return s.export()
}

// Export runs the process of exporting
func (e *state) export() io.Reader {
	defer puck.Stopwatch("Exported", "page", e.page.File).Record()

	// Initialize the html mapping after yunyun built regexes.
	markupHtmlMappingSetOnce.Do(func() {
		markupHtmlMapping = map[*regexp.Regexp]string{
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

	if e.page.Accoutrement.Toc.IsEnabled() {
		e.page.Contents = append(e.toc(), e.page.Contents...)
	}

	if e.page.Accoutrement.PreviewGenerate.IsEnabled() {
		e.page.Accoutrement.PreviewWidth = puck.PagePreviewWidthString
		e.page.Accoutrement.PreviewHeight = puck.PagePreviewHeightString
		akane.RequestPagePreview(e.page.Location, e.page.Title, e.page.Date)
	}

	// Build the HTML (string) representation of each content
	content := make([]string, 0, len(e.page.Contents))
	for i, v := range e.page.Contents {
		e.currentContentIndex = i
		e.currentContent = v
		content = append(content, e.buildContent(v))
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
func (e *state) buildContent(content *yunyun.Content) string {
	// Build the HTML (string) representation of each content.
	built := e.contentFunctions[e.currentContent.Type](e.currentContent)

	// Set the content flags, like whether it's in writing mode or not.
	e.setContentFlags(e.currentContent)

	// If the content is in writing mode, wrap it in a writing div.
	// otherwise, wrap it in other divs, depending on the content type.
	return e.resolveDivTags(built)
}

// leftHeading leaves the heading.
func (e *state) leftHeading() {
	e.inHeading = false
}

func (e *state) combineAndFilterHtmlHead() string {
	// Build the array of all head elements (except page's specific head options).
	allHead := [][]string{e.linkTags(), e.metaTags(), e.styleTags(), e.scriptTags(), e.conf.Website.ExtraHead}
	// Go through all the head elements and filter them out depending on page's specific exclusion rules.
	finalHead := ""
	for _, head := range allHead {
		finalHead += strings.Join(gana.Filter(e.page.Accoutrement.ExcludeHtmlHeadContains.ShouldKeep, head), "\n")
	}
	// Page's specific html head elements are not filtered out.
	return finalHead + "\n" + strings.Join(e.page.HtmlHead, "\n")
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

func (e *state) authorName() string {
	if !e.conf.Author.NameEnable {
		return ""
	}
	return `<span id="author" class="author">` + e.conf.Author.Name + `</span><br>` + "\n"
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
%s%s%s`,
		e.authorImage(), processTitle(e.page.Title),
		e.rssLink(), e.authorName(), e.authorEmail(),
	)
	content += `<span id="revdate">` + "\n"

	// Build the navigation links.
	navLinks := make([]string, 0, len(e.conf.Navigation))

	// Go through elements.
	for i := 1; i <= len(e.conf.Navigation); i++ {
		// Get the navigation element read from Darkness' toml.
		v := e.conf.Navigation[strconv.FormatInt(int64(i), 10)]
		// If the nav element wants to hide in this location, then skip it.
		if e.page.Location == v.Hide {
			continue
		}
		// Build each of the navlinks and concat the hrefs.
		navLinks = append(navLinks, fmt.Sprintf(`<a href="%s">%s</a>`,
			e.conf.Runtime.Join(yunyun.RelativePathFile(v.Link)),
			v.Title))
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
func (e *state) toc() []*yunyun.Content {
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
		// Finally, add the horizontal line.
		{
			Type: yunyun.TypeHorizontalLine,
		},
	}
}
