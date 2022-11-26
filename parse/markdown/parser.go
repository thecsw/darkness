package markdown

import "github.com/thecsw/darkness/parse"

// ParserMarkdownBuilder builds `ParserMarkdown`.
type ParserMarkdownBuilder struct{}

// ParserMarkdown implements `parse.Parser`.
type ParserMarkdown struct {
	// Filename is the filename of the source input.
	Filename string
	// Data is the contents that need to be parsed.
	Data string
}

// BuildParser builds the `Parser` interface object.
func (ParserMarkdownBuilder) BuildParser(filename, data string) parse.Parser {
	return &ParserMarkdown{
		Filename: filename,
		Data:     data,
	}
}
