package orgmode

import (
	"darkness/internals"
	"regexp"
)

var (
	listPrefixes = []string{"* ", "- ", "+ "}

	LinkRegexp     = internals.LinkRegexp
	BoldText       = internals.BoldText
	ItalicText     = internals.ItalicText
	VerbatimText   = internals.VerbatimText
	KeyboardRegexp = internals.KeyboardRegexp

	SourceCodeRegexp = regexp.MustCompile(`(?s)#\+begin_src ?([[:print:]]+)?`)
)
