package emilia

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/thecsw/darkness/yunyun"
)

// GenerateTableOfContents generates a table of contents for a page.
func GenerateTableOfContents(page *yunyun.Page) []yunyun.ListItem {
	toc := make([]yunyun.ListItem, len(page.Contents.Headings()))
	for i, heading := range page.Contents.Headings() {
		toc[i] = yunyun.ListItem{
			Level: uint8(heading.HeadingLevelAdjusted),
			Text:  fmt.Sprintf("[[%s][%s]]", "#"+ExtractID(heading.Heading), heading.Heading),
		}
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

	res := ""
	for _, c := range heading {
		if unicode.IsSpace(c) || unicode.IsPunct(c) || unicode.IsSymbol(c) {
			res += "-"
			continue
		}
		if c <= unicode.MaxASCII {
			res += string(unicode.ToLower(c))
		}
	}
	return strings.TrimRight(res, "-")
}
