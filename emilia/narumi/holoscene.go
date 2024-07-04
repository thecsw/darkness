package narumi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/thecsw/darkness/v3/emilia/puck"
)

const (
	// RfcEmily is the RFC3339 format for the emily time format
	RfcEmily = "Mon, 02 Jan 2006"
)

// HEParagraphRegex is a regex for matching Holoscene times in paragraphs
var HEParagraphRegex = regexp.MustCompile(`>\s*(\d+);\s*(\d+)\s*H.E.`)

// ConvertHoloscene takes a Holoscene time (127; 12022 H.E.) to a time struct.
func ConvertHoloscene(HEtime string) (time.Time, bool) {
	return getHoloscene(extractHoloscene(HEtime))
}

// extractHoloscene extracts the holoscene time from a string, the return
// values are day, year, hour, minute.
func extractHoloscene(data string) (string, string, string, string) {
	matches := puck.HEregex.FindAllStringSubmatch(data, 1)
	// Not a good match, nothing found
	if len(matches) < 1 {
		return "", "", "", ""
	}
	day := matches[0][puck.HEregex.SubexpIndex("day")]
	year := matches[0][puck.HEregex.SubexpIndex("year")]
	hourMinute := matches[0][puck.HEregex.SubexpIndex("hour_minute")]
	// If there is no hour/minute, set it to 00:00
	if hourMinute == "" {
		hourMinute = "0000"
	}
	hour := hourMinute[:2]
	minute := hourMinute[2:]
	return day, year, hour, minute
}

// getHoloscene returns a time struct for a given holoscene time, second
// return value is true if the time is valid.
func getHoloscene(dayS, yearS, hourS, minuteS string) (time.Time, bool) {
	day, year, hour, minute := a0(dayS), a0(yearS), a0(hourS), a0(minuteS)
	// Check if the time is valid
	if day == 0 && year == 0 && hour == 0 && minute == 0 {
		return time.Time{}, false
	}
	// Subtract the 10k holoscene years
	if year > 10000 {
		year -= 10000
	}
	tt := time.Date(year, time.January, day, hour, minute, 0, 0, time.Local)
	return tt, true
}

// AddHolosceneTitles adds the titles of the Holoscene to the page, and
// also how many to replace for the page, -1 for everything.
func AddHolosceneTitles(data string, num int) string {
	// Match all paragraphs with holoscene time
	matches := HEParagraphRegex.FindAllStringSubmatch(data, num)
	for _, match := range matches {
		tt, _ := getHoloscene(match[1], match[2], "", "")
		// Add the title to the paragraph
		data = strings.Replace(data,
			match[0],
			fmt.Sprintf(` title="%s"%s`, tt.Format(RfcEmily), match[0]),
			len(matches),
		)
	}
	return data
}

// a0 is a wrapper for atoiButDefault that returns 0.
func a0(s string) int {
	return atoiButDefault(s, 0)
}

// atoiButDefault is a wrapper for strconv.Atoi that returns a default value.
func atoiButDefault(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}
