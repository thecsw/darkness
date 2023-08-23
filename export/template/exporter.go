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

// BuildExporter builds the exporter to generate output from yunyun internals.
func (ExporterTemplateBuilder) BuildExporter(page *yunyun.Page) export.Exporter {
	what := &ExporterTemplate{
		page: page,
	}
	what.contentFunctions = export.ContentBuilder(what)
	return what
}
