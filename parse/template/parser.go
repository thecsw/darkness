package template

import (
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
)

// ParserTemplateBuilder builds `ParserTemplate`.
type ParserTemplateBuilder struct{}

// ParserTemplate implements `parse.Parser`.
type ParserTemplate struct {
	// Filename is the filename of the source input.
	Filename yunyun.RelativePathFile
	// Data is the contents that need to be parsed.
	Data string
}

// BuildParser will create a new parser object and return it.
func (ParserTemplateBuilder) BuildParser(
	filename yunyun.RelativePathFile, data string,
) parse.Parser {
	return &ParserTemplate{
		Filename: filename,
		Data:     data,
	}
}
