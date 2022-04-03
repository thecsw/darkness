package html

import (
	"html"
	"strings"
	"unicode"

	"github.com/thecsw/darkness/internals"
)

// processText returns a properly formatted HTML of a text
func processText(text string) string {
	// To make bold itolics, it has to be wrapped in /*...*/
	// instead of */.../*
	text = internals.BoldItalicTextBegin.ReplaceAllString(text, `$1/*`)
	text = internals.BoldItalicTextEnd.ReplaceAllString(text, `*/$1`)
	text = strings.ReplaceAll(text, "'s", "’s")
	text = strings.ReplaceAll(text, "n't", "n’t")
	text = strings.ReplaceAll(text, "'re", "’re")
	//text = strings.ReplaceAll(text, "`", "‘")
	text = strings.ReplaceAll(text, "``", "“")
	text = strings.ReplaceAll(text, "''", "”")
	text = strings.ReplaceAll(text, "--", "—")

	text = html.EscapeString(text)
	text = internals.ItalicText.ReplaceAllString(text, `$1<em>$2</em>$3`)
	text = internals.BoldText.ReplaceAllString(text, `$1<strong>$2</strong>$3`)
	text = internals.VerbatimText.ReplaceAllString(text, `$1<code>$2</code>$3`)
	text = internals.KeyboardRegexp.ReplaceAllString(text, `<kbd>$1</kbd>`)
	text = internals.NewLineRegexp.ReplaceAllString(text, `$1<br>`)
	text = strings.ReplaceAll(text, "◼", `<b style="color:#ba3925">◼︎</b>`)

	text = internals.LinkRegexp.ReplaceAllString(text, `<a href="$1">$2</a>`)

	//text = internals.MathRegexp.ReplaceAllString(text, `\($1\)`)

	text = internals.FootnotePostProcessingRegexp.ReplaceAllString(text, `
<sup class="footnote">[<a id="_footnoteref_$1" class="footnote" href="#_footnotedef_$1" title="View footnote.">$1</a>]</sup>
`)

	return strings.TrimSpace(text)
}

// processTitle returns a properly formatted HTML of a title
func processTitle(title string) string {
	title = strings.ReplaceAll(title, "'s", "’s")
	title = internals.MathRegexp.ReplaceAllString(title, `\($1\)`)
	return title
}

// extractID returns a properly formatted ID for a heading title
func extractID(heading string) string {
	// Check if heading is a link
	match := internals.LinkRegexp.FindStringSubmatch(heading)
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
