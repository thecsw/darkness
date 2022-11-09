package orgmode

import (
	"strings"

	"github.com/thecsw/darkness/internals"
)

// isHeader returns a non-nil object if the line is a header
func isHeader(line string) *internals.Content {
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
	return &internals.Content{
		Type:         internals.TypeHeading,
		HeadingLevel: level,
		Heading:      line[level+1:],
	}
}

// isComment returns true if the line is a comment
func isComment(line string) bool {
	return strings.HasPrefix(line, commentPrefix)
}

// isOption returns true if the line is an option
func isOption(line string) bool {
	return strings.HasPrefix(line, optionPrefix)
}

// isLink returns a non-nil object if the line is a link
func isLink(line string) *internals.Content {
	line = strings.TrimSpace(line)
	// Not a link
	if !linkRegexp.MatchString(line) {
		return nil
	}
	submatches := linkRegexp.FindAllStringSubmatch(line, 1)
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

func formParagraph(text, extra string, options internals.Bits) *internals.Content {
	val := &internals.Content{
		Type:      internals.TypeParagraph,
		Paragraph: strings.TrimSpace(text),
		Options:   options,
	}
	if internals.HasFlag(&options, internals.InDetailsFlag) {
		val.Summary = extra
	}
	return val
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
	return strings.HasPrefix(strings.ToLower(line), optionPrefix+optionBeginSource)
}

func isSourceCodeEnd(line string) bool {
	return strings.ToLower(line) == optionPrefix+optionEndSource
}

func sourceExtractLang(line string) string {
	return sourceCodeRegexp.FindAllStringSubmatch(strings.ToLower(line), 1)[0][1]
}

func detailsExtractSummary(line string) string {
	return detailsRegexp.FindAllStringSubmatch(strings.ToLower(line), 1)[0][1]
}

func isHTMLExportBegin(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), optionPrefix+optionBeginExport+" html")
}

func isHTMLExportEnd(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), optionPrefix+optionEndExport)
}

func isHorizonalLine(line string) bool {
	return strings.TrimSpace(line) == horizontalLine
}

func isAttentionBlock(line string) *internals.Content {
	matches := attentionBlockRegexp.FindAllStringSubmatch(line, 1)
	if len(matches) < 1 {
		return nil
	}
	return &internals.Content{
		Type:           internals.TypeAttentionText,
		AttentionTitle: matches[0][1],
		AttentionText:  matches[0][2],
	}
}
