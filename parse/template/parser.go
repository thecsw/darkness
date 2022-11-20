package template

import "github.com/thecsw/darkness/parse"

// ParserTemplateBuilder builds `ParserTemplate`.
type ParserTemplateBuilder struct{}

// ParserTemplate implements `parse.Parser`.
type ParserTemplate struct {
	// Filename is the filename of the source input.
	Filename string
	// Data is the contents that need to be parsed.
	Data string
}

// BuildParser will create a new parser object and return it.
func (ParserTemplateBuilder) BuildParser(filename, data string) parse.Parser {
	return &ParserTemplate{
		Filename: filename,
		Data:     data,
	}
}
