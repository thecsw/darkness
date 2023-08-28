package alpha

import (
	"fmt"
	"os"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

const (
	highlightJsThemeDefaultPath yunyun.RelativePathFile = `scripts/highlight/styles/agate.min.css`
)

// setupHighlightJsLanguages logs all the languages we support through
// the directory included in the config.
func (conf *DarknessConfig) setupHighlightJsLanguages() error {
	dir := conf.Website.SyntaxHighlightingLanguages
	languages, err := os.ReadDir(string(dir))
	if err != nil {
		return fmt.Errorf("opening %s: %v", dir, err)
	}
	conf.Runtime.HtmlHighlightLanguages = map[string]struct{}{}
	for _, language := range languages {
		// Skip if it's not a proper minified js from highlight.js
		if !strings.HasSuffix(language.Name(), ".min.js") {
			continue
		}
		langName := strings.TrimSuffix(language.Name(), ".min.js")
		conf.Runtime.HtmlHighlightLanguages[langName] = struct{}{}
	}
	return nil
}
