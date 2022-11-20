package orgmode

import "github.com/thecsw/darkness/parse"

// ParserOrgmode implements `parse.Parser`.
type ParserOrgmode struct {
	// Filename is the filename of the source input.
	Filename string
	// Data is the contents that need to be parsed.
	Data string
}

func (p ParserOrgmode) WithFilenameData(filename, data string) parse.Parser {
	return &ParserOrgmode{
		Filename: filename,
		Data:     data,
	}
}
