package ichika

import (
	"fmt"
	"os"
	"time"
	"unicode"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/ichika/hizuru"
	"github.com/thecsw/darkness/yunyun"
)

// if true, darkness cleans with no output
var isQuietMegumin = false

// MeguminCommandFunc blows up the directory.
func MeguminCommandFunc() {
	options := getAlphaOptions(darknessFlagset(meguminCommand))
	options.Dev = true
	conf := alpha.BuildConfig(options)
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
	removeOutputFiles(conf)
	delayedLinesPrint([]string{
		"Wahahahahaha!",
		"My name is Megumin, the number one mage of Axel!",
		"Come, you shall all become my experience points today!",
	})
}

// CleanCommandFunc cleans the files like `megumin` but without any output (except for errors).
func CleanCommandFunc() {
	options := getAlphaOptions(darknessFlagset(cleanCommand))
	options.Dev = true
	conf := alpha.BuildConfig(options)
	isQuietMegumin = true
	removeOutputFiles(conf)
}

// removeOutputFiles is the low-level command to be used when cleaning data.
func removeOutputFiles(conf *alpha.DarknessConfig) {
	inputFilenames := hizuru.FindFilesByExtSimple(conf)
	for _, inputFilename := range inputFilenames {
		toRemove := conf.Project.InputFilenameToOutput(inputFilename)
		toPrint := conf.Runtime.WorkDir.Rel(yunyun.FullPathFile(toRemove))
		if err := os.Remove(toRemove); err != nil && !os.IsNotExist(err) {
			fmt.Println(toPrint, "failed to blow up!!")
		}
		if !isQuietMegumin {
			fmt.Println(toPrint, "went boom!")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// delayedLinesPrint prints lines with a delay.
func delayedLinesPrint(lines []string) {
	if isQuietMegumin {
		return
	}
	for _, line := range lines {
		time.Sleep(200 * time.Millisecond)
		delayedSentencePrint(line)
		time.Sleep(900 * time.Millisecond)
		fmt.Printf("\n")
	}
}

// delayedSentencePrint prints a sentence with a delay.
func delayedSentencePrint(line string) {
	if isQuietMegumin {
		return
	}
	for _, c := range line {
		fmt.Printf("%c", c)
		time.Sleep(60 * time.Millisecond)
		if unicode.IsPunct(c) {
			time.Sleep(400 * time.Millisecond)
		}
	}
}
