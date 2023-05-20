package emilia

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/lipgloss"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/export"
	"github.com/thecsw/darkness/parse"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
	"github.com/thecsw/rei"
)

var (
	// Config is the global darkness config.
	Config *DarknessConfig

	// ParserBuilder returns the parser builder.
	ParserBuilder parse.ParserBuilder

	// ExporterBuilder returns the exporter builder.
	ExporterBuilder export.ExporterBuilder
)

// EmiliaOptions is used for passing options when initiating emilia.
type EmiliaOptions struct {
	// DarknessConfig is the location of darkness's toml config file.
	DarknessConfig string

	// URL is a custom website url, usually used for serving from localhost.
	URL string

	// OutputExtension overrides whatever is in the file.
	OutputExtension string

	// WorkDir is the working directory of the darkness project.
	WorkDir string

	// Dev enables URL generation through local paths.
	Dev bool

	// Test enables test environment, where darkness config is not needed.
	Test bool

	// VendorGalleries dictates whether we should stub in local gallery images.
	VendorGalleries bool
}

var Logger = puck.Logger.WithPrefix("Emilia ❄️ ")

// InitDarkness initializes the darkness config.
func InitDarkness(options *EmiliaOptions) {
	defer puck.Stopwatch("Initialized options").Record()
	Config = &DarknessConfig{}
	Config.WorkDir = options.WorkDir
	if isZero(Config.WorkDir) {
		Config.WorkDir = filepath.Dir(options.DarknessConfig)
		Logger.Info("Guessing working directory", "result", Config.WorkDir)
	}
	data, err := ioutil.ReadFile(options.DarknessConfig)
	if err != nil && !options.Test {
		Logger.Error("Opening config", "path", options.DarknessConfig, "err", err)
		os.Exit(1)
	}
	_, err = toml.Decode(string(data), Config)
	if err != nil {
		Logger.Error("Decoding config", "path", options.DarknessConfig, "err", err)
		os.Exit(1)
	}
	// If input/output formats are empty, default to .org/.html respectively.
	if isZero(Config.Project.Input) {
		Logger.Info("Input format not found, using a default", "ext", puck.ExtensionOrgmode)
		Config.Project.Input = puck.ExtensionOrgmode
	}
	// Output section.
	if isZero(Config.Project.Output) {
		Logger.Info("Input format not found, using a default", "ext", puck.ExtensionHtml)
		Config.Project.Output = puck.ExtensionHtml
	}
	if !isZero(options.OutputExtension) && Config.Project.Output != options.OutputExtension {
		Logger.Warn("Output extension was overwritten", "ext", options.OutputExtension)
		Config.Project.Output = options.OutputExtension
	}
	// Build the parser and exporter builders.
	ParserBuilder = getParserBuilder()
	ExporterBuilder = getExporterBuilder()

	// Define the preview filename.
	if isZero(Config.Website.Preview) {
		Config.Website.Preview = puck.DefaultPreviewFile
	}
	if Config.Website.DescriptionLength < 1 {
		Config.Website.DescriptionLength = 100
	}
	// If the URL is empty, then plug in the current directory.
	if len(Config.URL) < 1 || options.Dev {
		Config.URL, err = os.Getwd()
		if err != nil {
			Logger.Error("Getting working directory, no config url found", "err", err)
			os.Exit(1)
		}
	}
	Config.URLIsLocal = !yunyun.URLRegexp.MatchString(Config.URL)
	// Check if custom URL has been passed
	if len(options.URL) > 0 {
		Config.URL = options.URL
	}
	// URL must end with a trailing forward slash
	if !strings.HasSuffix(Config.URL, "/") {
		Config.URL += "/"
	}
	if !Config.URLIsLocal {
		Config.URLPath, err = url.Parse(Config.URL)
		if err != nil {
			Logger.Error("Parsing url from config", "url", Config.URL, "err", err)
			os.Exit(1)
		}
	}
	Config.URLSlice = []string{Config.URL}
	// Set up the custom highlight languages
	if !isZero(Config.Website.SyntaxHighlightingLanguages) {
		setupHighlightJsLanguages(Config.Website.SyntaxHighlightingLanguages)
	}
	// Set the default syntax highlighting theme
	if isZero(Config.Website.SyntaxHighlightingTheme) {
		Config.Website.SyntaxHighlightingTheme = highlightJsThemeDefaultPath
	}
	// Set the default vendor directory.
	if isZero(Config.Project.DarknessVendorDirectory) {
		Config.Project.DarknessVendorDirectory = defaultVendorDirectory
	}
	// Set the default preview directory.
	if isZero(Config.Project.DarknessPreviewDirectory) {
		Config.Project.DarknessPreviewDirectory = defaultPreviewDirectory
	}
	// Build the regex that will be used to exclude files that
	// have been denoted in emilia darkness config.
	if len(Config.Project.Exclude) > 0 {
		Config.Project.ExcludeEnabled = true
	}
	excludePattern := fmt.Sprintf("(?mU)(%s)/.*",
		strings.Join(yunyun.AnyPathsToStrings(Config.Project.Exclude), "|"))
	Config.Project.ExcludeRegex, err = regexp.Compile(excludePattern)
	if err != nil {
		Logger.Fatalf("bad exclude regex passed ('%s'): %v", excludePattern, err)
	}

	// Work through the vendored galleries.
	Config.VendorGalleries = options.VendorGalleries
	if Config.VendorGalleries {
		cmdColor := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff50a2"))
		yellow := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffff00"))
		fmt.Println("I'm going to vendor all gallery paths!")
		fmt.Println("If this is the first time, it will take a while... otherwise,",
			yellow.Render("an instant"))
		fmt.Printf("Please add %s to your .gitignore, so you don't pollute your git objects.\n",
			cmdColor.Render(string(Config.Project.DarknessVendorDirectory)))
		fmt.Println()
		if err := rei.Mkdir(filepath.Join(Config.WorkDir,
			string(Config.Project.DarknessVendorDirectory))); err != nil {
			Logger.Warnf("creating vendor directory %s: %v", Config.Project.DarknessVendorDirectory, err)
			Logger.Warn("disabling vendoring by force")
			Config.VendorGalleries = false
		}
	}

	// Check whether the author image is full or not by running
	// a url regexp and just hardcode the emilia path. If it's
	// already a URL, then do nothing.
	if !yunyun.URLRegexp.MatchString(string(Config.Author.Image)) {
		Config.Author.ImagePreComputed = JoinPath(Config.Author.Image)
	} else {
		Config.Author.ImagePreComputed = yunyun.FullPathFile(Config.Author.Image)
	}
	// Monkey patch the function if we're using the roman footnotes.
	if Config.Website.RomanFootnotes {
		FootnoteLabeler = rei.NumberToRoman
	}

	// Init the math script
	InitMathJS()
}

