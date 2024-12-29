package alpha

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/ichika/kuroko"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
)

// BuildConfig builds the config from the passed options.
func BuildConfig(options Options) *DarknessConfig {
	conf := &DarknessConfig{}
	conf.Runtime.Logger = puck.NewLogger("Alpha ☕")
	conf.Runtime.WorkDir = WorkingDirectory(options.WorkDir)

	// Record the time it takes to initialize the options.
	defer puck.Stopwatch("Initialized options").Record(conf.Runtime.Logger)

	// If the working directory is empty, then default to the directory of the config file.
	if isUnset(conf.Runtime.WorkDir) {
		conf.Runtime.WorkDir = WorkingDirectory(filepath.Dir(options.DarknessConfig))
		conf.Runtime.Logger.Info("Guessing working directory", "result", conf.Runtime.WorkDir)
	}

	// Read the config file.
	data, err := os.ReadFile(options.DarknessConfig)
	if err != nil && !options.Test {
		conf.Runtime.Logger.Fatal("Opening config", "path", options.DarknessConfig, "err", err)
	}

	// If we can't decode the config, then exit.
	_, err = toml.Decode(string(data), &conf)
	if err != nil {
		conf.Runtime.Logger.Fatal("Decoding config", "path", options.DarknessConfig, "err", err)
	}

	// Define the preview filename.
	if isUnset(conf.Website.Preview) {
		conf.Website.Preview = puck.DefaultPreviewFile
	}
	if conf.Website.DescriptionLength < 1 {
		conf.Website.DescriptionLength = 100
	}

	// If the Url is empty, then plug in the current directory.
	if len(conf.Url) < 1 || options.Dev {
		conf.Url, err = os.Getwd()
		if err != nil {
			conf.Runtime.Logger.Error("Getting working directory, no config url found", "err", err)
			os.Exit(1)
		}
	}

	// Check if custom Url has been passed
	if len(options.Url) > 0 {
		conf.Url = options.Url
	}

	conf.Runtime.isUrlLocal = !yunyun.UrlRegexp.MatchString(conf.Url)

	// Url must end with a trailing forward slash
	if !strings.HasSuffix(conf.Url, "/") {
		conf.Url += "/"
	}

	// If the Url is not local, then parse it. Otherwise, just use the Url as is.
	if !conf.Runtime.isUrlLocal {
		conf.Runtime.UrlPath, err = url.Parse(conf.Url)
		if err != nil {
			conf.Runtime.Logger.Error("Parsing url from config", "url", conf.Url, "err", err)
			os.Exit(1)
		}
	}
	conf.Runtime.urlSlice = []string{conf.Url}

	// Set up the custom highlight languages if they exist.
	conf.setupHighlightJsLanguages()

	// Set up the highlight theme if it's not given.
	if isUnset(conf.Website.SyntaxHighlightingTheme) {
		conf.Website.SyntaxHighlightingTheme = highlightJsThemeDefaultPath
	}

	// Set the default vendor directory if it's not set.
	if isUnset(conf.Project.DarknessVendorDirectory) {
		conf.Project.DarknessVendorDirectory = puck.DefaultVendorDirectory
	}

	// Set the default preview directory if it's not set.
	if isUnset(conf.Project.DarknessPreviewDirectory) {
		conf.Project.DarknessPreviewDirectory = puck.DefaultPreviewDirectory
	}

	// Build the regex that will be used to exclude files that
	// have been denoted in emilia darkness config.
	if len(conf.Project.Exclude) > 0 {
		conf.Project.ExcludeEnabled = true
	}

	// Build the regex that will be used to exclude files that
	// have been denoted in emilia darkness config.
	excludePattern := fmt.Sprintf("(?mU)(%s)/.*",
		strings.Join(yunyun.AnyPathsToStrings(conf.Project.Exclude), "|"))
	conf.Project.ExcludeRegex, err = regexp.Compile(excludePattern)
	if err != nil {
		conf.Runtime.Logger.Fatalf("bad exclude regex passed ('%s'): %v", excludePattern, err)
	}

	// Check whether the author image is full or not by running
	// a url regexp and just hardcode the emilia path. If it's
	// already a Url, then do nothing.
	if !yunyun.UrlRegexp.MatchString(string(conf.Author.Image)) {
		conf.Author.ImagePreComputed = conf.Runtime.Join(conf.Author.Image)
	} else {
		conf.Author.ImagePreComputed = yunyun.FullPathFile(conf.Author.Image)
	}

	// Set up the project extensions.
	conf.setupProjectExtensions(options)

	// Set up the gallery vendoring.
	conf.setupGalleryVendoring(options)

	// Only if we need LFS, do we need to enable the Git integration workflow.x
	if kuroko.LfsEnabled {
		// Last but not least, let's try to set up the git remote.
		if isUnset(conf.External.GitRemotePath) || isUnset(conf.External.GitRemoteService) {
			service, path, err := ExtractGitRemote(conf)
			if err != nil {
				conf.Runtime.Logger.Warnf("could not get the git remote info: %v", err)
			} else {
				conf.External.GitRemoteService = service
				conf.External.GitRemotePath = path
			}
		}

		// Handle the git branch as well.
		if isUnset(conf.External.GitBranch) {
			branch, err := ExtractGitBranch(conf)
			if err != nil {
				conf.Runtime.Logger.Warnf("could not get the current git branch: %v", err)
			} else {
				conf.External.GitBranch = branch
			}
		}
	}

	// If all git values are set, then put a quick and cheap marker to signify so.
	if allAreSet(conf.External.GitRemoteService, conf.External.GitRemotePath, conf.External.GitBranch) {
		conf.External.GitRemotesAreValid = true
	}

	// Set up the gallery vendoring.
	return conf
}

// isUnset returns true if the passed value is a zero value of its type.
func isUnset[T comparable](what T) bool {
	return what == gana.ZeroValue[T]()
}

// isSet returns true if the passed value is not a zero value of its type.
func isSet[T comparable](what T) bool {
	return what != gana.ZeroValue[T]()
}

// allAreSet will return true iff all the values given are not default.
func allAreSet[T comparable](whats ...T) bool {
	for _, what := range whats {
		if isUnset(what) {
			return false
		}
	}
	return true
}

// trimExt trims extension of a file (only top level, so `file.min.js`
// will be `file.min`)
func trimExt(s string) string {
	return strings.TrimSuffix(s, filepath.Ext(s))
}
