package export

import (
	"io"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/export/html"
	"github.com/thecsw/darkness/yunyun"
)

type Exporter interface {
	Do(*yunyun.Page) io.Reader
}

func BuildExporter(conf alpha.DarknessConfig) Exporter {
	var exporter Exporter
	switch conf.Project.Output {
	case puck.ExtensionHtml:
		exporter = html.ExporterHTML{Conf: conf}
	}
	return exporter
}
