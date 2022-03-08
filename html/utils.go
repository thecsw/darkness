package html

import (
	"darkness/internals"
	"html"
)

func htmlize(text string) string {
	text = html.EscapeString(text)
	text = internals.LinkRegexp.ReplaceAllString(text, `<a href="$1">$2</a>`)
	text = internals.BoldText.ReplaceAllString(text, ` <strong>$2</strong>$3`)
	text = internals.ItalicText.ReplaceAllString(text, ` <em>$2</em>$3`)
	return text
}
