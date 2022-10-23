package emilia

import "net/url"

// DarknessConfig is the global darkness config
type DarknessConfig struct {
	// Title is the title of the site
	Title string `toml:"title"`
	// URL is the URL of the site
	URL string `toml:"url"`
	// Project is the project section of the config
	Project ProjectConfig
	// Website is the website section of the config
	Website WebsiteConfig
	// Author is the author section of the config
	Author AuthorConfig
	// Navigation is the navigation section of the config
	Navigation map[string]NavigationConfig
	// URLPath is the parsed URL of the site
	URLPath *url.URL `toml:"-"`
}

// ProjectConfig is the project section of the config
type ProjectConfig struct {
	// Input is the input format (default ".org")
	Input string `toml:"input"`
	// Output is the output format (defaulte ".html")
	Output string `toml:"output"`
	// Excludes is the list of relative paths to exclude from the project
	Exclude []string `toml:"exclude"`
}

// WebsiteConfig is the website section of the config
type WebsiteConfig struct {
	// Locale is the locale of the site
	Locale string `toml:"locale"`
	// Color is the color of the site
	Color string `toml:"color"`
	// Twitter is the twitter handle of the site
	Twitter string `toml:"twitter"`
	// Styles is the list of relative paths to css files for the site
	Styles []string `toml:"styles"`
	// Tombs is the list of relative paths where to include the tombstones
	Tombs []string `toml:"tombs"`
	// Preview is the filename of the picture in the
	// same directory to use as a page preview. defaults to preview.png
	Preview string `toml:"preview"`
	// Description length dictates on how many characters do we extract
	// from the page to show in the web prewies, like OpenGraph and Twitter
	DescriptionLength int `toml:"description_length"`
	// Normalize headings will shift heading levels if enabled
	NormalizeHeadings bool `toml:"normalize_headings"`
	// RomanFootnotes tells if we have to use roman numerals for footnotes
	RomanFootnotes bool `toml:"roman_footnotes"`
	// FootnoteBrackets decides whether to use brackets on footnotes
	FootnoteBrackets bool `toml:"footnote_brackets"`
	// SyntaxHighlighting decides whether to enable code blocks'
	// syntax highlighting with highlight.js
	SyntaxHighlighting bool `toml:"syntax_highlighting"`
	// SyntaxHighlightingTheme decides what theme to use from highlight.js
	SyntaxHighlightingTheme string `toml:"syntax_highlighting_theme"`
}

// AuthorConfig is the author section of the config
type AuthorConfig struct {
	// AuthorImage is the header image (can be empty)
	Image string `toml:"image"`
	// Name is the name of the author
	Name string `toml:"name"`
	// Email is the email of the author
	Email string `toml:"email"`
}

// NavigationConfig is the navigation section of the config
type NavigationConfig struct {
	// Link is the link of the navigation item
	Link string `toml:"link"`
	// Title is the title of the navigation item
	Title string `toml:"title"`
	// Hide is the path where to hide the element
	Hide string `toml:"hide"`
}
