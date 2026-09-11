package html

import (
	"os"
	"strings"
	"testing"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/memento"
	"github.com/thecsw/darkness/v3/yunyun"
)

func testLinkState(mode string) *state {
	yunyun.ActiveMarkings.BuildRegex()
	return &state{conf: &alpha.DarknessConfig{
		Website: alpha.WebsiteConfig{EmbedMode: mode},
	}}
}

func TestStaticEmbedCardsRenderSupportedProviders(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		class    string
		kind     string
		fallback string
	}{
		{
			name: "youtube", link: "https://youtu.be/SGZc72g_aPI",
			class: "video", kind: "moving picture · youtube", fallback: "YouTube video",
		},
		{
			name: "spotify track", link: "https://open.spotify.com/track/54BW4qpq5ms4bnzBgiWVOo",
			class: "track", kind: "record · spotify", fallback: "Spotify track",
		},
		{
			name: "spotify playlist", link: "https://open.spotify.com/playlist/59bal0ZSNOlRC6jwhu0ocw",
			class: "playlist", kind: "mixtape · spotify", fallback: "Spotify playlist",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content := &yunyun.Content{Link: test.link, Caption: "an authored note"}
			got := testLinkState("card").link(content)

			for _, want := range []string{
				`class="media embed-card embed-card--` + test.class + `"`,
				`href="` + test.link + `"`,
				test.kind,
				test.fallback,
				`target="_blank" rel="external noopener noreferrer"`,
				`class="embed-card__description">an authored note</span>`,
			} {
				if !strings.Contains(got, want) {
					t.Errorf("card does not contain %q:\n%s", want, got)
				}
			}
			if strings.Contains(got, "<iframe") || strings.Contains(got, "<script") {
				t.Errorf("static card contains remote executable markup:\n%s", got)
			}
		})
	}
}

func TestRemoteEmbedsRemainTheDefault(t *testing.T) {
	content := &yunyun.Content{
		Link:      "https://youtu.be/SGZc72g_aPI",
		LinkTitle: "The Day Before You Came",
	}
	got := testLinkState("").link(content)

	if !strings.Contains(got, `<iframe src="https://www.youtube.com/embed/SGZc72g_aPI"`) {
		t.Fatalf("default embed did not retain the remote iframe:\n%s", got)
	}
}

func TestEmbedRemoteAttributeOverridesCardMode(t *testing.T) {
	content := &yunyun.Content{
		Link:       "https://open.spotify.com/track/54BW4qpq5ms4bnzBgiWVOo",
		Attributes: "embed-remote",
	}
	got := testLinkState("card").link(content)

	if !strings.Contains(got, "<iframe") || strings.Contains(got, "embed-card") {
		t.Fatalf("embed-remote did not restore the provider player:\n%s", got)
	}
}

