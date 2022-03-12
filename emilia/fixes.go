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
	// Mark the first heading
	for i := range page.Contents {
		c := &page.Contents[i]
		if c.Type == internals.TypeHeader {
			c.HeaderFirst = true
			break
		}
	}
	// Mark the last heading
	for i := len(page.Contents) - 1; i >= 0; i-- {
		c := &page.Contents[i]
		if c.Type == internals.TypeHeader {
			c.HeaderLast = true
			break
		}
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
