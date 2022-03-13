package orgmode

import (
	"darkness/internals"
	"regexp"
)

var (
	SourceCodeRegexp    = regexp.MustCompile(`(?s)#\+begin_src ?([[:print:]]+)?`)
	LinkRegexp          = internals.LinkRegexp
	UnorderedListRegexp = regexp.MustCompile(`(?mU)- (.+) ∆`)
)
