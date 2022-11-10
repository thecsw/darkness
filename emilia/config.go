package emilia

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	// Config is the global darkness config
	Config *DarknessConfig
)

// InitDarkness initializes the darkness config
func InitDarkness(file string) {
	Config = &DarknessConfig{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("failed to open the config %s: %s", file, err.Error())
		os.Exit(1)
	}
	_, err = toml.Decode(string(data), Config)
	if err != nil {
		fmt.Printf("failed to decode the config %s: %s", file, err.Error())
		os.Exit(1)
	}
	// If input/output formats are empty, default to .org/.html respectively
	if undef(Config.Project.Input) {
		Config.Project.Input = ".org"
	}
	if undef(Config.Project.Output) {
		Config.Project.Output = ".html"
	}
	if undef(Config.Website.Preview) {
		Config.Website.Preview = "preview.png"
	}
	if Config.Website.DescriptionLength < 1 {
		Config.Website.DescriptionLength = 100
	}
	// If the URL is empty, then plug in the current directory
	if len(Config.URL) < 1 {
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
	// Set the default syntax highlighting theme
	if undef(Config.Website.SyntaxHighlightingTheme) {
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
	// Monkey patch the function if we're using the roman footnotes
	if Config.Website.RomanFootnotes {
		FootnoteLabeler = numberToRoman
	}
}

// JoinPath joins a path to the darkness config URL
func JoinPath(elem string) string {
	u, _ := url.Parse(elem)
	return Config.URLPath.ResolveReference(u).String()
}

func undef(what string) bool {
	return what == ""
}
