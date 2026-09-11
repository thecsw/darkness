package misa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/memento"
	"github.com/thecsw/darkness/v3/ichika/hizuru"
	"github.com/thecsw/darkness/v3/yunyun"
)

const skipEmbedVendoringAttribute = "no-embed-vendor"

type vendorLinkFunc func(pageLocation string, link *yunyun.ExtractedLink, request memento.Request) (string, error)

type embedVendorStats struct {
	found    int
	vendored int
	skipped  int
	manual   int
	failed   int
}

// VendorEmbeds refreshes existing manifest references before turning new
// standalone provider links into colocated TOML-backed macro calls.
func VendorEmbeds(conf *alpha.DarknessConfig, dryRun bool) {
	refreshVendoredEmbeds(conf, dryRun)
	vendorEmbedLinks(conf, dryRun)
}

// vendorEmbedLinks performs the source-rewriting half of the unified workflow.
// Retrieval failures still produce an editable manual template.
func vendorEmbedLinks(conf *alpha.DarknessConfig, dryRun bool) {
	initLog()
	yunyun.ActiveMarkings.BuildRegex()
	importer := memento.NewImporter()
	ctx := context.Background()
	total := embedVendorStats{}

	vendor := func(pageLocation string, link *yunyun.ExtractedLink, request memento.Request) (string, error) {
		fetched, err := importer.Fetch(ctx, request)

		title := strings.TrimSpace(link.Text)
		if err == nil && strings.TrimSpace(fetched.Snapshot.Title) != "" {
			title = fetched.Snapshot.Title
		}
		base := embedManifestSlug(title, request)
		reference, exists, errPath := availableEmbedReference(conf, pageLocation, base, request)
		if errPath != nil {
			return "", errPath
		}
		if exists {
			return macroForEmbed(request, reference), nil
		}
		if dryRun {
			if err != nil {
				total.manual++
			}
			return macroForEmbed(request, reference), nil
		}

		if err == nil {
			manifestPath, pathErr := memento.ResolveReferencePath(conf, pageLocation, reference)
			if pathErr != nil {
				return "", pathErr
			}
			if writeErr := memento.WriteFetchedPath(manifestPath, fetched); writeErr != nil {
				return "", writeErr
			}
		} else {
			total.manual++
			manual := manualEmbedTemplate(request, title)
			if writeErr := memento.SaveReference(conf, pageLocation, reference, manual); writeErr != nil {
				return "", writeErr
			}
			logger.Warn("Wrote manual embed template after retrieval failed",
				"url", request.URL, "manifest", reference, "err", err)
		}
		return macroForEmbed(request, reference), nil
	}

	for _, sourcePath := range hizuru.FindFilesByExtSimple(conf) {
		data, err := os.ReadFile(filepath.Clean(string(sourcePath)))
		if err != nil {
			total.failed++
			logger.Warn("Could not read source while vendoring embeds", "file", sourcePath, "err", err)
			continue
		}
		relativeSource := conf.Runtime.WorkDir.Rel(sourcePath)
		pageLocation := filepath.ToSlash(filepath.Dir(string(relativeSource)))
		if pageLocation == "." {
			pageLocation = ""
		}
		rewritten, stats := rewriteEmbedLinks(string(data), pageLocation, vendor)
		total.found += stats.found
		total.vendored += stats.vendored
		total.skipped += stats.skipped
		total.failed += stats.failed
		if rewritten == string(data) || dryRun {
			continue
		}
		if err := writeSourceAtomically(string(sourcePath), []byte(rewritten)); err != nil {
			total.failed++
			logger.Warn("Could not rewrite vendored embed links", "file", sourcePath, "err", err)
		}
	}

	logger.Info("Finished vendoring embed links",
		"found", total.found,
		"vendored", total.vendored,
		"skipped", total.skipped,
		"manual", total.manual,
		"failures", total.failed,
		"dry_run", dryRun,
	)
}

