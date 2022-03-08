package orgmode

import "regexp"

var (
	listPrefixes = []string{"* ", "- ", "+ "}

	LinkRegexp = regexp.MustCompile(`\[\[([^][]+)\]\[([^][]+)\]\]`)
	BoldText   = regexp.MustCompile(`(^| )\*([^* ][^*]+[^* ]|[^*])\*([^\w.]|$)`)
	ItalicText = regexp.MustCompile(`(^| )/([^/ ][^/]+[^/ ]|[^/])/([^\w.])`)

	SourceCodeRegexp = regexp.MustCompile(`#\+begin_src ?(.+)?`)
)
