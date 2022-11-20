package template

import (
	"github.com/thecsw/darkness/parse"
)

var (
	// Make sure that this parser implements `parse.Parser`.
	parser                     = &ParserTemplate{}
	_      parse.Parser        = parser
	_      parse.ParserBuilder = parser
)

// This init registers the parser with the root module.
func init() {
	parse.Register("TEMPLATE PARSER", parser)
}
