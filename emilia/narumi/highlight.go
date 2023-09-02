package narumi

import (
	"fmt"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/yunyun"
)

const (
	highlightJsTheme                                     = `<link rel="stylesheet" href="%s">`
	highlightJsScript                                    = `<script src="%s"></script>`
	highlightJsScriptDefaultPath yunyun.RelativePathFile = `scripts/highlight/highlight.min.js`
	highlightJsAction                                    = `<script>hljs.highlightAll();</script>`
)

// WithSyntaxHighlighting adds syntax highlighting to the page.
func WithSyntaxHighlighting(conf *alpha.DarknessConfig) yunyun.PageOption {
	return func(page *yunyun.Page) {
		// If Emilia disabled the syntax highlighting, don't even bother.
		if !conf.Website.SyntaxHighlighting {
			return
		}
		// Find all the code blocks.
		sourceCodes := page.Contents.SourceCodeBlocks()
		// If there are none, the page doesn't require highlighting.
		if len(sourceCodes) < 1 {
			return
		}

		// Add the basic processing scripts.
		page.Stylesheets = append(page.Stylesheets,
			fmt.Sprintf(highlightJsTheme, conf.Runtime.Join(conf.Website.SyntaxHighlightingTheme)))
		page.Scripts = append(page.Scripts,
			fmt.Sprintf(highlightJsScript, conf.Runtime.Join(highlightJsScriptDefaultPath)))

		// Trigger the action after all the highlight scripts are imported.
		defer func() {
			page.Scripts = append(page.Scripts, highlightJsAction)
		}()

		// If language lookup table was not filled, skip the next step.
		if conf.Runtime.HtmlHighlightLanguages == nil {
			return
		}
		addedLanguages := map[string]struct{}{}

		// For each codeblock, look up the language and see if we have a
		// highlight.js processor for it. If we don't simply skip, otherwise,
		// build a new script import and inject it into the page.
		for _, sourceCode := range sourceCodes {
			lang := MapSourceCodeLang(sourceCode.SourceCodeLang)
			if _, ok := conf.Runtime.HtmlHighlightLanguages[lang]; !ok {
				lang = defaultHighlightLanguage
			}
			if _, alreadyAdded := addedLanguages[lang]; alreadyAdded {
				continue
			}
			page.Scripts = append(page.Scripts, fmt.Sprintf(highlightJsScript,
				conf.Runtime.Join(yunyun.JoinRelativePaths(
					conf.Website.SyntaxHighlightingLanguages,
					yunyun.RelativePathFile(lang+".min.js")),
				),
			))
			addedLanguages[lang] = struct{}{}
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
	// IsDefault to whatever was passed.
	return s
}