// JoinPathGeneric joins the target path with the final root path (url or local).
func JoinPathGeneric[
	relative yunyun.RelativePath,
	full yunyun.FullPath,
](what ...relative) full {
	if !Config.URLIsLocal {
		return full(Config.URLPath.JoinPath(yunyun.AnyPathsToStrings(what)...).String())
	}
	return full(filepath.Join(append(Config.URLSlice, yunyun.AnyPathsToStrings(what)...)...))
}

// JoinPath joins the target path with the final root path (url or local).
func JoinPath(relative ...yunyun.RelativePathFile) yunyun.FullPathFile {
	return JoinPathGeneric[yunyun.RelativePathFile, yunyun.FullPathFile](relative...)
}

// JoinWorkdirGeneric joins target path with the working directory.
func JoinWorkdirGeneric[R yunyun.RelativePath, F yunyun.FullPath](target R) F {
	return F(filepath.Join(Config.WorkDir, string(target)))
}

// JoinWorkdir joins target path with the working directory.
func JoinWorkdir(target yunyun.RelativePathFile) yunyun.FullPathFile {
	return JoinWorkdirGeneric[yunyun.RelativePathFile, yunyun.FullPathFile](target)
}

// isZero returns true if the passed value is a zero value of its type.
func isZero[T comparable](what T) bool {
	return what == gana.ZeroValue[T]()
}

// getParserBuilder returns a new parser object.
func getParserBuilder() parse.ParserBuilder {
	if v, ok := parse.ParserMap[Config.Project.Input]; ok {
		return v
	}
	Logger.Warnf("parser %s not found, defaulting to %s", Config.Project.Input, puck.ExtensionOrgmode)
	Config.Project.Input = puck.ExtensionOrgmode
	return parse.ParserMap[puck.ExtensionOrgmode]
}

// getExporterBuilder returns a new exporter object.
func getExporterBuilder() export.ExporterBuilder {
	if v, ok := export.ExporterMap[Config.Project.Output]; ok {
		return v
	}
	Logger.Warnf("exporter %s not found, defaulting to %s", Config.Project.Output, puck.ExtensionHtml)
	Config.Project.Output = puck.ExtensionHtml
	return export.ExporterMap[puck.ExtensionHtml]
}

// setupHighlightJsLanguages logs all the languages we support through
// the directory included in the config.
func setupHighlightJsLanguages(dir yunyun.RelativePathDir) {
	languages, err := ioutil.ReadDir(string(dir))
	if err != nil {
		Logger.Warnf("failed to open %s: %v", dir, err)
		AvailableLanguages = nil
		return
	}
	for _, language := range languages {
		if !strings.HasSuffix(language.Name(), ".min.js") {
			continue
		}
		AvailableLanguages[strings.TrimSuffix(language.Name(), ".min.js")] = true
	}
}

// trimExt trims extension of a file (only top level, so `file.min.js`
// will be `file.min`)
func trimExt(s string) string {
	return strings.TrimSuffix(s, filepath.Ext(s))
}
