package html

import (
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/yunyun"
)

// ExporterHTML will consume a `Page` and emit final HTML representation of it.
type ExporterHTML struct {
	// page is the source data that will be used for HTML building.
	page *yunyun.Page
	// contentsNum is a pre-computed value of how many contents there are in Page.
	contentsNum int
	// currentDiv is used as a state variable for internal processing.
	currentDiv divType
	// inHeading is used as a state variable for internal processing.
	inHeading bool
	// inWriting is used as a state variable for internal processing.
	inWriting bool
	// contentFunctions is dictionary of rules to execute on content types.
	contentFunctions []func(*yunyun.Content) string
}

// SetPage sets the internal page and creates internal content mappers.
func (e *ExporterHTML) SetPage(page *yunyun.Page) {
	// Set up the content functions.
	e.contentFunctions = export.ContentBuilder(e)
	// If nil passed, it's used as a test and should be ignored.
	if page == nil {
		return
	}
	e.page = page
	e.contentsNum = len(page.Contents)
}
