package memento

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/thecsw/darkness/v3/emilia/alpha"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		url  string
		kind Kind
		id   string
		ok   bool
	}{
		{"https://youtu.be/SGZc72g_aPI", KindYouTube, "SGZc72g_aPI", true},
		{"https://www.youtube.com/watch?v=SGZc72g_aPI", KindYouTube, "SGZc72g_aPI", true},
		{"https://open.spotify.com/track/54BW4qpq5ms4bnzBgiWVOo?si=hello", KindSpotifyTrack, "54BW4qpq5ms4bnzBgiWVOo", true},
		{"https://open.spotify.com/playlist/59bal0ZSNOlRC6jwhu0ocw", KindSpotifyPlaylist, "59bal0ZSNOlRC6jwhu0ocw", true},
		{"http://youtu.be/SGZc72g_aPI", "", "", false},
		{"ftp://open.spotify.com/track/54BW4qpq5ms4bnzBgiWVOo", "", "", false},
		{"//youtu.be/SGZc72g_aPI", "", "", false},
		{"javascript:alert(1)", "", "", false},
		{"https://example.com/video", "", "", false},
	}
	for _, test := range tests {
		request, ok := ParseURL(test.url)
		if ok != test.ok || request.Kind != test.kind || request.ID != test.id {
			t.Errorf("ParseURL(%q) = %#v, %v; want kind=%q id=%q ok=%v", test.url, request, ok, test.kind, test.id, test.ok)
		}
	}
}

func TestDecodeSpotifyPlaylistEntity(t *testing.T) {
	document := []byte(`<html><script id="__NEXT_DATA__" type="application/json">{
  "props":{"pageProps":{"state":{"data":{"entity":{
    "type":"playlist","id":"playlist-id","title":"night drive","subtitle":"sandy",
    "trackList":[
      {"uri":"spotify:track:first","title":"First","subtitle":"Artist One","duration":61000},
      {"uri":"spotify:track:second","title":"Second","subtitle":"Artist Two","duration":122000,"isExplicit":true}
    ],
    "visualIdentity":{"backgroundBase":{"red":12,"green":34,"blue":56},"image":[{"url":"small.jpg","maxWidth":64,"maxHeight":64},{"url":"large.jpg","maxWidth":640,"maxHeight":640}]}
  }}}}}}</script></html>`)
	entity, err := decodeSpotifyEntity(document)
	if err != nil {
		t.Fatal(err)
	}
	if entity.Title != "night drive" || entity.Subtitle != "sandy" || len(entity.TrackList) != 2 {
		t.Fatalf("unexpected entity: %#v", entity)
	}
	if got := bestSpotifyImage(entity); got != "large.jpg" {
		t.Errorf("best image = %q, want large.jpg", got)
	}
	if got := spotifyColors(entity.VisualIdentity).Background; got != "#0c2238" {
		t.Errorf("background = %q, want #0c2238", got)
	}
}

