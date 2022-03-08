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
<meta charset="UTF-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta name="generator" content="Darkness">
<meta name="author" content="email">
`

	finalContents += addOpenGraph(page)

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

type meta struct {
	Name     string
	Property string
	Content  string
}

func addOpenGraph(page *Page) string {
	content := ""
	toAdd := []meta{
		{
			Name:     "og:title",
			Property: "og:title",
			Content:  page.Title,
		},
		{
			Name:     "og:site_name",
			Property: "og:site_name",
			Content:  html.EscapeString(conf.Title),
		},
		{
			Name:     "og:url",
			Property: "og:url",
			Content:  page.URL,
		},
		{
			Name:     "og:locale",
			Property: "og:locale",
			Content:  conf.Website.Locale,
		},
		{
			Name:     "og:type",
			Property: "og:type",
			Content:  "website",
		},
	}
	for _, add := range toAdd {
		content += metaTag(add)
	}
	return content
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

func metaTag(val meta) string {
	return fmt.Sprintf(
		`<meta name="%s" property="%s" content="%s">`+"\n",
		val.Name, val.Property, html.EscapeString(val.Content),
	)
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
<div class="sect%d>"
<h%d id="%s">%s</h%d>
<div class="sectionbody">
`,
		content.HeaderLevel,
		content.HeaderLevel, content.Header, content.Header, content.HeaderLevel,
	)
}

func addParagraph(content *Content) string {
	text := html.EscapeString(content.Paragraph)
	text = OrgLinkRegexp.ReplaceAllString(text, `<a href="$1">$2</a>`)
	text = OrgBoldText.ReplaceAllString(text, ` <strong>$1</strong>$2`)
	text = OrgItalicText.ReplaceAllString(text, ` <em>$1</em>$2`)

	return fmt.Sprintf(
		`
<div class="paragraph">
<p>%s</p>
</div>`,
		text,
	)
}

func addList(content *Content) string {
	return ""
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
<hr>
`, content.Link, content.Link, content.LinkTitle, content.LinkTitle)
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
