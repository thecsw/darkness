package export

import "github.com/thecsw/darkness/yunyun"

// Exporter is a generic interface that other output extensions should
// implement.
type Exporter interface {
	SetPage(*yunyun.Page)
	Export() string

	Heading(*yunyun.Content) string
	Paragraph(*yunyun.Content) string
	List(*yunyun.Content) string
	ListNumbered(*yunyun.Content) string
	Link(*yunyun.Content) string
	SourceCode(*yunyun.Content) string
	RawHTML(*yunyun.Content) string
	HorizontalLine(*yunyun.Content) string
	AttentionBlock(*yunyun.Content) string
	Table(*yunyun.Content) string
	Details(*yunyun.Content) string
}

func ContentBuilder(exporter Exporter) []func(*yunyun.Content) string {
	return []func(*yunyun.Content) string{
		exporter.Heading,
		exporter.Paragraph,
		exporter.List,
		exporter.ListNumbered,
		exporter.Link,
		exporter.SourceCode,
		exporter.RawHTML,
		exporter.HorizontalLine,
		exporter.AttentionBlock,
		exporter.Table,
		exporter.Details,
	}
}
