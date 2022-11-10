package html

import (
	"github.com/thecsw/darkness/internals"
)

// ExporterHTML will consume a `Page` and emit final HTML representation of it.
type ExporterHTML struct {
	// Page is the source data that will be used for HTML building.
	Page *internals.Page
	// contentsNum is a pre-computed value of how many contents there are in Page.
	contentsNum int
	// currentDiv is used as a state variable for internal processing.
	currentDiv divType
	// inHeading is used as a state variable for internal processing.
	inHeading bool
	// inWriting is used as a state variable for internal processing.
	inWriting bool
	// contentFunctions is dictionary of rules to execute on content types.
	contentFunctions []func(*internals.Content) string
}

// ExporterOption defines options that can be passed to `NewExporterHTML`.
type ExporterOption func(*ExporterHTML)

// NewExporterHTML creates `ExporterHTML`, should be run after `emilia` consumed darkness.toml
func NewExporterHTML(options ...ExporterOption) *ExporterHTML {
	e := &ExporterHTML{}
	// Run given `ExporterOption` options.
	for _, option := range options {
		option(e)
	}

	// contentFunctions is a map of functions that process content
	e.contentFunctions = []func(*internals.Content) string{
		e.headings, e.paragraph, e.list, e.listNumbered, e.link,
		e.sourceCode, e.rawHTML, e.horizontalLine, e.attentionBlock,
		e.table, e.details,
	}
	return e
}

// WithPage is an option that gives `ExporterHTML` the `Page` to consume
// and run some processes to pre-compute frequently used values.
func WithPage(page *internals.Page) ExporterOption {
	return func(e *ExporterHTML) {
		e.Page = page
		e.contentsNum = len(page.Contents)
	}
}
