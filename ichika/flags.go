package ichika

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thecsw/darkness/emilia"
)

var (
	// workDir is the directory to look for files.
	workDir = "."
	// darknessToml is the location of `darkness.toml`.
	darknessToml = "darkness.toml"
	// filename is the file to build.
	filename = "index.org"
	// defaultNumOfWorkers gives us the number of workers to
	// spin up in each stage: parsing and processing.
	defaultNumOfWorkers = 14
	// disableParallel sets the number of workers to 1.
	disableParallel bool
	// customNumWorkers sets the custom number of workers.
	customNumWorkers int
	// customChannelCapacity low-level sets the capacity of
	// workers' input/output capacity, defaults to the default
	// number of workers.
	customChannelCapacity int
	// useCurrentDirectory is used for development and local
	// serving, such that you can browse the url files locally.
	useCurrentDirectory bool
)

// getEmiliaOptions takes a cmd subcommand and parses general flags
// and returns emilia options that should be used when initializing emilia.
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

	// Find the absolute path of the work directory to stub in the files.
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
		WorkDir:        workDir,
	}
}

// darknessFlagset returns flagset based on darkness command.
func darknessFlagset(command DarknessCommand) *flag.FlagSet {
	return flag.NewFlagSet(string(command), flag.ExitOnError)

}
