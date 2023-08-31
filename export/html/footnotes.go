package html

import (
	"fmt"
	"strings"

	"github.com/thecsw/darkness/emilia/narumi"
)

// addFootnotes adds the footnotes
func (e *state) addFootnotes() string {
	if len(e.page.Footnotes) < 1 {
		return ""
	}
	footnotes := make([]string, len(e.page.Footnotes))
	for i, footnote := range e.page.Footnotes {
		footnotes[i] = fmt.Sprintf(`
<div class="footnote" id="_footnotedef_%d">
<a href="#_footnoteref_%d">%s</a>
%s
</div>
`,
			i+1, i+1, narumi.FootnoteLabeler(i+1), processText(footnote))
	}
	return fmt.Sprintf(`
<div id="footnotes">
<hr>
%s
</div>
`, strings.Join(footnotes, ""))
}
