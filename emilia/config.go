package emilia

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
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
	if notDefined(Config.Project.Input) {
		Config.Project.Input = ".org"
	}
	if notDefined(Config.Project.Output) {
		Config.Project.Output = ".html"
	}
	if notDefined(Config.Website.Preview) {
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
	if notDefined(Config.Website.SyntaxHighlightingTheme) {
		Config.Website.SyntaxHighlightingTheme = highlightJsThemeDefaultPath
	}
}

// JoinPath joins a path to the darkness config URL
func JoinPath(elem string) string {
	u, _ := url.Parse(elem)
	return Config.URLPath.ResolveReference(u).String()
}

func notDefined(what string) bool {
	return what == ""
}
