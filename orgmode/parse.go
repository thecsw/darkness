package orgmode

import (
	"darkness/emilia"
	"darkness/internals"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ParseFile(workDir, file string) *internals.Page {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	page := Parse(string(data))
	page.URL = emilia.JoinPath(strings.TrimPrefix(filepath.Dir(file), workDir))
	return page
}

func Parse(data string) *internals.Page {
	// Pad a newline so that last elements can be processed
	// properly before an EOF is encountered during parsing
	data += "\n"
	lines := strings.Split(data, "\n")
	page := &internals.Page{}
	page.Contents = make([]internals.Content, 0, 16)

	inList := false
	inSourceCode := false
	inRawHTML := false
	sourceCodeLang := ""

	// Our context is a parody of a state machine
	currentContext := ""
	addContent := func(content internals.Content) {
		page.Contents = append(page.Contents, content)
		currentContext = ""
	}

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		previousContext := currentContext
		currentContext = currentContext + line

		// If we are in a raw html envoronment
		if inRawHTML {
			// Maybe it's time to leave it?
			if isHTMLExportEnd(line) {
				// Mark the leave
				inRawHTML = false
				// Save the raw html
				addContent(internals.Content{
					Type:    internals.TypeRawHTML,
					RawHTML: previousContext,
				})
				continue
			}
			// Otherwise, continue saving the context
			currentContext = previousContext + rawLine + "\n"
			continue
		}
		// Now, check if we can enter a raw html environment
		if isHTMLExportBegin(line) {
			inRawHTML = true
			currentContext = previousContext
			continue
		}
		// If we are in a source code block?
		if inSourceCode {
			// Check if it's time to leave
			if isSourceCodeEnd(line) {
				// Mark the leave
				inSourceCode = false
				// Save the source code
				addContent(internals.Content{
					Type:           internals.TypeSourceCode,
					SourceCodeLang: sourceCodeLang,
					SourceCode:     strings.TrimRight(previousContext, "\n"),
				})
				continue
			}
			// Save the context and continue
			currentContext = previousContext + rawLine + "\n"
			continue
		}
		// Should we enter a source code environment?
		if isSourceCodeBegin(line) {
			sourceCodeLang = sourceExtractLang(line)
			inSourceCode = true
			currentContext = ""
			continue
		}
		if isComment(line) || isOption(line) {
			currentContext = previousContext
			continue
		}
		if isHorizonalLine(line) {
			page.Contents = append(page.Contents, internals.Content{
				Type: internals.TypeHorizontalLine,
			})
			currentContext = previousContext
			continue
		}

		// Now, we need to parse headings here
		if header := isHeader(line); header != nil {
			if header.HeaderLevel == 1 {
				page.Title = header.Header
				currentContext = ""
				continue
			}
			addContent(*header)
			continue
		}
		// If we hit an empty line, end the whatever context we had
		if line == "" {
			// If we were in a list, save it as a list
			if inList {
				matches := UnorderedListRegexp.FindAllStringSubmatch(previousContext, -1)
				// Shouldn't happen, continue as a failure
				if len(matches) < 1 {
					continue
				}
				currentList := make([]string, 0, len(matches))
				for _, match := range matches {
					currentList = append(currentList, match[1])
				}
				addContent(internals.Content{
					Type: internals.TypeList,
					List: currentList,
				})
				continue
			}
			// Otherwise, save as a paragraph if not empty
			if len(previousContext) < 1 {
				continue
			}
			addContent(*formParagraph(previousContext))
			continue
		}
		currentContext += " "
		// Mark if the current line is a list
		inList = isList(line)
		// Add a delimeter
		if inList {
			currentContext += "∆"
		}
	}
	return page
}

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
		Type:        internals.TypeHeader,
		HeaderLevel: level,
		Header:      line[level+1:],
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
