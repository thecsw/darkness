package parse

import (
	"log"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/parse/orgmode"
	"github.com/thecsw/darkness/v3/yunyun"
)

// Parser is the interface for all parsers.
type Parser interface {
	// Do parses the file and returns a Page.
	Do(yunyun.RelativePathFile, string) *yunyun.Page
}

// BuildParser builds a parser based on the config.
func BuildParser(conf *alpha.DarknessConfig) Parser {
	var parser Parser
	switch conf.Project.Input {
	case puck.ExtensionOrgmode: // orgmode
		parser = orgmode.ParserOrgmode{Config: conf}
	default: // unknown
		log.Fatalf("unknown input format: %s", conf.Project.Input)
	}
	return parser
}
