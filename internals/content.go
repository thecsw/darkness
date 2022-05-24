package internals

// TypeContent is the type of content, used for enums
type TypeContent uint8

const (
	// TypeHeading is the type of heading
	TypeHeading TypeContent = iota
	// TypeParagraph is the type of paragraph, which is just text
	TypeParagraph
	// TypeList is the type of unordered list
	TypeList
	// TypeListNumbered is the type of numbered list
	TypeListNumbered
	// TypeLink is the type of link
	TypeLink
	// TypeSourceCode is the type of source code block
	TypeSourceCode
	// TypeRawHTML is the type of raw HTML block
	TypeRawHTML
	// TypeHorizontalLine is the type of horizontal line
	TypeHorizontalLine
	// TypeAttentionText is the type of attention text block
	TypeAttentionText
	// TypeTable is the type of a table
	TypeTable
)

// Content is a piece of content of a page
type Content struct {
	// Type is the type of content
	Type TypeContent

	// HeadingLevel is the heading level of the content (1 being the title, starts at 2)
	HeadingLevel int
	// HeadingChild tells us if the current heading is a child of some previous heading
	HeadingChild bool
	// HeadingFirst tells us if the current heading is the first heading on the page
	HeadingFirst bool
	// HeadingLast tells us if the current heading is the last heading on the page
	HeadingLast bool
	// Heading is the heading text
	Heading string
	// Paragraph is the paragraph text
	Paragraph string
	// List is the list of items, unordered
	List []string
	// ListNumbered is the list of items, numbered
	ListNumbered []string
	// Link is the link text
	Link string
	// LinkTitle is the link title
	LinkTitle string
	// SourceCode is the source code
	SourceCode string
	// SourceCodeLanguage is the language of the source code
	SourceCodeLang string
	// RawHTML is the raw HTML
	RawHTML string
	// AttentionTitle is the attention text title (IMPORTANT, WARNING, etc.)
	AttentionTitle string
	// AttentionText is the attention text
	AttentionText string
	// IsQuote defines if this part is a quote or not
	IsQuote bool
	// IsCentered define if this part should be centered
	IsCentered bool
	// IsDropCap defines whether to enlarge the first letter
	IsDropCap bool
	// Table is the table of items
	Table [][]string
	// TableHeaders tell us whether the table has headers
	// (use the first row as headers instead of data)
	TableHeaders bool
	// Caption is the current caption
	Caption string
}

// IsHeading tells us if the content is a heading
func (c Content) IsHeading() bool { return c.Type == TypeHeading }

// IsParagraph tells us if the content is a paragraph
func (c Content) IsParagraph() bool { return c.Type == TypeParagraph }

// IsList tells us if the content is a list
func (c Content) IsList() bool { return c.Type == TypeList }

// IsListNumbered tells us if the content is a numbered list
func (c Content) IsListNumbered() bool { return c.Type == TypeListNumbered }

// IsLink tells us if the content is a link
func (c Content) IsLink() bool { return c.Type == TypeLink }

// IsSourceCode tells us if the content is a source code block
func (c Content) IsSourceCode() bool { return c.Type == TypeSourceCode }

// IsRawHTML tells us if the content is a raw HTML block
func (c Content) IsRawHTML() bool { return c.Type == TypeRawHTML }

// IsHorizontalLine tells us if the content is a horizontal line
func (c Content) IsHorizontalLine() bool { return c.Type == TypeHorizontalLine }

// IsAttentionBlock tells us if the content is an attention text block
func (c Content) IsAttentionBlock() bool { return c.Type == TypeAttentionText }

// IsTable tells us if the content block is a table
func (c Content) IsTable() bool { return c.Type == TypeTable }
