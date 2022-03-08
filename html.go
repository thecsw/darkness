package main

import (
	"fmt"
	"html"
)

var (
	contentFunctions = []func(*Content) string{
		addHeader, addParagraph, addList, addListNumbered,
		addLink, addImage, addYoutube, addSpotifyTrack,
		addSpotifyPlaylist,
	}
)

func buildHTML(page *Page) string {
	finalContents := `
<!DOCTYPE html>
<html lang="en">
<head>
`
	finalContents += addRelTags(page)
	finalContents += addMetaTags(page)
	finalContents += addScripts(page)

	finalContents += `<title>` + page.Title + `</title>` + "\n"
	finalContents += `</head>` + "\n"

	finalContents += addStyles()
	finalContents += `<body class="article">`

	finalContents += addAuthorHeader(page)
	finalContents += `<div id="content">`

	for _, content := range page.Contents {
		finalContents += contentFunctions[content.Type](&content)
	}

	finalContents += "\n" + `</div>` + "\n" + `</body>`
	return finalContents + "\n" + `</html>`
}

func addStyles() string {
	content := ""
	for _, style := range conf.Website.Styles {
		content += fmt.Sprintf(
			`<link rel="stylesheet" type="text/css" href="%s">`+"\n",
			style,
		)
	}
	return content
}

func addScripts(page *Page) string {
	scripts := `<script type="module">document.documentElement.classList.remove("no-js");document.documentElement.classList.add("js");</script>
<script async src="https://sandyuraz.com//scripts/time.js"></script>
`
	return scripts
}

func addAuthorHeader(page *Page) string {
	content := fmt.Sprintf(`
<div id="header">
<h1><img id="myface" src="%s" width="100">%s</h1>
<div class="details">
<span id="author" class="author">%s</span><br>
<span id="email" class="email">%s</span><br>
`,
		conf.Author.Header, html.EscapeString(page.Title),
		conf.Author.Name, conf.Author.Email,
	)

	content += `<span id="revdate">` + "\n"
	for i := 1; i <= len(conf.Navigation); i++ {
		v := conf.Navigation[fmt.Sprintf("%d", i)]
		content += fmt.Sprintf(
			`<a href="%s">%s</a>`,
			v.Link, v.Title,
		)
		if i < len(conf.Navigation) {
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

var headerCounter = 0

func addHeader(content *Content) string {
	start := ``
	if headerCounter > 0 {
		start = `</div>` + "\n" + `</div>`
	}
	headerCounter++
	return fmt.Sprintf(start+`
<div class="sect%d">
<h%d id="%s">%s</h%d>
<div class="sectionbody">`,
		content.HeaderLevel,
		content.HeaderLevel, content.Header, content.Header, content.HeaderLevel,
	)
}

func addParagraph(content *Content) string {
	text := orgToHTML(content.Paragraph)
	return fmt.Sprintf(
		`
<div class="paragraph">
<p>%s</p>
</div>`,
		text,
	)
}

func addList(content *Content) string {
	start := `<div class="ulist">
<ul>`
	for _, item := range content.List {
		start += `
<li>
<p>
` + orgToHTML(item) + `
</p>
</li>`
	}
	start += `</ul>
</div>`
	return start
}

func addListNumbered(content *Content) string {
	return ""
}

func addLink(content *Content) string {
	return ""
}

func addImage(content *Content) string {
	return fmt.Sprintf(`
<hr>
<div class="imageblock">
<div class="content">
<a class="image" href="%s"><img src="%s" alt="%s"></a>
</div>
<div class="title">%s</div>
</div>
<hr>`, content.Link, content.Link, content.LinkTitle, content.LinkTitle)
}

func addYoutube(content *Content) string {
	return fmt.Sprintf(`
<iframe width="100%%" height="330px" src="https://www.youtube.com/embed/%s" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
<hr>`, content.Youtube)
}

func addSpotifyTrack(content *Content) string {
	return fmt.Sprintf(`
<iframe src="https://open.spotify.com/embed/track/%s" width="79%%" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
<hr>`, content.SpotifyTrack)
}

func addSpotifyPlaylist(content *Content) string {
	return fmt.Sprintf(`
<iframe src="https://open.spotify.com/embed/playlist/%s" width="79%%" height="380" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
<hr>`, content.SpotifyPlaylist)
}

func orgToHTML(text string) string {
	text = html.EscapeString(text)
	text = OrgLinkRegexp.ReplaceAllString(text, `<a href="$1">$2</a>`)
	text = OrgBoldText.ReplaceAllString(text, ` <strong>$1</strong>$2`)
	text = OrgItalicText.ReplaceAllString(text, ` <em>$1</em>$2`)
	return text
}
