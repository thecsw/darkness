package emilia

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/gana"
)

var (
	// Config is the global darkness config.
	Config *DarknessConfig
	// ParserBuilder returns the parser builder.
	ParserBuilder parse.ParserBuilder
	// ExporterBuilder returns the exporter builder.
	ExporterBuilder export.ExporterBuilder
)

// EmiliaOptions is used for passing options when initiating emilia.
type EmiliaOptions struct {
	// DarknessConfig is the location of darkness's toml config file.
	DarknessConfig string
	// Dev enables URL generation through local paths.
	Dev bool
	// Test enables test environment, where darkness config is not needed.
	Test bool
}

// InitDarkness initializes the darkness config.
func InitDarkness(options *EmiliaOptions) {
	Config = &DarknessConfig{}
	data, err := ioutil.ReadFile(options.DarknessConfig)
	if err != nil && !options.Test {
		fmt.Printf("failed to open the config %s: %s", options.DarknessConfig, err.Error())
		os.Exit(1)
	}
	_, err = toml.Decode(string(data), Config)
	if err != nil {
		fmt.Printf("failed to decode the config %s: %s", options.DarknessConfig, err.Error())
		os.Exit(1)
	}
	// If input/output formats are empty, default to .org/.html respectively.
	if isZero(Config.Project.Input) {
		Config.Project.Input = puck.ExtensionOrgmode
	}
	if isZero(Config.Project.Output) {
		Config.Project.Output = puck.ExtensionHtml
	}
	// Build the parser and exporter builders.
	ParserBuilder = getParserBuilder()
	ExporterBuilder = getExporterBuilder()

	// Define the preview filename.
	if isZero(Config.Website.Preview) {
		Config.Website.Preview = puck.DefaultPreviewFile
	}
	if Config.Website.DescriptionLength < 1 {
		Config.Website.DescriptionLength = 100
	}
	// If the URL is empty, then plug in the current directory.
	if len(Config.URL) < 1 || options.Dev {
		Config.URL, err = os.Getwd()
		if err != nil {
			fmt.Printf("failed to get current directory because config url was not given: %s", err.Error())
			os.Exit(1)
		}
	}
	// URL must end with a trailing forward slash
	if !strings.HasSuffix(Config.URL, "/") {
		Config.URL += "/"
	}
	Config.URLPath, err = url.Parse(Config.URL)
	if err != nil {
		fmt.Printf("failed to parse url from config %s: %s", Config.URL, err.Error())
		os.Exit(1)
	}
	// Set up the custom highlight languages
	if !isZero(Config.Website.SyntaxHighlightingLanguages) {
		setupHighlightJsLanguages(Config.Website.SyntaxHighlightingLanguages)
	}
	// Set the default syntax highlighting theme
	if isZero(Config.Website.SyntaxHighlightingTheme) {
		Config.Website.SyntaxHighlightingTheme = highlightJsThemeDefaultPath
	}
	// Build the regex that will be used to exclude files that
	// have been denoted in emilia darkness config.
	excludePattern := fmt.Sprintf("(?mU)(%s)/.*",
		strings.Join(Config.Project.Exclude, "|"))
	Config.Project.ExcludeRegex, err = regexp.Compile(excludePattern)
	if err != nil {
		fmt.Println("Bad exclude regex, made", excludePattern,
			"\nFailed with error:", err.Error())
		os.Exit(1)
	}
	// Monkey patch the function if we're using the roman footnotes.
	if Config.Website.RomanFootnotes {
		FootnoteLabeler = numberToRoman
	}
}

// JoinPath joins a path to the darkness config URL.
func JoinPath(elem string) string {
	u, _ := url.Parse(elem)
	// if err != nil {
	// 	fmt.Printf("Failde to parse target %s: %s", elem, err.Error())
	// 	os.Exit(1)
	// }
	return Config.URLPath.ResolveReference(u).String()
}

// isZero returns true if the passed value is a zero value of its type.
func isZero[T comparable](what T) bool {
	return what == gana.ZeroValue[T]()
}

// getParserBuilder returns a new parser object.
func getParserBuilder() parse.ParserBuilder {
	if v, ok := parse.ParserMap[Config.Project.Input]; ok {
		return v
	}
	fmt.Printf("No'%s parser, defaulting to %s\n", Config.Project.Input, puck.ExtensionOrgmode)
	Config.Project.Input = puck.ExtensionOrgmode
	return parse.ParserMap[puck.ExtensionOrgmode]
}

// getExporterBuilder returns a new exporter object.
func getExporterBuilder() export.ExporterBuilder {
	if v, ok := export.ExporterMap[Config.Project.Output]; ok {
		return v
	}
	fmt.Printf("No %s exporter, defaulting to %s\n", Config.Project.Output, puck.ExtensionHtml)
	Config.Project.Output = puck.ExtensionHtml
	return export.ExporterMap[puck.ExtensionHtml]
}

// setupHighlightJsLanguages logs all the languages we support through
// the directory included in the config.
func setupHighlightJsLanguages(dir string) {
	languages, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Failed to open %s: %s", dir, err.Error())
		AvailableLanguages = nil
		return
	}
	for _, language := range languages {
		if !strings.HasSuffix(language.Name(), ".min.js") {
			continue
		}
		AvailableLanguages[strings.TrimSuffix(language.Name(), ".min.js")] = true
	}
}

// trimExt trims extension of a file (only top level, so `file.min.js`
// will be `file.min`)
func trimExt(s string) string {
	return strings.TrimSuffix(s, filepath.Ext(s))
}
