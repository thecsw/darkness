package html

import (
	"unicode"

	"github.com/thecsw/darkness/yunyun"
)

// extractID returns a properly formatted ID for a heading title
func extractID(heading string) string {
	// Check if heading is a link
	match := yunyun.LinkRegexp.FindStringSubmatch(heading)
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

func (e ExporterHTML) whatDivType(content *yunyun.Content) divType {
	dt := divTypes[int(content.Type)]
	if dt != divSpecial {
		return dt
	}
	// If the list has the gallery flag on, do not wrap it writing.
	if content.IsList() {
		if yunyun.HasFlag(&content.Options, yunyun.InGalleryFlag) {
			return divOutside
		}
		return divWriting
	}
	// If the link was not an embed, wrap it in writing.
	if content.IsLink() {
		if e.linkWasNotEmbed {
			return divWriting
		}
		return divOutside
	}
	// default to writing div
	return divWriting
}
