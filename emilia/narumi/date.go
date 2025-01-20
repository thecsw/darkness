package narumi

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/thecsw/darkness/v3/yunyun"
)

var (
	// Some emojis are compound, like lime, so they don't fit in a single rune.
	randomDateEmojis = []string{
		"🍓", "🍒", "🍋", "🍋‍🟩", "🍸", "🥧", "🍊", "☕️",
		"🍑", "🥑", "🍐", "🥥", "🍈", "🫐", "🪵", "🍌", "🍉",
		"🍤", "🍇", "🥞", "🥗", "🍯", "🥐", "🥭", "🍙", "🧀",
	}
)

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
		daysAgo := int64(-1)
		regular, isHoloscene := ConvertHoloscene(page.Date)
		dateString := strings.TrimSpace(page.Date)
		if isHoloscene {
			// I was thinking of having something like
			// > 364; 12023 H.E. 1400 - Sat, 30 Dec 2023 14:00 CST
			// but just the simple
			// > 364; 12023 H.E. 1400
			// looks so much cleaner. We already have misa holoscene
			// post-processing, which will add the actual time as a tooltip.
			daysAgo = int64(time.Since(regular).Hours() / 24)
			dateString = fmt.Sprintf(`%s %s ^{{(at least %s days ago)}}`,
				randomDateEmojis[rand.IntN(len(randomDateEmojis))],
				page.Date, humanize.Comma(daysAgo))
		}
		dateContents[0] = &yunyun.Content{
			CustomHtmlTags: `id="date-section"`,
			Paragraph:      dateString,
			Type:           yunyun.TypeParagraph,
		}
		page.Contents = append(dateContents, page.Contents...)
	}
}
