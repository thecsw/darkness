package emilia

import "github.com/thecsw/darkness/internals"

func EnrichHeadings(page *internals.Page) {
	minHeadingLevel := 999
	// Find the smallest heanding
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
