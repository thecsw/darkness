package orgmode

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
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

// isOption returns true if the line is an option
func isOption(line string) bool {
	return strings.HasPrefix(line, optionPrefix)
}

// getLink returns a non-nil object if the line is a link
func getLink(line string) *yunyun.Content {
	line = strings.TrimSpace(line)
	extractedLink := yunyun.ExtractLink(line)
	// Extraction didn't yield any results.
	if extractedLink == nil {
		return nil
	}
	// Check if this is a standalone link (just by itself on a line)
	// If it's not, then it's a simple link in a paragraph, deal with
	// it later in `htmlize`
	if extractedLink.MatchLength != len(line) {
		return nil
	}
	return &yunyun.Content{
		Type:            yunyun.TypeLink,
		Link:            extractedLink.Link,
		LinkTitle:       extractedLink.Text,
		LinkDescription: extractedLink.Description,
	}
}

// formParagraph builds a proper paragraph-oriented `Content` object.
func formParagraph(text, extra string, options yunyun.Bits) *yunyun.Content {
	val := &yunyun.Content{
		Type:      yunyun.TypeParagraph,
		Paragraph: strings.TrimSpace(text),
		Options:   options,
	}
	if val.IsDetails() {
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
	return strings.HasPrefix(strings.ToLower(line), optionPrefix+optionBeginSource)
}

// isSourceCodeEnd returns true if we are currently reading the end of a
// source code block, false otherwise.
func isSourceCodeEnd(line string) bool {
	return strings.ToLower(line) == optionPrefix+optionEndSource
}

// isHTMLExportBegin returns true if we are currently reading the start
// of an html export block, false otherwise.
func isHTMLExportBegin(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), optionPrefix+optionBeginExport+" html")
}

// isHTMLExportEnd returns true if we are currently reading the end of an
// html export block, false otherwise.
func isHTMLExportEnd(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), optionPrefix+optionEndExport)
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

// extractOptionLabel is a utility function used to extract option values.
func extractOptionLabel(given string, option string) string {
	return strings.TrimSpace(gana.SkipString(uint(len(optionPrefix)+len(option)), given))
}

// extractSourceCodeLanguage extracts language `LANG` from `#+begin_src LANG`.
func extractSourceCodeLanguage(line string) string {
	return extractOptionLabel(line, optionBeginSource)
}

// extractDetailsSummary extracts summary `SUMMARY` from `#+begin_details SUMMARY`.
func extractDetailsSummary(line string) string {
	return extractOptionLabel(line, optionBeginDetails)
}

// extractHtmlHead extracts the html to inject in the head, like custom CSS.
func extractHtmlHead(line string) string {
	return extractOptionLabel(line, optionHtmlHead)
}

// extractOptions extracts the line of options.
func extractOptions(line string) string {
	return extractOptionLabel(line, optionOptions)
}

// extractCaptionTitle extracts caption `TITLE` from `#+caption: TITLE`.
func extractCaptionTitle(line string) string {
	return extractOptionLabel(line, optionCaption)
}

// extractDate extracts date `DATE` from `#+date: DATE`.
func extractDate(line string) string {
	return extractOptionLabel(line, optionDate)
}

// extractGalleryFolder extracts gallery `FOLDER` from `#+begin_gallery FOLDER`.
func extractGalleryFolder(line string) string {
	path, err := extractCustomBlockOption(line, `path`, regexpPatternNoWhitespace)
	if err != nil {
		if err != errNoMatches {
			fmt.Println("gallery path extraction failed:", err.Error())
		}
		return ""
	}
	return *path
}

func extractGalleryImagesPerRow(line string) uint {
	num, err := extractCustomBlockOption(line, `num`, regexpPatternOnlyDigits)
	if err != nil {
		if err != errNoMatches {
			fmt.Println("gallery path extraction failed:", err.Error())
		}
		return defaultGalleryImagesPerRow
	}
	ans, err := strconv.Atoi(*num)
	if err != nil {
		fmt.Println("failed to format gallery width of", line, ", defaulting to", defaultGalleryImagesPerRow)
		return defaultGalleryImagesPerRow
	}
	if ans < 1 {
		fmt.Println("gallery width should be at least 1, defaulting to", defaultGalleryImagesPerRow)
		return defaultGalleryImagesPerRow
	}
	return uint(ans)
}

type regexpPattern string

const (
	regexpPatternNoWhitespace regexpPattern = `([^\s]+)`
	regexpPatternOnlyDigits   regexpPattern = `(\d+)`
	regexpPatternNumber       regexpPattern = `(-?\d+)`
)

var (
	errNoMatches = errors.New(`no matches found`)
)

func extractCustomBlockOption(target, optionName string, pattern regexpPattern) (*string, error) {
	optP := fmt.Sprintf(`:%s %s`, optionName, pattern)
	optR, err := regexp.Compile(optP)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make regex "+optP)
	}
	matches := optR.FindAllStringSubmatch(target, 1)
	if len(matches) < 1 {
		return nil, errNoMatches
	}
	if len(matches[0]) < 1 {
		return nil, errNoMatches
	}
	return &matches[0][1], nil
}
