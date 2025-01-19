package yunyun

import (
	"regexp"
	"strings"
)

// Markings is used to store the regex patterns for
// our formatting, like bold, italics, etc.
type Markings struct {
	// To prevent unkeyed literars.
	_ struct{}
	// Bold defaults to orgmode's `*` (need to escape).
	Bold string
	// Italic defaults to orgmode's `/`.
	Italic string
	// */bold italic/*
	BoldItalic string
	// /*italic bold*/
	ItalicBold string
	// Verbatim defaults to orgmode's `[~=]`.
	Verbatim string
	// Strikethrough defaults to orgmode's `+` (need to escape).
	Strikethrough string
	// Underline defaults to orgmode's `_`.
	Underline string
	// Link is the whole regexp pattern.
	Link string
	// SuperscriptStart is the pattern to start superscript.
	SuperscriptStart string
	// SuperscriptEnd is the pattern to end superscript.
	SuperscriptEnd string
	// SubscriptStart is the pattern to start subscript.
	SubscriptStart string
	// SubscriptEnd is the pattern to end subscript.
	SubscriptEnd string
}

var (
	// ActiveMarkings can be set by parser to define custom markings.
	ActiveMarkings = Markings{
		Bold:             `\*`,
		Italic:           `/`,
		BoldItalic:       `\*/`,
		ItalicBold:       `/\*`,
		Verbatim:         `[~=]`,
		Strikethrough:    `\+`,
		Underline:        `_`,
		Link:             `(?mU)\[\[(?P<link>[^][]+)\](?:\[(?P<text>[^][]+)(?: "(?P<desc>[^"]+)")?\])?\]`,
		SuperscriptStart: `\^\{\{`,
		SuperscriptEnd:   `\}\}`,
		SubscriptStart:   `_\{\{`,
		SubscriptEnd:     `\}\}`,
	}

	// Pre-computed group indexes to use for group extraction.
	linkLinkIndex, linkTextIndex, linkDescIndex = -1, -1, -1
)

// BuildRegex uses patterns from `ActiveMarkings` to build Yunyun's regexes.
func (m Markings) BuildRegex() {
	BoldText = SymmetricEmphasis(m.Bold)
	ItalicText = SymmetricEmphasis(m.Italic)
	BoldItalicText = AsymmetricEmphasis(m.BoldItalic, m.ItalicBold)
	ItalicBoldText = AsymmetricEmphasis(m.ItalicBold, m.BoldItalic)
	VerbatimText = SymmetricEmphasis(m.Verbatim)
	StrikethroughText = SymmetricEmphasis(m.Strikethrough)
	UnderlineText = SymmetricEmphasis(m.Underline)
	SuperscriptText = AsymmetricEmphasis(m.SuperscriptStart, m.SuperscriptEnd)
	SubscriptText = AsymmetricEmphasis(m.SubscriptStart, m.SubscriptEnd)

	// Compile the link regexp and pre-compute named groups' indices.
	LinkRegexp = regexp.MustCompile(m.Link)
	linkLinkIndex = LinkRegexp.SubexpIndex("link")
	linkTextIndex = LinkRegexp.SubexpIndex("text")
	linkDescIndex = LinkRegexp.SubexpIndex("desc")

	SpecialTextMarkups = []*regexp.Regexp{
		BoldItalicText, ItalicBoldText, BoldText, ItalicText,
		VerbatimText, StrikethroughText, UnderlineText,
		SuperscriptText, SubscriptText,
	}
}

// ExtractedLink represents the link that was extracted with
// `ExtractLink` or `ExtractLinks`.
type ExtractedLink struct {
	Link        string
	Text        string
	Description string
	MatchLength int
}

