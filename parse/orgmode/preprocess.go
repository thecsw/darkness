package orgmode

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/rei"
)

const (
	// setup file should be on its own line like an #include directive.
	specialSetupFileDirective = `(?i)^#[+]setupfile:\s*([^\s]+)\s*\n`
)

var (
	// headingRegexp is the regexp for matching headlines.
	headingRegexp = regexp.MustCompile(`(?m)^(\*{1,6} )`)

	// specialSetupFileDirectivePattern is just compiled specialSetupFileDirective
	specialSetupFileDirectivePattern = regexp.MustCompile(specialSetupFileDirective)

	// These will be surrounded by newlines to acommodate parsing.
	surroundWithNewlines = []string{
		optionBeginQuote, optionEndQuote,
		optionBeginCenter, optionEndCenter,
		optionBeginDetails, optionEndDetails,
		optionBeginGallery, optionEndGallery,
	}
)

// preprocess preprocesses the input string to be parser-friendly
func (p ParserOrgmode) preprocess(filename yunyun.RelativePathFile, what string) string {
	// Add a newline before every heading just in case if
	// there is no terminating empty line before each one
	what = headingRegexp.ReplaceAllString(what, "\n$1")
	what = preprocessBySurroundingWithNewline(what)

	// Let's fin any SETUPFILE directives and treat it as a special case,
	// since those operate almost exactly as C-style #include directives would.
	what = preprocessByExpandingSetupFile(p.Config, filename, what)

	// Pad a newline so that last elements can be processed
	// properly before an EOF is encountered during parsing
	what += "\n"
	return what
}

// will replace the #+setupfile directive with the referenced org and fatally fail if the
// target file is not found or we can't read the contents, for some reason, latter will panic.
func preprocessByExpandingSetupFile(conf *alpha.DarknessConfig,
	inputFilename yunyun.RelativePathFile, what string) string {
	if !specialSetupFileDirectivePattern.MatchString(what) {
		return what
	}
	matches := specialSetupFileDirectivePattern.FindAllStringSubmatch(what, -1)
	for _, match := range matches {
		setupFileDirectiveMatch := strings.TrimSpace(match[0])
		setupFileTargetFilename := yunyun.RelativePathFile(strings.TrimSpace(match[1]))

		// If the file doesn't exist, we must fail.
		currentDirectory := yunyun.RelativePathFile(filepath.Dir(string(inputFilename)))
		relativeImportFilename := yunyun.RelativePathFile(filepath.Join(string(currentDirectory), string(setupFileTargetFilename)))
		absoluteImportFilename := conf.Runtime.WorkDir.Join(relativeImportFilename)
		if !rei.FileMustExist(string(absoluteImportFilename)) {
			conf.Runtime.Logger.Fatal("setupfile target not found",
				"orgfile", inputFilename, "target", absoluteImportFilename)
		}

		// Read the data and splash it into the input.
		setupFileTargetContents := rei.Must(os.ReadFile(filepath.Clean(string(absoluteImportFilename))))
		what = strings.Replace(what, setupFileDirectiveMatch, string(setupFileTargetContents), 1)
	}
	return what
}

// preprocessBySurroundingWithNewline just surrounds what's needed with newlines.
func preprocessBySurroundingWithNewline(what string) string {
	// Center and quote delimeters need a new line around
	for _, v := range surroundWithNewlines {
		what = strings.ReplaceAll(what,
			optionPrefix+v,
			"\n"+optionPrefix+v)
	}
	return what
}
