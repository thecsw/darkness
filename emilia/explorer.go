package emilia

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

var (
	NumFoundFiles int32 = 0
)

// FindFilesByExt finds all files with a given extension.
func FindFilesByExt(inputFilenames chan<- string, workDir string, wg *sync.WaitGroup) {
	NumFoundFiles = 0
	if err := godirwalk.Walk(workDir, &godirwalk.Options{
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			fmt.Printf("Encountered an error while traversing %s: %s\n", osPathname, err.Error())
			return godirwalk.SkipNode
		},
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if filepath.Ext(osPathname) != Config.Project.Input {
				return nil
			}
			if (Config.Project.ExcludeEnabled && Config.Project.ExcludeRegex.MatchString(osPathname)) ||
				gana.First([]rune(de.Name())) == rune('.') {
				return filepath.SkipDir
			}
			wg.Add(1)
			relPath, err := filepath.Rel(workDir, osPathname)
			inputFilenames <- filepath.Join(workDir, relPath)
			atomic.AddInt32(&NumFoundFiles, 1)
			return err
		},
	}); err != nil {
		fmt.Printf("File traversal returned an error: %s\n", err.Error())
	}
	close(inputFilenames)
}

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
	)
}
