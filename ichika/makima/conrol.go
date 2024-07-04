package makima

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/export"
	"github.com/thecsw/darkness/v3/ichika/chiho"
	"github.com/thecsw/darkness/v3/ichika/misaka"
	"github.com/thecsw/darkness/v3/parse"
	"github.com/thecsw/darkness/v3/yunyun"
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
	defer puck.
		Stopwatch("Read", "input", c.InputFilename).
		RecordWithFile(misaka.RecordReadTime, c.InputFilename)
	file, err := os.ReadFile(filepath.Clean(string(c.InputFilename)))
	if err != nil {
		return nil, fmt.Errorf("reading input file %s: %v", c.InputFilename, err)
	}
	c.Input = string(file)
	return c, nil
}

// Parse parses the input file and returns the Control.
func (c *Control) Parse() Woof {
	defer puck.
		Stopwatch("Parsed", "input", c.InputFilename).
		RecordWithFile(misaka.RecordParseTime, c.InputFilename)
	c.Page = c.Parser.Do(c.Conf.Runtime.WorkDir.Rel(c.InputFilename), c.Input)
	return c
}

// Export exports the parsed page and returns the Control.
func (c *Control) Export() Woof {
	defer puck.
		Stopwatch("Exported", "input", c.InputFilename).
		RecordWithFile(misaka.RecordExportTime, c.InputFilename)
	c.OutputFilename = c.Conf.Project.InputFilenameToOutput(c.InputFilename)
	c.Output = c.Exporter.Do(chiho.EnrichPage(c.Conf, c.Page))
	return c
}

// Write copies the exported contents onto the output file.
func (c *Control) Write() error {
	defer puck.
		Stopwatch("Wrote", "output", c.OutputFilename).
		RecordWithFile(misaka.RecordWriteTime, c.InputFilename)
	file, err := os.Create(c.OutputFilename)
	if err != nil {
		return fmt.Errorf("creating output file %s: %v", c.OutputFilename, err)
	}
	if _, err := io.Copy(file, c.Output); err != nil {
		return fmt.Errorf("writing to output file %s: %v", c.OutputFilename, err)
	}
	return nil
}
