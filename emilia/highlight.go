package emilia

import (
	"fmt"

	"github.com/thecsw/darkness/internals"
)

const (
	highlightJsTheme             = `<link rel="stylesheet" href="%s">`
	highlightJsThemeDefaultPath  = `scripts/highlight/styles/agate.min.css`
	highlightJsScript            = `<script src="%s"></script>`
	highlightJsScriptDefaultPath = `scripts/highlight/highlight.min.js`
	highlightJsAction            = `<script>hljs.highlightAll();</script>`
)

func WithSyntaxHighlighting() internals.PageOption {
	return func(page *internals.Page) {
		if !Config.Website.SyntaxHighlighting {
			return
		}
		if !hasCodeBlocks(page) {
			return
		}
		page.Stylesheets = append(page.Stylesheets,
			fmt.Sprintf(highlightJsTheme, JoinPath(Config.Website.SyntaxHighlightingTheme)))
		page.Scripts = append(page.Scripts,
			fmt.Sprintf(highlightJsScript, JoinPath(highlightJsScriptDefaultPath)),
			highlightJsAction)
	}
}

func hasCodeBlocks(page *internals.Page) bool {
	for _, content := range page.Contents {
		if !content.IsSourceCode() && len(content.SourceCodeLang) > 0 {
			return true
		}
	}
	return false
}
