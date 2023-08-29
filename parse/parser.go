package parse

import (
	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/parse/orgmode"
	"github.com/thecsw/darkness/yunyun"
)

type Parser interface {
	Do(yunyun.RelativePathFile, string) *yunyun.Page
}

func BuildParser(conf *alpha.DarknessConfig) Parser {
	var parser Parser
	switch conf.Project.Input {
	case puck.ExtensionOrgmode:
		parser = orgmode.ParserOrgmode{Conf: conf}
	}
	return parser
}
