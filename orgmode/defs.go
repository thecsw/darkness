package orgmode

import (
	"regexp"

	"github.com/thecsw/darkness/internals"
)

const (
	CommentPrefix      = "# "
	OptionPrefix       = "#+"
	OptionDropCap      = "drop_cap"
	OptionBeginSource  = "begin_src"
	OptionEndSource    = "end_src"
	OptionBeginExport  = "begin_export"
	OptionEndExport    = "end_export"
	OptionBeginQuote   = "begin_quote"
	OptionEndQuote     = "end_quote"
	OptionBeginCenter  = "begin_center"
	OptionEndCenter    = "end_center"
	OptionBeginDetails = "begin_details"
	OptionEndDetails   = "end_details"
	OptionCaption      = "caption"
	OptionDate         = "date"
	HorizontalLine     = "---"

	SectionLevelOne   = "* "
	SectionLevelTwo   = "** "
	SectionLevelThree = "*** "
	SectionLevelFour  = "**** "
	SectionLevelFive  = "***** "
	listSeparator    = string(rune(30))
	listSeparatorWS  = " " + listSeparator
	tableSeparator   = string(rune(29))
	tableSeparatorWS = " " + tableSeparator
)

var (
	SurroundWithNewlines = []string{
		OptionBeginQuote, OptionEndQuote,
		OptionBeginCenter, OptionEndCenter,
		OptionBeginDetails, OptionEndDetails,
	}
	// SourceCodeRegexp is the regexp for matching source blocks
	SourceCodeRegexp = regexp.MustCompile(`(?s)#\+begin_src ?([[:print:]]+)?`)
	// DetailsRegexp is the regexp for matching details
	DetailsRegexp = regexp.MustCompile(`(?s)#\+begin_details ?([[:print:]]+)?`)
	// LinkRegexp is the regexp for matching links
	LinkRegexp = internals.LinkRegexp
	// AttentionBlockRegexp is the regexp for matching attention blocks
	AttentionBlockRegexp = regexp.MustCompile(`^(WARNING|NOTE|TIP|IMPORTANT|CAUTION):\s*(.+)`)
	// UnorderedListRegexp is the regexp for matching unordered lists
	UnorderedListRegexp = regexp.MustCompile(`(?mU)- (.+) ∆`)
	// HeadingRegexp is the regexp for matching headlines
	HeadingRegexp = regexp.MustCompile(`(?m)^(\*\*\*\*\*|\*\*\*\*|\*\*\*|\*\*|\*\s+)`)
)
