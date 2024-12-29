package kuroko

import "github.com/charmbracelet/log"

var (
	// WorkDir is the directory to look for files.
	WorkDir = "."

	// DarknessToml is the location of `darkness.toml`.
	DarknessToml = "darkness.toml"

	// DisableParallel sets the number of workers to 1.
	DisableParallel bool

	// CustomNumWorkers sets the custom number of workers.
	CustomNumWorkers int

	// DebugEnabled tells us whether to show debug logs.
	DebugEnabled bool

	// InfoEnabled tells us whether to show info logs.
	InfoEnabled bool

	// LfsEnabled turns on the "lfs:" image path expansions.
	LfsEnabled bool

	// UseCurrentDirectory is used for development and local
	// serving, such that you can browse the url files locally.
	UseCurrentDirectory bool

	// Akaneless will skip akane post-processing if enabled.
	Akaneless bool

	// Force will command some of the post-processors, like akane
	// and misa to generate something like page previews even if
	// they already exist.
	Force bool

	// VendorGalleryImages is a flag that dictates whether we should
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
	VendorGalleryImages bool

	// BuildReport will output a timestamped file in the local
	// project's .darkness directory with then files discovered, duration,
	// and the output file that they reached.
	BuildReport bool
)

// LogLevel returns the log level as defined in kuroko
func LogLevel() log.Level {
	if DebugEnabled {
		return log.DebugLevel
	}
	if InfoEnabled {
		return log.InfoLevel
	}
	return log.WarnLevel
}
