package html

import (
	"strconv"

	"github.com/thecsw/darkness/emilia"
)

var footnoteLabel = strconv.Itoa

// InitializeExporter initializes the constant tags
func InitializeExporter() {
	// Monkey patch the function if we're using the roman footnotes
	if emilia.Config.Website.RomanFootnotes {
		footnoteLabel = numberToRoman
	}
}
