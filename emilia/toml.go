package emilia

import "net/url"

// DarknessConfig is the global darkness config
type DarknessConfig struct {
	// Title is the title of the site
	Title string
	// URL is the URL of the site
	URL string
	// Project is the project section of the config
	Project ProjectConfig
	// Website is the website section of the config
	Website WebsiteConfig
	// Author is the author section of the config
	Author AuthorConfig
	// Navigation is the navigation section of the config
	Navigation map[string]NavigationConfig
	// URLPath is the parsed URL of the site
	URLPath *url.URL
}

// ProjectConfig is the project section of the config
type ProjectConfig struct {
	// Input is the input format (default ".org")
	Input string
	// Output is the output format (defaulte ".html")
	Output string
	// Excuses is the list of relative paths to exclude from the project
	Exclude []string
}

// WebsiteConfig is the website section of the config
type WebsiteConfig struct {
	// Locale is the locale of the site
	Locale string
	// Color is the color of the site
	Color string
	// Twitter is the twitter handle of the site
	Twitter string
	// Styles is the list of relative paths to css files for the site
	Styles []string
	// Tombs is the list of relative paths where to include the tombstones
	Tombs []string
}

// AuthorConfig is the author section of the config
type AuthorConfig struct {
	// AuthorImage is the header image (can be empty)
	Image string
	// Name is the name of the author
	Name string
	// Email is the email of the author
	Email string
}

// NavigationConfig is the navigation section of the config
type NavigationConfig struct {
	// Link is the link of the navigation item
	Link string
	// Title is the title of the navigation item
	Title string
	// Hide is the path where to hide the element
	Hide string
}
