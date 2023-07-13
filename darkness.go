package main

import (
	"fmt"
	"os"

	"github.com/pkg/profile"
	"github.com/thecsw/darkness/ichika"
)

const (
	PROFILE_CPU   = false
	PROFILE_MEM   = false
	PROFILE_CLOCK = false
)

// main is the entry point for the program.
func main() {
	// debug.SetGCPercent(-1)
	// debug.SetMemoryLimit(math.MaxInt64)

	if PROFILE_CPU {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	} else if PROFILE_MEM {
		defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()
	} else if PROFILE_CLOCK {
		defer profile.Start(profile.ClockProfile, profile.ProfilePath(".")).Stop()
	}

	// Darkness needs something, if nothing given, then show help.
	if len(os.Args) < 2 {
		ichika.HelpCommandFunc()
		return
	}

	// Find the supplied command...
	if commandFunc := ichika.GetDarknessFunc(os.Args[1]); commandFunc != nil {
		commandFunc()
		return
	}

	// or show a snarky error message
	fmt.Println("command not found?")
	fmt.Println("see help, you pathetic excuse of a man")
}
