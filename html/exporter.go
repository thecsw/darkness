package html

import (
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
)

type ExporterHTML struct {
	Page             *internals.Page
	contentsNum      int
	currentDiv       divType
	inHeading        bool
	inWriting        bool
	contentFunctions []func(*internals.Content) string
}

type ExporterOption func(*ExporterHTML)

func InitPackage() {
	// Monkey patch the function if we're using the roman footnotes
	if emilia.Config.Website.RomanFootnotes {
		footnoteLabel = numberToRoman
	}
}

func NewExporterHTML(options ...ExporterOption) *ExporterHTML {
	e := &ExporterHTML{}
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

func WithPage(page *internals.Page) ExporterOption {
	return func(e *ExporterHTML) {
		e.Page = page
		e.contentsNum = len(page.Contents)
	}
}
