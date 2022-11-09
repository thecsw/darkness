package html

import (
	"unicode"

	"github.com/thecsw/darkness/internals"
)

// extractID returns a properly formatted ID for a heading title
func extractID(heading string) string {
	// Check if heading is a link
	match := internals.LinkRegexp.FindStringSubmatch(heading)
	if len(match) > 0 {
		heading = match[2] // 0 is whole match, 1 is link, 2 is title
	}
	res := "_"
	for _, c := range heading {
		if unicode.IsSpace(c) || unicode.IsPunct(c) || unicode.IsSymbol(c) {
			res += "_"
			continue
		}
		if c <= unicode.MaxASCII {
			res += string(unicode.ToLower(c))
		}
	}
	return res
}

type divType uint8

const (
	// This element should be wrapped in "writing" div.
	divWriting divType = 1 + iota
	// This element should just be in the root scope.
	divOutside
	// Not meant to be used outside of processing.
	divSpecial
)

var (
	divTypes = []divType{
		divWriting, // internals.TypeHeading
		divWriting, // internals.TypeParagraph
		divSpecial, // internals.TypeList
		divWriting, // internals.TypeListNumbered
		divOutside, // internals.TypeLink
		divWriting, // internals.TypeSourceCode
		divOutside, // internals.TypeRawHTML
		divOutside, // internals.TypeHorizontalLine
		divWriting, // internals.TypeAttentionText
		divWriting, // internals.TypeTable
		divWriting, // internals.TypeDetails
	}
)

func whatDivType(content *internals.Content) divType {
	dt := divTypes[int(content.Type)]
	if dt != divSpecial {
		return dt
	}
	if content.IsList() {
		if internals.HasFlag(&content.Options, internals.InGalleryFlag) {
			return divOutside
		}
		return divWriting
	}
	// default to writing div
	return divWriting
}
