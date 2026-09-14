// Package memento stores durable, first-party snapshots of remote embeds.
package memento

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/thecsw/darkness/v3/emilia/alpha"
)

const (
	SnapshotVersion = 1

	manifestControlHeader = `# Misa control flags (uncomment only when needed):
# locked = true       # never rewrite this file, even with -force
# needs_manual = true # retry incomplete metadata on the next normal vendor run
# If both are true, locked wins.

`
)

type Kind string

const (
	KindYouTube         Kind = "youtube"
	KindSpotifyTrack    Kind = "spotify-track"
	KindSpotifyPlaylist Kind = "spotify-playlist"
)

type Request struct {
	Kind Kind
	ID   string
	URL  string
}

type Colors struct {
	Background       string `toml:"background,omitempty"`
	BackgroundTinted string `toml:"background_tinted,omitempty"`
	Text             string `toml:"text,omitempty"`
	TextSubdued      string `toml:"text_subdued,omitempty"`
}

type Track struct {
	ID         string `toml:"id"`
	URL        string `toml:"url"`
	Title      string `toml:"title"`
	Artist     string `toml:"artist"`
	DurationMS int64  `toml:"duration_ms,omitzero"`
	Explicit   bool   `toml:"explicit,omitempty"`
}

type Snapshot struct {
	Version     int     `toml:"version"`
	Locked      bool    `toml:"locked,omitempty"`
	NeedsManual bool    `toml:"needs_manual,omitempty"`
	Provider    string  `toml:"provider"`
	Kind        Kind    `toml:"kind"`
	ID          string  `toml:"id"`
	URL         string  `toml:"url"`
	Title       string  `toml:"title"`
	Creator     string  `toml:"creator,omitempty"`
	Class       string  `toml:"class,omitempty"`
	Artwork     string  `toml:"artwork,omitempty"`
	DurationMS  int64   `toml:"duration_ms,omitzero"`
	Explicit    bool    `toml:"explicit,omitempty"`
	Colors      Colors  `toml:"colors,omitempty"`
	Tracks      []Track `toml:"tracks,omitempty"`

	// ResolvedArtwork is the project-root-relative path used by the renderer.
	// Artwork itself remains relative to the TOML file and pleasant to edit.
	ResolvedArtwork string `toml:"-"`
}

func ParseURL(raw string) (Request, bool) {
	clean := strings.TrimSpace(raw)
	parsed, err := url.Parse(clean)
	if err != nil || parsed.Scheme != "https" {
		return Request{}, false
	}

	host := strings.ToLower(parsed.Hostname())
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	request := Request{URL: clean}
	switch {
	case host == "youtu.be" && len(parts) == 1:
		request.Kind, request.ID = KindYouTube, parts[0]
	case (host == "youtube.com" || host == "www.youtube.com") && parsed.Path == "/watch":
		request.Kind, request.ID = KindYouTube, parsed.Query().Get("v")
	case host == "open.spotify.com" && len(parts) >= 2 && parts[0] == "track":
		request.Kind, request.ID = KindSpotifyTrack, parts[1]
	case host == "open.spotify.com" && len(parts) >= 2 && parts[0] == "playlist":
		request.Kind, request.ID = KindSpotifyPlaylist, parts[1]
	default:
		return Request{}, false
	}
	if !validID(request.ID) {
		return Request{}, false
	}
	return request, true
}

func validID(id string) bool {
	if id == "" {
		return false
	}
	for _, char := range id {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '-' || char == '_' {
			continue
		}
		return false
	}
	return true
}

// ResolveReferencePath resolves a TOML manifest relative to the current Org
// page and guarantees that the result remains inside that page directory.
func ResolveReferencePath(conf *alpha.DarknessConfig, pageLocation, reference string) (string, error) {
	cleanReference := strings.TrimSpace(reference)
	if cleanReference == "" || strings.Contains(cleanReference, "://") || filepath.Ext(cleanReference) != ".toml" {
		return "", errors.New("embed reference must be a local TOML file")
	}
	if filepath.IsAbs(cleanReference) {
		return "", errors.New("embed reference must be relative to its page")
	}
	pageDirectory := filepath.Clean(conf.Runtime.WorkDir.JoinGeneric(pageLocation))
	manifestPath := filepath.Clean(filepath.Join(pageDirectory, cleanReference))
	if !insideDirectory(pageDirectory, manifestPath) {
		return "", errors.New("embed reference leaves its page directory")
	}
	return manifestPath, nil
}

