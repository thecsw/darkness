package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"

	_ "github.com/thecsw/darkness/export/html"
	_ "github.com/thecsw/darkness/export/template"
	_ "github.com/thecsw/darkness/parse/orgmode"
	_ "github.com/thecsw/darkness/parse/template"
)

// findFilesByExt finds all files with a given extension
func findFilesByExt(inputFilenames chan<- string, wg *sync.WaitGroup) {
	if err := godirwalk.Walk(workDir, &godirwalk.Options{
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			fmt.Printf("Encountered an error while traversing %s: %s\n", osPathname, err.Error())
			return godirwalk.SkipNode
		},
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if filepath.Ext(osPathname) != emilia.Config.Project.Input {
				return nil
			}
			if emilia.Config.Project.ExcludeRegex.MatchString(osPathname) || gana.First([]rune(de.Name())) == rune('.') {
				return filepath.SkipDir
			}
			wg.Add(1)
			relPath, err := filepath.Rel(workDir, osPathname)
			inputFilenames <- filepath.Join(workDir, relPath)
			return err
		},
	}); err != nil {
		fmt.Printf("File traversal returned an error: %s\n", err.Error())
	}
	close(inputFilenames)
}

// getTarget returns the target file name
func getTarget(file string) string {
	htmlFilename := strings.Replace(filepath.Base(file),
		emilia.Config.Project.Input, emilia.Config.Project.Output, 1)
	return filepath.Join(filepath.Dir(file), htmlFilename)
}

// inputToOutput converts a single input file to its output.
func inputToOutput(file string) string {
	data, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		panic(err)
	}
	page := emilia.ParserBuilder.BuildParser(file, string(data)).Parse()
	return exportAndEnrich(applyEmilia(page))
}

// exportAndEnrich automatically applies all the emilia enhancements
// and converts Page into an html document.
func exportAndEnrich(page *yunyun.Page) string {
	result := emilia.AddHolosceneTitles(
		emilia.ExporterBuilder.BuildExporter(applyEmilia(page)).Export(),
		func() int {
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
