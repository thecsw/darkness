package orgmode

import (
	"regexp"

	"github.com/thecsw/darkness/internals"
)

var (
	SourceCodeRegexp     = regexp.MustCompile(`(?s)#\+begin_src ?([[:print:]]+)?`)
	LinkRegexp           = internals.LinkRegexp
	AttentionBlockRegexp = regexp.MustCompile(`^(WARNING|NOTE|TIP):\s*(.+)`)
	UnorderedListRegexp  = regexp.MustCompile(`(?mU)- (.+) ∆`)
	HeadingRegexp        = regexp.MustCompile(`^(\*{1,5})[ ]`)
)
