package html

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/thecsw/darkness/v3/yunyun"
)

// GenerateTableOfContents generates a table of contents for a page.
func GenerateTableOfContents(page *yunyun.Page) []yunyun.ListItem {
	toc := make([]yunyun.ListItem, 0, len(page.Contents.Headings()))
	for _, heading := range page.Contents.Headings() {
		// The user could have excluded the heading from appearing in the index.
		if yunyun.HasFlag(&heading.Options, yunyun.HeadingNoIndexFlag) {
			continue
		}

		// Bound check to prevent integer overflow when converting uint32 to uint8
		var level uint8
		if heading.HeadingLevelAdjusted > 255 {
			level = 255
		} else {
			level = uint8(heading.HeadingLevelAdjusted)
		}

		toc = append(toc, yunyun.ListItem{
			Level: level,
			Text:  fmt.Sprintf("[[%s][%s]]", "#"+ExtractID(heading.Heading), heading.Heading),
		})
	}
	return toc
}

// ExtractID returns a properly formatted ID for a heading title
func ExtractID(heading string) string {
	// Check if heading is a link
	extractedLink := yunyun.ExtractLink(heading)
	if extractedLink != nil {
		heading = extractedLink.Text // 0 is whole match, 1 is link, 2 is title
	}

	var res strings.Builder
	for _, c := range heading {
		if unicode.IsSpace(c) || unicode.IsPunct(c) || unicode.IsSymbol(c) {
			res.WriteString("-")
			continue
		}
		if c <= unicode.MaxASCII {
			res.WriteString(string(unicode.ToLower(c)))
		}
	}
	return strings.TrimRight(res.String(), "-")
}