func TestExplicitGenericEmbedCardIsStaticAndEscaped(t *testing.T) {
	content := &yunyun.Content{
		Link:       `https://example.com/read?a=1&b=2`,
		LinkTitle:  `<A quiet & durable web>`,
		Caption:    `<script>alert("no")</script>`,
		Attributes: "embed-card",
	}
	got := testLinkState("").link(content)

	for _, want := range []string{
		`embed-card--link`,
		`href="https://example.com/read?a=1&amp;b=2"`,
		`example.com`,
		`&lt;A quiet &amp; durable web&gt;`,
		`&lt;script&gt;alert(&#34;no&#34;)&lt;/script&gt;`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("generic card does not contain %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "<script") || strings.Contains(got, "<iframe") {
		t.Errorf("generic card contains executable markup:\n%s", got)
	}
}

func TestStaticEmbedCardCaptionCannotCreateNestedLinks(t *testing.T) {
	content := &yunyun.Content{
		Link:    "https://youtu.be/SGZc72g_aPI",
		Caption: `Read [[https://example.com][the note]]`,
	}
	got := testLinkState("card").link(content)
	if strings.Count(got, "<a ") != 1 || !strings.Contains(got, "Read the note") {
		t.Fatalf("card caption created nested interactive markup:\n%s", got)
	}
}

func TestEmbedCardAttributeRequiresAnExactTokenAndWebURL(t *testing.T) {
	tests := []*yunyun.Content{
		{Link: "https://example.com", LinkTitle: "Example", Attributes: "not-an-embed-card"},
		{Link: "mailto:hello@example.com", LinkTitle: "Email", Attributes: "embed-card"},
	}

	for _, content := range tests {
		got := testLinkState("").link(content)
		if strings.Contains(got, "embed-card__link") {
			t.Errorf("unexpected static card for %#v:\n%s", content, got)
		}
	}
}

func TestEmbedReferenceRendersAColocatedTOMLManifest(t *testing.T) {
	state := testLinkState("card")
	state.conf.Runtime.WorkDir = alpha.WorkingDirectory(t.TempDir())
	state.page = &yunyun.Page{Location: "blogs/music"}
	pageDirectory := state.conf.Runtime.WorkDir.JoinGeneric("blogs/music")
	if err := os.MkdirAll(pageDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `version = 1
provider = "Spotify"
kind = "spotify-track"
id = "54BW4qpq5ms4bnzBgiWVOo"
url = "https://open.spotify.com/track/54BW4qpq5ms4bnzBgiWVOo"
title = "My own NIGHT RUNNING"
creator = "Shin Sakiura, AAAMYYY"
class = "night-running feature"
artwork = "running.jpg"
duration_ms = 233666
`
	if err := os.WriteFile(pageDirectory+"/running.toml", []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	got := state.link(&yunyun.Content{Link: "embed:running.toml"})
	for _, want := range []string{
		"embed-spotify--track night-running feature", "My own NIGHT RUNNING",
		`src="/blogs/music/running.jpg"`, `href="https://open.spotify.com/track/54BW4qpq5ms4bnzBgiWVOo"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("referenced embed does not contain %q:\n%s", want, got)
		}
	}
}

func TestRichSpotifyPlaylistContainsEveryRealTrack(t *testing.T) {
	snapshot := &memento.Snapshot{
		Version: memento.SnapshotVersion, Provider: "Spotify", Kind: memento.KindSpotifyPlaylist,
		ID: "playlist", URL: "https://open.spotify.com/playlist/playlist",
		Title: "Night drive", Creator: "sandy", Artwork: "embeds/spotify-playlist/playlist/artwork.jpg",
		Colors: memento.Colors{Background: "#123456", BackgroundTinted: "#081522", Text: "#ffffff", TextSubdued: "#cccccc"},
		Tracks: []memento.Track{
			{ID: "first", URL: "https://open.spotify.com/track/first", Title: "First song", Artist: "One", DurationMS: 61000},
			{ID: "second", URL: "https://open.spotify.com/track/second", Title: "Second song", Artist: "Two", DurationMS: 122000, Explicit: true},
		},
	}
	got := renderSpotifyPlaylistSnapshot(snapshot, "")
	for _, want := range []string{
		"Night drive", "sandy · 2 songs", "First song", "Second song", "1:01", "2:02",
		`href="https://open.spotify.com/track/first"`, `href="https://open.spotify.com/track/second"`,
		`style="--embed-bg:#123456;--embed-bg-deep:#081522;--embed-fg:#ffffff;--embed-muted:#cccccc;"`,
		`target="_blank" rel="external noopener noreferrer"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("rich Spotify playlist does not contain %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "<iframe") || strings.Contains(got, "spotifycdn") {
		t.Errorf("rich Spotify playlist contains a remote resource:\n%s", got)
	}
}

func TestStaticCardsDoNotEmitProviderPreconnects(t *testing.T) {
	content := &yunyun.Content{
		Type: yunyun.TypeLink,
		Link: "https://youtu.be/SGZc72g_aPI",
	}
	cardState := testLinkState("card")
	cardState.page = &yunyun.Page{Contents: []*yunyun.Content{content}}
	if got := strings.Join(cardState.resourceHints(), "\n"); strings.Contains(got, "youtube.com") {
		t.Errorf("static card emitted a provider preconnect: %s", got)
	}

	content.Attributes = "embed-remote"
	if got := strings.Join(cardState.resourceHints(), "\n"); !strings.Contains(got, "youtube.com") {
		t.Errorf("remote override did not emit a provider preconnect: %s", got)
	}
}
