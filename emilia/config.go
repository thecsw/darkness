package emilia

import (
	"io/ioutil"
	"net/url"

	"github.com/BurntSushi/toml"
)

var (
	Config *DarknessConfig
)

func InitDarkness(file string) {
	Config = &DarknessConfig{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	_, err = toml.Decode(string(data), Config)
	if err != nil {
		panic(err)
	}
	Config.URLPath, err = url.Parse(Config.URL)
	if err != nil {
		panic(err)
	}
}

func JoinPath(elem string) string {
	u, _ := url.Parse(elem)
	return Config.URLPath.ResolveReference(u).String()
}

type DarknessConfig struct {
	Title      string
	URL        string
	Website    WebsiteConfig
	Author     AuthorConfig
	Navigation map[string]NavigationConfig

	URLPath *url.URL
}

type WebsiteConfig struct {
	Locale  string
	Color   string
	Twitter string
	Styles  []string
}

type AuthorConfig struct {
	Header string
	Name   string
	Email  string
}

type NavigationConfig struct {
	Link  string
	Title string
}
