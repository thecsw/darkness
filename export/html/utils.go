package html

import (
	"strings"
	"unicode"

	"github.com/thecsw/darkness/yunyun"
)

// extractID returns a properly formatted ID for a heading title
func extractID(heading string) string {
	// Check if heading is a link
	matchLen, _, title, _ := yunyun.ExtractLink(heading)
	if matchLen > 0 {
		heading = title // 0 is whole match, 1 is link, 2 is title
	}

	res := ""
	for _, c := range heading {
		if unicode.IsSpace(c) || unicode.IsPunct(c) || unicode.IsSymbol(c) {
			res += "-"
			continue
		}
		if c <= unicode.MaxASCII {
			res += string(unicode.ToLower(c))
		}
	}
	return strings.TrimRight(res, "-")
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
		divWriting, // yunyun.TypeHeading
		divWriting, // yunyun.TypeParagraph
		divSpecial, // yunyun.TypeList
		divWriting, // yunyun.TypeListNumbered
		divSpecial, // yunyun.TypeLink
		divWriting, // yunyun.TypeSourceCode
		divOutside, // yunyun.TypeRawHTML
		divOutside, // yunyun.TypeHorizontalLine
		divWriting, // yunyun.TypeAttentionText
		divOutside, // yunyun.TypeTable
		divWriting, // yunyun.TypeDetails
	}
)

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
