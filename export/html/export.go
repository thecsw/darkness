package html

import (
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
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

// Export runs the process of exporting
func (e ExporterHTML) Export() string {
	if e.page == nil {
		fmt.Println("Export should be called after SetPage")
		os.Exit(1)
	}

	// Initialize the html mapping after yunyun built regexes.
	markupHTMLMapping = map[*regexp.Regexp]string{
		yunyun.ItalicText:        `$l<em>$text</em>$r`,
		yunyun.BoldText:          `$l<strong>$text</strong>$r`,
		yunyun.VerbatimText:      `$l<code>$text</code>$r`,
		yunyun.StrikethroughText: `$l<s>$text</s>$r`,
		yunyun.UnderlineText:     `$l<u>$text</u>$r`,
		yunyun.SuperscriptText:   `$l<sup>$text</sup>$r`,
		yunyun.SubscriptText:     `$l<sub>$text</sub>$r`,
	}

	// Add the red tomb to the last paragraph on given directories.
	// Only trigger if the tombs were manually flipped.
	if e.page.Accoutrement.Tomb.IsEnabled() {
		e.addTomb()
	}
	// Build the HTML (string) representation of each content
	content := make([]string, e.contentsNum)
	for i, v := range e.page.Contents {
		e.currentContentIndex = i
		e.currentContent = v
		content[i] = e.buildContent()
	}

	return fmt.Sprintf(`%s<!DOCTYPE html>
<html lang="en">
<head>
<!-- Links -->
%s
<!-- Meta -->
%s
<!-- Styling -->
%s
<!-- Scripts -->
%s
<!-- Extra -->
%s
<!-- Title -->
<title>%s</title>
</head>
<body class="article">
<!-- Header -->
%s
<!-- Content -->
%s
<!-- Footnotes -->
%s
</body>
</html>`,
		darknessBanner, e.linkTags(), e.metaTags(), e.styleTags(),
		e.scriptTags(), e.htmlHead(),
		processTitle(flattenFormatting(e.page.Title)), e.authorHeader(),
		strings.Join(content, ""), e.addFootnotes(),
	)
}

// leftHeading leaves the heading.
func (e *ExporterHTML) leftHeading() {
	e.inHeading = false
}

// htmlHead builds the HTML head by also excluding any user-defined contains rules.
func (e ExporterHTML) htmlHead() string {
	htmlHeadElements := ""
	for _, v := range append(emilia.Config.Website.ExtraHead, e.page.HtmlHead...) {
		if e.page.Accoutrement.ExcludeHtmlHeadContains.ShouldExclude(v) {
			continue
		}
		htmlHeadElements += v + "\n"
	}
	return htmlHeadElements
}

// styleTags is the processed style tags.
func (e ExporterHTML) styleTags() string {
	content := make([]string, len(emilia.Config.Website.Styles)+len(e.page.Stylesheets))
	for i, style := range emilia.Config.Website.Styles {
		content[i] = fmt.Sprintf(
			`<link rel="stylesheet" type="text/css" href="%s">`+"\n", emilia.JoinPath(style),
		)
	}
	content = append(content, e.page.Stylesheets...)
	return strings.Join(content, "")
}

// defaultScripts are the default scripts.
var defaultScripts = []string{
	`<script type="module">document.documentElement.classList.remove("no-js");document.documentElement.classList.add("js");</script>`,
	`<script async src="https://sandyuraz.com/scripts/time.js"></script>`,
}

// scriptTags returns the script tags.
func (e ExporterHTML) scriptTags() string {
	allScripts := append(defaultScripts, e.page.Scripts...)
	return strings.Join(allScripts, "\n")
}

// authorHeader returns the author header.
func (e ExporterHTML) authorHeader() string {
	content := fmt.Sprintf(`
<div class="header">
<h1 class="section-1">%s%s</h1>
<div class="menu">
<span id="author" class="author">%s</span><br>
<span id="email" class="email">%s</span><br>
`,
		authorImage(e.page.Accoutrement.AuthorImage), processTitle(e.page.Title),
		emilia.Config.Author.Name, emilia.Config.Author.Email,
	)
	content += `<span id="revdate">` + "\n"

	// Build the navigation links.
	navLinks := make([]string, 0, len(emilia.Config.Navigation))

	// Go through elements.
	for i := 1; i <= len(emilia.Config.Navigation); i++ {
		// Get the navigation element read from Darkness' toml.
		v := emilia.Config.Navigation[strconv.FormatInt(int64(i), 10)]
		// If the nav element wants to hide in this location, then skip it.
		if e.page.Location == v.Hide {
			continue
		}
		// Build each of the navlinks and concat the hrefs.
		navLinks = append(navLinks, fmt.Sprintf(`<a href="%s">%s</a>`,
			emilia.JoinPathGeneric[yunyun.RelativePathDir, yunyun.FullPathDir](v.Link),
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
func authorImage(authorImageFlag yunyun.AccoutrementFlip) string {
	// Return nothing if it's not provided.
	if emilia.Config.Author.Image == "" || authorImageFlag.IsDisabled() {
		return ""
	}
	return fmt.Sprintf(`<img id="myface" src="%s" alt="avatar">`,
		emilia.Config.Author.ImagePreComputed)
}

// addTomb adds the tomb to the last paragraph.
func (e ExporterHTML) addTomb() {
	// Empty???
	if e.contentsNum < 1 {
		return
	}
	// Find the last paragrapd and attached the tomb.
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
