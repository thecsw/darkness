package html

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

const (
	// imageEmbedTemplateWithHref is the template for image embeds that's clickable.
	imageEmbedTemplateWithHref = `
<div class="media" %s>
<a class="image" href="%s"><img class="image" src="%s" title="%s" alt="%s"></a>
<div class="title">%s</div>
<hr>
</div>`

	// imageEmbedTemplateNoHref is the template of image embeds that is not clickable
	imageEmbedTemplateNoHref = `
<div class="media" %s>
<a class="image" ><img class="image" src="%s" title="%s" alt="%s"></a>
<div class="title">%s</div>
<hr>
</div>`

	// audioEmbedTemplate is the template for audio embeds.
	audioEmbedTemplate = `
<div class="media" %s>
<audio controls><source src="%s" type="audio/mpeg">music is good for the soul</audio>
</div>`

	// videoEmbedTemplate is the template for video embeds.
	videoEmbedTemplate = `
<div class="media" %s>
<video controls class="responsive-iframe">
<source src="%s" type="video/%s">
Sorry, your browser doesn't support embedded videos.
</video>
<div class="title">%s</div>
<hr>
</div>
`

	// rawHTMLTemplate wraps raw html in `mediablock`.
	rawHTMLTemplate = `
<div class="media" %s>
%s
<div class="title">%s</div>
</div>`

	// tableTemplate is the template for image embeds.
	tableTemplate = `
<div class="media" %s>
<div class="title centered">%s</div>
%s
</div>`

	// youtubeEmbedPrefix is the prefix for youtube embeds.
	youtubeEmbedPrefix = "https://youtu.be/"
	// youtubeEmbedTemplate is the template for youtube embeds.
	youtubeEmbedTemplate = `
<div class="media" %s>
<div class="yt-container">
<iframe src="https://www.youtube.com/embed/%s" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>
<hr>
</div>`

	// Put iframes here to have a youtube-embed-like experience.
	responsiveIFrameHTMLTemplate = `
<div class="media" %s>
<div class="yt-container">
%s
</div>
<hr>
</div>`

	// spotifyTrackEmbedPrefix is the prefix for spotify track embeds.
	spotifyTrackEmbedPrefix = "https://open.spotify.com/track/"
	// spotifyTrackEmbedTemplate is the template for spotify track embeds.
	spotifyTrackEmbedTemplate = `
<div class="media" %s>
<iframe class="spotify-embed-track" style="border-radius:12px" src="https://open.spotify.com/embed/track/%s?utm_source=generator" width="69%%" height="152" frameBorder="0" allowfullscreen="" allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture" loading="lazy"></iframe>
</div>`

	// spotifyPlaylistEmbedPrefix is the prefix for spotify playlist embeds.
	spotifyPlaylistEmbedPrefix = "https://open.spotify.com/playlist/"
	// spotifyPlaylistEmbedTemplate is the template for spotify playlist embeds.
	spotifyPlaylistEmbedTemplate = `
<div class="media" %s>
<iframe class="spotify-embed-playlist" style="border-radius:12px" src="https://open.spotify.com/embed/playlist/%s?utm_source=generator" width="69%%" height="550" frameBorder="0" allowfullscreen="" allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture" loading="lazy"></iframe>
</div>`
)

// Link returns an html representation of a link even if it's an embed command
func (e state) Link(content *yunyun.Content) string {
	cleanLink := strings.TrimSpace(content.Link)
	switch {
	case yunyun.ImageExtRegexp.MatchString(cleanLink) || strings.Contains(content.Attributes, "image"):
		// Put imageblocks.
		return linkImage(content, e.conf.Website.ClickableImages)
	case yunyun.AudioFileExtRegexp.MatchString(cleanLink):
		// Audiofiles
		return fmt.Sprintf(audioEmbedTemplate,
			content.CustomHtmlTags,
			cleanLink,
		)
	case yunyun.VideoFileExtRegexp.MatchString(cleanLink):
		// Raw videofiles
		return fmt.Sprintf(videoEmbedTemplate,
			content.CustomHtmlTags,
			cleanLink, func(v string) string {
				return yunyun.VideoFileExtRegexp.FindAllStringSubmatch(v, 1)[0][1]
			}(cleanLink),
			processText(content.LinkTitle),
		)
	case strings.HasPrefix(cleanLink, youtubeEmbedPrefix):
		// Youtube videos
		return fmt.Sprintf(youtubeEmbedTemplate,
			content.CustomHtmlTags,
			gana.SkipString(uint(len(youtubeEmbedPrefix)), cleanLink),
		)
	case strings.HasPrefix(cleanLink, spotifyTrackEmbedPrefix):
		// Spotify songs
		return fmt.Sprintf(spotifyTrackEmbedTemplate,
			content.CustomHtmlTags,
			gana.SkipString(uint(len(spotifyTrackEmbedPrefix)), cleanLink),
		)
	case strings.HasPrefix(cleanLink, spotifyPlaylistEmbedPrefix):
		return fmt.Sprintf(spotifyPlaylistEmbedTemplate,
			content.CustomHtmlTags,
			gana.SkipString(uint(len(spotifyPlaylistEmbedPrefix)), cleanLink),
		)
	default:
		yunyun.AddFlag(&content.Options, linkWasNotSpecialFlag)
		return fmt.Sprintf(`<a href="%s" title="%s">%s</a>`,
			cleanLink,
			yunyun.RemoveFormatting(content.LinkDescription),
			processText(content.LinkTitle),
		)
	}
}

func linkImage(content *yunyun.Content, isClickable bool) string {
	// User can elect in darkness.toml to make images clickable.
	if isClickable {
		return fmt.Sprintf(imageEmbedTemplateWithHref,
			content.CustomHtmlTags,
			content.Link,
			content.Link,
			yunyun.RemoveFormatting(content.LinkDescription),
			yunyun.RemoveFormatting(content.LinkTitle),
			processText(content.LinkTitle),
		)
	}
	// Send the embed with no clickable images. Default behavior.
	return fmt.Sprintf(imageEmbedTemplateNoHref,
		content.CustomHtmlTags,
		content.Link,
		yunyun.RemoveFormatting(content.LinkDescription),
		yunyun.RemoveFormatting(content.LinkTitle),
		processText(content.LinkTitle),
	)
}
