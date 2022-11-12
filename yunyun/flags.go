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

// Bits is aliased to `uint16` to store flags.
type Bits uint16

const (
	// InListFlag is used internally to mark list states.
	InListFlag Bits = 1 << iota
	// InTableFlag is used internally to mark table states.
	InTableFlag
	// InTableHasHeadersFlag is used internally to mark table
	// delimiter states.
	InTableHasHeadersFlag
	// InSourceCodeFlag is used internally to mark source code
	// states.
	InSourceCodeFlag
	// InRawHTMLFlag is used internally to mark raw html states.
	InRawHTMLFlag
	// InRawHtmlFlagUnsafe is used internally to mark unsafe html states.
	InRawHtmlFlagUnsafe
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
)

var (
	// LatchFlags returns four functions: add, remove, flip, and has
	// for the flag container.
	LatchFlags = func(v *Bits) (
		func(f Bits), func(f Bits), func(f Bits), func(f Bits) bool) {
		return func(f Bits) { AddFlag(v, f) },
			func(f Bits) { RemoveFlag(v, f) },
			func(f Bits) { FlipFlag(v, f) },
			func(f Bits) bool { return HasFlag(v, f) }
	}
	// AddFlag marks the flag in the first argument.
	AddFlag = func(v *Bits, f Bits) { *v |= f }
	// RemoveFlag removes the flag in the first argument.
	RemoveFlag = func(v *Bits, f Bits) { *v &^= f }
	// FlipFlag flips the flag in the first argument.
	FlipFlag = func(v *Bits, f Bits) { *v ^= f }
	// HasFlag returns true if the flag is found, false otherwise.
	HasFlag = func(v *Bits, f Bits) bool { return *v&f != 0 }
)
