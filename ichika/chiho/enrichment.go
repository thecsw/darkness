package chiho

import (
	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/narumi"
	"github.com/thecsw/darkness/v3/yunyun"
)

// EnrichPage enriches the page with the following:
// - Resolved comments
// - Enriched headings
// - Footnotes
// - Math support
// - Source code trimmed left whitespace
// - Syntax highlighting
// - Lazy galleries
func EnrichPage(conf *alpha.DarknessConfig, page *yunyun.Page) *yunyun.Page {
	return page.Options(
		narumi.WithResolvedComments(),
		narumi.WithEnrichedHeadings(),
		narumi.WithFootnotes(),
		narumi.WithMathSupport(),
		narumi.WithSourceCodeTrimmedLeftWhitespace(),
		narumi.WithSyntaxHighlighting(conf),
		narumi.WithLazyGalleries(conf),
	)
}
