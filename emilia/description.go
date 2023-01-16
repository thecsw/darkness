package emilia

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

func GetDescription(page *yunyun.Page) string {
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
		description = yunyun.RemoveFormatting(
			paragraph[:gana.Min(len(paragraph), Config.Website.DescriptionLength)]) + "..."
		break
	}
	return description
}
