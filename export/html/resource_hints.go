package html

import (
	"fmt"
	"strings"
)

// resourceHints returns `<link rel="preconnect">` and `<link rel="preload">`
// hints for the page. They are emitted before the rest of the head so the
// browser can warm up connections and start fetching critical assets as early
// as possible.
func (e *state) resourceHints() []string {
	hints := make([]string, 0, 8)

	// Preconnect to the embed origins this page actually uses, so scrolling to
	// a YouTube or Spotify embed doesn't pay the connection cost at that moment.
	usesYouTube, usesSpotify := false, false
	for _, c := range e.page.Contents {
		if !c.IsLink() {
			continue
		}
		link := strings.TrimSpace(c.Link)
		usesYouTube = usesYouTube || strings.HasPrefix(link, youtubeEmbedPrefix)
		usesSpotify = usesSpotify ||
			strings.HasPrefix(link, spotifyTrackEmbedPrefix) ||
			strings.HasPrefix(link, spotifyPlaylistEmbedPrefix)
	}
	if usesYouTube {
		hints = append(hints,
			`<link rel="preconnect" href="https://www.youtube.com">`,
			`<link rel="preconnect" href="https://i.ytimg.com">`,
		)
	}
	if usesSpotify {
		hints = append(hints,
			`<link rel="preconnect" href="https://open.spotify.com">`,
			`<link rel="preconnect" href="https://i.scdn.co">`,
		)
	}

	// Preload the primary (first, local) stylesheet so it's fetched at high
	// priority as early as possible.
	for _, style := range e.conf.Website.Styles {
		if strings.HasPrefix(string(style), "http") {
			continue
		}
		hints = append(hints, fmt.Sprintf(
			`<link rel="preload" as="style" href="%s">`, e.conf.Runtime.Join(style)))
		break
	}

	return hints
}
