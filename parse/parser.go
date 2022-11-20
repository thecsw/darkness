package parse

import "github.com/thecsw/darkness/yunyun"

// Parser is an interface used to define import packages,
// which convert source data into a yunyun `Page`.
type Parser interface {
	// Parse returns `*yunyunPage`.
	Parse() *yunyun.Page
	// WithFilenameData returns a new `Parser` object with
	// filename and data set.
	WithFilenameData(string, string) Parser
}

// ParserMap stores mappings of extensions to their parsers.
var ParserMap = make(map[string]Parser)

// Register is called by parsers to register themselves.
func Register(ext string, p Parser) {
	ParserMap[ext] = p
}
