package html

import (
	"slices"
	"strings"

	"github.com/thecsw/darkness/v3/yunyun"
)

type divType uint8

const (
	// This element should be wrapped in "writing" div.
	divWriting divType = 1 + iota
	// This element should just be in the root scope.
	divOutside
	// Not meant to be used outside of processing.
	divSpecial
)

// TODO: To whoever came up with this---me---is there a better way?
// TODO: That was of course me, Sandy. Hi. This is ugly as it has to
// connect directly with yunyun/flags.go. Travesty, but whatever.
var divTypes = []divType{
	divWriting, // yunyun.TypeHeading
	divWriting, // yunyun.TypeParagraph
	divSpecial, // yunyun.TypeList
	divWriting, // yunyun.TypeListNumbered
	divSpecial, // yunyun.TypeLink
	divOutside, // yunyun.TypeSourceCode
	divOutside, // yunyun.TypeRawHtml
	divOutside, // yunyun.TypeHorizontalLine
	divWriting, // yunyun.TypeAttentionText
	divOutside, // yunyun.TypeTable
	divWriting, // yunyun.TypeDetails
	divWriting, // yunyun.TypeTableOfContents (since it's just a list).
}

func whatDivType(content *yunyun.Content) divType {
	dt := divTypes[int(content.Type)]
	if dt != divSpecial {
		return dt
	}
	// If the list has the gallery flag on, do not wrap it writing.
	if content.IsList() {
		if content.IsGallery() {
			return divOutside
		}
		return divWriting
	}
	// If the link was not an embed, wrap it in writing.
	if content.IsLink() {
		if yunyun.HasFlag(&content.Options, linkWasNotSpecialFlag) {
			return divWriting
		}
		return divOutside
	}
	// default to writing div
	return divWriting
}

func filterByLatestMetaName(heads []string) []string {
	// The deduped list from the tail.
	res := make([]string, 0, len(heads))

	// Marking the meta names we had already seen.
	seen := make(map[string]struct{})

	// Going through it in reverse.
	for i := len(heads) - 1; i >= 0; i-- {
		name, isApplicable := extractMetaName(heads[i])
		// If name couldn't be extracted, stay safe.
		if len(name) < 1 || !isApplicable {
			res = append(res, heads[i])
			continue
		}
		// If seen, then skip.
		if _, ok := seen[name]; ok {
			continue
		}
		// Mark and add.
		seen[name] = struct{}{}
		res = append(res, heads[i])
	}
	// Original order, since if the user say applied stylesheet.css and override.css,
	// if we don't preserve the original order, the override.css would end up doing
	// nothing. Found this the hard way.
	slices.Reverse(res)
	return res
}

// extractMetaName returns the name and flag if this method is even applicable.
func extractMetaName(head string) (string, bool) {
	if !strings.HasPrefix(head, "<meta") {
		return "", false
	}
	for split := range strings.SplitSeq(head, " ") {
		// Poor man's pattern-matching.
		if !strings.HasPrefix(split, `name="`) {
			continue
		}
		// Get the nice value out of it.
		return strings.Trim(split[5:], `"`), true
	}
	return "", false
}
