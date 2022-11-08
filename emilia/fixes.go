package emilia

import (
	"strings"

	"github.com/thecsw/darkness/internals"
)

// WithEnrichedHeadings shifts heading levels to their correct layouts and
// adds some additional information to the headings for later export
func WithEnrichedHeadings() internals.PageOption {
	return func(page *internals.Page) {
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
		// Mark the first heading
		for i := range page.Contents {
			c := &page.Contents[i]
			if c.IsHeading() {
				c.HeadingFirst = true
				break
			}
		}
		// Mark the last heading
		for i := len(page.Contents) - 1; i >= 0; i-- {
			c := &page.Contents[i]
			if c.IsHeading() {
				c.HeadingLast = true
				break
			}
		}
		// Mark headings that are children
		currentLevel := 0
		for i := range page.Contents {
			c := &page.Contents[i]
			if !c.IsHeading() {
				continue
			}
			if c.HeadingLevel > currentLevel {
				c.HeadingChild = true
			}
			currentLevel = c.HeadingLevel
		}
	}
}

// WithResolvedComments resolves heading comments and cleans up the page if
// COMMENT headings are encountered
func WithResolvedComments() internals.PageOption {
	return func(page *internals.Page) {
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
