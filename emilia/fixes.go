package emilia

import "darkness/internals"

func EnrichHeadings(page *internals.Page) {
	minHeadingLevel := 999
	// Find the smallest heanding
	for i := range page.Contents {
		c := &page.Contents[i]
		if c.Type != internals.TypeHeader {
			continue
		}
		if c.HeaderLevel < minHeadingLevel {
			minHeadingLevel = c.HeaderLevel
		}
	}
	// Shift everything over
	for i := range page.Contents {
		c := &page.Contents[i]
		if c.Type != internals.TypeHeader {
			continue
		}
		c.HeaderLevel -= (minHeadingLevel - 2)
	}
	// Mark headings that are children
	currentLevel := 0
	for i := range page.Contents {
		c := &page.Contents[i]
		if c.Type != internals.TypeHeader {
			continue
		}
		if c.HeaderLevel > currentLevel {
			c.HeaderChild = true
		}
		currentLevel = c.HeaderLevel
	}
}
