package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thecsw/darkness/emilia"
)

func getEmiliaOptions(cmd *flag.FlagSet) *emilia.EmiliaOptions {
	cmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	cmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	cmd.BoolVar(&disableParallel, "disable-parallel", false, "disable parallel build (only use one worker)")
	cmd.IntVar(&customNumWorkers, "workers", defaultNumOfWorkers, "number of workers to spin up")
	cmd.IntVar(&customChannelCapacity, "capacity", defaultNumOfWorkers, "worker channels' capacity")
	cmd.BoolVar(&useCurrentDirectory, "dev", false, "use local path for urls (development)")
	if err := cmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("failed to parse build arguments, fatal: %s", err.Error())
		os.Exit(1)
	}

	var err error
	workDir, err = filepath.Abs(workDir)
	if err != nil {
		fmt.Println("Couldn't determine absolute path of", workDir)
		os.Exit(1)
	}

	// If parallel processing is disabled, only provision one workers
	// per each processing stage.
	if disableParallel {
		customNumWorkers = 1
	}

	// Read the config and initialize emilia settings.
	return &emilia.EmiliaOptions{
		DarknessConfig: darknessToml,
		Dev:            useCurrentDirectory,
	}
}