// ExtractLinks uses `linkRegexp` and returns an array of
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
		// Commented out because if there is no text provided, just the link,
		// then the author might want the image to appear clean without text.
		// if len(extractedLinks[i].Text) < 1 {
		// 	extractedLinks[i].Text = extractedLinks[i].Link
		// }
		if len(extractedLinks[i].Description) < 1 {
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
	// UrlRegexp is yoinked from https://ihateregex.io/expr/url/
	UrlRegexp = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*))`)
	// BoldText is the regexp for matching bold text.
	BoldText *regexp.Regexp
	// ItalicText is the regexp for matching italic text.
	ItalicText *regexp.Regexp
	// matches "*/bold italic/*"
	BoldItalicText *regexp.Regexp
	// matches "/*italic bold*/"
	ItalicBoldText *regexp.Regexp
	// VerbatimText is the regexp for matching verbatim text.
	VerbatimText *regexp.Regexp
	// StrikethroughText is the regexp for matching strikethrough text.
	StrikethroughText *regexp.Regexp
	// UnderlineText is the regexp for matching underline text.
	UnderlineText *regexp.Regexp
	// SuperscriptText is the regex for matching superscript text.
	SuperscriptText *regexp.Regexp
	// SubscriptText is the regexp for matching subscript text.
	SubscriptText *regexp.Regexp
	// SpecialTextMarkups simply combines some common formatting options above.
	SpecialTextMarkups []*regexp.Regexp
	// KeyboardRegexp is the regexp for matching keyboard text.
	KeyboardRegexp = regexp.MustCompile(`kbd:\[([^][]+)\]`)
	// MathRegexp is the regexp for matching math text, it's handled in a bit of a
	// special way compared to regular markups, because we are capturing group between
	// dollar signs and modifying it. Maybe I'll merge it later with the rest of the
	// markups, but not planning on it at the moment.
	MathRegexp = regexp.MustCompile(`(?mU)` +
		darknessPunctLeftWithSpace +
		`\$(?P<text>.+)\$` +
		darknessPunctRightWithSpace)
	// ImageExtRegexp is the regexp for matching images (png, gif, jpg, jpeg, svg, webp).
	ImageExtRegexp = regexp.MustCompile(`\.(png|gif|jpg|jpeg|svg|webp)$`)
	// AudioFileExtRegexp is the regexp for matching audio (mp3, flac, midi).
	AudioFileExtRegexp = regexp.MustCompile(`\.(mp3|flac|midi)$`)
	// VideoFileExtRegexp matches commonly used video file formats.
	VideoFileExtRegexp = regexp.MustCompile(`\.(mp4|mkv|mov|flv|webm)$`)
	// NewLineRegexp matches a new line for non-math environments.
	NewLineRegexp = regexp.MustCompile(`(?mU)([^\\ ])(?:[ ]|^)?(?:[\\])(?:[ ]|$)`)
	// FootnoteRegexp is the regexp for matching footnotes.
	FootnoteRegexp = regexp.MustCompile(`(?mU)\[fn:: (.+)\]([:;!?\t\n. ]|$)`)
	// FootnotePostProcessingRegexp is the regexp for matching footnotes references.
	FootnotePostProcessingRegexp = regexp.MustCompile(`!(\d+)!`)
)

// RemoveFormatting will remove all special markup symbols.
func RemoveFormatting(what string) string {
	for _, source := range SpecialTextMarkups {
		what = source.ReplaceAllString(what, `$l$text$r`)
	}
	// only show the text and not the link
	what = LinkRegexp.ReplaceAllString(what, `$text`)
	what = KeyboardRegexp.ReplaceAllString(what, `$1`)
	what = NewLineRegexp.ReplaceAllString(what, `$1`)
	// don't even show the footnotes
	what = FootnoteRegexp.ReplaceAllString(what, ` `)
	return strings.TrimSpace(what)
}

const (
	// darknessPunctLeft is our alternative to [[:punct:]] re2
	// class for matching left punctuation symbols.
	darknessPunctLeft = `(?:[()\[\]_%“”—–-]|[>])`
	// darknessPunctRight is our alternative to [[:punct:]] re2
	// class for matching right punctuation symbols.
	darknessPunctRight = `(?:[()\[\],.!?:;&_%“”’—–-]|[<])`
	// darknessPunctLeftWithSpace also includes spaces and left capturing group.
	darknessPunctLeftWithSpace = `(?P<l>[[:space:]]|` + darknessPunctLeft + `|^)`
	// darknessPunctRightWithSpace also includes spaces and right capturing group.
	darknessPunctRightWithSpace = `(?P<r>[[:space:]]|` + darknessPunctRight + `|$)`
)

// SymmetricEmphasis is a useful tool to create simple text markups.
func SymmetricEmphasis(delimeter string) *regexp.Regexp {
	return AsymmetricEmphasis(delimeter, delimeter)
}

// AsymmetricEmphasis returns regexp with asymmetric borders.
func AsymmetricEmphasis(left, right string) *regexp.Regexp {
	return regexp.MustCompile(emphasisPattern(left, right))
}

// emphasisPattern returns pattern given left and right delimeters.
func emphasisPattern(left, right string) string {
	return `(?mU)` +
		darknessPunctLeftWithSpace +
		`(?:` + left + `)` +
		`(?P<text>\S|\S\S|\S.+\S)` +
		`(?:` + right + `)` +
		darknessPunctRightWithSpace
}
