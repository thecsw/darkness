package html

import (
	"fmt"
	stdhtml "html"
	"path"
	"regexp"
	"strings"

	"github.com/thecsw/darkness/v3/emilia/memento"
	"github.com/thecsw/darkness/v3/yunyun"
)

var (
	snapshotColorRegexp = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	snapshotClassRegexp = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)
)

const (
	spotifyMark = `<svg class="embed-provider-mark" viewBox="0 0 24 24" aria-hidden="true"><path d="M12 1a11 11 0 1 0 0 22 11 11 0 0 0 0-22Zm5.04 15.86a.69.69 0 0 1-.95.23c-2.6-1.59-5.88-1.95-9.74-1.07a.69.69 0 1 1-.31-1.34c4.22-.96 7.85-.55 10.77 1.23.32.2.43.62.23.95Zm1.35-3a.86.86 0 0 1-1.18.28c-2.98-1.83-7.52-2.36-11.04-1.29a.86.86 0 1 1-.5-1.65c4.03-1.22 9.03-.63 12.44 1.46.4.25.53.78.28 1.2Zm.12-3.12C14.94 8.62 9.04 8.42 5.63 9.45a1.03 1.03 0 1 1-.6-1.97c3.92-1.19 10.44-.95 14.54 1.48a1.03 1.03 0 0 1-1.06 1.78Z"/></svg>`
	youtubeMark = `<svg class="embed-provider-mark" viewBox="0 0 28 20" aria-hidden="true"><path d="M27.4 3.1A3.5 3.5 0 0 0 25 .6C22.8 0 18.5 0 14 0S5.2 0 3 .6A3.5 3.5 0 0 0 .6 3.1C0 5.3 0 10 0 10s0 4.7.6 6.9A3.5 3.5 0 0 0 3 19.4c2.2.6 6.5.6 11 .6s8.8 0 11-.6a3.5 3.5 0 0 0 2.4-2.5c.6-2.2.6-6.9.6-6.9s0-4.7-.6-6.9Z"/><path class="embed-provider-mark__cutout" d="m11.2 14.3 7.3-4.3-7.3-4.3v8.6Z"/></svg>`
)

func (e *state) renderSnapshotReference(content *yunyun.Content, reference string) string {
	pageLocation := ""
	if e.page != nil {
		pageLocation = string(e.page.Location)
	}
	snapshot, err := memento.LoadReference(e.conf, pageLocation, reference)
	if err != nil {
		if e.conf != nil && e.conf.Runtime.Logger != nil {
			e.conf.Runtime.Logger.Warn("Could not load embed manifest", "manifest", reference, "err", err)
		}
		return fmt.Sprintf(`
<div class="media embed-card embed-card--missing" %s>
<span class="embed-card__copy">
<span class="embed-card__kind">embed manifest unavailable</span>
<strong class="embed-card__title">%s</strong>
</span>
</div>`, content.CustomHtmlTags, stdhtml.EscapeString(reference))
	}
	return renderSnapshot(snapshot, content.CustomHtmlTags)
}

func renderSnapshot(snapshot *memento.Snapshot, customTags string) string {
	switch snapshot.Kind {
	case memento.KindYouTube:
		return renderYouTubeSnapshot(snapshot, customTags)
	case memento.KindSpotifyTrack:
		return renderSpotifyTrackSnapshot(snapshot, customTags)
	case memento.KindSpotifyPlaylist:
		return renderSpotifyPlaylistSnapshot(snapshot, customTags)
	default:
		return ""
	}
}

func renderYouTubeSnapshot(snapshot *memento.Snapshot, customTags string) string {
	return fmt.Sprintf(`
<div class="media embed-snapshot embed-youtube%s" %s>
<a class="embed-snapshot__link embed-youtube__link" href="%s" target="_blank" rel="external noopener noreferrer">
<div class="embed-youtube__frame">
%s
<span class="embed-youtube__play" aria-hidden="true"><span></span></span>
</div>
<div class="embed-youtube__meta">
<span class="embed-youtube__brand">%s <span>YouTube</span></span>
<strong class="embed-youtube__title">%s</strong>
<span class="embed-youtube__creator">%s</span>
</div>
</a>
</div>`,
		snapshotClasses(snapshot),
		customTags,
		stdhtml.EscapeString(snapshot.URL),
		renderSnapshotArtwork(snapshot, "embed-youtube__thumbnail", snapshot.Title),
		youtubeMark,
		stdhtml.EscapeString(snapshot.Title),
		stdhtml.EscapeString(snapshot.Creator),
	)
}

func renderSpotifyTrackSnapshot(snapshot *memento.Snapshot, customTags string) string {
	return fmt.Sprintf(`
<div class="media embed-snapshot embed-spotify embed-spotify--track%s" %s>
<a class="embed-snapshot__link embed-spotify__track-link" href="%s" target="_blank" rel="external noopener noreferrer"%s>
%s
<div class="embed-spotify__track-copy">
<span class="embed-spotify__brand">%s <span>Spotify</span></span>
<strong class="embed-spotify__title">%s%s</strong>
<span class="embed-spotify__creator">%s</span>
<div class="embed-spotify__controls" aria-hidden="true">
<span class="embed-spotify__progress"><i></i></span>
<span class="embed-spotify__time">%s</span>
<span class="embed-spotify__play"><i></i></span>
</div>
</div>
</a>
</div>`,
		snapshotClasses(snapshot),
		customTags,
		stdhtml.EscapeString(snapshot.URL),
		snapshotColorStyle(snapshot.Colors),
		renderSnapshotArtwork(snapshot, "embed-spotify__artwork", snapshot.Title+" cover"),
		spotifyMark,
		stdhtml.EscapeString(snapshot.Title),
		explicitBadge(snapshot.Explicit),
		stdhtml.EscapeString(snapshot.Creator),
		formatDuration(snapshot.DurationMS),
	)
}

