package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thecsw/darkness/emilia"
)

var (
	//go:embed ishmael/ishmael.tar
	defaultDarknessTemplate []byte
)

// main is the entry point for the program
func main() {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()

	if len(os.Args) == 1 {
		help()
		return
	}

	switch os.Args[1] {
	case "new":
		newDarkness()
	case "file":
		oneFile()
	case "build":
		buildCommand()
	case "serve":
		serve()
	case "clean":
		isQuietMegumin = true
		megumin()
	case "megumin":
		megumin()
	case "tools":
		tools()
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

// help shows default darkness help message
func help() {
	fmt.Println(`My name is Darkness.
My calling is that of a crusader.
Do Shometing Gwazy!

If you don't have a darkness website yet, start with
creating it with new followed by the directory name

  $> darkness new axel

Here are the commands you can use, -help is supported:
  file - build a single input file and output to stdout
  build - build the entire directory
  megumin - blow up the directory
  lalatina - DO NOT
  aqua - ...

Don't hold back! You have no choice!`)
}

// newDarkness creates a default darkness config in the current directory
// if one already exists, nothing will happen, except a notice of that
func newDarkness() {
	if len(os.Args) != 3 {
		fmt.Println("you forgot to add a directory name after new")
		return
	}
	dirName := os.Args[2]
	f, err := os.Open(filepath.Clean(dirName))
	if err == nil {
		fmt.Println("this directory already exists, bailing")
		if err := f.Close(); err != nil {
			fmt.Printf("failed to close directory: %s\n", err.Error())
		}
		return
	}
	if err := os.Mkdir(dirName, os.FileMode(0777)); err != nil {
		fmt.Println("couldn't create a directory for you:", err.Error())
		return
	}
	// Create a buffer for the tar file so we can start untarring it
	tarBuf := bytes.NewReader(defaultDarknessTemplate)
	if err := emilia.Untar(tarBuf, dirName); err != nil {
		fmt.Println("failed at flushing the template files:", err.Error())
		return
	}
	// Done
	fmt.Printf("Done! Go to %s and start creating!\n(run darkness build in there)\n", dirName)
}

func aqua() {
	// KAZUMAAA-SAAAAAAAAN
}
