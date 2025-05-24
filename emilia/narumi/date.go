package narumi

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/thecsw/darkness/v3/yunyun"
)

var (
	// Some emojis are compound, like lime, so they don't fit in a single rune.
	randomDateEmojis = []string{
		"🍓", "🍒", "🍋", "🍋‍🟩", "🍸", "🥧", "🍊", "☕️", "🥧",
		"🍑", "🥑", "🍐", "🥥", "🍈", "🫐", "🪵", "🍌", "🍉",
		"🍤", "🍇", "🥞", "🥗", "🍯", "🥐", "🥭", "🍙", "🧀",
	}
)

// secureRandIntn returns a cryptographically secure random integer in the range [0, n).
// It panics if there's an error reading from the random source.
func secureRandIntn(n int) int {
	if n <= 0 {
		panic("invalid argument to secureRandIntn")
	}
	
	// For small values of n, we can simply use a single byte
	if n <= 256 {
		buf := make([]byte, 1)
		_, err := rand.Read(buf)
		if err != nil {
			panic(fmt.Sprintf("failed to generate random number: %v", err))
		}
		return int(buf[0]) % n
	}
	
	// For larger values, use 4 bytes
	buf := make([]byte, 4)
	_, err := rand.Read(buf)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random number: %v", err))
	}
	
	// Use a more explicit approach to avoid the integer overflow
	// Converting n to uint32 first to safely perform modulo
	randomInt := binary.BigEndian.Uint32(buf)
	// #nosec G115 - This is safe because n is positive and already verified to be > 0
	result := randomInt % uint32(n)
	return int(result) // Safe conversion from uint32 to int since result < n
}

// WithDate is a PageOption that adds the date to the page.
func WithDate() yunyun.PageOption {
	return func(page *yunyun.Page) {
		// The user explicitly opted out.
		if page.Accoutrement.Date.IsDisabled() {
			return
		}
		// Nonsense dates.
		if len(page.Date) < 1 || page.Date == "0; 0 H.E." {
			return
		}
		// The user manually put the date in.
		if len(page.Contents) > 0 &&
			page.Contents[0].Type == yunyun.TypeParagraph &&
			strings.TrimSpace(page.Contents[0].Paragraph) == strings.TrimSpace(page.Date) {
			return

		}
		dateContents := make(yunyun.Contents, 1)
		regular, isHoloscene := ConvertHoloscene(page.Date)
		dateString := strings.TrimSpace(page.Date)
		if isHoloscene {
			dateString = fmt.Sprintf(`%s At least %s ago`,
				randomDateEmojis[secureRandIntn(len(randomDateEmojis))],
				formatSince(time.Since(regular)))
		}
		dateContents[0] = &yunyun.Content{
			CustomHtmlTags: fmt.Sprintf(`id="date-section" title="%s"`,
				strings.TrimSpace(regular.Format(RfcEmily))),
			Paragraph: dateString,
			Type:      yunyun.TypeParagraph,
			Options:   yunyun.NotADescriptionFlag,
		}
		page.Contents = append(dateContents, page.Contents...)
	}
}

// formatSince formats a time.Duration into a human-readable string.
func formatSince(since time.Duration) string {
	months, days := sinceToMonthsAndDays(since)
	years := months / 12
	months = months % 12
	if years == 0 && months == 0 && days == 0 {
		return "today"
	}
	sb := strings.Builder{}

	if years > 0 {
		sb.WriteString(fmt.Sprintf("%d year", years))
		if years > 1 {
			sb.WriteString("s")
		}
	}
	if months > 0 {
		if years > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d month", months))
		if months > 1 {
			sb.WriteString("s")
		}
	}
	if days > 0 {
		// We will use the Oxford comma if we have both years and months.
		if years > 0 && months > 0 {
			sb.WriteString(", and ")
		} else if years > 0 || months > 0 {
			// But, if we only have one of them, we will use "and".
			sb.WriteString(" and ")
		}
		sb.WriteString(fmt.Sprintf("%d day", days))
		if days > 1 {
			sb.WriteString("s")
		}
	}
	return sb.String()
}

// sinceToMonthsAndDays converts a time.Duration to months and days.
func sinceToMonthsAndDays(since time.Duration) (months, days int64) {
	// 30 days in a month.
	days = int64(since.Hours() / 24)
	months = days / 30
	days = days % 30
	return
}
