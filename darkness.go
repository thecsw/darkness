package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"darkness/emilia"
	"darkness/html"
	"darkness/orgmode"
)

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	initFlags()
	flag.Parse()
	emilia.InitDarkness(darknessToml)
	html.InitConstantTags()

	if len(os.Args) == 1 {
		fmt.Println("hm? I didn't get anything, see -help")
		return
	}

	if os.Args[1] == "build" {
		build()
		return
	}
	if os.Args[1] == "lalatina" {
		fmt.Println("DONT CALL ME THAT >.<")
		return
	}
}

func build() {
	orgfiles, err := findFilesByExt(workDir, sourceExt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d files\n", len(orgfiles))
	fmt.Printf("Working on them... ")
	toSave := make(map[string]string)
	for _, file := range orgfiles {
		page := orgmode.ParseFile(workDir, file)
		// Ask emilia to work over the page a little
		emilia.EnrichHeadings(page)
		emilia.ResolveFootnotes(page)
		htmlFilename := strings.Replace(filepath.Base(file), sourceExt, targetExt, 1)
		targetFile := filepath.Join(filepath.Dir(file), htmlFilename)
		htmlPage := html.ExportPage(page)
		htmlPage = emilia.AddHolosceneTitles(file, htmlPage)
		toSave[targetFile] = htmlPage
	}
	fmt.Println("done")
	fmt.Printf("Flushing files... ")
	for k, v := range toSave {
		ioutil.WriteFile(k, []byte(v), 0644)
	}
	fmt.Println("done")
}

var (
	workDir      = "."
	darknessToml = "darkness.toml"
	sourceExt    = ".org"
	targetExt    = ".html"
)

func initFlags() {
	flag.StringVar(&workDir, "dir", ".", "where do I look for files")
	flag.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	flag.StringVar(&sourceExt, "source", ".org", "source extension [default: .org]")
	flag.StringVar(&targetExt, "target", ".html", "target extension [default: .html]")
}

func findFilesByExt(dir, ext string) ([]string, error) {
	files := make([]string, 0, 32)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
