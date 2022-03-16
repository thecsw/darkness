package html

var (
	// styleTagsProcessed is the processed style tags
	styleTagsProcessed string
)

// InitConstantTags initializes the constant tags
func InitConstantTags() {
	styleTagsProcessed = styleTags()
}
