package html

import (
	"strconv"

	"github.com/thecsw/darkness/emilia"
)

var (
	// styleTagsProcessed is the processed style tags
	styleTagsProcessed string
)

var footnoteLabel = strconv.Itoa

// InitializeExporter initializes the constant tags
func InitializeExporter() {
	styleTagsProcessed = styleTags()
	// Monkey patch the function if we're using the roman footnotes
	if emilia.Config.Website.RomanFootnotes {
		footnoteLabel = numberToRoman
	}
}
