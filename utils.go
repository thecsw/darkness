package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/karrick/godirwalk"
	"github.com/sanity-io/litter"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/html"
	"github.com/thecsw/darkness/internals"
	"github.com/thecsw/darkness/orgmode"
)

// bundle is a struct that hold filename and contents -- used for
// reading files and passing context or writing them too.
type bundle struct {
	File string
	Data string
}

// findFilesByExt finds all files with a given extension
func findFilesByExt(orgfiles chan<- string, wg *sync.WaitGroup) {
	godirwalk.Walk(workDir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if filepath.Ext(osPathname) != emilia.Config.Project.Input {
				return nil
			}
			if emilia.Config.Project.ExcludeRegex.MatchString(osPathname) ||
				internals.First([]rune(de.Name())) == rune('.') {
				return filepath.SkipDir
			}
			wg.Add(1)
			relPath, err := filepath.Rel(workDir, osPathname)
			orgfiles <- relPath
			return err
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	close(orgfiles)
}

// getTarget returns the target file name
func getTarget(file string) string {
	htmlFilename := strings.Replace(filepath.Base(file),
		emilia.Config.Project.Input, emilia.Config.Project.Output, 1)
	return filepath.Join(filepath.Dir(file), htmlFilename)
}

// orgToHTML converts an org file to html
func orgToHTML(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	page := orgmode.Parse(string(data), file)
	litter.Dump(*page)
	return exportAndEnrich(page)
}

// exportAndEnrich automatically applies all the emilia enhancements
// and converts Page into an html document.
func exportAndEnrich(page *internals.Page) string {
	result := html.NewExporterHTML(html.WithPage(applyEmilia(page))).Export()
	result = emilia.AddHolosceneTitles(result, func() int {
		if strings.HasSuffix(page.URL, "quotes") {
			return -1
		}
		return 1
	}())
	return result
}

// applyEmilia applies common emilia enhancements.
func applyEmilia(page *internals.Page) *internals.Page {
	return page.Options(
		emilia.WithResolvedComments(),
		emilia.WithEnrichedHeadings(),
		emilia.WithFootnotes(),
		emilia.WithMathSupport(),
		emilia.WithSourceCodeTrimmedLeftWhitespace(),
		emilia.WithSyntaxHighlighting(),
	)
}
