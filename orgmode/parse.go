package orgmode

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
)

// ParseFile parses a single file and returns a list of elements
func ParseFile(workDir, file string) *internals.Page {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("failed to open the file %s: %s", file, err.Error())
		os.Exit(1)
	}
	page := Parse(string(data))
	page.URL = emilia.JoinPath(strings.TrimPrefix(filepath.Dir(file), workDir))
	return page
}

// Preprocess preprocesses the input string to be parser-friendly
func Preprocess(data string) string {
	// Add a newline before every heading just in case if
	// there is no terminating empty line before each one
	data = HeadingRegexp.ReplaceAllString(data, "\n$1")
	// Debug stuff
	// fmt.Println(data)
	// fmt.Println("---------------------")
	// Pad a newline so that last elements can be processed
	// properly before an EOF is encountered during parsing
	data += "\n"
	return data
}

// Parse parses the input string and returns a list of elements
func Parse(data string) *internals.Page {
	// Split the data into lines
	lines := strings.Split(Preprocess(data), "\n")
	page := &internals.Page{
		Title:     "",
		URL:       "",
		MetaTags:  []internals.MetaTag{},
		Links:     []internals.Link{},
		Contents:  []internals.Content{},
		Footnotes: []string{},
		Scripts:   []string{},
	}
	page.Contents = make([]internals.Content, 0, 16)

	// inList is true if we are in a list
	inList := false
	// inSourceCode is true if we are in a source code block
	inSourceCode := false
	// inRaw is true if we are in a raw block
	inRawHTML := false
	// sourceCodeLanguage is the language of the source code block
	sourceCodeLang := ""

	// Our context is a parody of a state machine
	currentContext := ""
	// addContent is a helper function to add content to the page
	addContent := func(content internals.Content) {
		page.Contents = append(page.Contents, content)
		currentContext = ""
	}

	// Loop through the lines
	for _, rawLine := range lines {
		// Trimp the line from whitespaces
		line := strings.TrimSpace(rawLine)
		// Save the previous state and update the current
		// one with the newly read line
		previousContext := currentContext
		currentContext = currentContext + line

		// If we are in a raw html envoronment
		if inRawHTML {
			// Maybe it's time to leave it?
			if isHTMLExportEnd(line) {
				// Mark the leave
				inRawHTML = false
				// Save the raw html
				addContent(internals.Content{
					Type:    internals.TypeRawHTML,
					RawHTML: previousContext,
				})
				continue
			}
			// Otherwise, continue saving the context
			currentContext = previousContext + rawLine + "\n"
			continue
		}
		// Now, check if we can enter a raw html environment
		if isHTMLExportBegin(line) {
			inRawHTML = true
			currentContext = previousContext
			continue
		}
		// If we are in a source code block?
		if inSourceCode {
			// Check if it's time to leave
			if isSourceCodeEnd(line) {
				// Mark the leave
				inSourceCode = false
				// Save the source code
				addContent(internals.Content{
					Type:           internals.TypeSourceCode,
					SourceCodeLang: sourceCodeLang,
					SourceCode:     strings.TrimRight(previousContext, "\n\t\r\f\b"),
				})
				continue
			}
			// Save the context and continue
			currentContext = previousContext + rawLine + "\n"
			continue
		}
		// Should we enter a source code environment?
		if isSourceCodeBegin(line) {
			sourceCodeLang = sourceExtractLang(line)
			inSourceCode = true
			currentContext = ""
			continue
		}
		// Ignore orgmode comments and options, where source code blocks
		// and export block options are exceptions to this rule
		if isComment(line) || isOption(line) {
			currentContext = previousContext
			continue
		}
		// Now, we need to parse headings here
		if header := isHeader(line); header != nil {
			if header.HeadingLevel == 1 {
				page.Title = header.Heading
				currentContext = ""
				continue
			}
			addContent(*header)
			continue
		}
		// If we hit an empty line, end the whatever context we had
		if line == "" {
			// Empty context gets us nowhere
			if len(previousContext) < 1 {
				continue
			}
			// Add a horizontal line divider
			if isHorizonalLine(previousContext) {
				addContent(internals.Content{Type: internals.TypeHorizontalLine})
				continue
			}
			// If we were in a list, save it as a list
			if inList {
				matches := UnorderedListRegexp.FindAllStringSubmatch(previousContext[2:]+" ∆ ", -1)
				// Shouldn't happen, continue as a failure
				if len(matches) < 1 {
					continue
				}
				currentList := make([]string, len(matches))
				for i, match := range matches {
					currentList[i] = match[1]
				}
				// Add the list
				addContent(internals.Content{
					Type: internals.TypeList,
					List: currentList,
				})
				inList = false
				continue
			}
			// Let's see if our context is a standalone link
			if link := isLink(previousContext); link != nil {
				addContent(*link)
				continue
			}
			// Also check if this is an attention block, like "NOTE:..." or "WARNING:..."
			if attention := isAttentionBlack(previousContext); attention != nil {
				addContent(*attention)
				continue
			}
			// By default, save whatever we have as a paragraph
			addContent(*formParagraph(previousContext))
			continue
		}
		if isList(line) {
			inList = true
			currentContext = previousContext + " ∆" + line
		}
		currentContext += " "
	}
	return page
}
