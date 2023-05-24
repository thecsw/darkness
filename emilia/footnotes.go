package emilia

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// FootnoteLabeler will take an integer and return string representation as
// defined in the darkness config, with either Roman or Arabic numerals.
var FootnoteLabeler = strconv.Itoa

// WithFootnotes resolves footnotes and cleans up the page if necessary
func WithFootnotes() yunyun.PageOption {
	return func(page *yunyun.Page) {
		footnotes := make([]string, 0, 4)
		for i := range page.Contents {
			c := page.Contents[i]
			// Replace footnotes in paragraphs
			if c.IsParagraph() {
				c.Paragraph = findFootnotes(c.Paragraph, &footnotes)
			}
			// Footnotes can also appear in lists
			if c.IsList() {
				for i := 0; i < len(c.List); i++ {
					c.List[i] = findFootnotes(c.List[i], &footnotes)
				}
			}
		}
		page.Footnotes = footnotes
	}
}

// findFootnotes finds footnotes in a paragraph and replaces them with a footnote reference
func findFootnotes(text string, footnotes *[]string) string {
	matches := yunyun.FootnoteRegexp.FindAllStringSubmatch(text, -1)
	// no footnotes found
	if len(matches) < 1 {
		return text
	}
	newText := text
	for _, match := range matches {
		*footnotes = append(*footnotes, match[1])
		newText = strings.Replace(newText, match[0], fmt.Sprintf("!%d!", len(*footnotes)), 1)
	}
	return newText
}
