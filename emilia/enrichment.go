package emilia

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// EnrichExportPage automatically applies all the emilia enhancements
// and converts Page into an html document.
func EnrichExportPage(page *yunyun.Page) string {
	return ExporterBuilder.BuildExporter(EnrichPage(page)).Export()
}

// EnrichExportPageAsBufio is the same as `EnrichExportPage` but returns a
// bufio-based buffered reader.
func EnrichExportPageAsBufio(page *yunyun.Page) *bufio.Reader {
	return bufio.NewReader(bytes.NewBufferString(EnrichExportPage(page)))
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
	return strings.Replace(string(file), Config.Project.Input, Config.Project.Output, 1)
}

// InputToOutput converts a single input file to its output.
func InputToOutput(file yunyun.FullPathFile) string {
	data, err := ioutil.ReadFile(filepath.Clean(string(file)))
	if err != nil {
		panic(err)
	}
	page := ParserBuilder.BuildParser(Pack(file, string(data))).Parse()
	return EnrichExportPage(EnrichPage(page))
}
