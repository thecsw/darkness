package html

import (
	"darkness/internals"
	"fmt"
	"html"
)

var (
	contentFunctions = []func(*internals.Content) string{
		headings, paragraph, list, listNumbered,
		link, image, youtube, spotifyTrack,
		spotifyPlaylist,
	}
)

var headerCounter = 0

func headings(content *internals.Content) string {
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
		content.HeaderLevel,
		html.EscapeString(content.Header),
		htmlize(content.Header),
		content.HeaderLevel,
	)
}

func paragraph(content *internals.Content) string {
	text := htmlize(content.Paragraph)
	return fmt.Sprintf(
		`
<div class="paragraph">
<p>%s</p>
</div>`,
		text,
	)
}

func list(content *internals.Content) string {
	start := `<div class="ulist">
<ul>`
	for _, item := range content.List {
		start += `
<li>
<p>
` + htmlize(item) + `
</p>
</li>`
	}
	start += `</ul>
</div>`
	return start
}

func listNumbered(content *internals.Content) string {
	return ""
}

func link(content *internals.Content) string {
	return ""
}

func image(content *internals.Content) string {
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

func youtube(content *internals.Content) string {
	return fmt.Sprintf(`
<iframe width="100%%" height="330px" src="https://www.youtube.com/embed/%s" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
<hr>`, content.Youtube)
}

func spotifyTrack(content *internals.Content) string {
	return fmt.Sprintf(`
<iframe src="https://open.spotify.com/embed/track/%s" width="79%%" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
<hr>`, content.SpotifyTrack)
}

func spotifyPlaylist(content *internals.Content) string {
	return fmt.Sprintf(`
<iframe src="https://open.spotify.com/embed/playlist/%s" width="79%%" height="380" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
<hr>`, content.SpotifyPlaylist)
}
