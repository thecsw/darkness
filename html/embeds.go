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
<a class="image" href="%s"><img src="%s" alt="%s"></a>
</div>
<div class="title">%s</div>
</div>
<hr>`
	// AudioEmbedTemplate is the template for audio embeds
	AudioEmbedTemplate = `
<audio controls><source src="%s" type="audio/mpeg">music is good for the soul</audio>`

	// YoutubeEmbedPrefix is the prefix for youtube embeds
	YoutubeEmbedPrefix = "https://youtu.be/"
	// YoutubeEmbedTemplate is the template for youtube embeds
	YoutubeEmbedTemplate = `
<iframe width="100%%" height="330px" src="https://www.youtube.com/embed/%s" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
<hr>`

	// SpotifyTrackEmbedPrefix is the prefix for spotify track embeds
	SpotifyTrackEmbedPrefix = "https://open.spotify.com/track/"
	// SpotifyTrackEmbedTemplate is the template for spotify track embeds
	SpotifyTrackEmbedTemplate = `
<iframe src="https://open.spotify.com/embed/track/%s" width="79%%" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
<hr>`

	// SpotifyPlaylistEmbedPrefix is the prefix for spotify playlist embeds
	SpotifyPlaylistEmbedPrefix = "https://open.spotify.com/playlist/"
	// SpotifyPlaylistEmbedTemplate is the template for spotify playlist embeds
	SpotifyPlaylistEmbedTemplate = `
<iframe src="https://open.spotify.com/embed/playlist/%s" width="79%%" height="380" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>
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
