package orgmode

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
	"github.com/thecsw/rei"
)

const (
	// setup file should be on its own line like an #include directive.
	specialSetupFileDirective = `^#[+](setupfile|SETUPFILE):\s*([^\s]+)$`

	// this is the syntax for evaluating macros
	specialMacroEvalDirective = `{{{([a-z0-9]+)(\([^)(]+\))?}}}`

	macroDefinition = `macro:`
	macroPrefix     = optionPrefix + macroDefinition + " "
)

var (
	// headingRegexp is the regexp for matching headlines.
	headingRegexp = regexp.MustCompile(`(?m)^(\*{1,6} )`)

	// specialSetupFileDirectivePattern is just compiled specialSetupFileDirective
	specialSetupFileDirectivePattern = regexp.MustCompile(specialSetupFileDirective)

	// specialMacroEvalDirectivePattern dictates on how we match and eval macros
	specialMacroEvalDirectivePattern = regexp.MustCompile(specialMacroEvalDirective)

	shouldBeSurroundedWithNewLines = map[string]struct{}{
		optionBeginQuote: {}, optionEndQuote: {},
		optionBeginCenter: {}, optionEndCenter: {},
		optionBeginDetails: {}, optionEndDetails: {},
		optionBeginGallery: {}, optionEndGallery: {},
	}

	expandedFiles = sync.Map{}
)

func (p ParserOrgmode) preprocess(filename yunyun.RelativePathFile, what string) string {
	// We will do everything in one pass here and build the final input file using
	// a string builder for performance.
	sb := &strings.Builder{}

	// Here we will store the macro definitions.
	macrosLookupTable := make(map[string]string)

	// We add a newline before lists start
	previousLine := ""

	// We will read the original input line by line and build the final input same way.
	// Be ready for a very greedy loop.
	lines := strings.Split(what, "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, commentPrefix) {
			continue
		}

		// Let's see if we have any macros to expand on this line.
		if updatedLine, expandedMacros := expandMacros(p.Config,
			filename, macrosLookupTable, line); expandedMacros {
			line = updatedLine
		}

		if headingRegexp.MatchString(line) {
			sb.WriteRune('\n')
			sb.WriteString(line)
			sb.WriteRune('\n')
			continue
		}

		// We are in the realm of options now.
		lowercase := strings.ToLower(line)
		if isOption(line) {
			option := gana.SkipString(uint(optionPrefixLen), lowercase)
			parts := strings.SplitN(option, " ", 2)
			// Don't know what this is, don't let it reach the parser.
			if len(parts) < 1 {
				continue
			}
			key := strings.TrimSpace(parts[0])

			// Check if it needs to be surrounded by a newline.
			if _, ok := shouldBeSurroundedWithNewLines[key]; ok {
				sb.WriteRune('\n')
				sb.WriteString(line)
				sb.WriteRune('\n')
				continue
			}

			// What if it's a setup file?
			if setupFile, found := expandSetupFile(p.Config, filename, line); found {
				sb.WriteRune('\n')
				sb.WriteString(setupFile)
				sb.WriteRune('\n')

				// Expanded files may and will contain macro definitions that we
				// will need to evaluate later in the file.
				collectMacros(p.Config, filename, macrosLookupTable, setupFile)
				continue
			}

			// If we found any macros on a single line, then the whole line was a macro
			// definition and gets consumed.
			if collectMacros(p.Config, filename, macrosLookupTable, line) {
				continue
			}
		}

		// We add a newline before listings, so the parser has an easier time.
		if isList(line) && !isList(previousLine) {
			sb.WriteRune('\n')
		}

		// By default, if we reached the end of the iteration, write the line as is.
		sb.WriteString(line)
		sb.WriteRune('\n') // regular linefeed

		// Save it for the next iteration.
		previousLine = line
	}

	// Pad a newline so that last elements can be processed
	// properly before an EOF is encountered during parsing
	sb.WriteRune('\n')
	return sb.String()
}

