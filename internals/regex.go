package internals

import "regexp"

var (
	LinkRegexp = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	// URLRegexp is yoinked from https://ihateregex.io/expr/url/
	URLRegexp      = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*))`)
	BoldText       = regexp.MustCompile(`(?mU)(^|[ ()_%<>])\*(\S|\S\S|\S.+\S)\*($|[ (),.!?;&_%<>])`)
	ItalicText     = regexp.MustCompile(`(?mU)(^|[ ()_%<>])/(\S|\S\S|\S.+\S)/($|[ (),.!?;&_%<>])`)
	VerbatimText   = regexp.MustCompile(`(?mU)(^|[ ()_%<>])=(\S|\S\S|\S.+\S)=($|[ (),.!?;&_%<>])`)
	KeyboardRegexp = regexp.MustCompile(`kbd:\[([^][]+)\]`)
	MathRegexp     = regexp.MustCompile(`(?mU)\$(.+)\$`)

	FootnoteRegexp               = regexp.MustCompile(`\[fn::([^][]+)\]`)
	FootnotePostProcessingRegexp = regexp.MustCompile(`!(\d+)!`)
)
