package html

import (
	"fmt"

	"github.com/thecsw/darkness/yunyun"
)

func init() {
	if yunyun.TypeContent(len(divTypes)) !=
		yunyun.TypeShouldBeLastDoNotTouch {
		panic(fmt.Sprintf("len(html.divTypes) should be %d",
			yunyun.TypeShouldBeLastDoNotTouch))
	}
}

func init() {
	e := NewExporterHTML()
	if yunyun.TypeContent(len(e.contentFunctions)) !=
		yunyun.TypeShouldBeLastDoNotTouch {
		panic(fmt.Sprintf("len(html.contentFunctions) should be %d",
			yunyun.TypeShouldBeLastDoNotTouch))
	}
}
