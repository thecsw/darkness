package alpha

import (
	"regexp"

	"github.com/thecsw/darkness/v3/yunyun"
)

// DarknessConfig is the global darkness config.
type DarknessConfig struct {
	// Navigation is the navigation section of the config.
	Navigation map[string]NavigationConfig `toml:"navigation"`

	// Title is the title of the site.
	Title string `toml:"title"`

	// Url is the Url of the site.
	Url string `toml:"url"`

	// RSS is the rss config.
	RSS RSSConfig `toml:"rss"`

	// Author is the author section of the config.
	Author AuthorConfig `toml:"author"`

	// Project is the project section of the config.
	Project ProjectConfig `toml:"project"`

	// Website is the website section of the config.
	Website WebsiteConfig `toml:"website"`

	// Runtime holds the state we use during the runtime.
	Runtime RuntimeConfig `toml:"-"`

	// External is config for external services.
	External ExternalConfig `toml:"external"`
}

// ProjectConfig is the project section of the config
type ProjectConfig struct {
	ExcludeRegex *regexp.Regexp `toml:"-"`
	// Input is the input format (default ".org")
	Input string `toml:"input"`

	// Output is the output format (defaulte ".html")
	Output string `toml:"output"`

	// DarknessVendorDirectory where to vendor, default to `darkness_vendor`.
	DarknessVendorDirectory yunyun.RelativePathDir `toml:"vendor_directory"`

	// DarknessPreviewDirectory where to store previews, default to `darkness_preview`.
	DarknessPreviewDirectory yunyun.RelativePathDir `toml:"preview_directory"`

	// Excludes is the list of relative paths to exclude from the project
	Exclude []yunyun.RelativePathDir `toml:"exclude"`

	ExcludeEnabled bool `toml:"-"`
}

// WebsiteConfig is the website section of the config
type WebsiteConfig struct {
	// SyntaxHighlightingLanguages is the location of highlight.js languages
	SyntaxHighlightingLanguages yunyun.RelativePathDir `toml:"syntax_highlighting_languages"`

	// Color is the color of the site
	Color string `toml:"color"`

	// Twitter is the twitter handle of the site
	Twitter string `toml:"twitter"`

	// Preview is the filename of the picture in the
	// same directory to use as a page preview. defaults to preview.png
	Preview yunyun.RelativePathFile `toml:"preview"`

	// Locale is the locale of the site
	Locale string `toml:"locale"`

	// SyntaxHighlightingTheme decides what theme to use from highlight.js
	SyntaxHighlightingTheme yunyun.RelativePathFile `toml:"syntax_highlighting_theme"`

	// Styles is the list of relative paths to css files for the site
	Styles []yunyun.RelativePathFile `toml:"styles"`

	// Tombs is the list of relative paths where to include the tombstones
	Tombs []yunyun.RelativePathDir `toml:"tombs"`

	// ExtraHead is the list of html directives to insert in <head>
	ExtraHead []string `toml:"extra_head"`

	// Description length dictates on how many characters do we extract
	// from the page to show in the web prewies, like OpenGraph and Twitter
	DescriptionLength int `toml:"description_length"`

	// SyntaxHighlighting decides whether to enable code blocks'
	// syntax highlighting with highlight.js
	SyntaxHighlighting bool `toml:"syntax_highlighting"`

	// ClickableImages marks whether the images should href to the img link.
	ClickableImages bool `toml:"clickable_images"`

	// FootnoteBrackets decides whether to use brackets on footnotes
	FootnoteBrackets bool `toml:"footnote_brackets"`

	// RomanFootnotes tells if we have to use roman numerals for footnotes
	RomanFootnotes bool `toml:"roman_footnotes"`
}

// AuthorConfig is the author section of the config
type AuthorConfig struct {
	// AuthorImage is the header image (can be empty)
	Image yunyun.RelativePathFile `toml:"image"`

	// ImagePreComputed is what actually gets stubbed in.
	ImagePreComputed yunyun.FullPathFile `toml:"-"`

	// Name is the name of the author
	Name string `toml:"name"`

	// Email is the email of the author
	Email string `toml:"email"`

	// NameEnable will not show the name in the menu
	// if disabled, but will use it for metadata.
	NameEnable bool `toml:"name_enable"`

	// EmailEnable will hide the email if false.
	EmailEnable bool `toml:"email_enable"`
}

// NavigationConfig is the navigation section of the config
type NavigationConfig struct {
	// Link is the link of the navigation item
	Link yunyun.RelativePathDir `toml:"link"`

	// Title is the title of the navigation item
	Title string `toml:"title"`

	// Hide is the path where to hide the element
	Hide yunyun.RelativePathDir `toml:"hide"`

	// Will hide the element if this value is encountered.
	HideIf yunyun.RelativePathDir `toml:"hide_if"`
}

// RSSConfig is for filling out rss stuff.
type RSSConfig struct {
	// The language the channel is written in. This allows
	// aggregators to group all Italian language sites, for
	// example, on a single page. A list of allowable values
	// for this element, as provided by Netscape, is here.
	// You may also use values defined by the W3C.
	//
	// Example: "en-us"
	Language string `toml:"language"`

	// Phrase or sentence describing the channel.
	//
	// The latest news from GoUpstate.com, a Spartanburg
	// Herald-Journal Web site.
	Description string `toml:"description"`

	// Copyright notice for content in the channel.
	//
	// Example: "Copyright 2002, Spartanburg Herald-Journal"
	Copyright string `toml:"copyright"`

	// Email address for person responsible for editorial content.
	//
	// Example: "geo@herald.com (George Matesky)"
	ManagingEditor string `toml:"managing_editor"`

	// DefaultAuthor name to use for posts. Overwritten by posts'
	// author directives.
	//
	// Example: "Sandy Urazayev"
	DefaultAuthor string `toml:"default_author"`

	// Email address for person responsible for technical issues
	// relating to channel.
	//
	// Example: "betty@herald.com (Betty Guernsey)"
	WebMaster string `toml:"web_master"`

	// Specify one or more categories that the channel belongs to.
	// Follows the same rules as the <item>-level category element.
	//
	// Example: "<category>Newspapers</category>"
	Category string `toml:"category"`

	// If true, darkness will add the rss icon to the menu.
	Enable bool `toml:"enable"`

	// Timezone sets the default timezone for RSS timestamps.
	Timezone string `toml:"timezone"`

	// DefaultHour defines the hour value in RSS timestamp if one
	// is not provided. Use the 24 hrs.
	DefaultHour int `toml:"default_hour"`
}

// ExternalConfig is config for external services.
type ExternalConfig struct {
	// SearchEngines is the list of search engines to notify of
	// new and refreshed content.
	SearchEngines []string `toml:"search_engines"`

	// GitRemoteService would be like github.com and
	// GitRemotePath would be thecsw/whatever.
	GitRemoteService   string `toml:"git_remote_service"`
	GitRemotePath      string `toml:"git_remote_path"`
	GitBranch          string `toml:"git_branch"`
	GitRemotesAreValid bool   `toml:"-"`
}
