package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/html"
	"github.com/thecsw/darkness/orgmode"
)

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	if len(os.Args) == 1 {
		help()
		return
	}

	switch os.Args[1] {
	case "file":
		oneFile()
	case "build":
		build()
	case "megumin":
		megumin()
	case "lalatina":
		fmt.Println("DONT CALL ME THAT (╥︣﹏᷅╥)")
	case "aqua":
		os.Exit(1)
	case "-h", "--help", "-help", "help":
		help()
	default:
		fmt.Println("see help, you pathetic excuse of a man")
	}
}

func help() {
	fmt.Println(`My name is Darkness.
My calling is that of a crusader.
Do Shometing Gwazy!

Here are the commands you can use, -help is supported:
  file - build a single input file and output to stdout
  build - build the entire directory
  megumin - blow up the directory
  lalatina - DO NOT
  aqua - ...

Don't hold back! You have no choice!`)
}

func oneFile() {
	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	fileCmd.StringVar(&filename, "i", "index.org", "file on input")
	fileCmd.Parse(os.Args[2:])
	fmt.Println(orgToHTML(filename))
}

func build() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	buildCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	buildCmd.StringVar(&sourceExt, "source", ".org", "source extension")
	buildCmd.StringVar(&targetExt, "target", ".html", "target extension")
	buildCmd.Parse(os.Args[2:])

	emilia.InitDarkness(darknessToml)
	html.InitConstantTags()
	orgfiles, err := findFilesByExt(workDir, sourceExt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d files\n", len(orgfiles))
	fmt.Printf("Working on them... ")
	toSave := make(map[string]string)
	for _, file := range orgfiles {
		toSave[getTarget(file)] = orgToHTML(file)
	}
	fmt.Println("done")
	fmt.Printf("Flushing files... ")
	for k, v := range toSave {
		ioutil.WriteFile(k, []byte(v), 0644)
	}
	fmt.Println("done")
}

func aqua() {
	// KAZUMAAA-SAAAAAAAAN
}

func megumin() {
	explosionCmd := flag.NewFlagSet("megumin", flag.ExitOnError)
	explosionCmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	explosionCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	explosionCmd.StringVar(&sourceExt, "source", ".org", "source extension")
	explosionCmd.StringVar(&targetExt, "target", ".html", "target extension")
	emilia.InitDarkness(darknessToml)

	orgfiles, err := findFilesByExt(workDir, sourceExt)
	if err != nil {
		panic(err)
	}
	delayedLinesPrint([]string{
		"Darker than black, darker than darkness, combine with my intense crimson.",
		"Time to wake up, descend to these borders and appear as an intangible distortion.",
		"Dance, dance, dance!",
		"May a destructive force flood my torrent of power, a destructive force like no other!",
		"Send all creation to its source!",
		"Come out of your abyss!",
		"Humanity knows no other more powerful offensive technique!",
		"It is the ultimate magical attack!",
		"Explosion!",
	})
	for _, orgfile := range orgfiles {
		toRemove := getTarget(orgfile)
		fmt.Println(toRemove, "went boom!")
		time.Sleep(50 * time.Millisecond)
		if err := os.Remove(toRemove); err != nil {
			fmt.Println(toRemove, "failed to blow up!!")
		}
	}
	delayedLinesPrint([]string{
		"Wahahahahaha!",
		"My name is Megumin, the number one mage of Axel!",
		"Come, you shall all become my experience points today!",
	})
}

func orgToHTML(file string) string {
	page := orgmode.ParseFile(workDir, file)
	// Debug line to show the current page
	//litter.Dump(page)
	// Ask emilia to work over the page a little
	emilia.EnrichHeadings(page)
	emilia.ResolveFootnotes(page)
	emilia.AddMathSupport(page)
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
			// Check if it is not excluded
			isExcluded := false
			for _, excludedPath := range emilia.Config.Website.Exclude {
				if strings.HasPrefix(path, excludedPath) {
					isExcluded = true
					break
				}
			}
			if !isExcluded {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

func getTarget(file string) string {
	htmlFilename := strings.Replace(filepath.Base(file), sourceExt, targetExt, 1)
	return filepath.Join(filepath.Dir(file), htmlFilename)
}

func delayedLinesPrint(lines []string) {
	for _, line := range lines {
		time.Sleep(200 * time.Millisecond)
		delayedSentencePrint(line)
		time.Sleep(900 * time.Millisecond)
		fmt.Printf("\n")
	}
}

func delayedSentencePrint(line string) {
	for _, c := range line {
		fmt.Printf("%c", c)
		time.Sleep(60 * time.Millisecond)
		if c == ',' {
			time.Sleep(400 * time.Millisecond)
		}
	}
}
