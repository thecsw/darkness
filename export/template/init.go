package template

import "github.com/thecsw/darkness/export"

var (
	// Make sure that this exporter implements `export.Exporter`.
	exporter                        = &ExporterTemplate{}
	_        export.Exporter        = exporter
	_        export.ExporterBuilder = exporter
)

// This init registers this exporter with the root module.
func init() {
	export.Register("TEMPLATE EXPORTER", exporter)
}
