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
	"unicode"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/html"
	"github.com/thecsw/darkness/orgmode"
)

// main is the entry point for the program
func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	//defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()

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
  new - create darkness.toml in the current directory
  file - build a single input file and output to stdout
  build - build the entire directory
  megumin - blow up the directory
  lalatina - DO NOT
  aqua - ...

Don't hold back! You have no choice!`)
}

// oneFile builds a single file
func oneFile() {
	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	fileCmd.StringVar(&filename, "i", "index.org", "file on input")
	fileCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	fileCmd.Parse(os.Args[2:])
	emilia.InitDarkness(darknessToml)
	fmt.Println(orgToHTML(filename))
}

// build builds the entire directory
func build() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	buildCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	buildCmd.StringVar(&sourceExt, "source", ".org", "source extension")
	buildCmd.StringVar(&targetExt, "target", ".html", "target extension")
	buildCmd.Parse(os.Args[2:])

	emilia.InitDarkness(darknessToml)
	html.InitConstantTags()
	fmt.Printf("Looking for files... ")
	start := time.Now()
	orgfiles, err := findFilesByExt(workDir, sourceExt)
	if err != nil {
		fmt.Printf("failed to find files by extension %s: %s", sourceExt, err.Error())
		os.Exit(1)
	}
	fmt.Printf("found %d in %d ms\n", len(orgfiles), time.Since(start).Milliseconds())
	fmt.Printf("Building and flushing... ")
	start = time.Now()
	for _, file := range orgfiles {
		ioutil.WriteFile(getTarget(file), []byte(orgToHTML(file)), 0644)
	}
	fmt.Printf("done in %d ms\n", time.Since(start).Milliseconds())
	fmt.Println("farewell")
}

func aqua() {
	// KAZUMAAA-SAAAAAAAAN
}

// megumin blows up the directory
func megumin() {
	explosionCmd := flag.NewFlagSet("megumin", flag.ExitOnError)
	explosionCmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	explosionCmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	explosionCmd.StringVar(&sourceExt, "source", ".org", "source extension")
	explosionCmd.StringVar(&targetExt, "target", ".html", "target extension")
	explosionCmd.Parse(os.Args[2:])
	emilia.InitDarkness(darknessToml)

	orgfiles, err := findFilesByExt(workDir, sourceExt)
	if err != nil {
		fmt.Printf("failed to find files by extension %s: %s", sourceExt, err.Error())
		os.Exit(1)
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

// orgToHTML converts an org file to html
func orgToHTML(file string) string {
	page := orgmode.ParseFile(workDir, file)
	// Debug line to show the current page
	//litter.Dump(page)
	// Ask emilia to work over the page a little
	emilia.ResolveComments(page)
	emilia.EnrichHeadings(page)
	emilia.ResolveFootnotes(page)
	emilia.AddMathSupport(page)
	emilia.SourceCodeTrimLeftWhitespace(page)
	htmlPage := html.ExportPage(page)
	htmlPage = emilia.AddHolosceneTitles(htmlPage)
	return htmlPage
}

var (
	// workDir is the directory to look for files
	workDir = "."
	// darknessToml is the location of darkness.toml
	darknessToml = "darkness.toml"
	// sourceExt is the extension to look for
	sourceExt = ".org"
	// targetExt is the extension to output
	targetExt = ".html"
	// filename is the file to build
	filename = "index.org"
)

// findFilesByExt finds all files with a given extension
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

// getTarget returns the target file name
func getTarget(file string) string {
	htmlFilename := strings.Replace(filepath.Base(file), sourceExt, targetExt, 1)
	return filepath.Join(filepath.Dir(file), htmlFilename)
}

// delayedLinesPrint prints lines with a delay
func delayedLinesPrint(lines []string) {
	for _, line := range lines {
		time.Sleep(200 * time.Millisecond)
		delayedSentencePrint(line)
		time.Sleep(900 * time.Millisecond)
		fmt.Printf("\n")
	}
}

// delayedSentencePrint prints a sentence with a delay
func delayedSentencePrint(line string) {
	for _, c := range line {
		fmt.Printf("%c", c)
		time.Sleep(60 * time.Millisecond)
		if unicode.IsPunct(c) {
			time.Sleep(400 * time.Millisecond)
		}
	}
}