func rewriteEmbedLinks(source, pageLocation string, vendor vendorLinkFunc) (string, embedVendorStats) {
	lines := strings.Split(source, "\n")
	stats := embedVendorStats{}
	pendingSkip := false
	blockDepth := 0
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "#+begin_") {
			blockDepth++
			pendingSkip = false
			continue
		}
		if blockDepth > 0 {
			if strings.HasPrefix(lower, "#+end_") {
				blockDepth--
			}
			continue
		}
		if strings.HasPrefix(lower, "#+attr_darkness:") {
			attributes := strings.Fields(strings.TrimSpace(strings.TrimPrefix(lower, "#+attr_darkness:")))
			pendingSkip = pendingSkip || containsString(attributes, skipEmbedVendoringAttribute) ||
				containsString(attributes, "embed-remote")
			continue
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "#+") || strings.HasPrefix(trimmed, "# ") {
			continue
		}

		link := yunyun.ExtractLink(trimmed)
		if link == nil || link.MatchLength != len(trimmed) || !standaloneOrgElement(lines, index) {
			pendingSkip = false
			continue
		}
		request, ok := memento.ParseURL(link.Link)
		if !ok {
			pendingSkip = false
			continue
		}
		stats.found++
		if pendingSkip {
			stats.skipped++
			pendingSkip = false
			continue
		}

		replacement, err := vendor(pageLocation, link, request)
		if err != nil {
			stats.failed++
			pendingSkip = false
			continue
		}
		indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		lines[index] = indent + replacement
		stats.vendored++
		pendingSkip = false
	}
	return strings.Join(lines, "\n"), stats
}

func standaloneOrgElement(lines []string, index int) bool {
	previousBoundary := index == 0
	for previous := index - 1; previous >= 0; previous-- {
		trimmed := strings.TrimSpace(lines[previous])
		if strings.HasPrefix(strings.ToLower(trimmed), "#+") {
			continue
		}
		previousBoundary = orgParagraphBoundary(trimmed)
		break
	}
	if !previousBoundary {
		return false
	}
	if index == len(lines)-1 {
		return true
	}
	return orgParagraphBoundary(strings.TrimSpace(lines[index+1]))
}

func orgParagraphBoundary(trimmed string) bool {
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return true
	}
	stars := 0
	for stars < len(trimmed) && trimmed[stars] == '*' {
		stars++
	}
	return stars > 0 && stars < len(trimmed) && trimmed[stars] == ' '
}

func embedManifestSlug(title string, request memento.Request) string {
	var slug strings.Builder
	previousDash := false
	for _, char := range strings.ToLower(strings.TrimSpace(title)) {
		switch {
		case unicode.IsLetter(char) || unicode.IsDigit(char):
			slug.WriteRune(char)
			previousDash = false
		case !previousDash && slug.Len() > 0:
			slug.WriteByte('-')
			previousDash = true
		}
		if slug.Len() >= 64 {
			break
		}
	}
	result := strings.Trim(slug.String(), "-")
	generic := result == "" || result == "song" || result == "spotify-track" ||
		result == "spotify-playlist" || result == "youtube-video"
	if generic {
		result = string(request.Kind) + "-" + request.ID
	}
	return result
}

func availableEmbedReference(conf *alpha.DarknessConfig, pageLocation, base string, request memento.Request) (string, bool, error) {
	candidates := []string{base + ".toml", base + "-" + shortEmbedID(request.ID) + ".toml"}
	for suffix := 2; suffix < 100; suffix++ {
		candidates = append(candidates, fmt.Sprintf("%s-%d.toml", base, suffix))
	}
	for _, reference := range candidates {
		manifestPath, err := memento.ResolveReferencePath(conf, pageLocation, reference)
		if err != nil {
			return "", false, err
		}
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			return reference, false, nil
		}
		existing, err := memento.ReadReference(conf, pageLocation, reference)
		if err == nil {
			existingRequest, ok := memento.ParseURL(existing.URL)
			if ok && existingRequest.Kind == request.Kind && existingRequest.ID == request.ID {
				return reference, true, nil
			}
		}
	}
	return "", false, fmt.Errorf("could not find an available embed manifest name for %s", request.URL)
}

func manualEmbedTemplate(request memento.Request, title string) *memento.Snapshot {
	provider := "Spotify"
	fallbackTitle := "Spotify media — fill out manually"
	if request.Kind == memento.KindYouTube {
		provider = "YouTube"
		fallbackTitle = "YouTube video — fill out manually"
	}
	if strings.TrimSpace(title) == "" {
		title = fallbackTitle
	}
	return &memento.Snapshot{
		Version: memento.SnapshotVersion, NeedsManual: true,
		Provider: provider, Kind: request.Kind, ID: request.ID,
		URL: request.URL, Title: title,
	}
}

func macroForEmbed(request memento.Request, reference string) string {
	if request.Kind == memento.KindYouTube {
		return "{{{youtube_embed(" + reference + ")}}}"
	}
	return "{{{spotify_embed(" + reference + ")}}}"
}

func shortEmbedID(id string) string {
	if len(id) <= 8 {
		return id
	}
	return id[:8]
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func writeSourceAtomically(path string, data []byte) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	temporary := path + ".embed-vendor.tmp"
	if err := os.WriteFile(temporary, data, info.Mode().Perm()); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}
