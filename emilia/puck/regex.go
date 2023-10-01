package puck

import "regexp"

// HEregex is a regex for matching Holoscene times
var HEregex = regexp.MustCompile(`(\d+);\s*(\d+)\s*H.E.`)
