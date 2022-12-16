package ichika

import (
	"fmt"
	"os"
	"sync"
	"time"
	"unicode"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
)

var (
	// if true, darkness cleans with no output
	isQuietMegumin = false
)

// MeguminCommandFunc blows up the directory.
func MeguminCommandFunc() {
	options := getEmiliaOptions(darknessFlagset(meguminCommand))
	options.Dev = true
	emilia.InitDarkness(options)
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
	removeOutputFiles()
	delayedLinesPrint([]string{
		"Wahahahahaha!",
		"My name is Megumin, the number one mage of Axel!",
		"Come, you shall all become my experience points today!",
	})
	// Also clean gallery previews
	removeGalleryFiles()
}

// CleanCommandFunc cleans the files like `megumin` but without any output (except for errors).
func CleanCommandFunc() {
	options := getEmiliaOptions(darknessFlagset(cleanCommand))
	options.Dev = true
	emilia.InitDarkness(options)
	isQuietMegumin = true
	removeOutputFiles()
	removeGalleryFiles()
}

// removeOutputFiles is the low-level command to be used when cleaning data.
func removeOutputFiles() {
	orgfiles := make(chan yunyun.FullPathFile, defaultNumOfWorkers)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go emilia.FindFilesByExt(orgfiles, emilia.Config.Project.Input, wg)
	for orgfile := range orgfiles {
		toRemove := emilia.InputFilenameToOutput(orgfile)
		if err := os.Remove(toRemove); err != nil && !os.IsNotExist(err) {
			fmt.Println(toRemove, "failed to blow up!!")
		}
		if !isQuietMegumin {
			fmt.Println(toRemove, "went boom!")
			time.Sleep(50 * time.Millisecond)
		}
		wg.Done()
	}
	wg.Done()
	wg.Wait()
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
