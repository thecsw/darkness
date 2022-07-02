package orgmode

import (
	"strings"

	"github.com/thecsw/darkness/internals"
)

// isHeader returns a non-nil object if the line is a header
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
		Heading:      line[level+1:],
	}
}

// isComment returns true if the line is a comment
func isComment(line string) bool {
	return strings.HasPrefix(line, "# ")
}

// isOption returns true if the line is an option
func isOption(line string) bool {
	return strings.HasPrefix(line, "#+")
}

// isLink returns a non-nil object if the line is a link
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
	return &internals.Content{
		Type:      internals.TypeLink,
		Link:      link,
		LinkTitle: text,
	}
}

func formParagraph(text string, inQuote bool, inCenter bool, inDropCap bool) *internals.Content {
	return &internals.Content{
		Type:       internals.TypeParagraph,
		Paragraph:  strings.TrimSpace(text),
		IsCentered: inCenter,
		IsQuote:    inQuote,
		IsDropCap:  inDropCap,
	}
}

func isList(line string) bool {
	return strings.HasPrefix(line, "- ")
}

func isTable(line string) bool {
	return strings.HasPrefix(line, "| ") || strings.HasPrefix(line, "|-")
}

func isTableHeaderDelimeter(line string) bool {
	return strings.HasPrefix(line, "|-")
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
	return strings.HasPrefix(strings.ToLower(line), "#+begin_export html")
}

func isHTMLExportEnd(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), "#+end_export")
}

func isHorizonalLine(line string) bool {
	return strings.TrimSpace(line) == "---"
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
