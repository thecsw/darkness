package emilia

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

const (
	descriptionMinLength = 14
)

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
		if paragraph == "" || HEregex.MatchString(paragraph) {
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
