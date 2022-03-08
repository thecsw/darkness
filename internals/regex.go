package internals

import "regexp"

var (
	LinkRegexp = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	BoldText   = regexp.MustCompile(`(^| )\*([^* ][^*]+[^* ]|[^*])\*([^\w.]|$)`)
	ItalicText = regexp.MustCompile(`(^| )/([^/ ][^/]+[^/ ]|[^/])/([^\w.]|$)`)
)
