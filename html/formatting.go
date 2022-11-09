package html

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
)

// quotesReplace is the map to replace
var quotesReplace = map[string]string{
	"'s":  "’s",
	"'m":  "’m",
	"n't": "n’t",
	"'re": "’re",
	"'ll": "’ll",
	"``":  "“",
	"''":  "”",
	"--":  "—",
}

// fancyQuotes replaces boring single and double quotes with fancier Unicode versions
func fancyQuotes(text string) string {
	for k, v := range quotesReplace {
		text = strings.ReplaceAll(text, k, v)
	}
	return text
}

// markupHTMLMapping maps the regex markup to html replacements
var markupHTMLMapping = map[*regexp.Regexp]string{
	internals.ItalicText:        `$1<em>$2</em>$3`,
	internals.BoldText:          `$1<strong>$2</strong>$3`,
	internals.VerbatimText:      `$1<code>$2</code>$3`,
	internals.StrikethroughText: `$1<s>$2</s>$3`,
	internals.UnderlineText:     `$1<u>$2</u>$3`,
}

// markupHTML replaces the markup regexes defined in internal with HTML tags
func markupHTML(text string) string {
	// To make bold italics, it has to be wrapped in /*...*/
	// instead of */.../*
	text = internals.BoldItalicTextBegin.ReplaceAllString(text, `$1/*`)
	text = internals.BoldItalicTextEnd.ReplaceAllString(text, `*/$1`)
	for source, replacement := range markupHTMLMapping {
		text = source.ReplaceAllString(text, replacement)
	}
	// Double pass for giggles
	for source, replacement := range markupHTMLMapping {
		text = source.ReplaceAllString(text, replacement)
	}
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

	text = internals.FootnotePostProcessingRegexp.ReplaceAllStringFunc(text, func(what string) string {
		num, _ := strconv.Atoi(strings.ReplaceAll(what, "!", ""))
		// get the footnote HTML body
		footnote := fmt.Sprintf(
			`<a id="_footnoteref_%d" class="footnote" href="#_footnotedef_%d" title="View footnote.">%s</a>`,
			num, num, footnoteLabel(num))
		// Decide if we need to wrap the footnote in square brackets
		if emilia.Config.Website.FootnoteBrackets {
			footnote = "[" + footnote + "]"
		}
		return `
<sup class="footnote">` + footnote + `</sup>
`
	})

	return strings.TrimSpace(text)
}

// processTitle returns a properly formatted HTML of a title
func processTitle(title string) string {
	title = fancyQuotes(title)
	title = markupHTML(title)
	title = internals.MathRegexp.ReplaceAllString(title, `\($1\)`)
	return title
}

// flattenFormatting returns a plain-text to be fit into the description
func flattenFormatting(what string) string {
	what = fancyQuotes(what)
	// To make bold italics, it has to be wrapped in /*...*/
	// instead of */.../*
	what = internals.BoldItalicTextBegin.ReplaceAllString(what, `$1/*`)
	what = internals.BoldItalicTextEnd.ReplaceAllString(what, `*/$1`)
	for source := range markupHTMLMapping {
		what = source.ReplaceAllString(what, `$1$2$3`)
	}
	what = internals.KeyboardRegexp.ReplaceAllString(what, `$1`)
	what = internals.NewLineRegexp.ReplaceAllString(what, `$1`)
	return what
}
