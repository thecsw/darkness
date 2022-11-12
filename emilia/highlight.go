package emilia

import (
	"fmt"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

const (
	highlightJsTheme             = `<link rel="stylesheet" href="%s">`
	highlightJsThemeDefaultPath  = `scripts/highlight/styles/agate.min.css`
	highlightJsScript            = `<script src="%s"></script>`
	highlightJsScriptDefaultPath = `scripts/highlight/highlight.min.js`
	highlightJsAction            = `<script>hljs.highlightAll();</script>`
)

func WithSyntaxHighlighting() yunyun.PageOption {
	return func(page *yunyun.Page) {
		// If Emilia disabled the syntax highlighting, don't even bother.
		if !Config.Website.SyntaxHighlighting {
			return
		}
		// Check if any content has a source code and if not, no highlighting.
		if !gana.Anyf(func(v yunyun.Content) bool { return v.IsSourceCode() }, page.Contents) {
			return
		}
		page.Stylesheets = append(page.Stylesheets,
			fmt.Sprintf(highlightJsTheme, JoinPath(Config.Website.SyntaxHighlightingTheme)))
		page.Scripts = append(page.Scripts,
			fmt.Sprintf(highlightJsScript, JoinPath(highlightJsScriptDefaultPath)),
			highlightJsAction)
	}
}
