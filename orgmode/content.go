package orgmode

import (
	"darkness/internals"
	"strings"
)

func isHeader(line string) *internals.Content {
	level := 0
	switch {
	case strings.HasPrefix(line, "* "):
		level = 1
	case strings.HasPrefix(line, "** "):
		level = 2
	case strings.HasPrefix(line, "*** "):
		level = 3
	case strings.HasPrefix(line, "**** "):
		level = 4
	case strings.HasPrefix(line, "***** "):
		level = 5
	default:
		level = 0
	}
	// Not a header
	if level < 1 {
		return nil
	}
	// Is a header
	return &internals.Content{
		Type:         internals.TypeHeading,
		HeadingLevel: level,
		Header:       line[level+1:],
	}
}

func isComment(line string) bool {
	return strings.HasPrefix(line, "# ")
}

func isOption(line string) bool {
	return strings.HasPrefix(line, "#+")
}

func isLink(line string) *internals.Content {
	line = strings.TrimSpace(line)
	// Not a link
	if !LinkRegexp.MatchString(line) {
		return nil
	}
	submatches := LinkRegexp.FindAllStringSubmatch(line, 1)
	// Sanity check
	if len(submatches) < 1 {
		return nil
	}
	match := strings.TrimSpace(submatches[0][0])
	link := strings.TrimSpace(submatches[0][1])
	text := strings.TrimSpace(submatches[0][2])
	// Check if this is a standalone link (just by itself on a line)
	// If it's not, then it's a simple link in a paragraph, deal with
	// it later in `htmlize`
	if len(match) != len(line) {
		return nil
	}
	content := &internals.Content{
		Type:      internals.TypeLink,
		Link:      link,
		LinkTitle: text,
	}
	// Our link is standalone. Check if it's an image
	if strings.HasSuffix(link, ".png") {
		content.Type = internals.TypeImage
		content.ImageSource = link
		content.ImageCaption = text
		return content
	}
	// Check if it's a youtube video embed
	if strings.HasPrefix(link, "https://youtu.be/") {
		content.Type = internals.TypeYoutube
		content.Youtube = link[17:]
		return content
	}
	// Check if it's a spotify track link
	if strings.HasPrefix(link, "https://open.spotify.com/track/") {
		content.Type = internals.TypeSpotifyTrack
		content.SpotifyTrack = link[31:]
		return content
	}
	// Check if it's a spotify playlist link
	if strings.HasPrefix(link, "https://open.spotify.com/playlist/") {
		content.Type = internals.TypeSpotifyPlaylist
		content.SpotifyPlaylist = link[34:]
		return content
	}
	return nil
}

func formParagraph(text string) *internals.Content {
	return &internals.Content{
		Type:      internals.TypeParagraph,
		Paragraph: text,
	}
}

func isList(line string) bool {
	return strings.HasPrefix(line, "- ")
}

func isSourceCodeBegin(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), "#+begin_src")
}

func isSourceCodeEnd(line string) bool {
	return strings.ToLower(line) == "#+end_src"
}

func sourceExtractLang(line string) string {
	return SourceCodeRegexp.FindAllStringSubmatch(strings.ToLower(line), 1)[0][1]
}

func isHTMLExportBegin(line string) bool {
	return line == "#+begin_export html"
}

func isHTMLExportEnd(line string) bool {
	return line == "#+end_export"
}

func isHorizonalLine(line string) bool {
	return line == "---"
}

func isAttentionBlack(line string) *internals.Content {
	matches := AttentionBlockRegexp.FindAllStringSubmatch(line, 1)
	if len(matches) < 1 {
		return nil
	}
	return &internals.Content{
		Type:           internals.TypeAttentionText,
		AttentionTitle: matches[0][1],
		AttentionText:  matches[0][2],
	}
}
