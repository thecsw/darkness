package main

import (
	"regexp"
	"strings"
)

var (
	OrgLinkRegexp = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	OrgBoldText   = regexp.MustCompile(` \*([^* ][^*]+[^* ]|[^*])\*([^\w.])`)
	OrgItalicText = regexp.MustCompile(` /([^/ ][^/]+[^/ ]|[^/])/([^\w.])`)
)

func Parse(lines []string) *Page {
	page := &Page{}
	page.Contents = make([]Content, 0, 16)

	currentContext := ""
	inList := false
	currentList := make([]string, 0, 8)
	for i, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		currentContext = currentContext + line
		// If it's an empty line, then process current text
		if line == "" {
			// Empty context
			if currentContext == "" {
				continue
			}
			// Let's see if our context is a standalone link
			if link := isLink(currentContext); link != nil {
				page.Contents = append(page.Contents, *link)
				currentContext = ""
				continue
			}
			// New line break means we have to save the paragraph
			// we just read if we're not currently reading a list
			if !inList {
				page.Contents = append(
					page.Contents,
					*formParagraph(strings.TrimSpace(currentContext)))
			}
			currentContext = ""
		}
		// We are in a list now
		if isList(line) && i != 0 {
			// If we were not in a list, save the current context
			if !inList && len(currentContext) != len(line)+1 {
				page.Contents = append(
					page.Contents,
					*formParagraph(strings.TrimSpace(currentContext[:len(currentContext)-len(line)])))
				currentContext = ""
			}
			inList = true
			currentList = append(currentList, line[2:])
		} else if inList {
			page.Contents = append(page.Contents, Content{
				Type: TypeList,
				List: currentList,
			})
			currentList = []string{}
			inList = false
			currentContext = currentContext[len(currentContext)-len(line):]
		}
		// Find whether the current line is a part of a list
		// A header is found, append and continue
		if header := isHeader(line); header != nil &&
			(((i == 0) && header.HeaderLevel == 1) || header.HeaderLevel > 1) {
			currentContext = ""
			// Level 1 is the page title
			if header.HeaderLevel == 1 {
				page.Title = header.Header
				continue
			}
			page.Contents = append(page.Contents, *header)
			continue
		}
		currentContext += " "
	}
	return page
}

func isHeader(line string) *Content {
	level := 0
	for _, c := range line {
		if c != '*' {
			break
		}
		level++
	}
	// Not a header
	if level < 1 {
		return nil
	}
	// Is a header
	return &Content{
		Type:        TypeHeader,
		HeaderLevel: level,
		Header:      line[level+1:],
	}
}

func isLink(line string) *Content {
	line = strings.TrimSpace(line)
	// Not a link
	if !OrgLinkRegexp.MatchString(line) {
		return nil
	}
	submatches := OrgLinkRegexp.FindAllStringSubmatch(line, 1)
	// Sanity check
	if len(submatches) < 1 {
		return nil
	}
	match := strings.TrimSpace(submatches[0][0])
	link := strings.TrimSpace(submatches[0][1])
	text := strings.TrimSpace(submatches[0][2])
	// Check if this is a standalone link (just by itself on a line)
	// If it's not, then it's a simple link in a paragraph, deal with
	// it in a different way
	if len(match) != len(line) {
		return nil
	}
	content := &Content{
		Type:      TypeLink,
		Link:      link,
		LinkTitle: text,
	}
	// Our link is standalone. Check if it's an image
	if strings.HasSuffix(link, ".png") {
		content.Type = TypeImage
		content.ImageSource = link
		content.ImageCaption = text
		return content
	}
	// Check if it's a youtube video embed
	if strings.HasPrefix(link, "https://youtu.be/") {
		content.Type = TypeYoutube
		content.Youtube = link[17:]
		return content
	}
	// Check if it's a spotify track link
	if strings.HasPrefix(link, "https://open.spotify.com/track/") {
		content.Type = TypeSpotifyTrack
		content.SpotifyTrack = link[31:]
		return content
	}
	// Check if it's a spotify playlist link
	if strings.HasPrefix(link, "https://open.spotify.com/playlist/") {
		content.Type = TypeSpotifyPlaylist
		content.SpotifyPlaylist = link[34:]
		return content
	}
	return nil
}

func formParagraph(text string) *Content {
	return &Content{
		Type:      TypeParagraph,
		Paragraph: text,
	}
}

var (
	listPrefixes = []string{"* ", "- ", "+ "}
)

func isList(line string) bool {
	for _, prefix := range listPrefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}
