package yunyun

import (
	"regexp"
)

// Markings is used to store the regex patterns for
// our formatting, like bold, italics, etc.
type Markings struct {
	// Bold defaults to orgmode's `*` (need to escape).
	Bold string
	// Italic defaults to orgmode's `/`.
	Italic string
	// Verbatim defaults to orgmode's `[~=]`.
	Verbatim string
	// Strikethrough defaults to orgmode's `+` (need to escape).
	Strikethrough string
	// Underline defaults to orgmode's `_`.
	Underline string
	// Link is the whole regexp pattern.
	Link string
}

var (
	// ActiveMarkings can be set by parser to define custom markings.
	ActiveMarkings = Markings{
		Bold:          `\*`,
		Italic:        `/`,
		Verbatim:      `[~=]`,
		Strikethrough: `\+`,
		Underline:     `_`,
		Link:          `(?mU)\[\[(?P<link>[^][]+)\]\[(?P<text>[^][]+)(?: "(?P<desc>[^"]+)")?\]\]`,
	}

	linkLinkIndex, linkTextIndex, linkDescIndex = -1, -1, -1
)

// BuildRegex uses patterns from `ActiveMarkings` to build Yunyun's regexes.
func (m Markings) BuildRegex() {
	BoldText = MarkupRegex(m.Bold)
	ItalicText = MarkupRegex(m.Italic)
	BoldItalicTextBegin = regexp.MustCompile(`(?mU)(^|[ ()_%<>])` + m.Bold + m.Italic)
	BoldItalicTextEnd = regexp.MustCompile(`(?mU)` + m.Italic + m.Bold + `($|[ (),.!?;&_%><])`)
	VerbatimText = MarkupRegex(m.Verbatim)
	StrikethroughText = MarkupRegex(m.Strikethrough)
	UnderlineText = MarkupRegex(m.Underline)

	LinkRegexp = regexp.MustCompile(m.Link)
	linkLinkIndex = LinkRegexp.SubexpIndex("link")
	linkTextIndex = LinkRegexp.SubexpIndex("text")
	linkDescIndex = LinkRegexp.SubexpIndex("desc")

	SpecialTextMarkups = []*regexp.Regexp{
		BoldText, ItalicText, VerbatimText, StrikethroughText, UnderlineText,
	}
}

// ExtractedLink represents the link that was extracted with
// `ExtractLink` or `ExtractLinks`.
type ExtractedLink struct {
	MatchLength int
	Link        string
	Text        string
	Description string
}

// ExtractLink uses `linkRegexp` and returns an array of
// extracted links, if any.
func ExtractLinks(line string) []*ExtractedLink {
	if !LinkRegexp.MatchString(line) {
		return nil
	}
	submatches := LinkRegexp.FindAllStringSubmatch(line, -1)
	// Sanity check
	if len(submatches) < 1 {
		return nil
	}
	extractedLinks := make([]*ExtractedLink, len(submatches))
	for i, submatch := range submatches {
		extractedLinks[i] = &ExtractedLink{
			MatchLength: len(submatch[0]),
			Link:        submatch[linkLinkIndex],
			Text:        submatch[linkTextIndex],
			Description: submatch[linkDescIndex],
		}
		if extractedLinks[i].Description == "" {
			extractedLinks[i].Description = extractedLinks[i].Text
		}
	}
	return extractedLinks
}

// ExtractLink uses `linkRegexp` and returns first extracted link.
func ExtractLink(line string) *ExtractedLink {
	if links := ExtractLinks(line); len(links) > 0 {
		return links[0]
	}
	return nil
}

var (
	// LinkRegexp is the regexp for matching links.
	LinkRegexp *regexp.Regexp
	// URLRegexp is yoinked from https://ihateregex.io/expr/url/
	URLRegexp = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*))`)
	// BoldText is the regexp for matching bold text.
	BoldText *regexp.Regexp
	// ItalicText is the regexp for matching italic text.
	ItalicText *regexp.Regexp
	// BoldItalicText is the regexp for matching bold-italic text from the left.
	BoldItalicTextBegin *regexp.Regexp
	// BoldItalicTextEnd is the regexp for matching bold-italic text from the right.
	BoldItalicTextEnd *regexp.Regexp
	// VerbatimText is the regexp for matching verbatim text.
	VerbatimText *regexp.Regexp
	// StrikethroughText is the regexp for matching strikethrough text.
	StrikethroughText *regexp.Regexp
	// UnderlineText is the regexp for matching underline text.
	UnderlineText *regexp.Regexp
	// SpecialTextMarkups simply combines some common formatting options above.
	SpecialTextMarkups []*regexp.Regexp
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
		what = source.ReplaceAllString(what, `$l$text$r`)
	}
	what = KeyboardRegexp.ReplaceAllString(what, `$1`)
	what = NewLineRegexp.ReplaceAllString(what, `$1`)
	return what
}

// MarkupRegex is a useful tool to create simple text markups.
func MarkupRegex(delimeter string) *regexp.Regexp {
	return regexp.MustCompile(
		`(?mU)(?P<l>^|[ ()\[\]_%>])` + delimeter +
			`(?P<text>\S|\S\S|\S.+\S)` + delimeter +
			`(?P<r>$|[ ()\[\],.!?:;&_%<“”])`)
}
