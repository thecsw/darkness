package orgmode

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
)

// isHeader returns a non-nil object if the line is a header
func isHeader(line string) *yunyun.Content {
	level := uint32(0)
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
func isOption(line string) (string, bool) {
	if !strings.HasPrefix(line, optionPrefix) {
		return "", false
	}
	option := gana.SkipString(uint(optionPrefixLen), line)
	parts := strings.SplitN(option, " ", 2)

	// Don't know what this is, don't let it reach the parser.
	if len(parts) < 1 {
		return "", false
	}
	val := strings.TrimSpace(parts[0])
	return val, len(val) > 0
}

// isOptionLine returns true if the line is an option (wrapper for isOption)
func isOptionLine(line string) bool {
	_, ok := isOption(line)
	return ok
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

// isOrderedListStart returns true if we are starting an ordered list.
func isOrderedListStart(line string) bool {
	return strings.HasPrefix(line, "1. ")
}

// listAnyRegex checks whether the start of the line signifies we are in an ordered list.
var listAnyRegex = regexp.MustCompile(`^[0-9]+[.] `)

// isOrderedListAny returns true if we are anywhere within the ordered list.
func isOrderedListAny(line string) bool {
	return listAnyRegex.MatchString(line)
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

// isHtmlExportBegin returns true if we are currently reading the start
// of an html export block, false otherwise.
func isHtmlExportBegin(line string) bool {
	return strings.HasPrefix(strings.ToLower(line), optionPrefix+optionBeginExport+" html")
}

// isHtmlExportEnd returns true if we are currently reading the end of an
// html export block, false otherwise.
func isHtmlExportEnd(line string) bool {
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

// safeIntToUint safely converts an int to uint by handling negative values
func safeIntToUint(val int) uint {
	if val < 0 {
		return 0
	}
	return uint(val)
}

// extractOptionLabel is a utility function used to extract option values.
func extractOptionLabel(given string, option string) string {
	skipLen := len(optionPrefix) + len(option)
	return strings.TrimSpace(gana.SkipString(safeIntToUint(skipLen), given))
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

// extractAttributes extracts the line of attributes.
func extractAttributes(line string) string {
	return extractOptionLabel(line, optionAttributes)
}

// extractHtmlTags extracts the line of html styles.
func extractHtmlTags(line string) string {
	return extractOptionLabel(line, optionHtmlTags)
}

// extractAttributes extracts the line of html attributes.
func extractHtmlAttributes(line string) string {
	return extractOptionLabel(line, optionAttrHtml)
}

// extractCaptionTitle extracts caption `TITLE` from `#+caption: TITLE`.
func extractCaptionTitle(line string) string {
	return extractOptionLabel(line, optionCaption)
}

// extractDate extracts date `DATE` from `#+date: DATE`.
func extractDate(line string) string {
	return extractOptionLabel(line, optionDate)
}

// extractDate extracts author `AUTHOR` from `#+author: AUTHOR`.
func extractAuthor(line string) string {
	return extractOptionLabel(line, optionAuthor)
}

// extractGalleryFolder extracts gallery `FOLDER` from `#+begin_gallery FOLDER`.
func extractGalleryFolder(line string) string {
	path, err := extractCustomBlockOption(line, `path`, regexpPatternNoWhitespace)
	if err != nil {
		if !errors.Is(err, errNoMatches) {
			puck.Logger.Errorf("gallery path extraction: %v", err)
		}
		return ""
	}
	return *path
}

func extractGalleryImagesPerRow(line string) uint {
	num, err := extractCustomBlockOption(line, `num`, regexpPatternOnlyDigits)
	if err != nil {
		if !errors.Is(err, errNoMatches) {
			puck.Logger.Errorf("gallery path extraction: %v", err)
		}
		return defaultGalleryImagesPerRow
	}
	ans, err := strconv.Atoi(*num)
	if err != nil {
		puck.Logger.Warnf("failed to format gallery width of %s, defaulting to %d", line, defaultGalleryImagesPerRow)
		return defaultGalleryImagesPerRow
	}
	if ans < 1 {
		puck.Logger.Warnf("gallery width should be at least 1, defaulting to %d", defaultGalleryImagesPerRow)
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

var errNoMatches = errors.New(`no matches found`)

func extractCustomBlockOption(target, optionName string, pattern regexpPattern) (*string, error) {
	optP := fmt.Sprintf(`:%s %s`, optionName, pattern)
	optR, err := regexp.Compile(optP)
	if err != nil {
		return nil, fmt.Errorf("compiling regex ('%s'): %v", optP, err)
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
