package alpha

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// // EnrichExportPage automatically applies all the emilia enhancements
// // and converts Page into an html document.
//
//	func (conf RuntimeConfig) EnrichExportPage(page *yunyun.Page) string {
//		return conf.ExporterBuilder.BuildExporter(EnrichPage(page)).Export()
//	}
//
// // EnrichExportPageAsBufio is the same as `EnrichExportPage` but returns a
// // bufio-based buffered reader.
//
//	func (conf RuntimeConfig) EnrichExportPageAsBufio(page *yunyun.Page) *bufio.Reader {
//		return bufio.NewReader(bytes.NewBufferString(conf.EnrichExportPage(page)))
//	}
//
// // EnrichPage applies common emilia enhancements.
//
//	func EnrichPage(page *yunyun.Page) *yunyun.Page {
//		defer puck.Stopwatch("Enriched", "page", page.File).Record()
//		return page
//		//return page.Options(
//		//	WithResolvedComments(),
//		//	WithEnrichedHeadings(),
//		//	WithFootnotes(),
//		//	WithMathSupport(),
//		//	WithSourceCodeTrimmedLeftWhitespace(),
//		//	WithSyntaxHighlighting(),
//		//	WithLazyGalleries(),
//		//)
//	}
//
// InputFilenameToOutput converts input filename to the filename to write.
func (conf *DarknessConfig) InputFilenameToOutput(file yunyun.FullPathFile) string {
	return strings.Replace(string(file), conf.Project.Input, conf.Project.Output, 1)
}

//
//// InputToOutput converts a single input file to its output.
//func (conf *DarknessConfig) InputToOutput(file yunyun.FullPathFile) string {
//	data := rei.Must(os.ReadFile(filepath.Clean(string(file))))
//	a, b := conf.Runtime.WorkDir.PackRel(file, string(data))
//	page := conf.Runtime.ParserBuilder.BuildParser(*conf, a, b).Parse()
//	return conf.Runtime.EnrichExportPage(EnrichPage(page))
//}
