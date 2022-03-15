package internals

type TypeContent uint8

const (
	TypeHeading TypeContent = iota
	TypeParagraph
	TypeList
	TypeListNumbered
	TypeLink
	TypeSourceCode
	TypeRawHTML
	TypeHorizontalLine
	TypeAttentionText
)

type Content struct {
	Type TypeContent

	HeadingLevel   int
	HeadingChild   bool
	HeadingFirst   bool
	HeadingLast    bool
	Heading        string
	Paragraph      string
	List           []string
	ListNumbered   []string
	Link           string
	LinkTitle      string
	SourceCode     string
	SourceCodeLang string
	RawHTML        string
	AttentionTitle string
	AttentionText  string
}

func (c Content) IsHeading() bool        { return c.Type == TypeHeading }
func (c Content) IsParagraph() bool      { return c.Type == TypeParagraph }
func (c Content) IsList() bool           { return c.Type == TypeList }
func (c Content) IsListNumbered() bool   { return c.Type == TypeListNumbered }
func (c Content) IsLink() bool           { return c.Type == TypeLink }
func (c Content) IsSourceCode() bool     { return c.Type == TypeSourceCode }
func (c Content) IsRawHTML() bool        { return c.Type == TypeRawHTML }
func (c Content) IsHorizontalLine() bool { return c.Type == TypeHorizontalLine }
func (c Content) IsAttentionBlock() bool { return c.Type == TypeAttentionText }
