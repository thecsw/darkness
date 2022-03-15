package orgmode

import (
	"regexp"

	"github.com/thecsw/darkness/internals"
)

var (
	SourceCodeRegexp     = regexp.MustCompile(`(?s)#\+begin_src ?([[:print:]]+)?`)
	LinkRegexp           = internals.LinkRegexp
	AttentionBlockRegexp = regexp.MustCompile(`^(WARNING|NOTE|TIP|IMPORTANT):\s*(.+)`)
	UnorderedListRegexp  = regexp.MustCompile(`(?mU)- (.+) ∆`)
	HeadingRegexp        = regexp.MustCompile(`(?m)^(\*\*\*\*\*|\*\*\*\*|\*\*\*|\*\*|\*\s+)`)
	ImageExtRegexp       = regexp.MustCompile(`\.(png|gif|jpg|jpeg|svg|webp)$`)
	AudioFileExtRegexp   = regexp.MustCompile(`\.(mp3|flac|midi)$`)
)
