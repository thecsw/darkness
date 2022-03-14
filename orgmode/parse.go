package orgmode

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
)

func ParseFile(workDir, file string) *internals.Page {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	page := Parse(string(data))
	page.URL = emilia.JoinPath(strings.TrimPrefix(filepath.Dir(file), workDir))
	return page
}

func Parse(data string) *internals.Page {
	// Add a newline before every heading just in case if
	// there is no terminating empty line before each one
	data = HeadingRegexp.ReplaceAllString(data, "\n"+`$1 `)
	// Pad a newline so that last elements can be processed
	// properly before an EOF is encountered during parsing
	data += "\n"
	// Split the data into lines
	lines := strings.Split(data, "\n")
	page := &internals.Page{}
	page.Contents = make([]internals.Content, 0, 16)

	inList := false
	inSourceCode := false
	inRawHTML := false
	sourceCodeLang := ""

	// Our context is a parody of a state machine
	currentContext := ""
	addContent := func(content internals.Content) {
		page.Contents = append(page.Contents, content)
		currentContext = ""
	}

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
					SourceCode:     previousContext,
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
				page.Title = header.Header
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
				matches := UnorderedListRegexp.FindAllStringSubmatch(previousContext, -1)
				// Shouldn't happen, continue as a failure
				if len(matches) < 1 {
					continue
				}
				currentList := make([]string, 0, len(matches))
				for _, match := range matches {
					currentList = append(currentList, match[1])
				}
				// Add the list
				addContent(internals.Content{
					Type: internals.TypeList,
					List: currentList,
				})
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
		currentContext += " "
		// Mark if the current line is a list
		inList = isList(line)
		// Add a delimeter so we can later regex out each item
		if inList {
			currentContext += "∆ "
		}
	}
	return page
}