func renderSpotifyPlaylistSnapshot(snapshot *memento.Snapshot, customTags string) string {
	var tracks strings.Builder
	for index, track := range snapshot.Tracks {
		fmt.Fprintf(&tracks, `
<li class="embed-spotify__row">
<a class="embed-snapshot__link" href="%s" target="_blank" rel="external noopener noreferrer">
<span class="embed-spotify__index">%d</span>
<span class="embed-spotify__song"><strong>%s%s</strong><small>%s</small></span>
<time datetime="PT%dS">%s</time>
</a>
</li>`,
			stdhtml.EscapeString(track.URL),
			index+1,
			stdhtml.EscapeString(track.Title),
			explicitBadge(track.Explicit),
			stdhtml.EscapeString(track.Artist),
			track.DurationMS/1000,
			formatDuration(track.DurationMS),
		)
	}
	return fmt.Sprintf(`
<section class="media embed-snapshot embed-spotify embed-spotify--playlist%s" %s>
<div class="embed-spotify__surface"%s>
<a class="embed-snapshot__link embed-spotify__playlist-head" href="%s" target="_blank" rel="external noopener noreferrer">
%s
<span class="embed-spotify__playlist-copy">
<span class="embed-spotify__brand">%s <span>Spotify</span></span>
<strong class="embed-spotify__title">%s</strong>
<span class="embed-spotify__creator">%s · %d songs</span>
</span>
<span class="embed-spotify__play" aria-hidden="true"><i></i></span>
</a>
<ol class="embed-spotify__tracklist" aria-label="Tracks in %s">%s
</ol>
<a class="embed-snapshot__link embed-spotify__footer" href="%s" target="_blank" rel="external noopener noreferrer">open in Spotify <span aria-hidden="true">↗</span></a>
</div>
</section>`,
		snapshotClasses(snapshot),
		customTags,
		snapshotColorStyle(snapshot.Colors),
		stdhtml.EscapeString(snapshot.URL),
		renderSnapshotArtwork(snapshot, "embed-spotify__artwork", snapshot.Title+" cover"),
		spotifyMark,
		stdhtml.EscapeString(snapshot.Title),
		stdhtml.EscapeString(snapshot.Creator),
		len(snapshot.Tracks),
		stdhtml.EscapeString(snapshot.Title),
		tracks.String(),
		stdhtml.EscapeString(snapshot.URL),
	)
}

func renderSnapshotArtwork(snapshot *memento.Snapshot, class, alt string) string {
	artworkSource := snapshot.ResolvedArtwork
	if artworkSource == "" {
		artworkSource = snapshot.Artwork
	}
	artwork, ok := localArtworkPath(artworkSource)
	if !ok {
		return `<span class="` + class + ` embed-snapshot__missing-art" aria-hidden="true"></span>`
	}
	return fmt.Sprintf(`<img class="%s" src="/%s" alt="%s" loading="lazy" decoding="async">`,
		class, stdhtml.EscapeString(artwork), stdhtml.EscapeString(alt))
}

func localArtworkPath(artwork string) (string, bool) {
	if strings.Contains(artwork, "://") || strings.HasPrefix(artwork, "/") {
		return "", false
	}
	clean := path.Clean(strings.TrimSpace(artwork))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", false
	}
	return clean, true
}

func snapshotColorStyle(colors memento.Colors) string {
	values := []struct {
		name  string
		value string
	}{
		{"--embed-bg", colors.Background},
		{"--embed-bg-deep", colors.BackgroundTinted},
		{"--embed-fg", colors.Text},
		{"--embed-muted", colors.TextSubdued},
	}
	var declarations strings.Builder
	for _, value := range values {
		if snapshotColorRegexp.MatchString(value.value) {
			fmt.Fprintf(&declarations, "%s:%s;", value.name, value.value)
		}
	}
	if declarations.Len() == 0 {
		return ""
	}
	return ` style="` + declarations.String() + `"`
}

func snapshotClasses(snapshot *memento.Snapshot) string {
	classes := make([]string, 0, 2)
	for _, class := range strings.Fields(snapshot.Class) {
		if snapshotClassRegexp.MatchString(class) {
			classes = append(classes, class)
		}
	}
	if len(classes) == 0 {
		return ""
	}
	return " " + strings.Join(classes, " ")
}

func explicitBadge(explicit bool) string {
	if !explicit {
		return ""
	}
	return ` <span class="embed-spotify__explicit" title="Explicit">E</span>`
}

func formatDuration(milliseconds int64) string {
	if milliseconds <= 0 {
		return ""
	}
	totalSeconds := milliseconds / 1000
	return fmt.Sprintf("%d:%02d", totalSeconds/60, totalSeconds%60)
}
