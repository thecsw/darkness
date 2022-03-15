package emilia

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	RFC_EMILY = "Mon, 02 Jan 2006"
)

var (
	HEregex          = regexp.MustCompile(`(\d+);\s*(\d+)\s*H.E.`)
	HEParagraphRegex = regexp.MustCompile(`>(\d+);\s*(\d+)\s*H.E.`)
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

func AddHolosceneTitles(data string) string {
	// Match all paragraphs with holoscene time
	matches := HEParagraphRegex.FindAllStringSubmatch(data, -1)
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
