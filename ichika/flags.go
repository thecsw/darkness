package ichika

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/thecsw/darkness/emilia/alpha"

	"github.com/charmbracelet/log"
	"github.com/thecsw/darkness/emilia/puck"
)

var (
	// workDir is the directory to look for files.
	workDir = "."

	// darknessToml is the location of `darkness.toml`.
	darknessToml = "darkness.toml"

	// disableParallel sets the number of workers to 1.
	disableParallel bool

	// customNumWorkers sets the custom number of workers.
	customNumWorkers int

	// debugEnabled tells us whether to show debug logs.
	debugEnabled bool

	// infoEnabled tells us whether to show info logs.
	infoEnabled bool

	// useCurrentDirectory is used for development and local
	// serving, such that you can browse the url files locally.
	useCurrentDirectory bool

	// vendorGalleryImages is a flag that dictates whether we should
	// store a local copy of all remote gallery images and stub them
	// in the gallery links instead of the remote links.
	//
	// Turning this option on would result in a VERY slow build the
	// first time, as it would need to retrieve however many images
	// from remote services.
	//
	// All images will be put in "darkness_vendor" directory, which
	// will be skipped in discovery process AND should be put it
	// .gitignore by user, so they don't pollute their git objects.
	vendorGalleryImages bool
)

// getAlphaOptions takes a cmd subcommand and parses general flags
// and returns emilia options that should be used when initializing emilia.
func getAlphaOptions(cmd *flag.FlagSet) alpha.Options {
	cmd.StringVar(&workDir, "dir", ".", "where do I look for files")
	cmd.StringVar(&darknessToml, "conf", "darkness.toml", "location of darkness.toml")
	cmd.BoolVar(&disableParallel, "disable-parallel", false, "disable parallel build (only use one worker)")
	cmd.IntVar(&customNumWorkers, "workers", 4, "number of workers to use")
	cmd.BoolVar(&infoEnabled, "info", false, "enable info logs")
	cmd.BoolVar(&debugEnabled, "debug", false, "enable debug logs")
	cmd.BoolVar(&useCurrentDirectory, "dev", false, "use local path for urls (development)")
	cmd.BoolVar(&vendorGalleryImages, "vendor-galleries", false, "stub in local copies of gallery links (SLOW)")
	if err := cmd.Parse(os.Args[2:]); err != nil {
		puck.Logger.Fatalf("parsing build arguments: %v", err)
	}

	// Find the absolute path of the work directory to stub in the files.
	var err error
	workDir, err = filepath.Abs(workDir)
	if err != nil {
		puck.Logger.Fatalf("determining absolute path of %s: %v", workDir, err)
	}

	// If parallel processing is disabled, only provision one workers
	// per each processing stage.
	if disableParallel {
		customNumWorkers = 1
	}

	// Set the proper log levels and default to warn.
	switch {
	case debugEnabled:
		puck.Logger.SetLevel(log.DebugLevel)
	case infoEnabled:
		puck.Logger.SetLevel(log.InfoLevel)
	}

	// Read the config and initialize emilia settings.
	return alpha.Options{
		DarknessConfig:  darknessToml,
		Dev:             useCurrentDirectory,
		WorkDir:         workDir,
		VendorGalleries: vendorGalleryImages,
	}
}

// darknessFlagset returns flagset based on darkness command.
func darknessFlagset(command DarknessCommand) *flag.FlagSet {
	return flag.NewFlagSet(string(command), flag.ExitOnError)
}
