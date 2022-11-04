package html

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/internals"
)

const (
	// ImageEmbedTemplate is the template for image embeds
	ImageEmbedTemplate = `
<hr>
<div class="imageblock">
<div class="content">
<a class="image" href="%s"><img class="image" src="%s" alt="%s"></a>
</div>
<div class="title">%s</div>
</div>
<hr>`
	// AudioEmbedTemplate is the template for audio embeds
	AudioEmbedTemplate = `
<audio controls><source src="%s" type="audio/mpeg">music is good for the soul</audio>`

	// VideoEmbedTemplate is the template for video embeds
	VideoEmbedTemplate = `
<div class="videoblock">
<div class="content">
<video controls class="responsive-iframe">
<source src="%s" type="video/%s">
Sorry, your browser doesn't support embedded videos.
</video>
</div>
<div class="title">%s</div>
</div>
<hr>
`
	// Table is the template for image embeds
	TableTemplate = `
<div class="imageblock">
<div class="title">%s</div>
<div class="content">
%s
</div>
</div>`

	// YoutubeEmbedPrefix is the prefix for youtube embeds
	YoutubeEmbedPrefix = "https://youtu.be/"
	// YoutubeEmbedTemplate is the template for youtube embeds
	YoutubeEmbedTemplate = `
<div class="videoblock">
<iframe class="responsive-iframe" src="https://www.youtube.com/embed/%s" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>
<hr>`

	// SpotifyTrackEmbedPrefix is the prefix for spotify track embeds
	SpotifyTrackEmbedPrefix = "https://open.spotify.com/track/"
	// SpotifyTrackEmbedTemplate is the template for spotify track embeds
	SpotifyTrackEmbedTemplate = `
<center><iframe src="https://open.spotify.com/embed/track/%s" width="79%%" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe></center>
<hr>`

	// SpotifyPlaylistEmbedPrefix is the prefix for spotify playlist embeds
	SpotifyPlaylistEmbedPrefix = "https://open.spotify.com/playlist/"
	// SpotifyPlaylistEmbedTemplate is the template for spotify playlist embeds
	SpotifyPlaylistEmbedTemplate = `
<center><iframe src="https://open.spotify.com/embed/playlist/%s" width="79%%" height="380" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe></center>
<hr>`
)

// link returns an html representation of a link even if it's an embed command
func link(content *internals.Content) string {
	switch {
	case internals.ImageExtRegexp.MatchString(content.Link):
		return fmt.Sprintf(ImageEmbedTemplate, content.Link, content.Link,
			content.LinkTitle, processText(content.LinkTitle))
	case internals.AudioFileExtRegexp.MatchString(content.Link):
		return fmt.Sprintf(AudioEmbedTemplate, content.Link)
	case internals.VideoFileExtRegexp.MatchString(content.Link):
		return fmt.Sprintf(VideoEmbedTemplate, content.Link, func(v string) string {
			return internals.VideoFileExtRegexp.FindAllStringSubmatch(v, 1)[0][1]
		}(content.Link), content.LinkTitle)
	case strings.HasPrefix(content.Link, YoutubeEmbedPrefix):
		return fmt.Sprintf(YoutubeEmbedTemplate, content.Link[len(YoutubeEmbedPrefix):])
	case strings.HasPrefix(content.Link, SpotifyTrackEmbedPrefix):
		return fmt.Sprintf(SpotifyTrackEmbedTemplate, content.Link[len(SpotifyTrackEmbedPrefix):])
	case strings.HasPrefix(content.Link, SpotifyPlaylistEmbedPrefix):
		return fmt.Sprintf(SpotifyPlaylistEmbedTemplate, content.Link[len(SpotifyPlaylistEmbedPrefix):])
	default:
		return fmt.Sprintf(`<a href="%s">%s</a>`, content.Link, processText(content.LinkTitle))
	}
}
