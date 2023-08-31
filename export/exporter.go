package export

import (
	"io"
	"log"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/export/html"
	"github.com/thecsw/darkness/yunyun"
)

// Exporter is the interface for all exporters.
type Exporter interface {
	// Do is the main function of the exporter.
	Do(*yunyun.Page) io.Reader
}

// BuildExporter builds the exporter based on the config.
func BuildExporter(conf *alpha.DarknessConfig) Exporter {
	var exporter Exporter
	switch conf.Project.Output {
	case puck.ExtensionHtml: // html
		exporter = html.ExporterHTML{Config: conf}
	default: // unknown
		log.Fatalf("unknown output type: %s", conf.Project.Output)
	}
	return exporter
}
