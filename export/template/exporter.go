package template

import (
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/yunyun"
)

// ExporterTemplateBuilder builds `ExporterTemplate`.
type ExporterTemplateBuilder struct{}

// ExporterTemplate is a template exporter
type ExporterTemplate struct {
	page             *yunyun.Page
	contentFunctions []func(*yunyun.Content) string
}

// SetPage creates a new Exporter object and returns it with data filled.
func (ExporterTemplateBuilder) BuildExporter(page *yunyun.Page) export.Exporter {
	what := &ExporterTemplate{
		page: page,
	}
	what.contentFunctions = export.ContentBuilder(what)
	return what
}
