package markdown

import (
	"regexp"
)

const (
	commentPrefix     = "// "
	optionBeginSource = "```"
	optionEndSource   = optionBeginSource
	horizontalLine    = "---"

	sectionLevelOne   = "# "
	sectionLevelTwo   = "## "
	sectionLevelThree = "### "
	sectionLevelFour  = "#### "
	sectionLevelFive  = "##### "

	listSeparator    = string(rune(30))
	listSeparatorWS  = " " + listSeparator
	tableSeparator   = string(rune(29))
	tableSeparatorWS = " " + tableSeparator
)

var (
	// linkRegexp is the regexp for matching links
	linkRegexp *regexp.Regexp
	// attentionBlockRegexp is the regexp for matching attention blocks
	attentionBlockRegexp = regexp.MustCompile(`^(WARNING|NOTE|TIP|IMPORTANT|CAUTION):\s*(.+)`)
	// unorderedListRegexp is the regexp for matching unordered lists
	unorderedListRegexp = regexp.MustCompile(`(?mU)- (.+) ` + listSeparator)
	// headingRegexp is the regexp for matching headlines
	headingRegexp = regexp.MustCompile(`(?m)^(#{1,5}\s+)`)
)
