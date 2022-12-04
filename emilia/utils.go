package emilia

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// InputFilenameToOutput converts input filename to the filename to write.
func InputFilenameToOutput(file string) string {
	outputFilename := strings.Replace(filepath.Base(file),
		Config.Project.Input, Config.Project.Output, 1)
	return filepath.Join(filepath.Dir(file), outputFilename)
}

// InputToOutput converts a single input file to its output.
func InputToOutput(file string) string {
	data, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		panic(err)
	}
	page := ParserBuilder.BuildParser(file, string(data)).Parse()
	return EnrichAndExportPage(EnrichPage(page))
}

// EnrichAndExportPage automatically applies all the emilia enhancements
// and converts Page into an html document.
func EnrichAndExportPage(page *yunyun.Page) string {
	result := AddHolosceneTitles(
		ExporterBuilder.BuildExporter(EnrichPage(page)).Export(),
		func() int {
			if strings.HasSuffix(page.URL, "quotes") {
				return -1
			}
			return 1
		}())
	return result
}

// EnrichPage applies common emilia enhancements.
func EnrichPage(page *yunyun.Page) *yunyun.Page {
	return page.Options(
		WithResolvedComments(),
		WithEnrichedHeadings(),
		WithFootnotes(),
		WithMathSupport(),
		WithSourceCodeTrimmedLeftWhitespace(),
		WithSyntaxHighlighting(),
		WithLazyGalleries(),
	)
}

// GalleryPreview takes an original image's path and returns
// the preview path of it.
func GalleryPreview(img string) string {
	ext := filepath.Ext(img)
	return strings.TrimSuffix(img, ext) + "_preview" + ext
}
