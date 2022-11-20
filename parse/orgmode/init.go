package orgmode

import (
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/parse"
)

var (
	// Make sure that this parser implements `parse.Parser`.
	parser                     = &ParserOrgmode{}
	_      parse.Parser        = parser
	_      parse.ParserBuilder = parser
)

// This init registers the parser with the root module.
func init() {
	parse.Register(puck.ExtensionOrgmode, parser)
}
