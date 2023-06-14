package markdown

import (
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// preprocess preprocesses the input string to be parser-friendly
func preprocess(data string) string {
	// Add a newline before every heading just in case if
	// there is no terminating empty line before each one
	data = headingRegexp.ReplaceAllString(data, "\n$1")
	// Pad a newline so that last elements can be processed
	// properly before an EOF is encountered during parsing
	data += "\n"
	return data
}

// Parse parses the input string and returns a list of elements
func (p ParserMarkdown) Parse() *yunyun.Page {
	// Split the data into lines
	lines := strings.Split(preprocess(p.Data), "\n")

	page := yunyun.NewPage(
		yunyun.WithFilename(p.Filename),
		yunyun.WithLocation(yunyun.RelativePathTrim(p.Filename)),
		yunyun.WithContents(make([]*yunyun.Content, 0, 16)),
	)

	// currentFlags uses flags to set options
	currentFlags := yunyun.Bits(0)
	// sourceCodeLanguage is the language of the source code block
	sourceCodeLang := ""
	// caption is the current caption we can read
	caption := ""
	// detailsSummary is the current details' summary
	additionalContext := ""
	// Our context is a parody of a state machine
	currentContext := ""
	// listItemInitialIndent is the initial indent of the list item
	listItemInitialIndent := uint8(0)

	addFlag, removeFlag, flipFlag, hasFlag := yunyun.LatchFlags(&currentFlags)
	// addContent is a helper function to add content to the page
	addContent := func(content *yunyun.Content) {
		content.Options = currentFlags
		content.Summary = additionalContext
		content.Caption = caption
		page.Contents = append(page.Contents, content)
		currentContext = ""
	}

	yunyun.ActiveMarkings.Bold = `\*\*`
	yunyun.ActiveMarkings.Italic = `[\*|_]`
	yunyun.ActiveMarkings.Verbatim = "`"
	yunyun.ActiveMarkings.Link = `!?\[(?P<text>[^\]]+)\]\((?P<link>[^\\]+)(?: "(?P<desc>[^"]+)")\)`
	yunyun.ActiveMarkings.BuildRegex()
	linkRegexp = yunyun.LinkRegexp

	// Loop through the lines
	for _, rawLine := range lines {
		// Trimp the line from whitespaces
		line := strings.TrimSpace(rawLine)
		// Save the previous state and update the current
		// one with the newly read line
		previousContext := currentContext
		currentContext = currentContext + line

		// If we are in a raw html envoronment
		if hasFlag(yunyun.InRawHTMLFlag) {
			// Maybe it's time to leave it?
			if isHTMLExportEnd(line) {
				// Mark the leave
				removeFlag(yunyun.InRawHTMLFlag)
				// Save the raw html
				addContent(&yunyun.Content{
					Type:    yunyun.TypeRawHTML,
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
			addFlag(yunyun.InRawHTMLFlag)
			if strings.HasSuffix(line, "unsafe") {
				addFlag(yunyun.InRawHtmlFlagUnsafe)
			}
			currentContext = previousContext
			continue
		}
		// If we are in a source code block?
		if hasFlag(yunyun.InSourceCodeFlag) {
			// Check if it's time to leave
			if isSourceCodeEnd(line) {
				// Mark the leave
				removeFlag(yunyun.InSourceCodeFlag)
				// Save the source code
				addContent(&yunyun.Content{
					Type:           yunyun.TypeSourceCode,
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
			addFlag(yunyun.InSourceCodeFlag)
			currentContext = ""
			continue
		}
		// Ignore orgmode comments and options, where source code blocks
		// and export block options are exceptions to this rule
		if isComment(line) {
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
			addContent(header)
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
				addContent(&yunyun.Content{
					Type: yunyun.TypeHorizontalLine,
				})
				continue
			}
			// If we were in a list, save it as a list
			if hasFlag(yunyun.InListFlag) {
				rawListItems := strings.Split(previousContext, listSeparatorWS)[1:]
				matches := make([]*yunyun.ListItem, len(rawListItems))
				for i, match := range rawListItems {
					listItemRaw := strings.Replace(match, "- ", "", 1)
					matches[i] = &yunyun.ListItem{
						Level: emilia.CountStringsLeft(listItemRaw, "  ") - listItemInitialIndent + 1,
						Text:  strings.TrimSpace(listItemRaw),
					}
				}
				// Shouldn't happen, continue as a failure
				if len(rawListItems) < 1 {
					continue
				}
				// Add the list
				addContent(&yunyun.Content{
					Type: yunyun.TypeList,
					List: matches,
				})
				flipFlag(yunyun.InListFlag)
				listItemInitialIndent = 0
				continue
			}
			// If we were in a table, save it as such
			if hasFlag(yunyun.InTableFlag) {
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
				addContent(&yunyun.Content{
					Type:         yunyun.TypeTable,
					Table:        tableData,
					TableHeaders: hasFlag(yunyun.InTableHasHeadersFlag),
				})
				removeFlag(yunyun.InTableFlag | yunyun.InTableHasHeadersFlag)
				continue
			}
			// Let's see if our context is a standalone link
			if link := isLink(previousContext); link != nil {
				addContent(link)
				continue
			}
			// Also check if this is an attention block, like "NOTE:..." or "WARNING:..."
			if attention := isAttentionBlock(previousContext); attention != nil {
				addContent(attention)
				continue
			}
			// By default, save whatever we have as a paragraph
			addContent(formParagraph(previousContext, additionalContext, currentFlags))
			// Reset the drop cap flag
			removeFlag(yunyun.InDropCapFlag)
			continue
		}
		if isList(line) {
			addFlag(yunyun.InListFlag)
			currentContext = previousContext + listSeparatorWS + line
		}
		if isTable(line) {
			addFlag(yunyun.InTableFlag)
			// If it's a delimeter, save it and move on
			if isTableHeaderDelimeter(line) {
				addFlag(yunyun.InTableHasHeadersFlag)
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

func fillHolosceneDate(page *yunyun.Page) {
	// No contents found?
	if len(page.Contents) < 1 {
		return
	}
	first := gana.First(page.Contents)
	// Needs to be a simple text
	if !first.IsParagraph() {
		return
	}
	if !strings.HasSuffix(first.Paragraph, "H.E.") {
		return
	}
	page.Date = first.Paragraph
	page.DateHoloscene = true
}
