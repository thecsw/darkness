package emilia

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// RFC_EMILY is the RFC3339 format for the emily time format
	RFC_EMILY = "Mon, 02 Jan 2006"
)

var (
	// HEregex is a regex for matching Holoscene times
	HEregex = regexp.MustCompile(`(\d+);\s*(\d+)\s*H.E.`)
	// HEParagraphRegex is a regex for matching Holoscene times in paragraphs
	HEParagraphRegex = regexp.MustCompile(`>\s*(\d+);\s*(\d+)\s*H.E.`)
)

// ConvertHoloscene takes a Holoscene time (127; 12022 H.E.) to a time struct.
func ConvertHoloscene(HEtime string) *time.Time {
	matches := HEregex.FindAllStringSubmatch(HEtime, 1)
	// Not a good match, nothing found
	if len(matches) < 1 {
		return nil
	}
	return getHoloscene(matches[0][1], matches[0][2])
}

// getHoloscene returns a time struct for a given holoscene time.
func getHoloscene(dayS, yearS string) *time.Time {
	// By the regex, we are guaranteed to have good numbers
	day, _ := strconv.Atoi(dayS)
	year, _ := strconv.Atoi(yearS)
	// Subtract the 10k holoscene years
	year -= 10000

	tt := time.Date(year, time.January, 0, 0, 0, 0, 0, time.Local)
	tt = tt.Add(time.Duration(day) * 24 * time.Hour)
	return &tt
}

// AddHolosceneTitles adds the titles of the Holoscene to the page, and
// also how many to replace for the page, -1 for everything.
func AddHolosceneTitles(data string, num int) string {
	// Match all paragraphs with holoscene time
	matches := HEParagraphRegex.FindAllStringSubmatch(data, num)
	for _, match := range matches {
		// Add the title to the paragraph
		data = strings.Replace(data,
			match[0],
			fmt.Sprintf(` title="%s"%s`,
				getHoloscene(match[1], match[2]).Format(RFC_EMILY), match[0]),
			len(matches),
		)
	}
	return data
}
