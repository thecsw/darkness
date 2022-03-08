package main

type DarknessConfig struct {
	Title      string
	URL        string
	Website    WebsiteConfig
	Author     AuthorConfig
	Navigation map[string]NavigationConfig
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
