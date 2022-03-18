package orgmode

import (
	"regexp"

	"github.com/thecsw/darkness/internals"
)

var (
	// SourceCodeRegexp is the regexp for matching source blocks
	SourceCodeRegexp = regexp.MustCompile(`(?s)#\+begin_src ?([[:print:]]+)?`)
	// LinkRegexp is the regexp for matching links
	LinkRegexp = internals.LinkRegexp
	// AttentionBlockRegexp is the regexp for matching attention blocks
	AttentionBlockRegexp = regexp.MustCompile(`^(WARNING|NOTE|TIP|IMPORTANT|CAUTION):\s*(.+)`)
	// UnorderedListRegexp is the regexp for matching unordered lists
	UnorderedListRegexp = regexp.MustCompile(`(?mU)- (.+) ∆`)
	// HeadingRegexp is the regexp for matching headlines
	HeadingRegexp = regexp.MustCompile(`(?m)^(\*\*\*\*\*|\*\*\*\*|\*\*\*|\*\*|\*\s+)`)
)
