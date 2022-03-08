package orgmode

import "regexp"

var (
	listPrefixes = []string{"* ", "- ", "+ "}

	LinkRegexp     = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	BoldText       = regexp.MustCompile(`(^| )\*([^* ][^*]+[^* ]|[^*])\*([^\w]|$)`)
	ItalicText     = regexp.MustCompile(`(^| )/([^/ ][^/]+[^/ ]|[^/])/([^\w]|$)`)
	VerbatimText   = regexp.MustCompile(`(^| )=([^= ][^=]+[^= ]|[^=])=([^\w]|$)`)
	KeyboardRegexp = regexp.MustCompile(`kbd:\[([^][]+)\]`)

	SourceCodeRegexp = regexp.MustCompile(`(?s)#\+begin_src ?([[:print:]]+)?`)
)
