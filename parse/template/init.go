package template

import (
	"github.com/thecsw/darkness/parse"
)

var (
	// Make sure this parser implements `parse.Parser`.
	parser              = &ParserTemplate{}
	_      parse.Parser = parser
	// Make sure this parser implements `parse.ParserBuilder`.
	parserBuilder                     = &ParserTemplateBuilder{}
	_             parse.ParserBuilder = parserBuilder
)

// This init registers the parser with the root module.
func init() {
	parse.Register("TEMPLATE PARSER", parserBuilder)
}
