package internals

// TypeContent is the type of content, used for enums
type TypeContent uint8

const (
	// TypeHeading is the type of heading
	TypeHeading TypeContent = iota
	// TypeParagraph is the type of paragraph, which is just text
	TypeParagraph
	// TypeList is the type of unordered list
	TypeList
	// TypeListNumbered is the type of numbered list
	TypeListNumbered
	// TypeLink is the type of link
	TypeLink
	// TypeSourceCode is the type of source code block
	TypeSourceCode
	// TypeRawHTML is the type of raw HTML block
	TypeRawHTML
	// TypeHorizontalLine is the type of horizontal line
	TypeHorizontalLine
	// TypeAttentionText is the type of attention text block
	TypeAttentionText
	// TypeTable is the type of a table
	TypeTable
	// TypeDetails is the type for html details
	TypeDetails
	// This is the last type that should not be touched --
	// It's used to verify consistency within darkness.
	TypeShouldBeLastDoNotTouch
)

type Bits uint16

const (
	InListFlag Bits = 1 << iota
	InTableFlag
	InTableHasHeadersFlag
	InSourceCodeFlag
	InRawHTMLFlag
	InQuoteFlag
	InCenterFlag
	InDetailsFlag
	InDropCapFlag
	InGalleryFlag
)

var (
	LatchFlags = func(v *Bits) (
		func(f Bits), func(f Bits), func(f Bits), func(f Bits) bool) {
		return func(f Bits) { AddFlag(v, f) },
			func(f Bits) { RemoveFlag(v, f) },
			func(f Bits) { FlipFlag(v, f) },
			func(f Bits) bool { return HasFlag(v, f) }
	}

	AddFlag    = func(v *Bits, f Bits) { *v |= f }
	RemoveFlag = func(v *Bits, f Bits) { *v &^= f }
	FlipFlag   = func(v *Bits, f Bits) { *v ^= f }
	HasFlag    = func(v *Bits, f Bits) bool { return *v&f != 0 }
)
