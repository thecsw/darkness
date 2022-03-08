package internals

import "regexp"

var (
	LinkRegexp     = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	BoldText       = regexp.MustCompile(`(^| )\*([^* ][^*]+[^* ]|[^*])\*([^\w]|$)`)
	ItalicText     = regexp.MustCompile(`(^| )/([^/ ][^/]+[^/ ]|[^/])/([^\w]|$)`)
	VerbatimText   = regexp.MustCompile(`(^| )=([^= ][^=]+[^= ]|[^=])=([^\w]|$)`)
	KeyboardRegexp = regexp.MustCompile(`kbd:\[([^][]+)\]`)
)
