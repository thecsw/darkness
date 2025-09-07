package alpha

import (
	"net/url"

	l "github.com/charmbracelet/log"
)

// WorkingDirectory is the directory of where darkness project lives.
type WorkingDirectory string

// RuntimeConfig the config that is not a part of the config file,
// but us building when processing stuff.
type RuntimeConfig struct {
	// WorkDir is the directory of where darkness project lives.
	WorkDir WorkingDirectory

	// Slice with just `Url` in it.
	urlSlice []string

	// UrlPath is the parsed Url of the site
	UrlPath *url.URL

	// isUrlLocal is true if Url is the file path, not url.
	isUrlLocal bool

	// VendorGalleries tells us if we need to stub local copies
	// of remote links in galleries.
	VendorGalleries bool

	// HtmlHighlightLanguages is a map of languages that we want to
	// highlight in HTML.
	HtmlHighlightLanguages map[string]struct{}

	// Logger is the logger that we use.
	Logger *l.Logger

	// WriteParsedPagesAsJson will flush parsed pages as json in the same dir.
	WriteParsedPagesAsJson bool
}