func TestRefreshReferenceHydratesAColocatedTOMLStub(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/oembed":
			writer.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(writer, `{"title":"NIGHT RUNNING","author_name":"sandy","thumbnail_url":%q}`, server.URL+"/thumbnail.jpg")
		case "/thumbnail.jpg":
			writer.Header().Set("Content-Type", "image/jpeg")
			_, _ = writer.Write([]byte("local artwork"))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	conf := &alpha.DarknessConfig{}
	conf.Runtime.WorkDir = alpha.WorkingDirectory(t.TempDir())
	pageDirectory := conf.Runtime.WorkDir.JoinGeneric("blogs/music")
	if err := os.MkdirAll(pageDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	manifestPath := pageDirectory + "/running.toml"
	if err := os.WriteFile(manifestPath, []byte("url = \"https://youtu.be/SGZc72g_aPI\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	importer := NewImporter()
	importer.Client = server.Client()
	importer.YouTubeOEmbedURL = server.URL + "/oembed"
	request, _ := ParseURL("https://youtu.be/SGZc72g_aPI")
	status, err := importer.RefreshReference(context.Background(), conf, request,
		"blogs/music", "running.toml", false, false)
	if err != nil || status != "refreshed" {
		t.Fatalf("reference refresh = %q, %v", status, err)
	}
	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`# locked = true       # never rewrite this file, even with -force`,
		`# needs_manual = true # retry incomplete metadata on the next normal vendor run`,
		`kind = "youtube"`, `title = "NIGHT RUNNING"`, `artwork = "running.jpg"`,
	} {
		if !strings.Contains(string(manifest), want) {
			t.Errorf("hydrated TOML does not contain %q:\n%s", want, manifest)
		}
	}
	loaded, err := LoadReference(conf, "blogs/music", "running.toml")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ResolvedArtwork != "blogs/music/running.jpg" {
		t.Errorf("resolved artwork = %q", loaded.ResolvedArtwork)
	}
}

func TestResolveReferenceAndArtworkStayInPageDirectory(t *testing.T) {
	conf := &alpha.DarknessConfig{}
	conf.Runtime.WorkDir = alpha.WorkingDirectory(t.TempDir())
	for _, reference := range []string{"../shared.toml", "/shared.toml"} {
		if _, err := ResolveReferencePath(conf, "blogs/music", reference); err == nil {
			t.Errorf("ResolveReferencePath accepted %q outside the page", reference)
		}
	}
	if _, err := ResolveReferencePath(conf, "blogs/music", "embeds/song.toml"); err != nil {
		t.Fatalf("ResolveReferencePath rejected a page-local subdirectory: %v", err)
	}

	pageDirectory := conf.Runtime.WorkDir.JoinGeneric("blogs/music")
	if err := os.MkdirAll(pageDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `version = 1
provider = "YouTube"
kind = "youtube"
id = "SGZc72g_aPI"
url = "https://youtu.be/SGZc72g_aPI"
title = "Example"
artwork = "../outside.jpg"
`
	if err := os.WriteFile(pageDirectory+"/example.toml", []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadReference(conf, "blogs/music", "example.toml")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ResolvedArtwork != "" {
		t.Errorf("artwork escaped its manifest directory: %q", loaded.ResolvedArtwork)
	}
}

func TestRefreshPrefillsAndPreservesManualSnapshots(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/oembed":
			writer.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(writer, `{"title":"Remote title","author_name":"Remote author","thumbnail_url":%q}`, server.URL+"/thumbnail.jpg")
		case "/thumbnail.jpg":
			writer.Header().Set("Content-Type", "image/jpeg")
			_, _ = writer.Write([]byte("local artwork"))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	conf := &alpha.DarknessConfig{}
	conf.Runtime.WorkDir = alpha.WorkingDirectory(t.TempDir())
	pageLocation := "blogs/music"
	manifestReference := "running.toml"
	if err := os.MkdirAll(conf.Runtime.WorkDir.JoinGeneric(pageLocation), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := SaveReference(conf, pageLocation, manifestReference, &Snapshot{
		URL: "https://youtu.be/SGZc72g_aPI",
	}); err != nil {
		t.Fatal(err)
	}

	importer := NewImporter()
	importer.Client = server.Client()
	importer.YouTubeOEmbedURL = server.URL + "/oembed"
	request, _ := ParseURL("https://youtu.be/SGZc72g_aPI")

	status, err := importer.RefreshReference(context.Background(), conf, request,
		pageLocation, manifestReference, false, false)
	if err != nil || status != "refreshed" {
		t.Fatalf("initial refresh = %q, %v", status, err)
	}
	snapshot, err := LoadReference(conf, pageLocation, manifestReference)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Title != "Remote title" || snapshot.Artwork != "running.jpg" ||
		!strings.HasSuffix(snapshot.ResolvedArtwork, "/running.jpg") {
		t.Fatalf("unexpected imported snapshot: %#v", snapshot)
	}
	if _, err := os.Stat(conf.Runtime.WorkDir.JoinGeneric(snapshot.ResolvedArtwork)); err != nil {
		t.Fatalf("local artwork was not written: %v", err)
	}

	snapshot.Title = "Temporary manual fallback"
	snapshot.NeedsManual = true
	if err := SaveReference(conf, pageLocation, manifestReference, snapshot); err != nil {
		t.Fatal(err)
	}
	status, err = importer.RefreshReference(context.Background(), conf, request,
		pageLocation, manifestReference, false, false)
	if err != nil || status != "refreshed" {
		t.Fatalf("manual placeholder refresh = %q, %v", status, err)
	}
	snapshot, _ = LoadReference(conf, pageLocation, manifestReference)
	if snapshot.NeedsManual || snapshot.Title != "Remote title" {
		t.Fatalf("manual placeholder was not hydrated: %#v", snapshot)
	}

	snapshot.Title = "My better title"
	if err := SaveReference(conf, pageLocation, manifestReference, snapshot); err != nil {
		t.Fatal(err)
	}
	status, err = importer.RefreshReference(context.Background(), conf, request,
		pageLocation, manifestReference, false, false)
	if err != nil || status != "cached" {
		t.Fatalf("ordinary refresh = %q, %v", status, err)
	}
	preserved, _ := LoadReference(conf, pageLocation, manifestReference)
	if preserved.Title != "My better title" {
		t.Errorf("ordinary refresh overwrote manual changes: %#v", preserved)
	}

	preserved.Locked = true
	if err := SaveReference(conf, pageLocation, manifestReference, preserved); err != nil {
		t.Fatal(err)
	}
	status, err = importer.RefreshReference(context.Background(), conf, request,
		pageLocation, manifestReference, true, false)
	if err != nil || status != "locked" {
		t.Fatalf("forced locked refresh = %q, %v", status, err)
	}
	locked, _ := LoadReference(conf, pageLocation, manifestReference)
	if locked.Title != "My better title" {
		t.Errorf("forced refresh overwrote locked snapshot: %#v", locked)
	}
}
