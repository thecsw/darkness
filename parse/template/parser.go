package template

import "github.com/thecsw/darkness/parse"

// ParserTemplate implements `parse.Parser`.
type ParserTemplate struct {
	// Filename is the filename of the source input.
	Filename string
	// Data is the contents that need to be parsed.
	Data string
}

// WithFilenameData will create a new parser object and return it
func (p ParserTemplate) WithFilenameData(filename, data string) parse.Parser {
	return &ParserTemplate{
		Filename: filename,
		Data:     data,
	}
}