// ReadReference reads a page-local TOML manifest without requiring it to be
// complete yet. The importer uses this to hydrate a stub containing only url.
func ReadReference(conf *alpha.DarknessConfig, pageLocation, reference string) (*Snapshot, error) {
	manifestPath, err := ResolveReferencePath(conf, pageLocation, reference)
	if err != nil {
		return nil, err
	}
	return loadPath(conf, manifestPath)
}

// LoadReference resolves and validates a TOML manifest relative to the current
// Org page for rendering.
func LoadReference(conf *alpha.DarknessConfig, pageLocation, reference string) (*Snapshot, error) {
	snapshot, err := ReadReference(conf, pageLocation, reference)
	if err != nil {
		return nil, err
	}
	request, ok := ParseURL(snapshot.URL)
	if !ok {
		return nil, errors.New("embed manifest has an unsupported provider URL")
	}
	if err := validateSnapshot(snapshot, request); err != nil {
		return nil, err
	}
	return snapshot, nil
}

func loadPath(conf *alpha.DarknessConfig, manifestPath string) (*Snapshot, error) {
	var snapshot Snapshot
	if _, err := toml.DecodeFile(manifestPath, &snapshot); err != nil {
		return nil, fmt.Errorf("decoding embed snapshot: %w", err)
	}
	if snapshot.Artwork != "" {
		manifestDirectory := filepath.Dir(manifestPath)
		artworkPath := filepath.Clean(filepath.Join(manifestDirectory, snapshot.Artwork))
		if insideDirectory(manifestDirectory, artworkPath) {
			relative, err := filepath.Rel(string(conf.Runtime.WorkDir), artworkPath)
			if err == nil {
				snapshot.ResolvedArtwork = filepath.ToSlash(relative)
			}
		}
	}
	return &snapshot, nil
}

func validateSnapshot(snapshot *Snapshot, request Request) error {
	snapshotRequest, snapshotURLValid := ParseURL(snapshot.URL)
	if snapshot.Version != SnapshotVersion || snapshot.ID != request.ID || snapshot.Kind != request.Kind ||
		snapshot.Provider != providerForKind(request.Kind) || !snapshotURLValid ||
		snapshotRequest.ID != request.ID || snapshotRequest.Kind != request.Kind {
		return errors.New("embed snapshot does not match its source URL")
	}
	return nil
}

func providerForKind(kind Kind) string {
	if kind == KindYouTube {
		return "YouTube"
	}
	if kind == KindSpotifyTrack || kind == KindSpotifyPlaylist {
		return "Spotify"
	}
	return ""
}

func insideDirectory(directory, target string) bool {
	relative, err := filepath.Rel(directory, target)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

// SaveReference writes a manually supplied snapshot beside an Org page.
func SaveReference(conf *alpha.DarknessConfig, pageLocation, reference string, snapshot *Snapshot) error {
	if snapshot == nil {
		return errors.New("snapshot is nil")
	}
	manifestPath, err := ResolveReferencePath(conf, pageLocation, reference)
	if err != nil {
		return err
	}
	snapshot.ResolvedArtwork = ""
	return savePath(manifestPath, snapshot)
}

func savePath(manifestPath string, snapshot *Snapshot) error {
	// #nosec G301 - content directory must be world-readable for site serving
	if err := os.MkdirAll(filepath.Dir(manifestPath), 0o755); err != nil {
		return err
	}
	var data bytes.Buffer
	data.WriteString(manifestControlHeader)
	if err := toml.NewEncoder(&data).Encode(snapshot); err != nil {
		return err
	}
	temporary := manifestPath + ".tmp"
	// #nosec G306 - content file must be world-readable for site serving
	if err := os.WriteFile(temporary, data.Bytes(), 0o644); err != nil {
		return err
	}
	return os.Rename(temporary, manifestPath)
}
