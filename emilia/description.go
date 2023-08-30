package emilia

import (
	"strings"

	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

const (
	// Minimum length of the description
	descriptionMinLength = 14
)

// GetDescription returns the description of the page
// It will return the first paragraph that is not empty and not a holoscene time
// If no such paragraph is found, it will return an empty string
// If the description is less than 14 characters, it will return an empty string
func GetDescription(page *yunyun.Page, length int) string {
	// Find the first paragraph for description
	description := ""
	for _, content := range page.Contents {
		// We are only looking for paragraphs
		if !content.IsParagraph() {
			continue
		}
		// Skip holoscene times
		paragraph := strings.TrimSpace(content.Paragraph)
		if paragraph == "" || puck.HEregex.MatchString(paragraph) {
			continue
		}

		cleanText := yunyun.RemoveFormatting(paragraph[:gana.Min(len(paragraph), length+10)])
		description = cleanText[:gana.Max(len(cleanText)-10, 0)] + "..."
		if len(description) < descriptionMinLength {
			continue
		}
		break
	}
	return description
}
