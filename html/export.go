package html

import (
	"darkness/emilia"
	"darkness/internals"
	"fmt"
	"html"
	"strings"
)

func ExportPage(page *internals.Page) string {
	// Add the red tomb to the last paragraph
	addTomb(page)

	content := make([]string, 0, len(page.Contents))
	for _, v := range page.Contents {
		content = append(content, contentFunctions[v.Type](&v))
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
%s
%s
%s
<title>%s</title>
</head>
%s
<body class="article">
%s
<div id="content">
%s
</div>
%s
</body>
</html>
`,
		linkTags(page), metaTags(page), scriptTags(page),
		processTitle(page.Title), styleTagsProcessed, authorHeader(page), strings.Join(content, ""),
		addFootnotes(page),
	)
}

func styleTags() string {
	content := make([]string, 0, len(emilia.Config.Website.Styles))
	for _, style := range emilia.Config.Website.Styles {
		content = append(content, fmt.Sprintf(
			`<link rel="stylesheet" type="text/css" href="%s">`+"\n",
			emilia.JoinPath(style),
		))
	}
	return strings.Join(content, "")
}

func scriptTags(page *internals.Page) string {
	scripts := `<script type="module">document.documentElement.classList.remove("no-js");document.documentElement.classList.add("js");</script>
<script async src="https://sandyuraz.com//scripts/time.js"></script>
`
	return scripts
}

func authorHeader(page *internals.Page) string {
	content := fmt.Sprintf(`
<div id="header">
<h1><img id="myface" src="%s" width="112">%s</h1>
<div class="details">
<span id="author" class="author">%s</span><br>
<span id="email" class="email">%s</span><br>
`,
		emilia.JoinPath(emilia.Config.Author.Header), html.EscapeString(processTitle(page.Title)),
		emilia.Config.Author.Name, emilia.Config.Author.Email,
	)

	content += `<span id="revdate">` + "\n"
	for i := 1; i <= len(emilia.Config.Navigation); i++ {
		v := emilia.Config.Navigation[fmt.Sprintf("%d", i)]
		homeIsOnHome := page.URL == emilia.Config.URL && v.Link == ""
		if homeIsOnHome {
			// Remove the extra pipe character
			content = content[:len(content)-2]
			continue
		}
		content += fmt.Sprintf(
			`<a href="%s">%s</a>`,
			emilia.JoinPath(v.Link), v.Title,
		)
		if i < len(emilia.Config.Navigation) {
			content += ` | `
		}
	}
	content += `</span>`

	content += `
</div>
<div id="hetime" class="details"></div>
</div>`

	return content
}

func addTomb(page *internals.Page) {
	// Empty???
	if len(page.Contents) < 1 {
		return
	}
	last := &page.Contents[len(page.Contents)-1]
	// Onnly add it to paragraphs
	if last.Type != internals.TypeParagraph {
		return
	}
	last.Paragraph += " ◼"
}

func addFootnotes(page *internals.Page) string {
	if len(page.Footnotes) < 1 {
		return ""
	}
	footnotes := make([]string, 0, len(page.Footnotes))
	for i, footnote := range page.Footnotes {
		footnotes = append(footnotes, fmt.Sprintf(`
<div class="footnote" id="_footnotedef_%d">
<a href="#_footnotedef_%d">%d</a>
%s
</div>
`,
			i+1, i+1, i+1, processText(footnote)))
	}
	return fmt.Sprintf(`
<div id="footnotes">
<hr>
%s
</div>
`, strings.Join(footnotes, ""))
}