func expandSetupFile(conf *alpha.DarknessConfig, filename yunyun.RelativePathFile, line string) (string, bool) {
	matches := specialSetupFileDirectivePattern.FindAllStringSubmatch(line, 1)
	if len(matches) < 1 {
		return "", false
	}
	setupFileTargetFilename := yunyun.RelativePathFile(strings.TrimSpace(matches[0][2]))

	// If the file doesn't exist, we must fail.
	currentDirectory := yunyun.RelativePathFile(filepath.Dir(string(filename)))
	relativeImportFilename := yunyun.RelativePathFile(filepath.Join(string(currentDirectory), string(setupFileTargetFilename)))
	absoluteImportFilename := conf.Runtime.WorkDir.Join(relativeImportFilename)

	// Check the hot cache.
	if expandedFile, alreadyExpanded := expandedFiles.Load(absoluteImportFilename); alreadyExpanded {
		// See if the type is right, if it's not, drop in to the slow IO retrieval.
		if stringified, isString := expandedFile.(string); isString {
			return stringified, true
		}
	}

	if !rei.FileMustExist(string(absoluteImportFilename)) {
		conf.Runtime.Logger.Fatal("setupfile target not found",
			"orgfile", filename, "target", absoluteImportFilename)
	}

	// Read the data and splash it into the input.
	setupFileTargetContents := string(rei.Must(os.ReadFile(filepath.Clean(string(absoluteImportFilename)))))
	expandedFiles.Store(absoluteImportFilename, setupFileTargetContents)
	return setupFileTargetContents, true
}

func collectMacros(
	conf *alpha.DarknessConfig,
	filename yunyun.RelativePathFile,
	macrosLookupTable map[string]string,
	what string) bool {
	macroDefsFound := false
	for _, line := range strings.Split(what, "\n") {
		if !strings.HasPrefix(line, macroPrefix) {
			continue
		}
		macroLine := gana.SkipString(uint(len(macroPrefix)), line)
		split := strings.SplitN(macroLine, " ", 2)
		if len(split) != 2 {
			conf.Runtime.Logger.Fatal("malformed macro definition found",
				"path", filename, "line", line)
		}
		macroDefsFound = true
		macroName := strings.TrimSpace(split[0])
		macroBody := strings.TrimSpace(split[1])
		macrosLookupTable[macroName] = macroBody
	}
	return macroDefsFound
}

func expandMacros(conf *alpha.DarknessConfig,
	filename yunyun.RelativePathFile,
	macrosLookupTable map[string]string,
	line string,
) (string, bool) {
	macroEvaluationMatches := specialMacroEvalDirectivePattern.FindAllStringSubmatch(line, -1)
	if len(macroEvaluationMatches) < 1 {
		return "", false
	}
	for _, match := range macroEvaluationMatches {
		fullMatch := match[0]
		macroName := strings.TrimSpace(match[1])
		if _, ok := macrosLookupTable[macroName]; !ok {
			conf.Runtime.Logger.Fatal("macro used but not defined",
				"path", filename, "macro", macroName)
		}
		macroBody := strings.ReplaceAll(macrosLookupTable[macroName], "\\n", "\n")
		macroParamsString := match[2]

		// what if it has no params?
		if len(macroParamsString) < 1 {
			line = strings.ReplaceAll(line, fullMatch, macroBody)
			continue
		}

		macroParamsDirty := strings.Split(strings.Trim(macroParamsString, ")("), ",")
		macroParams := make([]string, 0, len(macroParamsDirty))
		for _, macroParamDirty := range macroParamsDirty {
			macroParams = append(macroParams, strings.TrimSpace(macroParamDirty))
		}

		// Let's get the body and hydrate the parameters.
		for i, param := range macroParams {
			macroBody = strings.ReplaceAll(macroBody, `$`+strconv.Itoa(i+1), param)
		}

		// then replace the macro eval with the hydrated body
		line = strings.ReplaceAll(line, fullMatch, macroBody)
	}
	return line, true
}
