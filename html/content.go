package html

import (
	"fmt"
	"html"

	"github.com/thecsw/darkness/internals"
)

var (
	contentFunctions = []func(*internals.Content) string{
		headings, paragraph, list, listNumbered,
		link, image, youtube, spotifyTrack,
		spotifyPlaylist, sourceCode, rawHTML,
		horizontalLine, attentionBlock, audioFile,
	}
)

func headings(content *internals.Content) string {
	start := ``
	if !content.HeadingChild && !content.HeadingFirst {
		start = `</div>` + "\n" + `</div>`
	}
	if content.HeadingChild && !content.HeadingFirst {
		start = `</div>`
	}
	toReturn := fmt.Sprintf(start+`
<div class="sect%d">
<h%d id="%s">%s</h%d>
<div class="sectionbody">`,
		content.HeadingLevel-1,
		content.HeadingLevel,
		html.EscapeString(content.Heading),
		processText(content.Heading),
		content.HeadingLevel,
	)
	if content.HeadingLast {
		toReturn += "\n" + `</div>`
	}
	return toReturn
}

func paragraph(content *internals.Content) string {
	text := processText(content.Paragraph)
	return fmt.Sprintf(
		`
<div class="paragraph">
<p>%s</p>
</div>`,
		text,
	)
}

func list(content *internals.Content) string {
	elements := ""
	for _, item := range content.List {
		elements += fmt.Sprintf(`
<li>
<p>
%s
</p>
</li>
`, processText(item))
	}
	return fmt.Sprintf(`
<div class="ulist">
<ul>
%s
</ul>
</div>
`, elements)
}

func listNumbered(content *internals.Content) string {
	// TODO
	return ""
}

func link(content *internals.Content) string {
	// TODO
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

func sourceCode(content *internals.Content) string {
	return fmt.Sprintf(`
<div class="listingblock">
<div class="content">
<pre class="highlight">
<code class="language-%s" data-lang="%s">%s</code>
</pre>
</div>
</div>
`, content.SourceCodeLang, content.SourceCodeLang, content.SourceCode)
}

func rawHTML(content *internals.Content) string {
	return content.RawHTML
}

func horizontalLine(content *internals.Content) string {
	return `<hr>`
}

func attentionBlock(content *internals.Content) string {
	return fmt.Sprintf(`
<div class="admonitionblock note">
<table>
<tr>
<td class="icon">
<div class="title">%s</div>
</td>
<td class="content">
%s
</td>
</tr>
</table>
</div>`, content.AttentionTitle, processText(content.AttentionText))
}

func audioFile(content *internals.Content) string {
	return fmt.Sprintf(`
<audio controls><source src="%s" type="audio/mpeg">music is good for the soul</audio>
`, content.AudioFile)
}
