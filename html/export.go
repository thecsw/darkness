package html

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
)

// ExportPage exports the page to HTML
func ExportPage(page *internals.Page) string {
	// Add the red tomb to the last paragraph on given directories
	// from the config
	for _, tombPage := range emilia.Config.Website.Tombs {
		if strings.HasPrefix(page.URL, emilia.JoinPath(tombPage)) {
			addTomb(page)
			break
		}
	}

	content := make([]string, len(page.Contents))
	for i, v := range page.Contents {
		content[i] = contentFunctions[v.Type](&v)
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
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
		linkTags(page), metaTags(page), scriptTags(page),
		strings.Join(emilia.Config.Website.ExtraHead, "\n"),
		processTitle(flattenFormatting(page.Title)), styleTags(page),
		authorHeader(page), strings.Join(content, ""), addFootnotes(page),
	)
}

// styleTags is the processed style tags
func styleTags(page *internals.Page) string {
	content := make([]string, len(emilia.Config.Website.Styles)+len(page.Stylesheets))
	for i, style := range emilia.Config.Website.Styles {
		content[i] = fmt.Sprintf(
			`<link rel="stylesheet" type="text/css" href="%s">`+"\n",
			emilia.JoinPath(style),
		)
	}
	content = append(content, page.Stylesheets...)
	return strings.Join(content, "")
}

// defaultScripts are the default scripts
var defaultScripts = []string{
	`<script type="module">document.documentElement.classList.remove("no-js");document.documentElement.classList.add("js");</script>`,
	`<script async src="https://sandyuraz.com/scripts/time.js"></script>`,
}

// scriptTags returns the script tags
func scriptTags(page *internals.Page) string {
	allScripts := append(defaultScripts, page.Scripts...)
	return strings.Join(allScripts, "\n")
}

// authorHeader returns the author header
func authorHeader(page *internals.Page) string {
	content := fmt.Sprintf(`
<div class="header">
<h1>%s%s</h1>
<div class="details">
<span id="author" class="author">%s</span><br>
<span id="email" class="email">%s</span><br>
`,
		authorImage(), processTitle(page.Title),
		emilia.Config.Author.Name, emilia.Config.Author.Email,
	)

	content += `<span id="revdate">` + "\n"
	navLinks := make([]string, 0, len(emilia.Config.Navigation))
	for i := 1; i <= len(emilia.Config.Navigation); i++ {
		v := emilia.Config.Navigation[fmt.Sprintf("%d", i)]
		if page.URL == emilia.Config.URL && v.Link == v.Hide {
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

// authorHeader returns img element if author header image is given
func authorImage() string {
	// Return nothing if it's not provided
	if emilia.Config.Author.Image == "" {
		return ""
	}
	return fmt.Sprintf(`<img id="myface" src="%s" width="112" alt="Top Face" height="112">`,
		emilia.JoinPath(emilia.Config.Author.Image))
}

// addTomb adds the tomb to the last paragraph
func addTomb(page *internals.Page) {
	// Empty???
	if len(page.Contents) < 1 {
		return
	}
	last := &page.Contents[len(page.Contents)-1]
	// Only add it to paragraphs
	if !last.IsParagraph() {
		return
	}
	last.Paragraph += " ◼"
}

// addFootnotes adds the footnotes
func addFootnotes(page *internals.Page) string {
	if len(page.Footnotes) < 1 {
		return ""
	}
	footnotes := make([]string, len(page.Footnotes))
	for i, footnote := range page.Footnotes {
		footnotes[i] = fmt.Sprintf(`
<div class="footnote" id="_footnotedef_%d">
<a href="#_footnoteref_%d">%s</a>
%s
</div>
`,
			i+1, i+1, footnoteLabel(i+1), processText(footnote))
	}
	return fmt.Sprintf(`
<div id="footnotes">
<hr>
%s
</div>
`, strings.Join(footnotes, ""))
}
