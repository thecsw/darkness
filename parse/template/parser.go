package template

import (
	"bufio"
	"io"

	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
)

// ParserTemplateBuilder builds `ParserTemplate`.
type ParserTemplateBuilder struct{}

// ParserTemplate implements `parse.Parser`.
type ParserTemplate struct {
	// Filename is the filename of the source input.
	Filename yunyun.RelativePathFile
	// Data is the contents that need to be parsed.
	Data string
}

// BuildParser will create a new parser object and return it.
func (ParserTemplateBuilder) BuildParser(
	filename yunyun.RelativePathFile, data string,
) parse.Parser {
	return &ParserTemplate{
		Filename: filename,
		Data:     data,
	}
}

func (ParserTemplateBuilder) BuildParserReader(
	filename yunyun.RelativePathFile, reader io.ReadCloser,
) parse.Parser {
	data, _ := io.ReadAll(bufio.NewReader(reader))
	if err := reader.Close(); err != nil {
		puck.Logger.Errorf("closing file %s: %v", filename, err)
	}
	return &ParserTemplate{
		Filename: filename,
		Data:     string(data),
	}
}
