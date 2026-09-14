package memento

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thecsw/darkness/v3/emilia/alpha"
)

const (
	spotifyEmbedBase = "https://open.spotify.com/embed"
	youtubeOEmbedURL = "https://www.youtube.com/oembed"
	maxMetadataBytes = 8 << 20
	maxArtworkBytes  = 16 << 20
)

type Importer struct {
	Client           *http.Client
	SpotifyEmbedBase string
	YouTubeOEmbedURL string
}

type FetchedSnapshot struct {
	Snapshot         *Snapshot
	Artwork          []byte
	ArtworkExtension string
}

func NewImporter() *Importer {
	return &Importer{
		Client:           &http.Client{Timeout: 30 * time.Second},
		SpotifyEmbedBase: spotifyEmbedBase,
		YouTubeOEmbedURL: youtubeOEmbedURL,
	}
}

// RefreshReference updates a specifically named, page-local TOML manifest.
func (importer *Importer) RefreshReference(ctx context.Context, conf *alpha.DarknessConfig,
	request Request, pageLocation, reference string, force, dryRun bool,
) (string, error) {
	manifestPath, err := ResolveReferencePath(conf, pageLocation, reference)
	if err != nil {
		return "", err
	}
	return importer.refreshPath(ctx, conf, request, manifestPath, force, dryRun)
}

func (importer *Importer) refreshPath(ctx context.Context, conf *alpha.DarknessConfig,
	request Request, manifestPath string, force, dryRun bool,
) (string, error) {
	if current, err := loadPath(conf, manifestPath); err == nil {
		if validateSnapshot(current, request) == nil {
			if current.Locked {
				return "locked", nil
			}
			if !force && !current.NeedsManual {
				return "cached", nil
			}
		}
	}

	fetched, err := importer.Fetch(ctx, request)
	if err != nil {
		return "", err
	}
	if dryRun {
		return "would refresh", nil
	}

	if err := WriteFetchedPath(manifestPath, fetched); err != nil {
		return "", err
	}
	return "refreshed", nil
}

// Fetch retrieves a provider snapshot and its artwork without writing files.
func (importer *Importer) Fetch(ctx context.Context, request Request) (*FetchedSnapshot, error) {
	switch request.Kind {
	case KindYouTube:
		return importer.fetchYouTube(ctx, request)
	case KindSpotifyTrack, KindSpotifyPlaylist:
		return importer.fetchSpotify(ctx, request)
	default:
		return nil, errors.New("unsupported embed kind")
	}
}

// WriteFetchedPath writes fetched metadata and gives its artwork a name based
// on the destination TOML file.
func WriteFetchedPath(manifestPath string, fetched *FetchedSnapshot) error {
	manifestBase := strings.TrimSuffix(filepath.Base(manifestPath), filepath.Ext(manifestPath))
	artworkBase := manifestBase
	if manifestBase == "embed" {
		artworkBase = "artwork"
	}
	fetched.Snapshot.ResolvedArtwork = ""
	if len(fetched.Artwork) > 0 {
		artworkName := artworkBase + fetched.ArtworkExtension
		artworkPath := filepath.Join(filepath.Dir(manifestPath), artworkName)
		// #nosec G301 - content directory must be world-readable for site serving
		if err := os.MkdirAll(filepath.Dir(artworkPath), 0o755); err != nil {
			return err
		}
		if err := writeAtomically(artworkPath, fetched.Artwork); err != nil {
			return fmt.Errorf("writing artwork: %w", err)
		}
		fetched.Snapshot.Artwork = artworkName
	}
	if err := savePath(manifestPath, fetched.Snapshot); err != nil {
		return fmt.Errorf("writing snapshot: %w", err)
	}
	return nil
}

