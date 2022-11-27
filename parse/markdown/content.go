package markdown

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// isHeader returns a non-nil object if the line is a header
func isHeader(line string) *yunyun.Content {
	level := 0
	switch {
	case strings.HasPrefix(line, sectionLevelOne):
		level = 1
	case strings.HasPrefix(line, sectionLevelTwo):
		level = 2
	case strings.HasPrefix(line, sectionLevelThree):
		level = 3
	case strings.HasPrefix(line, sectionLevelFour):
		level = 4
	case strings.HasPrefix(line, sectionLevelFive):
		level = 5
	default:
		level = 0
	}
	// Not a header
	if level < 1 {
		return nil
	}
	// Is a header
	return &yunyun.Content{
		Type:         yunyun.TypeHeading,
		HeadingLevel: level,
		Heading:      line[level+1:],
	}
}

// isComment returns true if the line is a comment
func isComment(line string) bool {
	return strings.HasPrefix(line, commentPrefix)
}

// isLink returns a non-nil object if the line is a link
func isLink(line string) *yunyun.Content {
	line = strings.TrimSpace(line)
	matchLen, link, text, desc := yunyun.ExtractLink(line)
	// Extraction didn't yield any results.
	if matchLen < 0 {
		return nil
	}
	// Check if this is a standalone link (just by itself on a line)
	// If it's not, then it's a simple link in a paragraph, deal with
	// it later in `htmlize`
	if matchLen != len(line) {
		return nil
	}
	return &yunyun.Content{
		Type:            yunyun.TypeLink,
		Link:            link,
		LinkTitle:       text,
		LinkDescription: desc,
	}
}

// formParagraph builds a proper paragraph-oriented `Content` object.
func formParagraph(text, extra string, options yunyun.Bits) *yunyun.Content {
	val := &yunyun.Content{
		Type:      yunyun.TypeParagraph,
		Paragraph: strings.TrimSpace(text),
		Options:   options,
	}
	if yunyun.HasFlag(&options, yunyun.InDetailsFlag) {
		val.Summary = extra
	}
	return val
}

// isList returns true if we are currently reading a list, false otherwise.
func isList(line string) bool {
	return strings.HasPrefix(line, "- ")
}

// isTable returns true if we are currently reading a table, false otherwise.
func isTable(line string) bool {
	return strings.HasPrefix(line, "| ") || strings.HasPrefix(line, "|-")
}

// isTableHeaderDelimeter returns true if we are currently reading a table
// header delimiter, false otherwise.
func isTableHeaderDelimeter(line string) bool {
	return strings.HasPrefix(line, "|-")
}

// isSourceCodeBegin returns true if we are currently reading the start of
// a source code block, false otherwise.
func isSourceCodeBegin(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), optionBeginSource)
}

// isSourceCodeEnd returns true if we are currently reading the end of a
// source code block, false otherwise.
func isSourceCodeEnd(line string) bool {
	return strings.ToLower(line) == optionEndSource
}

// isHTMLExportBegin returns true if we are currently reading the start
// of an html export block, false otherwise.
func isHTMLExportBegin(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), "<")
}

// isHTMLExportEnd returns true if we are currently reading the end of an
// html export block, false otherwise.
func isHTMLExportEnd(line string) bool {
	return strings.HasPrefix(strings.ToLower(line)+" ", " ")
}

// isHorizonalLine returns true if we are currently reading a horizontal line,
// false otherwise.
func isHorizonalLine(line string) bool {
	return strings.TrimSpace(line) == horizontalLine
}

// isAttentionBlock returns *Content object if we have fonud an attention block
// with filled values, nil otherwise.
func isAttentionBlock(line string) *yunyun.Content {
	matches := attentionBlockRegexp.FindAllStringSubmatch(line, 1)
	if len(matches) < 1 {
		return nil
	}
	return &yunyun.Content{
		Type:           yunyun.TypeAttentionText,
		AttentionTitle: matches[0][1],
		AttentionText:  matches[0][2],
	}
}

// extractSourceCodeLanguage extracts language `LANG` from `#+begin_src LANG`.
func extractSourceCodeLanguage(line string) string {
	return strings.TrimPrefix(line, optionBeginSource)
}
