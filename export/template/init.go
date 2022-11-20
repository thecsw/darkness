package template

import "github.com/thecsw/darkness/export"

var (
	// Make sure this exporter implements `export.Exporter`.
	exporter                 = &ExporterTemplate{}
	_        export.Exporter = exporter
	// Make sure this exporter implements `exporter.ExporterBuilder` .
	exporterBuilder                        = &ExporterTemplateBuilder{}
	_               export.ExporterBuilder = exporterBuilder
)

// This init registers this exporter with the root module.
func init() {
	export.Register("TEMPLATE EXPORTER", exporterBuilder)
}
