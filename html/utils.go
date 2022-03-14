package html

import (
	"html"
	"strings"

	"github.com/thecsw/darkness/internals"
)

func processText(text string) string {
	text = html.EscapeString(text)
	text = internals.ItalicText.ReplaceAllString(text, `$1<em>$2</em>$3`)
	text = internals.BoldText.ReplaceAllString(text, `$1<strong>$2</strong>$3`)
	text = internals.VerbatimText.ReplaceAllString(text, `$1<code>$2</code>$3`)
	text = internals.KeyboardRegexp.ReplaceAllString(text, `<kbd>$1</kbd>`)
	text = strings.ReplaceAll(text, "◼", `<b style="color:#ba3925">◼︎</b>`)

	text = internals.LinkRegexp.ReplaceAllString(text, `<a href="$1">$2</a>`)

	text = internals.FootnotePostProcessingRegexp.ReplaceAllString(text, `
<sup class="footnote">[<a id="_footnoteref_$1" class="footnote" href="#_footnotedef_$1" title="View footnote.">$1</a>]</sup>
`)

	return text
}

func processTitle(title string) string {
	title = strings.ReplaceAll(title, "'s", "’s")
	return title
}