func (importer *Importer) fetchYouTube(ctx context.Context, request Request) (*FetchedSnapshot, error) {
	endpoint := importer.YouTubeOEmbedURL + "?format=json&url=" +
		urlQueryEscape(request.URL)
	data, _, err := importer.get(ctx, endpoint, maxMetadataBytes)
	if err != nil {
		return nil, err
	}
	var response struct {
		Title        string `json:"title"`
		AuthorName   string `json:"author_name"`
		ThumbnailURL string `json:"thumbnail_url"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("decoding YouTube oEmbed response: %w", err)
	}
	if response.Title == "" || response.ThumbnailURL == "" {
		return nil, errors.New("YouTube oEmbed response is missing title or thumbnail")
	}
	artwork, contentType, err := importer.get(ctx, response.ThumbnailURL, maxArtworkBytes)
	if err != nil {
		return nil, fmt.Errorf("downloading YouTube thumbnail: %w", err)
	}
	return &FetchedSnapshot{
		Snapshot: &Snapshot{
			Version: SnapshotVersion, Provider: "YouTube", Kind: request.Kind,
			ID: request.ID, URL: request.URL, Title: response.Title, Creator: response.AuthorName,
		},
		Artwork: artwork, ArtworkExtension: imageExtension(contentType, response.ThumbnailURL),
	}, nil
}

func (importer *Importer) fetchSpotify(ctx context.Context, request Request) (*FetchedSnapshot, error) {
	kind := "track"
	if request.Kind == KindSpotifyPlaylist {
		kind = "playlist"
	}
	endpoint := strings.TrimRight(importer.SpotifyEmbedBase, "/") + "/" + kind + "/" + request.ID
	data, _, err := importer.get(ctx, endpoint, maxMetadataBytes)
	if err != nil {
		return nil, err
	}
	entity, err := decodeSpotifyEntity(data)
	if err != nil {
		return nil, err
	}

	snapshot := &Snapshot{
		Version: SnapshotVersion, Provider: "Spotify", Kind: request.Kind,
		ID: request.ID, URL: request.URL, Title: entity.Title,
		Creator: entity.Subtitle, DurationMS: entity.Duration, Explicit: entity.IsExplicit,
		Colors: spotifyColors(entity.VisualIdentity),
	}
	if snapshot.Title == "" {
		snapshot.Title = entity.Name
	}
	if len(entity.Artists) > 0 {
		names := make([]string, 0, len(entity.Artists))
		for _, artist := range entity.Artists {
			names = append(names, artist.Name)
		}
		snapshot.Creator = strings.Join(names, ", ")
	}
	for _, item := range entity.TrackList {
		id := strings.TrimPrefix(item.URI, "spotify:track:")
		if !validID(id) {
			continue
		}
		snapshot.Tracks = append(snapshot.Tracks, Track{
			ID: id, URL: "https://open.spotify.com/track/" + id,
			Title: item.Title, Artist: item.Subtitle,
			DurationMS: item.Duration, Explicit: item.IsExplicit,
		})
	}

	artworkURL := bestSpotifyImage(entity)
	if artworkURL == "" {
		return nil, errors.New("spotify embed response is missing artwork")
	}
	artwork, contentType, err := importer.get(ctx, artworkURL, maxArtworkBytes)
	if err != nil {
		return nil, fmt.Errorf("downloading Spotify artwork: %w", err)
	}
	return &FetchedSnapshot{
		Snapshot: snapshot, Artwork: artwork,
		ArtworkExtension: imageExtension(contentType, artworkURL),
	}, nil
}

func (importer *Importer) get(ctx context.Context, endpoint string, limit int64) ([]byte, string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", err
	}
	request.Header.Set("User-Agent", "Darkness embed snapshot importer")
	response, err := importer.Client.Do(request)
	if err != nil {
		return nil, "", err
	}
	if response == nil {
		return nil, "", fmt.Errorf("GET %s: no response", endpoint)
	}
	defer response.Body.Close() //nolint:errcheck
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, "", fmt.Errorf("GET %s: %s", endpoint, response.Status)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, limit+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(data)) > limit {
		return nil, "", fmt.Errorf("GET %s: response exceeds %d bytes", endpoint, limit)
	}
	return data, response.Header.Get("Content-Type"), nil
}

type spotifyImage struct {
	URL       string `json:"url"`
	MaxWidth  int    `json:"maxWidth"`
	MaxHeight int    `json:"maxHeight"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type spotifyColor struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

type spotifyVisualIdentity struct {
	BackgroundBase       spotifyColor   `json:"backgroundBase"`
	BackgroundTintedBase spotifyColor   `json:"backgroundTintedBase"`
	TextBase             spotifyColor   `json:"textBase"`
	TextSubdued          spotifyColor   `json:"textSubdued"`
	Image                []spotifyImage `json:"image"`
}

type spotifyEntity struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	ID         string `json:"id"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Duration   int64  `json:"duration"`
	IsExplicit bool   `json:"isExplicit"`
	Artists    []struct {
		Name string `json:"name"`
	} `json:"artists"`
	CoverArt struct {
		Sources []spotifyImage `json:"sources"`
	} `json:"coverArt"`
	TrackList []struct {
		URI        string `json:"uri"`
		Title      string `json:"title"`
		Subtitle   string `json:"subtitle"`
		Duration   int64  `json:"duration"`
		IsExplicit bool   `json:"isExplicit"`
	} `json:"trackList"`
	VisualIdentity spotifyVisualIdentity `json:"visualIdentity"`
}

func decodeSpotifyEntity(document []byte) (*spotifyEntity, error) {
	const marker = `<script id="__NEXT_DATA__"`
	tagStart := bytes.Index(document, []byte(marker))
	if tagStart < 0 {
		return nil, errors.New("spotify embed response has no __NEXT_DATA__ snapshot")
	}
	tagEnd := bytes.IndexByte(document[tagStart:], '>')
	if tagEnd < 0 {
		return nil, errors.New("spotify __NEXT_DATA__ opening tag is not terminated")
	}
	start := tagStart + tagEnd + 1
	end := bytes.Index(document[start:], []byte("</script>"))
	if end < 0 {
		return nil, errors.New("spotify __NEXT_DATA__ snapshot is not terminated")
	}
	var payload struct {
		Props struct {
			PageProps struct {
				State struct {
					Data struct {
						Entity spotifyEntity `json:"entity"`
					} `json:"data"`
				} `json:"state"`
			} `json:"pageProps"`
		} `json:"props"`
	}
	if err := json.Unmarshal(document[start:start+end], &payload); err != nil {
		return nil, fmt.Errorf("decoding Spotify embed snapshot: %w", err)
	}
	entity := payload.Props.PageProps.State.Data.Entity
	if entity.ID == "" || entity.Title == "" {
		return nil, errors.New("spotify embed snapshot is missing identity metadata")
	}
	return &entity, nil
}

func bestSpotifyImage(entity *spotifyEntity) string {
	images := append([]spotifyImage(nil), entity.VisualIdentity.Image...)
	images = append(images, entity.CoverArt.Sources...)
	bestURL, bestArea := "", -1
	for _, image := range images {
		width, height := image.MaxWidth, image.MaxHeight
		if width == 0 {
			width = image.Width
		}
		if height == 0 {
			height = image.Height
		}
		area := width * height
		if image.URL != "" && area > bestArea {
			bestURL, bestArea = image.URL, area
		}
	}
	return bestURL
}

func spotifyColors(identity spotifyVisualIdentity) Colors {
	return Colors{
		Background:       colorHex(identity.BackgroundBase),
		BackgroundTinted: colorHex(identity.BackgroundTintedBase),
		Text:             colorHex(identity.TextBase),
		TextSubdued:      colorHex(identity.TextSubdued),
	}
}

func colorHex(color spotifyColor) string {
	if color.Red < 0 || color.Red > 255 || color.Green < 0 || color.Green > 255 || color.Blue < 0 || color.Blue > 255 {
		return ""
	}
	return fmt.Sprintf("#%02x%02x%02x", color.Red, color.Green, color.Blue)
}

func imageExtension(contentType, sourceURL string) string {
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	switch mediaType {
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "image/jpeg", "image/jpg":
		return ".jpg"
	}
	extension := strings.ToLower(filepath.Ext(strings.Split(sourceURL, "?")[0]))
	switch extension {
	case ".png", ".webp", ".gif", ".jpg", ".jpeg":
		if extension == ".jpeg" {
			return ".jpg"
		}
		return extension
	default:
		return ".jpg"
	}
}

func writeAtomically(path string, data []byte) error {
	temporary := path + ".tmp"
	// #nosec G306 - content file must be world-readable for site serving
	if err := os.WriteFile(temporary, data, 0o644); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}

func urlQueryEscape(value string) string {
	return url.QueryEscape(value)
}
