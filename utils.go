package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/export/html"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/parse/orgmode"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// bundle is a struct that hold filename and contents -- used for
// reading files and passing context or writing them too.
type bundle struct {
	File string
	Data string
}

// findFilesByExt finds all files with a given extension
func findFilesByExt(inputFilenames chan<- string, wg *sync.WaitGroup) {
	godirwalk.Walk(workDir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if filepath.Ext(osPathname) != emilia.Config.Project.Input {
				return nil
			}
			if emilia.Config.Project.ExcludeRegex.MatchString(osPathname) ||
				gana.First([]rune(de.Name())) == rune('.') {
				return filepath.SkipDir
			}
			wg.Add(1)
			relPath, err := filepath.Rel(workDir, osPathname)
			inputFilenames <- relPath
			return err
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	close(inputFilenames)
}

// getTarget returns the target file name
func getTarget(file string) string {
	htmlFilename := strings.Replace(filepath.Base(file),
		emilia.Config.Project.Input, emilia.Config.Project.Output, 1)
	return filepath.Join(filepath.Dir(file), htmlFilename)
}

// inputToOutput converts an org file to html
func inputToOutput(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	page := getParser().Parse(string(data), file)
	return exportAndEnrich(applyEmilia(page))
}

// exportAndEnrich automatically applies all the emilia enhancements
// and converts Page into an html document.
func exportAndEnrich(page *yunyun.Page) string {
	exporter := getExporter()
	exporter.SetPage(applyEmilia(page))
	result := emilia.AddHolosceneTitles(exporter.Export(), func() int {
		if strings.HasSuffix(page.URL, "quotes") {
			return -1
		}
		return 1
	}())
	return result
}

// applyEmilia applies common emilia enhancements.
func applyEmilia(page *yunyun.Page) *yunyun.Page {
	return page.Options(
		emilia.WithResolvedComments(),
		emilia.WithEnrichedHeadings(),
		emilia.WithFootnotes(),
		emilia.WithMathSupport(),
		emilia.WithSourceCodeTrimmedLeftWhitespace(),
		emilia.WithSyntaxHighlighting(),
	)
}

const (
	parserDefaultKey   = ".org"
	exporterDefaultKey = ".html"
)

var (
	parserMap = map[string]parse.Parser{
		".org": &orgmode.ParserOrgmode{},
	}
	exporterMap = map[string]export.Exporter{
		".html": &html.ExporterHTML{},
	}
)

func getParser() parse.Parser {
	if v, ok := parserMap[emilia.Config.Project.Input]; ok {
		return v
	}
	return parserMap[parserDefaultKey]
}

func getExporter() export.Exporter {
	if v, ok := exporterMap[emilia.Config.Project.Output]; ok {
		return v
	}
	return exporterMap[exporterDefaultKey]
}
