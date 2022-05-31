package internals

import "regexp"

var (
	// LinkRegexp is the regexp for matching links
	LinkRegexp = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	// URLRegexp is yoinked from https://ihateregex.io/expr/url/
	URLRegexp = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*))`)
	// BoldText is the regexp for matching bold text
	BoldText = markupRegex(`\*`)
	// ItalicText is the regexp for matching italic text
	ItalicText = markupRegex(`/`)
	// BoldItalicText is the regexp for matching bold-italic text from the left
	BoldItalicTextBegin = regexp.MustCompile(`(?mU)(^|[ ()_%<>])\*/`)
	// BoldItalicTextEnd is the regexp for matching bold-italic text from the right
	BoldItalicTextEnd = regexp.MustCompile(`(?mU)/\*($|[ (),.!?;&_%<>])`)
	// VerbatimText is the regexp for matching verbatim text
	VerbatimText = markupRegex(`[~=]`)
	// StrikethroughText is the regexp for matching strikethrough text
	StrikethroughText = markupRegex(`+`)
	// UnderlineText is the regexp for matching underline text
	UnderlineText = markupRegex(`_`)
	// KeyboardRegexp is the regexp for matching keyboard text
	KeyboardRegexp = regexp.MustCompile(`kbd:\[([^][]+)\]`)
	// MathRegexp is the regexp for matching math text
	MathRegexp = regexp.MustCompile(`(?mU)\$(.+)\$`)
	// ImageRegexp is the regexp for matching images (png, gif, jpg, jpeg, svg, webp)
	ImageExtRegexp = regexp.MustCompile(`\.(png|gif|jpg|jpeg|svg|webp)$`)
	// AudioRegexp is the regexp for matching audio (mp3, flac, midi)
	AudioFileExtRegexp = regexp.MustCompile(`\.(mp3|flac|midi)$`)
	// VideoFileExtRegexp matches commonly used video file formats
	VideoFileExtRegexp = regexp.MustCompile(`\.(mp4|mkv|mov|flv|webm)$`)

	// NewLineRegexp matches a new line for non-math environments
	NewLineRegexp = regexp.MustCompile(`(?m)([^\\])\\ `)

	// FootnoteRegexp is the regexp for matching footnotes
	FootnoteRegexp = regexp.MustCompile(`(?mU)\[fn:: (.+)\]([:;!?\t\n. ]|$)`)
	// FootnoteReferenceRegexp is the regexp for matching footnotes references
	FootnotePostProcessingRegexp = regexp.MustCompile(`!(\d+)!`)
)

func markupRegex(delimeter string) *regexp.Regexp {
	return regexp.MustCompile(
		`(?mU)(^|[ ()\[\]_%<>])` + delimeter +
			`(\S|\S\S|\S.+\S)` + delimeter +
			`($|[ ()\[\],.!?:;&_%<>“”])`)
}
