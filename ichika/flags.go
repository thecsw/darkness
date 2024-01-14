package ichika

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/ichika/kuroko"
)

// getAlphaOptions takes a cmd subcommand and parses general flags
// and returns emilia options that should be used when initializing emilia.
func getAlphaOptions(cmd *flag.FlagSet) alpha.Options {
	cmd.StringVar(&kuroko.WorkDir, "dir", ".", "where do I look for files")
	cmd.StringVar(&kuroko.DarknessToml, "conf", "darkness.toml", "location of darkness.toml")
	cmd.BoolVar(&kuroko.DisableParallel, "disable-parallel", false, "disable parallel build (only use one worker)")
	cmd.IntVar(&kuroko.CustomNumWorkers, "workers", 4, "number of workers to use")
	cmd.BoolVar(&kuroko.InfoEnabled, "info", false, "enable info logs")
	cmd.BoolVar(&kuroko.DebugEnabled, "debug", false, "enable debug logs")
	cmd.BoolVar(&kuroko.UseCurrentDirectory, "dev", false, "use local path for urls (development)")
	cmd.BoolVar(&kuroko.VendorGalleryImages, "vendor-galleries", false, "stub in local copies of gallery links (SLOW)")
	cmd.BoolVar(&kuroko.Akaneless, "akaneless", false, "skip akane processing")
	cmd.BoolVar(&kuroko.Force, "force", false, "force post-processing (akane or misa)")
	if len(os.Args) < 2 {
		puck.Logger.Fatalf("no command specified")
	}
	if err := cmd.Parse(os.Args[2:]); err != nil {
		puck.Logger.Fatalf("parsing build arguments: %v", err)
	}

	// Find the absolute path of the work directory to stub in the files.
	var err error
	kuroko.WorkDir, err = filepath.Abs(kuroko.WorkDir)
	if err != nil {
		puck.Logger.Fatalf("determining absolute path of %s: %v", kuroko.WorkDir, err)
	}

	// If parallel processing is disabled, only provision one workers
	// per each processing stage.
	if kuroko.DisableParallel {
		kuroko.CustomNumWorkers = 1
	}

	// Set the proper log levels and default to warn.
	switch {
	case kuroko.DebugEnabled:
		puck.Logger.SetLevel(log.DebugLevel)
	case kuroko.InfoEnabled:
		puck.Logger.SetLevel(log.InfoLevel)
	}

	// Read the config and initialize emilia settings.
	return alpha.Options{
		DarknessConfig:  kuroko.DarknessToml,
		Dev:             kuroko.UseCurrentDirectory,
		WorkDir:         kuroko.WorkDir,
		VendorGalleries: kuroko.VendorGalleryImages,
	}
}

// darknessFlagset returns flagset based on darkness command.
func darknessFlagset(command DarknessCommand) *flag.FlagSet {
	return flag.NewFlagSet(string(command), flag.ExitOnError)
}
