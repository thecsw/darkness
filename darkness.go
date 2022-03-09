package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"darkness/emilia"
	"darkness/html"
	"darkness/orgmode"
)

const (
	workDir      = "sandyuraz"
	darknessToml = "darkness.toml"
	sourceExt    = ".org"
	targetExt    = ".html"
)

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	emilia.InitDarkness(darknessToml)
	html.InitConstantTags()

	orgfiles, err := findFilesByExt(workDir, sourceExt)
	if err != nil {
		panic(err)
	}
	//litter.Dump(orgfiles)
	fmt.Printf("Found %d files\n", len(orgfiles))

	fmt.Println("Working on them...")
	for _, file := range orgfiles {
		page := orgmode.ParseFile(workDir, file)
		targetFile := filepath.Join(filepath.Dir(file),
			strings.Replace(filepath.Base(file), sourceExt, targetExt, 1))
		htmlPage := html.ExportPage(page)
		htmlPage = emilia.AddHolosceneTitles(file, htmlPage)
		ioutil.WriteFile(targetFile, []byte(htmlPage), 0644)
	}
	fmt.Println("done")
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
