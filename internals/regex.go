package internals

import "regexp"

var (
	LinkRegexp     = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	BoldText       = regexp.MustCompile(`(?mU)(^|[ ()_%<>])\*(\S|\S\S|\S.+\S)\*($|[ (),.!?;&_%<>])`)
	ItalicText     = regexp.MustCompile(`(?mU)(^|[ ()_%<>])/(\S|\S\S|\S.+\S)/($|[ (),.!?;&_%<>])`)
	VerbatimText   = regexp.MustCompile(`(?mU)(^|[ ()_%<>])=(\S|\S\S|\S.+\S)=($|[ (),.!?;&_%<>])`)
	KeyboardRegexp = regexp.MustCompile(`kbd:\[([^][]+)\]`)

	FootnoteRegexp               = regexp.MustCompile(`\[fn::([^][]+)\]`)
	FootnotePostProcessingRegexp = regexp.MustCompile(`!(\d+)!`)
)
