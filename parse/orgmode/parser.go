package orgmode

import (
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
)

// ParserOrgmodeBuilder builds `ParserOrgmode`.
type ParserOrgmodeBuilder struct{}

// ParserOrgmode implements `parse.Parser`.
type ParserOrgmode struct {
	// Filename is the filename of the source input.
	Filename yunyun.RelativePathFile
	// Data is the contents that need to be parsed.
	Data string
}

// BuildParser builds the `Parser` interface object.
func (ParserOrgmodeBuilder) BuildParser(
	filename yunyun.RelativePathFile, data string,
) parse.Parser {
	return &ParserOrgmode{
		Filename: filename,
		Data:     data,
	}
}
