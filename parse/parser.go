package parse

import (
	"io"

	"github.com/thecsw/darkness/yunyun"
)

type ParserBuilder interface {
	// BuildParser returns a new `Parser` object with
	// filename and data set.
	BuildParser(yunyun.RelativePathFile, string) Parser
	BuildParserReader(yunyun.RelativePathFile, io.ReadCloser) Parser
}

// Parser is an interface used to define import packages,
// which convert source data into a yunyun `Page`.
type Parser interface {
	// Parse returns `*yunyunPage`.
	Parse() *yunyun.Page
}

// ParserMap stores mappings of extensions to their parsers.
var ParserMap = make(map[string]ParserBuilder)

// Register is called by parsers to register themselves.
func Register(ext string, p ParserBuilder) {
	ParserMap[ext] = p
}
