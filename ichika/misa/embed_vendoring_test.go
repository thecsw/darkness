package misa

import (
	"os"
	"strings"
	"testing"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/memento"
	"github.com/thecsw/darkness/v3/yunyun"
)

func TestRewriteEmbedLinksVendorsStandaloneLinksAndHonorsSkipFlag(t *testing.T) {
	yunyun.ActiveMarkings.BuildRegex()
	source := `A paragraph with [[https://youtu.be/inline][an inline video]].

#+begin_src org
[[https://youtu.be/literal][a literal example]]
#+end_src

Introductory prose
[[https://youtu.be/paragraph][part of that paragraph]]

[[https://youtu.be/SGZc72g_aPI][A video]]

#+attr_darkness: no-embed-vendor
#+attr_darkness: decorative
[[https://open.spotify.com/track/skipped][Skip me]]

#+html_tags: style="width: 90%"
[[https://open.spotify.com/track/vendor][Vendor me]]
`
	calls := 0
	rewritten, stats := rewriteEmbedLinks(source, "blogs/music",
		func(_ string, _ *yunyun.ExtractedLink, request memento.Request) (string, error) {
			calls++
			return macroForEmbed(request, "snapshot.toml"), nil
		})

	if calls != 2 || stats.found != 3 || stats.vendored != 2 || stats.skipped != 1 || stats.failed != 0 {
		t.Fatalf("unexpected rewrite stats: calls=%d stats=%+v", calls, stats)
	}
	for _, want := range []string{
		`A paragraph with [[https://youtu.be/inline][an inline video]].`,
		`[[https://youtu.be/literal][a literal example]]`,
		`[[https://youtu.be/paragraph][part of that paragraph]]`,
		`{{{youtube_embed(snapshot.toml)}}}`,
		`[[https://open.spotify.com/track/skipped][Skip me]]`,
		`{{{spotify_embed(snapshot.toml)}}}`,
	} {
		if !strings.Contains(rewritten, want) {
			t.Errorf("rewritten source does not contain %q:\n%s", want, rewritten)
		}
	}
}

func TestEmbedManifestSlugUsesMetadataAndStableFallbacks(t *testing.T) {
	tests := []struct {
		title string
		kind  memento.Kind
		id    string
		want  string
	}{
		{"NIGHT RUNNING", memento.KindSpotifyTrack, "abc123", "night-running"},
		{"Spotify track", memento.KindSpotifyTrack, "abc123", "spotify-track-abc123"},
		{"", memento.KindYouTube, "SGZc72g_aPI", "youtube-SGZc72g_aPI"},
		{"空色デイズ", memento.KindSpotifyTrack, "japan", "空色デイズ"},
	}
	for _, test := range tests {
		got := embedManifestSlug(test.title, memento.Request{Kind: test.kind, ID: test.id})
		if got != test.want {
			t.Errorf("embedManifestSlug(%q) = %q, want %q", test.title, got, test.want)
		}
	}
}

func TestAvailableEmbedReferenceMarksMatchingManifestAsExisting(t *testing.T) {
	conf := &alpha.DarknessConfig{}
	conf.Runtime.WorkDir = alpha.WorkingDirectory(t.TempDir())
	if err := os.MkdirAll(conf.Runtime.WorkDir.JoinGeneric("blogs/music"), 0o755); err != nil {
		t.Fatal(err)
	}
	request, _ := memento.ParseURL("https://youtu.be/SGZc72g_aPI")
	if err := memento.SaveReference(conf, "blogs/music", "favorite.toml", &memento.Snapshot{
		Version: memento.SnapshotVersion, Locked: true, Provider: "YouTube",
		Kind: request.Kind, ID: request.ID, URL: request.URL, Title: "My title",
	}); err != nil {
		t.Fatal(err)
	}

	reference, exists, err := availableEmbedReference(conf, "blogs/music", "favorite", request)
	if err != nil || reference != "favorite.toml" || !exists {
		t.Fatalf("matching reference = %q, exists=%v, err=%v", reference, exists, err)
	}
}

func TestManualEmbedTemplateIsImmediatelyRenderable(t *testing.T) {
	request := memento.Request{
		Kind: memento.KindYouTube, ID: "missing", URL: "https://youtu.be/missing",
	}
	snapshot := manualEmbedTemplate(request, "An unavailable favorite")
	if !snapshot.NeedsManual || snapshot.Version != memento.SnapshotVersion ||
		snapshot.Kind != request.Kind || snapshot.ID != request.ID ||
		snapshot.URL != request.URL || snapshot.Title != "An unavailable favorite" {
		t.Fatalf("unexpected manual template: %#v", snapshot)
	}
}
