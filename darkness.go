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
	emilia.InitDarkness(darknessToml)
	html.InitConstantTags()

	if len(os.Args) == 1 {
		fmt.Println("hm? I didn't get anything, see -help")
		return
	}

	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	buildCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	buildCmd.StringVar(&sourceExt, "source", ".org", "source extension")
	buildCmd.StringVar(&targetExt, "target", ".html", "target extension")

	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	fileCmd.StringVar(&filename, "i", "index.org", "file on input")

	switch os.Args[1] {
	case "file":
		fileCmd.Parse(os.Args[2:])
		fmt.Println(orgToHTML(filename))
	case "build":
		buildCmd.Parse(os.Args[2:])
		build()
	case "lalatina":
		fmt.Println("DONT CALL ME THAT >.<")
	default:
		fmt.Println("see help, you pathetic excuse of a man")
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
		htmlFilename := strings.Replace(filepath.Base(file), sourceExt, targetExt, 1)
		targetFile := filepath.Join(filepath.Dir(file), htmlFilename)
		toSave[targetFile] = orgToHTML(file)
	}
	fmt.Println("done")
	fmt.Printf("Flushing files... ")
	for k, v := range toSave {
		ioutil.WriteFile(k, []byte(v), 0644)
	}
	fmt.Println("done")
}

func orgToHTML(file string) string {
	page := orgmode.ParseFile(workDir, file)
	// Debug line to show the current page
	//litter.Dump(page)
	// Ask emilia to work over the page a little
	emilia.EnrichHeadings(page)
	emilia.ResolveFootnotes(page)
	htmlPage := html.ExportPage(page)
	htmlPage = emilia.AddHolosceneTitles(file, htmlPage)
	return htmlPage
}

var (
	workDir      = "."
	darknessToml = "darkness.toml"
	sourceExt    = ".org"
	targetExt    = ".html"
	filename     = "index.org"
)

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
