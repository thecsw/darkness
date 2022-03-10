package internals

import "regexp"

var (
	LinkRegexp     = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	BoldText       = regexp.MustCompile(`(?mU)(^|[ ()])\*(\S|\S\S|\S[\w\W]+\S)\*($|[ (),.!?;&])`)
	ItalicText     = regexp.MustCompile(`(?mU)(^|[ ()])/(\S|\S\S|\S[\w\W]+\S)/($|[ (),.!?;&])`)
	VerbatimText   = regexp.MustCompile(`(?mU)(^|[ ()])=(\S|\S\S|\S[\w\W]+\S)=($|[ (),.!?;&])`)
	KeyboardRegexp = regexp.MustCompile(`kbd:\[([^][]+)\]`)
)
