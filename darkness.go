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
	emilia.InitDarkness(darknessToml)

	orgfiles, err := findFilesByExt(workDir, sourceExt)
	if err != nil {
		panic(err)
	}
	//litter.Dump(orgfiles)
	fmt.Printf("Found %d files\n", len(orgfiles))

	fmt.Println("Working on them...")
	for _, file := range orgfiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		lines := strings.Split(string(data), "\n")
		page := orgmode.Parse(lines)
		//		litter.Dump(page)
		page.URL = emilia.Config.URL + strings.TrimPrefix(filepath.Dir(file), workDir) + "/"
		targetFile := filepath.Join(filepath.Dir(file),
			strings.Replace(filepath.Base(file), sourceExt, targetExt, 1))
		//fmt.Println(targetFile)

		finalPage := html.ExportPage(page)
		finalPage = emilia.AddHolosceneTitles(file, finalPage)

		ioutil.WriteFile(targetFile, []byte(finalPage), 0644)
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
