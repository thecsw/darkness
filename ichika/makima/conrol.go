package makima

import (
	"io"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/ichika/chiho"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
)

type Control struct {
	Conf     *alpha.DarknessConfig
	Parser   parse.Parser
	Exporter export.Exporter

	InputFilename yunyun.FullPathFile
	Input         string

	Page *yunyun.Page

	OutputFilename string
	Output         io.Reader
}

func (c *Control) Parse() *Control {
	c.Page = c.Parser.Do(c.Conf.Runtime.WorkDir.Rel(c.InputFilename), c.Input)
	return c
}

func (c *Control) Export() *Control {
	c.OutputFilename = c.Conf.Project.InputFilenameToOutput(c.InputFilename)
	c.Output = c.Exporter.Do(chiho.EnrichPage(c.Conf, c.Page))
	return c
}
