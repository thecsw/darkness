package export

import "github.com/thecsw/darkness/yunyun"

// Exporter is a generic interface that other output extensions should
// implement.
type Exporter interface {
	// SetPage sets the import source for the exporter.
	SetPage(*yunyun.Page)
	// Export performs the exporting with options passed in.
	Export() string

	// Heading exports `TypeHeading` content.
	Heading(*yunyun.Content) string
	// Paragraph exports `TypeParagraph` content.
	Paragraph(*yunyun.Content) string
	// List exports `TypeList` content.
	List(*yunyun.Content) string
	// ListNumbered exports `TypeListNumbered` content.
	ListNumbered(*yunyun.Content) string
	// Link exports `TypeLink` content.
	Link(*yunyun.Content) string
	// SourceCode exports `TypeSourceCode` content.
	SourceCode(*yunyun.Content) string
	// RawHTML exports `TypeRawHTML` content.
	RawHTML(*yunyun.Content) string
	// HorizontalLine exports `TypeHorizontalLine` content.
	HorizontalLine(*yunyun.Content) string
	// AttentionBlock exports `TypeAttentionText` content.
	AttentionBlock(*yunyun.Content) string
	// Table exports `TypeTable` content.
	Table(*yunyun.Content) string
	// Details exports `TypeDetails` content.
	Details(*yunyun.Content) string
}

// ContentBuilder returns the map of type content builder functions.
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
