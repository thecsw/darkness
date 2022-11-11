package html

import (
	"fmt"

	"github.com/thecsw/darkness/yunyun"
)

func init() {
	if yunyun.TypeContent(len(divTypes)) !=
		yunyun.TypeShouldBeLastDoNotTouch {
		panic(fmt.Sprintf("len(html.divTypes) should be %d but it is %d",
			yunyun.TypeShouldBeLastDoNotTouch, len(divTypes)))
	}
}

func init() {
	e := ExporterHTML{}
	e.SetPage(nil)
	if yunyun.TypeContent(len(e.contentFunctions)) !=
		yunyun.TypeShouldBeLastDoNotTouch {
		panic(fmt.Sprintf("len(html.contentFunctions) should be %d but it is %d",
			yunyun.TypeShouldBeLastDoNotTouch, len(e.contentFunctions)))
	}
}
