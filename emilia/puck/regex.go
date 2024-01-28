package puck

import "regexp"

// HEregex is a regex for matching Holoscene times
var (
	HEregex              = regexp.MustCompile(`(?P<day>\d+);\s*(?P<year>\d+)\s*H.E.\s*(?P<hour_minute>[0-9]{4})?`)
	HERegexSubmatchNames = HEregex.SubexpNames()
)
