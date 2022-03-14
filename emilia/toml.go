package emilia

import "net/url"

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
	Tombs   []string
	Exclude []string
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
