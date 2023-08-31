package makima

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/ichika/chiho"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
)

// Control is the struct that is passed across darkness to build the site.
type Control struct {
	// Conf is the configuration for the site.
	Conf *alpha.DarknessConfig
	// Parser is the parser to use for the site.
	Parser parse.Parser
	// Exporter is the exporter to use for the site.
	Exporter export.Exporter

	// InputFilename is the filename of the input file.
	InputFilename yunyun.FullPathFile
	// Input is the input file's contents.
	Input string

	// Page is the parsed page.
	Page *yunyun.Page

	// OutputFilename is the filename of the output file.
	OutputFilename string
	// Output is the output file's contents.
	Output io.Reader
}

// Read reads the input file and returns the Control.
func (c *Control) Read() (Woof, error) {
	file, err := os.ReadFile(filepath.Clean(string(c.InputFilename)))
	if err != nil {
		return nil, fmt.Errorf("reading input file %s: %v", c.InputFilename, err)
	}
	c.Input = string(file)
	return c, nil
}

// Parse parses the input file and returns the Control.
func (c *Control) Parse() Woof {
	c.Page = c.Parser.Do(c.Conf.Runtime.WorkDir.Rel(c.InputFilename), c.Input)
	return c
}

// Export exports the parsed page and returns the Control.
func (c *Control) Export() Woof {
	c.OutputFilename = c.Conf.Project.InputFilenameToOutput(c.InputFilename)
	c.Output = c.Exporter.Do(chiho.EnrichPage(c.Conf, c.Page))
	return c
}

// Write copies the exported contents onto the output file.
func (c *Control) Write() error {
	file, err := os.Create(c.OutputFilename)
	if err != nil {
		return fmt.Errorf("creating output file %s: %v", c.OutputFilename, err)
	}
	if _, err := io.Copy(file, c.Output); err != nil {
		return fmt.Errorf("writing to output file %s: %v", c.OutputFilename, err)
	}
	return nil
}
