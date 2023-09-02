package html

import (
	"github.com/thecsw/darkness/yunyun"
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
