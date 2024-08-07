package export

import (
	"io"
	"log"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/export/html"
	"github.com/thecsw/darkness/v3/yunyun"
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
		exporter = html.ExporterHtml{Config: conf}
	default: // unknown
		log.Fatalf("unknown output type: %s", conf.Project.Output)
	}
	return exporter
}
