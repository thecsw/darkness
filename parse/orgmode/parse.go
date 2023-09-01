package orgmode

import (
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// preprocess preprocesses the input string to be parser-friendly
func preprocess(data string) string {
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

const (
	defaultGalleryImagesPerRow uint = 3
)

// Parse parses the input string and returns a list of elements
func (p ParserOrgmode) Do(
	filename yunyun.RelativePathFile,
	data string,
) *yunyun.Page {
	defer puck.Stopwatch("Parsed", "page", filename).Record()

	// Split the data into lines
	lines := strings.Split(preprocess(data), "\n")

	page := yunyun.NewPage(
		yunyun.WithFilename(filename),
		yunyun.WithLocation(yunyun.RelativePathTrim(filename)),
		yunyun.WithContents(make([]*yunyun.Content, 0, 32)),
	)
	page.Author = p.Config.RSS.DefaultAuthor

	// currentFlags uses flags to set options
	currentFlags := yunyun.Bits(0)
	// sourceCodeLanguage is the language of the source code block
	sourceCodeLang := ""
	// caption is the current caption we can read
	caption := ""
	// attributes is the attributes for the current content.
	attributes := ""
	// detailsSummary is the current details' summary
	additionalContext := ""
	// galleryPath stores the gallery's declared path
	galleryPath := ""
	// galleryWidth dictates on how many items per row we will have
	galleryWidth := defaultGalleryImagesPerRow
	// Our context is a parody of a state machine
	currentContext := ""
	// User can provide custom style for an image (like resizing).
	customHtmlTags := ""
	// listItemInitialIndent is the initial indent of the list item
	listItemInitialIndent := uint8(0)

	// optionsStrings will get populated as the page is being scanned
	// and then parsed out before leaving this parser.
	optionsStrings := ""
	defer emilia.FillAccoutrement(p.Config.Website.Tombs, &optionsStrings, page)

	// Optional parsing to see if H.E. has been left on the first line
	// as the date
	defer fillHolosceneDate(page)

	addFlag, removeFlag, flipFlag, hasFlag := yunyun.LatchFlags(&currentFlags)
	// addContent is a helper function to add content to the page
	addContent := func(content *yunyun.Content) {
		content.Options = currentFlags
		content.Summary = additionalContext
		content.GalleryPath = yunyun.RelativePathDir(galleryPath)
		content.GalleryImagesPerRow = galleryWidth
		content.Caption = caption
		content.Attributes = attributes
		content.CustomHtmlTags = customHtmlTags
		page.Contents = append(page.Contents, content)
		currentContext = ""
		galleryPath = ""
		galleryWidth = defaultGalleryImagesPerRow
		additionalContext = ""
		attributes = ""
		customHtmlTags = ""
	}
	optionsActions := map[string]func(line string){
		optionDropCap:     func(line string) { addFlag(yunyun.InDropCapFlag) },
		optionBeginQuote:  func(line string) { addFlag(yunyun.InQuoteFlag) },
		optionEndQuote:    func(line string) { removeFlag(yunyun.InQuoteFlag) },
		optionBeginCenter: func(line string) { addFlag(yunyun.InCenterFlag) },
		optionEndCenter:   func(line string) { removeFlag(yunyun.InCenterFlag) },
		optionBeginDetails: func(line string) {
			addFlag(yunyun.InDetailsFlag)
			additionalContext = extractDetailsSummary(line)
			if additionalContext == "" {
				additionalContext = "open for details"
			}
			addContent(&yunyun.Content{Type: yunyun.TypeDetails})
		},
		optionEndDetails: func(line string) {
			removeFlag(yunyun.InDetailsFlag)
			addContent(&yunyun.Content{Type: yunyun.TypeDetails})
		},
		optionBeginGallery: func(line string) {
			addFlag(yunyun.InGalleryFlag)
			galleryPath = extractGalleryFolder(line)
			galleryWidth = extractGalleryImagesPerRow(line)
		},
		optionEndGallery: func(line string) { removeFlag(yunyun.InGalleryFlag) },
		optionCaption:    func(line string) { caption = extractCaptionTitle(line) },
		optionDate:       func(line string) { page.Date = extractDate(line) },
		optionHtmlHead:   func(line string) { page.HtmlHead = append(page.HtmlHead, extractHtmlHead(line)) },
		optionOptions:    func(line string) { optionsStrings += extractOptions(line) + " " },
		optionAttributes: func(line string) { attributes = extractAttributes(line) },
		optionAuthor:     func(line string) { page.Author = extractAuthor(line) },
		optionHtmlTags:   func(line string) { customHtmlTags = extractHtmlTags(line) },
	}

	// Yunyun's markings default to orgmode
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
		if hasFlag(yunyun.InRawHtmlFlag) {
			// Maybe it's time to leave it?
			if isHtmlExportEnd(line) {
				// Mark the leave
				removeFlag(yunyun.InRawHtmlFlag)
				// Save the raw html
				addContent(&yunyun.Content{
					Type:    yunyun.TypeRawHtml,
					RawHtml: previousContext,
				})
				continue
			}
			// Otherwise, continue saving the context
			currentContext = previousContext + rawLine + "\n"
			continue
		}
		// Now, check if we can enter a raw html environment
		if isHtmlExportBegin(line) {
			addFlag(yunyun.InRawHtmlFlag)
			if strings.Contains(line, "unsafe") {
				addFlag(yunyun.InRawHtmlFlagUnsafe)
			} else if strings.Contains(line, "responsive") || strings.Contains(line, "iframe") {
				addFlag(yunyun.InRawHtmlFlagResponsive)
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
				matches := make([]yunyun.ListItem, len(rawListItems))
				for i, match := range rawListItems {
					listItemRaw := strings.Replace(match, "- ", "", 1)
					matches[i] = yunyun.ListItem{
						Level: gana.CountStringsLeft[uint8](listItemRaw, "  ") - listItemInitialIndent + 1,
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
			if link := getLink(previousContext); link != nil {
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
			if !hasFlag(yunyun.InListFlag) {
				listItemInitialIndent = gana.CountRunesLeft[uint8](rawLine, ' ')
			}
			addFlag(yunyun.InListFlag)
			currentContext = previousContext + listSeparatorWS + rawLine
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

	return page
}

// fillHolosceneDate tries to find a date in the format of "H.E." and
// saves it as the page's date.
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
