package orgmode

import (
	"darkness/internals"
	"regexp"
)

var (
	SourceCodeRegexp     = regexp.MustCompile(`(?s)#\+begin_src ?([[:print:]]+)?`)
	LinkRegexp           = internals.LinkRegexp
	AttentionBlockRegexp = regexp.MustCompile(`^(WARNING|NOTE|TIP):\s*(.+)`)
	UnorderedListRegexp  = regexp.MustCompile(`(?mU)- (.+) ∆`)
)
