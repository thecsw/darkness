package html

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
)

// quotesReplace is the map to replace
var quotesReplace = map[string]string{
	"'s":  "’s",
	"s'":  "s’",
	"'m":  "’m",
	"n't": "n’t",
	"'re": "’re",
	"'ve": "’ve",
	"'ll": "’ll",
	"``":  "“",
	"''":  "”",
}

// fancyQuotes replaces boring single and double quotes with fancier Unicode versions
func fancyQuotes(text string) string {
	for k, v := range quotesReplace {
		text = strings.ReplaceAll(text, k, v)
	}
	text = strings.ReplaceAll(text, "---", "—") // em-dash
	text = strings.ReplaceAll(text, "--", "–")  // en dash
	return text
}

// markupHTMLMapping maps the regex markup to html replacements
var markupHTMLMapping map[*regexp.Regexp]string

// markupHTML replaces the markup regexes defined in internal with HTML tags
func markupHTML(text string) string {
	for source, replacement := range markupHTMLMapping {
		text = source.ReplaceAllString(text, replacement)
	}
	// We only need to run bold text repacement again
	text = yunyun.BoldText.ReplaceAllString(text, markupHTMLMapping[yunyun.BoldText])
	text = yunyun.KeyboardRegexp.ReplaceAllString(text, `<kbd>$1</kbd>`)
	text = yunyun.NewLineRegexp.ReplaceAllString(text, `$1<br>`)
	return text
}

// processText returns a properly formatted HTML of a text
func processText(text string) string {
	text = markupHTML(html.EscapeString(fancyQuotes(text)))
	text = strings.ReplaceAll(text, "◼", `<b style="color:#ba3925">◼︎</b>`)
	text = yunyun.LinkRegexp.ReplaceAllString(text,
		fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, `$link`, `$desc`, `$text`))
	text = yunyun.MathRegexp.ReplaceAllString(text, `\($1\)`)
	text = yunyun.FootnotePostProcessingRegexp.ReplaceAllStringFunc(text, func(what string) string {
		num, _ := strconv.Atoi(strings.ReplaceAll(what, "!", ""))
		// get the footnote HTML body
		footnote := fmt.Sprintf(
			`<a id="_footnoteref_%d" class="footnote" href="#_footnotedef_%d" title="View footnote.">%s</a>`,
			num, num, emilia.FootnoteLabeler(num))
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
	return yunyun.MathRegexp.ReplaceAllString(markupHTML(fancyQuotes(title)), `\($1\)`)
}

// flattenFormatting returns a plain-text to be fit into the description
func flattenFormatting(what string) string {
	return yunyun.RemoveFormatting(fancyQuotes(what))
}
