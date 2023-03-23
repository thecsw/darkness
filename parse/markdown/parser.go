package markdown

import (
	"bufio"
	"fmt"
	"io"

	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
)

// ParserMarkdownBuilder builds `ParserMarkdown`.
type ParserMarkdownBuilder struct{}

// ParserMarkdown implements `parse.Parser`.
type ParserMarkdown struct {
	// Filename is the filename of the source input.
	Filename yunyun.RelativePathFile
	// Data is the contents that need to be parsed.
	Data string
}

// BuildParser builds the `Parser` interface object.
func (ParserMarkdownBuilder) BuildParser(
	filename yunyun.RelativePathFile, data string,
) parse.Parser {
	return &ParserMarkdown{
		Filename: filename,
		Data:     data,
	}
}
func (ParserMarkdownBuilder) BuildParserReader(
	filename yunyun.RelativePathFile, reader io.ReadCloser,
) parse.Parser {
	data, _ := io.ReadAll(bufio.NewReader(reader))
	if err := reader.Close(); err != nil {
		fmt.Printf("Failed to close file handler for %s: %s\n", filename, err)
	}
	return &ParserMarkdown{
		Filename: filename,
		Data:     string(data),
	}
}
