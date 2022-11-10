package orgmode

import (
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
)

// Preprocess preprocesses the input string to be parser-friendly
func Preprocess(data string) string {
	// Add a newline before every heading just in case if
	// there is no terminating empty line before each one
	data = headingRegexp.ReplaceAllString(data, "\n$1")
	// Center and quote delimeters need a new line around
	for _, v := range surroundWithNewlines {
		data = strings.ReplaceAll(data,
			optionPrefix+v,
			"\n"+optionPrefix+v)
	}
	// Pad a newline so that last elements can be processed
	// properly before an EOF is encountered during parsing
	data += "\n"
	return data
}

// Parse parses the input string and returns a list of elements
func Parse(data string, filename string) *internals.Page {
	// Split the data into lines
	lines := strings.Split(Preprocess(data), "\n")

	page := internals.NewPage(
		internals.WithFilename(filename),
		internals.WithURL(emilia.JoinPath(filepath.Dir(filename))),
		internals.WithContents(make([]internals.Content, 0, 16)),
	)

	// currentFlags uses flags to set options
	currentFlags := internals.Bits(0)
	// sourceCodeLanguage is the language of the source code block
	sourceCodeLang := ""
	// caption is the current caption we can read
	caption := ""
	// detailsSummary is the current details' summary
	additionalContext := ""
	// Our context is a parody of a state machine
	currentContext := ""

	addFlag, removeFlag, flipFlag, hasFlag := internals.LatchFlags(&currentFlags)
	// addContent is a helper function to add content to the page
	addContent := func(content internals.Content) {
		content.Options = currentFlags
		content.Summary = additionalContext
		content.Caption = caption
		page.Contents = append(page.Contents, content)
		currentContext = ""
	}
	optionsActions := map[string]func(line string){
		optionDropCap:     func(line string) { addFlag(internals.InDropCapFlag) },
		optionBeginQuote:  func(line string) { addFlag(internals.InQuoteFlag) },
		optionEndQuote:    func(line string) { removeFlag(internals.InQuoteFlag) },
		optionBeginCenter: func(line string) { addFlag(internals.InCenterFlag) },
		optionEndCenter:   func(line string) { removeFlag(internals.InCenterFlag) },
		optionBeginDetails: func(line string) {
			addFlag(internals.InDetailsFlag)
			additionalContext = extractDetailsSummary(line)
			if additionalContext == "" {
				additionalContext = "open for details"
			}
			addContent(internals.Content{Type: internals.TypeDetails})
		},
		optionEndDetails: func(line string) {
			removeFlag(internals.InDetailsFlag)
			addContent(internals.Content{Type: internals.TypeDetails})
		},
		optionBeginGallery: func(line string) {
			addFlag(internals.InGalleryFlag)
			additionalContext = extractGalleryFolder(line)
		},
		optionEndGallery: func(line string) { removeFlag(internals.InGalleryFlag) },
		optionCaption:    func(line string) { caption = extractCaptionTitle(line) },
		optionDate:       func(line string) { page.Date = extractDate(line) },
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
		if hasFlag(internals.InRawHTMLFlag) {
			// Maybe it's time to leave it?
			if isHTMLExportEnd(line) {
				// Mark the leave
				removeFlag(internals.InRawHTMLFlag)
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
			addFlag(internals.InRawHTMLFlag)
			currentContext = previousContext
			continue
		}
		// If we are in a source code block?
		if hasFlag(internals.InSourceCodeFlag) {
			// Check if it's time to leave
			if isSourceCodeEnd(line) {
				// Mark the leave
				removeFlag(internals.InSourceCodeFlag)
				// Save the source code
				addContent(internals.Content{
					Type:           internals.TypeSourceCode,
					SourceCodeLang: sourceCodeLang,
					SourceCode:     strings.TrimRight(previousContext, "\n\t\r\f\b"),
					Caption:        caption,
				})
				continue
			}
			// Save the context and continue
			currentContext = previousContext + rawLine + "\n"
			continue
		}
		// Should we enter a source code environment?
		if isSourceCodeBegin(line) {
			sourceCodeLang = extractSourceCodeLanguage(line)
			addFlag(internals.InSourceCodeFlag)
			currentContext = ""
			continue
		}
		// Ignore orgmode comments and options, where source code blocks
		// and export block options are exceptions to this rule
		if isComment(line) {
			currentContext = previousContext
			continue
		}
		// isOption is a sink for any options that darkness
		// does not support, hence will be ignored
		if isOption(line) {
			givenLine := line[2:]
			option := strings.Split(givenLine, " ")[0]
			if action, ok := optionsActions[option]; ok {
				action(rawLine)
			}
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
				addContent(internals.Content{
					Type: internals.TypeHorizontalLine,
				})
				continue
			}
			// If we were in a list, save it as a list
			if hasFlag(internals.InListFlag) {
				matches := strings.Split(previousContext, listSeparatorWS)[1:]
				for i, match := range matches {
					matches[i] = strings.Replace(match, "- ", "", 1)
				}
				// Shouldn't happen, continue as a failure
				if len(matches) < 1 {
					continue
				}
				// Add the list
				addContent(internals.Content{
					Type: internals.TypeList,
					List: matches,
				})
				flipFlag(internals.InListFlag)
				continue
			}
			// If we were in a table, save it as such
			if hasFlag(internals.InTableFlag) {
				rows := strings.Split(previousContext, tableSeparatorWS)[1:]
				tableData := make([][]string, len(rows))
				for i, row := range rows {
					row = strings.TrimSpace(row)
					if len(row) < 1 {
						continue
					}
					// Split by the item delimeter
					columns := strings.Split(row, "|")
					// Trim the array from the first and last element
					columns = columns[1 : len(columns)-1]
					// Trim each item by the left/right whitespace
					for j, item := range columns {
						columns[j] = strings.TrimSpace(item)
					}
					tableData[i] = columns
				}
				addContent(internals.Content{
					Type:         internals.TypeTable,
					Table:        tableData,
					TableHeaders: hasFlag(internals.InTableHasHeadersFlag),
				})
				removeFlag(internals.InTableFlag | internals.InTableHasHeadersFlag)
				continue
			}
			// Let's see if our context is a standalone link
			if link := isLink(previousContext); link != nil {
				addContent(*link)
				continue
			}
			// Also check if this is an attention block, like "NOTE:..." or "WARNING:..."
			if attention := isAttentionBlock(previousContext); attention != nil {
				addContent(*attention)
				continue
			}
			// By default, save whatever we have as a paragraph
			addContent(*formParagraph(previousContext, additionalContext, currentFlags))
			// Reset the drop cap flag
			removeFlag(internals.InDropCapFlag)
			continue
		}
		if isList(line) {
			addFlag(internals.InListFlag)
			currentContext = previousContext + listSeparatorWS + line
		}
		if isTable(line) {
			addFlag(internals.InTableFlag)
			// If it's a delimeter, save it and move on
			if isTableHeaderDelimeter(line) {
				addFlag(internals.InTableHasHeadersFlag)
				currentContext = previousContext
				continue
			}
			currentContext = previousContext + tableSeparatorWS + line
		}
		currentContext += " "
	}

	// Optional parsing to see if H.E. has been left on the first line
	// as the date
	fillHolosceneDate(page)
	return page
}

func fillHolosceneDate(page *internals.Page) {
	// No contents found?
	if len(page.Contents) < 1 {
		return
	}
	// Needs to be a simple text
	if !page.Contents[0].IsParagraph() {
		return
	}
	if !strings.HasSuffix(page.Contents[0].Paragraph, "H.E.") {
		return
	}
	page.Date = page.Contents[0].Paragraph
	page.DateHoloscene = true
}
