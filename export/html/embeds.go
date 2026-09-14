package html

import (
	"fmt"
	stdhtml "html"
	"net/url"
	"slices"
	"strings"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/kowloon"
	"github.com/thecsw/darkness/v3/yunyun"
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

	// rawHtmlTemplate wraps raw html in `mediablock`.
	rawHtmlTemplate = `
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

	// embedReferencePrefix marks a local TOML snapshot link produced by a macro.
	embedReferencePrefix = "embed:"

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
	responsiveIFrameHtmlTemplate = `
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

	// pdfEmbedTemplate is the template for showing the PDF file on the page through an embed.
	pdfEmbedTemplate = `
<div class="media" %s>
<div class="pdf-container">
<embed src="%s" type="application/pdf" />
</div>
</div>`

	// staticEmbedCardTemplate is a self-contained representation of an external
	// resource. It deliberately contains no remote images, scripts, or players.
	staticEmbedCardTemplate = `
<div class="media embed-card embed-card--%s" %s>
<a class="embed-card__link" href="%s" target="_blank" rel="external noopener noreferrer">
<span class="embed-card__sigil" aria-hidden="true">%s</span>
<span class="embed-card__copy">
<span class="embed-card__kind">%s</span>
<strong class="embed-card__title">%s</strong>
%s</span>
<span class="embed-card__action">%s <span aria-hidden="true">↗</span></span>
</a>
</div>`
)

type embedCardKind struct {
	class         string
	sigil         string
	label         string
	action        string
	fallbackTitle string
}

var (
	youtubeCard = embedCardKind{
		class: "video", sigil: "▷", label: "moving picture · youtube",
		action: "watch", fallbackTitle: "YouTube video",
	}
	spotifyTrackCard = embedCardKind{
		class: "track", sigil: "♫", label: "record · spotify",
		action: "listen", fallbackTitle: "Spotify track",
	}
	spotifyPlaylistCard = embedCardKind{
		class: "playlist", sigil: "≋", label: "mixtape · spotify",
		action: "listen", fallbackTitle: "Spotify playlist",
	}
)

func hasEmbedAttribute(attributes, attribute string) bool {
	return slices.Contains(strings.Fields(attributes), attribute)
}

func staticEmbedCardsEnabled(conf *alpha.DarknessConfig) bool {
	return conf != nil && strings.EqualFold(strings.TrimSpace(conf.Website.EmbedMode), "card")
}

func shouldRenderStaticCard(conf *alpha.DarknessConfig, content *yunyun.Content) bool {
	if hasEmbedAttribute(content.Attributes, "embed-remote") {
		return false
	}
	return staticEmbedCardsEnabled(conf) || hasEmbedAttribute(content.Attributes, "embed-card")
}

func isWebLink(link string) bool {
	parsed, err := url.Parse(link)
	if err != nil || parsed == nil {
		return false
	}
	return (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Hostname() != ""
}

func genericEmbedCardKind(link string) embedCardKind {
	parsed, _ := url.Parse(link)
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	return embedCardKind{
		class: "link", sigil: "↗", label: host,
		action: "visit", fallbackTitle: host,
	}
}

func renderStaticEmbedCard(content *yunyun.Content, cleanLink string, kind embedCardKind) string {
	title := strings.TrimSpace(content.LinkTitle)
	if title == "" {
		title = kind.fallbackTitle
	}

	description := ""
	if caption := strings.TrimSpace(content.Caption); caption != "" {
		description = fmt.Sprintf(`<span class="embed-card__description">%s</span>`+"\n",
			stdhtml.EscapeString(yunyun.RemoveFormatting(caption)))
	}

	return fmt.Sprintf(staticEmbedCardTemplate,
		kind.class,
		content.CustomHtmlTags,
		stdhtml.EscapeString(cleanLink),
		kind.sigil,
		stdhtml.EscapeString(kind.label),
		processText(title),
		description,
		stdhtml.EscapeString(kind.action),
	)
}

// link returns an html representation of a link even if it's an embed command
func (e *state) link(content *yunyun.Content) string {
	cleanLink := strings.TrimSpace(content.Link)
	switch {
	case yunyun.ImageExtRegexp.MatchString(cleanLink) || strings.Contains(content.Attributes, "image"):
		// Put imageblocks.
		return linkImage(e.conf, content, e.conf.Website.ClickableImages)
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
	case yunyun.PdfFileExtRegexp.MatchString(cleanLink):
		return fmt.Sprintf(pdfEmbedTemplate,
			content.CustomHtmlTags,
			cleanLink,
		)
	case strings.HasPrefix(cleanLink, embedReferencePrefix):
		return e.renderSnapshotReference(content,
			strings.TrimPrefix(cleanLink, embedReferencePrefix))
	case strings.HasPrefix(cleanLink, youtubeEmbedPrefix):
		if shouldRenderStaticCard(e.conf, content) {
			return renderStaticEmbedCard(content, cleanLink, youtubeCard)
		}
		// Youtube videos
		return fmt.Sprintf(youtubeEmbedTemplate,
			content.CustomHtmlTags,
			gana.SkipString(uint(len(youtubeEmbedPrefix)), cleanLink),
		)
	case strings.HasPrefix(cleanLink, spotifyTrackEmbedPrefix):
		if shouldRenderStaticCard(e.conf, content) {
			return renderStaticEmbedCard(content, cleanLink, spotifyTrackCard)
		}
		// Spotify songs
		return fmt.Sprintf(spotifyTrackEmbedTemplate,
			content.CustomHtmlTags,
			gana.SkipString(uint(len(spotifyTrackEmbedPrefix)), cleanLink),
		)
	case strings.HasPrefix(cleanLink, spotifyPlaylistEmbedPrefix):
		if shouldRenderStaticCard(e.conf, content) {
			return renderStaticEmbedCard(content, cleanLink, spotifyPlaylistCard)
		}
		return fmt.Sprintf(spotifyPlaylistEmbedTemplate,
			content.CustomHtmlTags,
			gana.SkipString(uint(len(spotifyPlaylistEmbedPrefix)), cleanLink),
		)
	case hasEmbedAttribute(content.Attributes, "embed-card") && isWebLink(cleanLink):
		return renderStaticEmbedCard(content, cleanLink, genericEmbedCardKind(cleanLink))
	default:
		yunyun.AddFlag(&content.Options, linkWasNotSpecialFlag)
		return fmt.Sprintf(`<div %s><a href="%s" title="%s">%s</a></div>`,
			content.CustomHtmlTags,
			cleanLink,
			yunyun.RemoveFormatting(content.LinkDescription),
			processText(content.LinkTitle),
		)
	}
}

func linkImage(conf *alpha.DarknessConfig, content *yunyun.Content, isClickable bool) string {
	loadableLink := kowloon.ConvertImageToLfsMediaLink(conf, content.Link)
	// User can elect in darkness.toml to make images clickable.
	if isClickable {
		return fmt.Sprintf(imageEmbedTemplateWithHref,
			content.CustomHtmlTags,
			loadableLink,
			loadableLink,
			yunyun.RemoveFormatting(content.LinkDescription),
			yunyun.RemoveFormatting(content.LinkTitle),
			processText(content.LinkTitle),
		)
	}
	// Send the embed with no clickable images. IsDefault behavior.
	return fmt.Sprintf(imageEmbedTemplateNoHref,
		content.CustomHtmlTags,
		loadableLink,
		yunyun.RemoveFormatting(content.LinkDescription),
		yunyun.RemoveFormatting(content.LinkTitle),
		processText(content.LinkTitle),
	)
}
