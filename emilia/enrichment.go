package emilia

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// EnrichAndExportPage automatically applies all the emilia enhancements
// and converts Page into an html document.
func EnrichAndExportPage(page *yunyun.Page) string {
	result := AddHolosceneTitles(
		ExporterBuilder.BuildExporter(EnrichPage(page)).Export(),
		func() int {
			if strings.HasSuffix(string(page.Location), "quotes") {
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

// InputFilenameToOutput converts input filename to the filename to write.
func InputFilenameToOutput(file yunyun.FullPathFile) string {
	outputFilename := strings.Replace(filepath.Base(string(file)),
		Config.Project.Input, Config.Project.Output, 1)
	return filepath.Join(filepath.Dir(string(file)), outputFilename)
}

// InputToOutput converts a single input file to its output.
func InputToOutput(file yunyun.FullPathFile) string {
	data, err := ioutil.ReadFile(filepath.Clean(string(file)))
	if err != nil {
		panic(err)
	}
	page := ParserBuilder.BuildParser(Pack(file, string(data))).Parse()
	return EnrichAndExportPage(EnrichPage(page))
}
