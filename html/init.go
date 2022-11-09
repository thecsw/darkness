package html

import (
	"fmt"

	"github.com/thecsw/darkness/internals"
)

func init() {
	if internals.TypeContent(len(divTypes)) !=
		internals.TypeShouldBeLastDoNotTouch {
		panic(fmt.Sprintf("len(html.divTypes) should be %d",
			internals.TypeShouldBeLastDoNotTouch))
	}
}

func init() {
	e := NewExporterHTML()
	if internals.TypeContent(len(e.contentFunctions)) !=
		internals.TypeShouldBeLastDoNotTouch {
		panic(fmt.Sprintf("len(html.contentFunctions) should be %d",
			internals.TypeShouldBeLastDoNotTouch))
	}
}
