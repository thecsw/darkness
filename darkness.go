package main

import (
	"fmt"
	"os"

	"github.com/thecsw/darkness/ichika"
)

// main is the entry point for the program.
func main() {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()

	// Darkness needs something, if nothing given, then show help.
	if len(os.Args) < 2 {
		ichika.HelpCommandFunc()
		return
	}

	// Find the supplied command or show a snarky error message otherwise.
	commandFunc, ok := ichika.CommandFuncs[ichika.DarknessCommand(os.Args[1])]
	if !ok {
		fmt.Println("command not found?")
		fmt.Println("see help, you pathetic excuse of a man")
		return
	}

	// Execute the retrieved command from Ichika.
	commandFunc()
}
