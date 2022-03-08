package html

import (
	"darkness/emilia"
	"darkness/internals"
	"fmt"
	"html"
)

func ExportPage(page *internals.Page) string {
	headerCounter = 0

	content := ""
	for _, v := range page.Contents {
		content += contentFunctions[v.Type](&v)
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
</body>
</html>
`,
		linkTags(page), metaTags(page), scriptTags(page),
		page.Title, styleTags(), authorHeader(page), content,
	)
}

func styleTags() string {
	content := ""
	for _, style := range emilia.Config.Website.Styles {
		content += fmt.Sprintf(
			`<link rel="stylesheet" type="text/css" href="%s">`+"\n",
			style,
		)
	}
	return content
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
<h1><img id="myface" src="%s" width="100">%s</h1>
<div class="details">
<span id="author" class="author">%s</span><br>
<span id="email" class="email">%s</span><br>
`,
		emilia.Config.Author.Header, html.EscapeString(page.Title),
		emilia.Config.Author.Name, emilia.Config.Author.Email,
	)

	content += `<span id="revdate">` + "\n"
	for i := 1; i <= len(emilia.Config.Navigation); i++ {
		v := emilia.Config.Navigation[fmt.Sprintf("%d", i)]
		content += fmt.Sprintf(
			`<a href="%s">%s</a>`,
			v.Link, v.Title,
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
