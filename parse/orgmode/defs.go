package orgmode

import (
	"regexp"
)

const (
	commentPrefix      = "# "
	optionPrefix       = "#+"
	optionPrefixLen    = len(optionPrefix)
	optionDropCap      = "drop_cap"
	optionBeginSource  = "begin_src"
	optionEndSource    = "end_src"
	optionBeginExport  = "begin_export"
	optionEndExport    = "end_export"
	optionBeginQuote   = "begin_quote"
	optionEndQuote     = "end_quote"
	optionBeginCenter  = "begin_center"
	optionEndCenter    = "end_center"
	optionBeginDetails = "begin_details"
	optionEndDetails   = "end_details"
	optionBeginGallery = "begin_gallery"
	optionEndGallery   = "end_gallery"
	optionCaption      = "caption:"
	optionDate         = "date:"
	optionHtmlHead     = "html_head:"
	optionOptions      = "options:"
	optionAttributes   = "attr_darkness:"
	optionHtmlTags     = "html_tags:"
	optionAttrHtml     = "attr_html:"
	optionAuthor       = "author:"
	horizontalLine     = "-----"

	sectionLevelOne   = "* "
	sectionLevelTwo   = "** "
	sectionLevelThree = "*** "
	sectionLevelFour  = "**** "
	sectionLevelFive  = "***** "

	listSeparator    = string(rune(30))
	listSeparatorWS  = " " + listSeparator
	tableSeparator   = string(rune(29))
	tableSeparatorWS = " " + tableSeparator
)

var (
	surroundWithNewlines = []string{
		optionBeginQuote, optionEndQuote,
		optionBeginCenter, optionEndCenter,
		optionBeginDetails, optionEndDetails,
		optionBeginGallery, optionEndGallery,
	}
	// linkRegexp is the regexp for matching links
	linkRegexp *regexp.Regexp
	// attentionBlockRegexp is the regexp for matching attention blocks
	attentionBlockRegexp = regexp.MustCompile(`^(WARNING|NOTE|TIP|IMPORTANT|CAUTION):\s*(.+)`)
	// unorderedListRegexp is the regexp for matching unordered lists
	unorderedListRegexp = regexp.MustCompile(`(?mU)- (.+) ` + listSeparator)
	// headingRegexp is the regexp for matching headlines
	headingRegexp = regexp.MustCompile(`(?m)^(\*{1,6} )`)
)
