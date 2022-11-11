package emilia

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// WithEnrichedHeadings shifts heading levels to their correct layouts and
// adds some additional information to the headings for later export
func WithEnrichedHeadings() yunyun.PageOption {
	return func(page *yunyun.Page) {
		// Normalizing headings
		if Config.Website.NormalizeHeadings {
			minHeadingLevel := 999
			// Find the smallest heading
			for i := range page.Contents {
				c := &page.Contents[i]
				if !c.IsHeading() {
					continue
				}
				if c.HeadingLevel < minHeadingLevel {
					minHeadingLevel = c.HeadingLevel
				}
			}
			// Shift everything over
			for i := range page.Contents {
				c := &page.Contents[i]
				if !c.IsHeading() {
					continue
				}
				c.HeadingLevel -= (minHeadingLevel - 2)
			}
		}
	}
}

// WithResolvedComments resolves heading comments and cleans up the page if
// COMMENT headings are encountered
func WithResolvedComments() yunyun.PageOption {
	return func(page *yunyun.Page) {
		start, headingLevel, searching := -1, -1, false
		for i, content := range page.Contents {
			if !content.IsHeading() {
				continue
			}
			if strings.HasPrefix(content.Heading, "COMMENT ") && !searching {
				start = i
				headingLevel = content.HeadingLevel
				searching = true
				continue
			}
			if searching && content.HeadingLevel <= headingLevel {
				page.Contents = append(page.Contents[:start], page.Contents[i:]...)
				start, headingLevel, searching = -1, -1, false
			}
		}
		// Still searching till the end? then set the finish to the last element
		if searching {
			page.Contents = page.Contents[:start]
		}
	}
}
