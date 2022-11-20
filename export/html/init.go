package html

import (
	"fmt"

	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/yunyun"
)

var (
	// Make sure that this exporter implements `export.Exporter`.
	exporter                        = &ExporterHTML{}
	_        export.Exporter        = exporter
	_        export.ExporterBuilder = exporter
)

// This init makes sure there are no discrepancies in html defs.
func init() {
	if yunyun.TypeContent(len(divTypes)) !=
		yunyun.TypeShouldBeLastDoNotTouch {
		panic(fmt.Sprintf("len(html.divTypes) should be %d but it is %d",
			yunyun.TypeShouldBeLastDoNotTouch, len(divTypes)))
	}
}

// This init registers the exporter with the root module.
func init() {
	export.Register(puck.ExtensionHtml, exporter)
}
