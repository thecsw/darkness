package internals

import "regexp"

var (
	// LinkRegexp is the regexp for matching links.
	LinkRegexp = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	// URLRegexp is yoinked from https://ihateregex.io/expr/url/
	URLRegexp = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*))`)
	// BoldText is the regexp for matching bold text.
	BoldText = markupRegex(`\*`)
	// ItalicText is the regexp for matching italic text.
	ItalicText = markupRegex(`/`)
	// BoldItalicText is the regexp for matching bold-italic text from the left.
	BoldItalicTextBegin = regexp.MustCompile(`(?mU)(^|[ ()_%<>])\*/`)
	// BoldItalicTextEnd is the regexp for matching bold-italic text from the right.
	BoldItalicTextEnd = regexp.MustCompile(`(?mU)/\*($|[ (),.!?;&_%><])`)
	// VerbatimText is the regexp for matching verbatim text.
	VerbatimText = markupRegex(`[~=]`)
	// StrikethroughText is the regexp for matching strikethrough text.
	StrikethroughText = markupRegex(`\+`)
	// UnderlineText is the regexp for matching underline text.
	UnderlineText = markupRegex(`_`)
	// SpecialTextMarkups simply combines some common formatting options above.
	SpecialTextMarkups = []*regexp.Regexp{
		BoldText, ItalicText, VerbatimText, StrikethroughText, UnderlineText,
	}
	// KeyboardRegexp is the regexp for matching keyboard text.
	KeyboardRegexp = regexp.MustCompile(`kbd:\[([^][]+)\]`)
	// MathRegexp is the regexp for matching math text.
	MathRegexp = regexp.MustCompile(`(?mU)\$(.+)\$`)
	// ImageRegexp is the regexp for matching images (png, gif, jpg, jpeg, svg, webp).
	ImageExtRegexp = regexp.MustCompile(`\.(png|gif|jpg|jpeg|svg|webp)$`)
	// AudioRegexp is the regexp for matching audio (mp3, flac, midi).
	AudioFileExtRegexp = regexp.MustCompile(`\.(mp3|flac|midi)$`)
	// VideoFileExtRegexp matches commonly used video file formats.
	VideoFileExtRegexp = regexp.MustCompile(`\.(mp4|mkv|mov|flv|webm)$`)
	// NewLineRegexp matches a new line for non-math environments.
	NewLineRegexp = regexp.MustCompile(`(?mU)([^\\])\\([ ]|$)`)
	// FootnoteRegexp is the regexp for matching footnotes.
	FootnoteRegexp = regexp.MustCompile(`(?mU)\[fn:: (.+)\]([:;!?\t\n. ]|$)`)
	// FootnoteReferenceRegexp is the regexp for matching footnotes references.
	FootnotePostProcessingRegexp = regexp.MustCompile(`!(\d+)!`)
)

// FixBoldItalicMarkups fixes bold italic text sources, such that all
// occurences of \*/.../\*  are switched to /\*...\*/
func FixBoldItalicMarkups(input string) string {
	return BoldItalicTextEnd.ReplaceAllString(BoldItalicTextBegin.ReplaceAllString(input, `$1/*`), `*/$1`)
}

// RemoveFormatting will remove all special markup symbols.
func RemoveFormatting(what string) string {
	what = FixBoldItalicMarkups(what)
	for _, source := range SpecialTextMarkups {
		what = source.ReplaceAllString(what, `$1$2$3`)
	}
	what = KeyboardRegexp.ReplaceAllString(what, `$1`)
	what = NewLineRegexp.ReplaceAllString(what, `$1`)
	return what
}

// markupRegex is a useful tool to create simple text markups.
func markupRegex(delimeter string) *regexp.Regexp {
	return regexp.MustCompile(
		`(?mU)(^|[ ()\[\]_%>])` + delimeter +
			`(\S|\S\S|\S.+\S)` + delimeter +
			`($|[ ()\[\],.!?:;&_%<“”])`)
}
