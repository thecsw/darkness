package alpha

import (
	"net/url"

	l "github.com/charmbracelet/log"
)

type WorkingDirectory string

// RuntimeConfig the config that is not a part of the config file,
// but us building when processing stuff.
type RuntimeConfig struct {
	// WorkDir is the directory of where darkness project lives.
	WorkDir WorkingDirectory

	// Slice with just `URL` in it.
	URLSlice []string

	// URLPath is the parsed URL of the site
	URLPath *url.URL

	// URLIsLocal is true if URL is the file path, not url.
	URLIsLocal bool

	// VendorGalleries tells us if we need to stub local copies
	// of remote links in galleries.
	VendorGalleries bool

	HtmlHighlightLanguages map[string]struct{}

	Logger *l.Logger
}
