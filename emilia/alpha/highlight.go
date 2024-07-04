package alpha

import (
	"fmt"
	"os"
	"strings"

	"github.com/thecsw/darkness/v3/yunyun"
)

const (
	// highlightJsThemeDefaultPath is the default path to the highlight.js
	highlightJsThemeDefaultPath yunyun.RelativePathFile = `scripts/highlight/styles/agate.min.css`
)

// setupHighlightJsLanguages logs all the languages we support through
// the directory included in the config.
func (conf *DarknessConfig) setupHighlightJsLanguages() error {
	// Don't do anything if it wasn't requested.
	if isUnset(conf.Website.SyntaxHighlightingLanguages) {
		return nil
	}

	// Read the directory and get all the languages.
	dir := conf.Website.SyntaxHighlightingLanguages
	languages, err := os.ReadDir(string(dir))
	if err != nil {
		return fmt.Errorf("opening %s: %v", dir, err)
	}

	// Add all the languages to the config.
	conf.Runtime.HtmlHighlightLanguages = map[string]struct{}{}
	for _, language := range languages {
		// Skip if it's not a proper minified js from highlight.js
		if !strings.HasSuffix(language.Name(), ".min.js") {
			continue
		}

		// Add the language to the config.
		langName := strings.TrimSuffix(language.Name(), ".min.js")
		conf.Runtime.HtmlHighlightLanguages[langName] = struct{}{}
	}
	return nil
}
