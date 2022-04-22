package html

import (
	"html"
	"strings"
	"unicode"

	"github.com/thecsw/darkness/internals"
)

// fancyQuotes replaces boring single and double quotes with fancier Unicode versions
func fancyQuotes(text string) string {
	text = strings.ReplaceAll(text, "'s", "’s")
	text = strings.ReplaceAll(text, "n't", "n’t")
	text = strings.ReplaceAll(text, "'re", "’re")
	text = strings.ReplaceAll(text, "'ll", "’ll")
	//text = strings.ReplaceAll(text, "`", "‘")
	text = strings.ReplaceAll(text, "``", "“")
	text = strings.ReplaceAll(text, "''", "”")
	text = strings.ReplaceAll(text, "--", "—")
	return text
}

// markupHTML replaces the markup regexes defined in internal with HTML tags
func markupHTML(text string) string {
	// To make bold itolics, it has to be wrapped in /*...*/
	// instead of */.../*
	text = internals.BoldItalicTextBegin.ReplaceAllString(text, `$1/*`)
	text = internals.BoldItalicTextEnd.ReplaceAllString(text, `*/$1`)
	text = internals.ItalicText.ReplaceAllString(text, `$1<em>$2</em>$3`)
	text = internals.BoldText.ReplaceAllString(text, `$1<strong>$2</strong>$3`)
	text = internals.VerbatimText.ReplaceAllString(text, `$1<code>$2</code>$3`)
	text = internals.KeyboardRegexp.ReplaceAllString(text, `<kbd>$1</kbd>`)
	text = internals.NewLineRegexp.ReplaceAllString(text, `$1<br>`)
	return text
}

// processText returns a properly formatted HTML of a text
func processText(text string) string {
	text = html.EscapeString(fancyQuotes(text))
	text = markupHTML(text)
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
	title = fancyQuotes(title)
	title = markupHTML(title)
	title = internals.MathRegexp.ReplaceAllString(title, `\($1\)`)
	return title
}

// processDescription returns a plain-text to be fit into the description
func processDescription(desc string) string {
	desc = fancyQuotes(desc)
	// To make bold itolics, it has to be wrapped in /*...*/
	// instead of */.../*
	desc = internals.BoldItalicTextBegin.ReplaceAllString(desc, `$1/*`)
	desc = internals.BoldItalicTextEnd.ReplaceAllString(desc, `*/$1`)
	desc = internals.ItalicText.ReplaceAllString(desc, `$1$2$3`)
	desc = internals.BoldText.ReplaceAllString(desc, `$1$2$3`)
	desc = internals.VerbatimText.ReplaceAllString(desc, `$1$2$3`)
	desc = internals.KeyboardRegexp.ReplaceAllString(desc, `$1`)
	desc = internals.NewLineRegexp.ReplaceAllString(desc, `$1`)
	return desc
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
