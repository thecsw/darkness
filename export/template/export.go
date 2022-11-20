package template

import "github.com/thecsw/darkness/yunyun"

// Export exports the input `*yunyun.Page` to string output.
func (e *ExporterTemplate) Export() string {
	panic("implement me")
}

// Heading exports `TypeHeading` content.
func (e *ExporterTemplate) Heading(*yunyun.Content) string {
	panic("implement me")
}

// Paragraph exports `TypeParagraph` content.
func (e *ExporterTemplate) Paragraph(*yunyun.Content) string {
	panic("implement me")
}

// List exports `TypeList` content.
func (e *ExporterTemplate) List(*yunyun.Content) string {
	panic("implement me")
}

// ListNumbered exports `TypeListNumbered` content.
func (e *ExporterTemplate) ListNumbered(*yunyun.Content) string {
	panic("implement me")
}

// Link exports `TypeLink` content.
func (e *ExporterTemplate) Link(*yunyun.Content) string {
	panic("implement me")
}

// SourceCode exports `TypeSourceCode` content.
func (e *ExporterTemplate) SourceCode(*yunyun.Content) string {
	panic("implement me")
}

// RawHTML exports `TypeRawHTML` content.
func (e *ExporterTemplate) RawHTML(*yunyun.Content) string {
	panic("implement me")
}

// HorizontalLine exports `TypeHorizontalLine` content.
func (e *ExporterTemplate) HorizontalLine(*yunyun.Content) string {
	panic("implement me")
}

// AttentionBlock exports `TypeAttentionText` content.
func (e *ExporterTemplate) AttentionBlock(*yunyun.Content) string {
	panic("implement me")
}

// Table exports `TypeTable` content.
func (e *ExporterTemplate) Table(*yunyun.Content) string {
	panic("implement me")
}

// Details exports `TypeDetails` content.
func (e *ExporterTemplate) Details(*yunyun.Content) string {
	panic("implement me")
}
