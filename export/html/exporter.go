package html

import (
	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/yunyun"
)

type ExporterHTML struct {
	Conf alpha.DarknessConfig
}

type state struct {
	// page is the source data that will be used for HTML building.
	page *yunyun.Page
	// currentContent is the pointer to the current `Content` object that is being processed.
	currentContent *yunyun.Content
	// contentFunctions is dictionary of rules to execute on content types.
	contentFunctions []func(*yunyun.Content) string
	// currentContentIndex is the index of the content that exporter is currently working on.
	currentContentIndex int
	// currentDiv is used as a state variable for internal processing.
	currentDiv divType
	// inHeading is used as a state variable for internal processing.
	inHeading bool
	// inWriting is used as a state variable for internal processing.
	inWriting bool
	conf      alpha.DarknessConfig
}
