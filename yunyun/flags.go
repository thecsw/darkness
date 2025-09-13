package yunyun

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
	// TypeRawHtml is the type of raw HTML block
	TypeRawHtml
	// TypeHorizontalLine is the type of horizontal line
	TypeHorizontalLine
	// TypeAttentionText is the type of attention text block
	TypeAttentionText
	// TypeTable is the type of a table
	TypeTable
	// TypeDetails is the type for html details
	TypeDetails
	// TypeTableOfContents is the type that splashes links to headings.
	TypeTableOfContents
	// TypeShouldBeLastDoNotTouch the last type that should not be touched --
	// It's used to verify consistency within darkness.
	TypeShouldBeLastDoNotTouch
)

// Bits is aliased to `uint16` to store flags.
type Bits uint32

const (
	// InListFlag is used internally to mark list states.
	InListFlag Bits = 1 << iota
	// InOrderedListFlag is used to internally mark ordered list states.
	InOrderedListFlag
	// InTableFlag is used internally to mark table states.
	InTableFlag
	// InTableHasHeadersFlag is used internally to mark table
	// delimiter states.
	InTableHasHeadersFlag
	// InSourceCodeFlag is used internally to mark source code
	// states.
	InSourceCodeFlag
	// InRawHtmlFlag is used internally to mark raw html states.
	InRawHtmlFlag
	// InRawHtmlFlagUnsafe is used internally to mark unsafe html states.
	InRawHtmlFlagUnsafe
	// InRawHtmlFlagResponsive is used internally to mark responsive html states.
	InRawHtmlFlagResponsive
	// InQuoteFlag is used internally to mark quote states.
	InQuoteFlag
	// InCenterFlag is used internally to mark center states.
	InCenterFlag
	// InDetailsFlag is used internally to mark details states.
	InDetailsFlag
	// InDropCapFlag is used internally to make drop cap states.
	InDropCapFlag
	// InGalleryFlag is used internally to mark gallery states.
	InGalleryFlag
	// Not description flag will mark to the exporter that if this paragaph
	// is the first content on the page, it should not be used for descriptions.
	NotADescriptionFlag
	// YunYunStartCustomFlags is used internally to mark last flag.
	YunYunStartCustomFlags
)

// LatchFlags returns four functions: add, remove, flip, and has
// for the flag container.
//
//go:inline
func LatchFlags(v *Bits) (
	func(f Bits), func(f Bits), func(f Bits), func(f Bits) bool,
) {
	return func(f Bits) { AddFlag(v, f) },
		func(f Bits) { RemoveFlag(v, f) },
		func(f Bits) { FlipFlag(v, f) },
		func(f Bits) bool { return HasFlag(v, f) }
}

// AddFlag marks the flag in the first argument.
//
//go:inline
func AddFlag(v *Bits, f Bits) { *v |= f }

// RemoveFlag removes the flag in the first argument.
//
//go:inline
func RemoveFlag(v *Bits, f Bits) { *v &^= f }

// FlipFlag flips the flag in the first argument.
//
//go:inline
func FlipFlag(v *Bits, f Bits) { *v ^= f }

// HasFlag returns true if the flag is found, false otherwise.
//
//go:inline
func HasFlag(v *Bits, f Bits) bool {
	if v == nil {
		return false
	}
	return *v&f != 0
}
