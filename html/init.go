package html

var (
	styleTagsProcessed string
)

func InitConstantTags() {
	styleTagsProcessed = styleTags()
}
