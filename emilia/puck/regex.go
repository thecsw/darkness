package puck

import "regexp"

var (
	// HEregex is a regex for matching Holoscene times
	HEregex = regexp.MustCompile(`(\d+);\s*(\d+)\s*H.E.`)
)
