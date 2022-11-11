package html

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

const (
	// imageEmbedTemplate is the template for image embeds
	imageEmbedTemplate = `
<div class="imageblock">
<hr>
<a class="image" href="%s"><img class="image" src="%s" alt="%s"></a>
<div class="title">%s</div>
<hr>
</div>`
	// audioEmbedTemplate is the template for audio embeds
	audioEmbedTemplate = `
<audio controls><source src="%s" type="audio/mpeg">music is good for the soul</audio>`

	// videoEmbedTemplate is the template for video embeds
	videoEmbedTemplate = `
<div class="videoblock">
<hr>
<video controls class="responsive-iframe">
<source src="%s" type="video/%s">
Sorry, your browser doesn't support embedded videos.
</video>
<div class="title">%s</div>
<hr>
</div>
`

	// rawHTMLTemplate wraps raw html in `mediablock`
	rawHTMLTemplate = `
<div class="mediablock">
<hr>
%s
<div class="title">%s</div>
<hr>
</div>`

	// tableTemplate is the template for image embeds
	tableTemplate = `
<div class="mediablock">
<div class="title">%s</div>
%s
</div>`

	// youtubeEmbedPrefix is the prefix for youtube embeds
	youtubeEmbedPrefix = "https://youtu.be/"
	// youtubeEmbedTemplate is the template for youtube embeds
	youtubeEmbedTemplate = `
<div class="videoblock">
<hr>
<iframe class="responsive-iframe" src="https://www.youtube.com/embed/%s" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
<hr>
</div>`

	// spotifyTrackEmbedPrefix is the prefix for spotify track embeds
	spotifyTrackEmbedPrefix = "https://open.spotify.com/track/"
	// spotifyTrackEmbedTemplate is the template for spotify track embeds
	spotifyTrackEmbedTemplate = `
<div class="mediablock">
<iframe src="https://open.spotify.com/embed/track/%s" width="79%%" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
<hr>
</div>`

	// spotifyPlaylistEmbedPrefix is the prefix for spotify playlist embeds
	spotifyPlaylistEmbedPrefix = "https://open.spotify.com/playlist/"
	// spotifyPlaylistEmbedTemplate is the template for spotify playlist embeds
	spotifyPlaylistEmbedTemplate = `
<div class="mediablock">
<iframe src="https://open.spotify.com/embed/playlist/%s" width="79%%" height="380" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
<hr>
</div>`
)

// link returns an html representation of a link even if it's an embed command
func (e *ExporterHTML) link(content *yunyun.Content) string {
	switch {
	case yunyun.ImageExtRegexp.MatchString(content.Link):
		return fmt.Sprintf(imageEmbedTemplate, content.Link, content.Link,
			content.LinkTitle, processText(content.LinkTitle))
	case yunyun.AudioFileExtRegexp.MatchString(content.Link):
		return fmt.Sprintf(audioEmbedTemplate, content.Link)
	case yunyun.VideoFileExtRegexp.MatchString(content.Link):
		return fmt.Sprintf(videoEmbedTemplate, content.Link, func(v string) string {
			return yunyun.VideoFileExtRegexp.FindAllStringSubmatch(v, 1)[0][1]
		}(content.Link), content.LinkTitle)
	case strings.HasPrefix(content.Link, youtubeEmbedPrefix):
		return fmt.Sprintf(youtubeEmbedTemplate, content.Link[len(youtubeEmbedPrefix):])
	case strings.HasPrefix(content.Link, spotifyTrackEmbedPrefix):
		return fmt.Sprintf(spotifyTrackEmbedTemplate, content.Link[len(spotifyTrackEmbedPrefix):])
	case strings.HasPrefix(content.Link, spotifyPlaylistEmbedPrefix):
		return fmt.Sprintf(spotifyPlaylistEmbedTemplate, content.Link[len(spotifyPlaylistEmbedPrefix):])
	default:
		return fmt.Sprintf(`<a href="%s">%s</a>`, content.Link, processText(content.LinkTitle))
	}
}
