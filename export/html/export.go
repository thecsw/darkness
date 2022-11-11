package html

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/emilia"
)

const (
	tombEnding = " ◼"
)

// Export runs the process of exporting
func (e ExporterHTML) Export() string {
	// Add the red tomb to the last paragraph on given directories
	// from the config
	for _, tombPage := range emilia.Config.Website.Tombs {
		if strings.HasPrefix(e.page.URL, emilia.JoinPath(tombPage)) {
			e.addTomb()
			break
		}
	}
	// Build the HTML (string) representation of each content
	content := make([]string, 0, e.contentsNum)
	for i, v := range e.page.Contents {
		content = append(content, e.buildContent(i, &v))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<!-- Links -->
%s
<!-- Meta -->
%s
<!-- Scripts -->
%s
<!-- Extra -->
%s
<!-- Title -->
<title>%s</title>
</head>
<!-- Styling -->
%s
<body class="article">
<!-- Header -->
%s
<!-- Content -->
%s
<!-- Footnotes -->
%s
</body>
</html>`,
		e.linkTags(), e.metaTags(), e.scriptTags(),
		strings.Join(emilia.Config.Website.ExtraHead, "\n"),
		processTitle(flattenFormatting(e.page.Title)), e.styleTags(),
		e.authorHeader(), strings.Join(content, ""), e.addFootnotes(),
	)
}

// leftHeading leaves the heading.
func (e ExporterHTML) leftHeading() {
	e.inHeading = false
}

// styleTags is the processed style tags.
func (e ExporterHTML) styleTags() string {
	content := make([]string, len(emilia.Config.Website.Styles)+len(e.page.Stylesheets))
	for i, style := range emilia.Config.Website.Styles {
		content[i] = fmt.Sprintf(
			`<link rel="stylesheet" type="text/css" href="%s">`+"\n",
			emilia.JoinPath(style),
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
<h1>%s%s</h1>
<div class="details">
<span id="author" class="author">%s</span><br>
<span id="email" class="email">%s</span><br>
`,
		authorImage(), processTitle(e.page.Title),
		emilia.Config.Author.Name, emilia.Config.Author.Email,
	)

	content += `<span id="revdate">` + "\n"
	navLinks := make([]string, 0, len(emilia.Config.Navigation))
	for i := 1; i <= len(emilia.Config.Navigation); i++ {
		v := emilia.Config.Navigation[fmt.Sprintf("%d", i)]
		if e.page.URL == emilia.Config.URL && v.Link == v.Hide {
			continue
		}
		navLinks = append(navLinks, fmt.Sprintf(`<a href="%s">%s</a>`, emilia.JoinPath(v.Link), v.Title))
	}
	content += strings.Join(navLinks, " | ") + `</span>`
	content += `
</div>
<div id="hetime" class="details"></div>
</div>`

	return content
}

// authorHeader returns img element if author header image is given.
func authorImage() string {
	// Return nothing if it's not provided
	if emilia.Config.Author.Image == "" {
		return ""
	}
	return fmt.Sprintf(`<img id="myface" src="%s" width="112" alt="Top Face" height="112">`,
		emilia.JoinPath(emilia.Config.Author.Image))
}

// addTomb adds the tomb to the last paragraph.
func (e ExporterHTML) addTomb() {
	// Empty???
	if e.contentsNum < 1 {
		return
	}
	last := &e.page.Contents[e.contentsNum-1]
	// Only add it to paragraphs
	if !last.IsParagraph() {
		return
	}
	last.Paragraph += tombEnding
}
