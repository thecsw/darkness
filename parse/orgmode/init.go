package orgmode

import (
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/parse"
)

var (
	// Make sure this parser implements `parse.Parser`.
	parser              = &ParserOrgmode{}
	_      parse.Parser = parser
	// Make sure this parser builder implements `parse.ParserBuilder`.
	parserBuilder                     = &ParserOrgmodeBuilder{}
	_             parse.ParserBuilder = parserBuilder
)

// This init registers the parser with the root module.
func init() {
	parse.Register(puck.ExtensionOrgmode, parserBuilder)
}
