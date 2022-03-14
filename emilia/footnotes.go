package emilia

import (
	"darkness/internals"
	"fmt"
	"strings"
)

func ResolveFootnotes(page *internals.Page) {
	footnotes := make([]string, 0, 4)
	fnCounter := 0
	for i := range page.Contents {
		c := &page.Contents[i]
		// Only replace footnotes in paragraphs
		if !c.IsParagraph() {
			continue
		}
		matches := internals.FootnoteRegexp.FindAllStringSubmatch(c.Paragraph, -1)
		// no footnotes found
		if len(matches) < 1 {
			continue
		}
		newParagraph := c.Paragraph
		for _, match := range matches {
			fnCounter++
			footnotes = append(footnotes, match[1])
			newParagraph = strings.Replace(newParagraph, match[0], fmt.Sprintf("!%d!", fnCounter), 1)
		}
		c.Paragraph = newParagraph
	}
	page.Footnotes = footnotes
}
