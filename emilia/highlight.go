package emilia

import (
	"fmt"
	"path/filepath"

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

var (
	// AvailableLanguages has the map lookup if a highlight.js language is found.
	AvailableLanguages = map[string]bool{}
)

func WithSyntaxHighlighting() yunyun.PageOption {
	return func(page *yunyun.Page) {
		// If Emilia disabled the syntax highlighting, don't even bother.
		if !Config.Website.SyntaxHighlighting {
			return
		}
		// Find all the code blocks.
		sourceCodes := gana.Filter(func(v *yunyun.Content) bool { return v.IsSourceCode() }, page.Contents)
		// If there are none, the page doesn't require highlighting.
		if len(sourceCodes) < 1 {
			return
		}
		// Add the basic processing scripts.
		page.Stylesheets = append(page.Stylesheets,
			fmt.Sprintf(highlightJsTheme, JoinPath(Config.Website.SyntaxHighlightingTheme)))
		page.Scripts = append(page.Scripts,
			fmt.Sprintf(highlightJsScript, JoinPath(highlightJsScriptDefaultPath)))
		// Trigger the action after all the highlight scripts are imported.
		defer func() {
			page.Scripts = append(page.Scripts, highlightJsAction)
		}()
		// If language lookup table was not filled, skip the next step.
		if AvailableLanguages == nil {
			return
		}
		addedLanguages := map[string]bool{}
		// For each codeblock, look up the language and see if we have a
		// highlight.js processor for it. If we don't simply skip, otherwise,
		// build a new script import and inject it into the page.
		for _, sourceCode := range sourceCodes {
			lang := MapSourceCodeLang(sourceCode.SourceCodeLang)
			if _, ok := AvailableLanguages[lang]; !ok {
				lang = defaultHighlightLanguage
			}
			if _, alreadyAdded := addedLanguages[lang]; alreadyAdded {
				continue
			}
			page.Scripts = append(page.Scripts, fmt.Sprintf(highlightJsScript,
				JoinPath(filepath.Join(
					Config.Website.SyntaxHighlightingLanguages,
					lang+".min.js"))))
			addedLanguages[lang] = true
		}
	}
}

const (
	defaultHighlightLanguage = "plaintext"
)

// sourceCodeLang maps our source code language name to the
// name that highlight.js will need when coloring code
var sourceCodeLang = map[string]string{
	"":           defaultHighlightLanguage,
	"sh":         "bash",
	"emacs-lisp": "lisp",
}

// MapSourceCodeLang tries to map the simple source code language
// to the one that highlight.js would accept.
func MapSourceCodeLang(s string) string {
	if v, ok := sourceCodeLang[s]; ok {
		return v
	}
	// Default to whatever was passed.
	return s
}
