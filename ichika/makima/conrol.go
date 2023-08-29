package makima

import (
	"io"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
)

type Control struct {
	Conf     alpha.DarknessConfig
	Parser   parse.Parser
	Exporter export.Exporter

	InputFilename yunyun.FullPathFile
	Input         string

	Page *yunyun.Page

	OutputFilename string
	Output         io.Reader
}
