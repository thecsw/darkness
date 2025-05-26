package html

import (
	"fmt"
	"html"
	"strings"

	"github.com/thecsw/darkness/v3/emilia/narumi"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
)

const (
	// linkWasNotSpecialFlag is used internally to mark that the read link
	// content was not an embed and should be treated as simple text.
	linkWasNotSpecialFlag = yunyun.YunYunStartCustomFlags << iota
	// thisContentOpensWritingFlag shows that we are entering the `.writing` tag.
	thisContentOpensWritingFlag
	// ThisContentClosesWriting shows that we have to close the `.writing` tag.
	thisContentClosesWritingFlag
	// ThisContentClosesDivSection shows that we have to close the previously
	// opened `.sectionbody`
	thisContentClosesDivSectionFlag
)

// setContentFlags sets the content flags.
func (e *state) setContentFlags(v *yunyun.Content) {
	contentDivType := whatDivType(v)
	// Mark situations when we have to leave writing
	if e.inWriting && (contentDivType != divWriting) {
		yunyun.AddFlag(&v.Options, thisContentClosesWritingFlag)
		e.inWriting = false
	}
	// Mark situations when we have to enter writing
	if contentDivType == divWriting && !e.inWriting {
		yunyun.AddFlag(&v.Options, thisContentOpensWritingFlag)
		e.inWriting = true
	}
}

// resolveDivTags applies results from `setContentFlags` by modifying the DOM.
func (e *state) resolveDivTags(built string) string {
	if yunyun.HasFlag(&e.currentContent.Options, thisContentOpensWritingFlag) {
		built = `<div class="writing">` + "\n" + built
	}
	if yunyun.HasFlag(&e.currentContent.Options, thisContentClosesWritingFlag) {
		built = "</div>\n" + built
	}
	if e.inWriting && e.currentContentIndex == len(e.page.Contents)-1 {
		built = built + "\n</div>\n"
	}
	return built
}

// heading gives us a heading html representation.
func (e *state) heading(content *yunyun.Content) string {
	toReturn := fmt.Sprintf(`
<h%d id="%s" class="section-%d">%s</h%d>`,
		content.HeadingLevelAdjusted, // HTML open tag
		ExtractID(content.Heading),   // ID
		content.HeadingLevel,         // section class
		processText(content.Heading), // Actual title
		content.HeadingLevelAdjusted, // HTML close tag
	)
	e.inHeading = true
	return toReturn
}

func paragraphClass(content *yunyun.Content) string {
	if content.IsQuote() {
		return " quote"
	}
	if content.IsCentered() {
		return " center"
	}
	if content.IsDropCap() {
		return " dropcap"
	}
	return ""
}

// paragraph gives us a paragraph html representation
func (e *state) paragraph(content *yunyun.Content) string {
	return fmt.Sprintf(
		`
<div class="paragraph%s" %s>
<p>
%s
</p>
</div>`,
		paragraphClass(content), content.CustomHtmlTags, processText(content.Paragraph),
	)
}

// makeListItem makes an html item
func makeListItem(item yunyun.ListItem) string {
	return fmt.Sprintf(`
<li class="l%d">
<p>
%s
</p>
</li>`, item.Level, processText(item.Text))
}

// list gives us a list html representation
func (e *state) list(content *yunyun.Content) string {
	// Hijack this type for galleries
	if content.IsGallery() {
		return e.gallery(content)
	}
	return fmt.Sprintf(`
<div class="ulist" %s>
<ul class="%s">
%s
</ul>
</div>
`,
		content.CustomHtmlTags,
		content.Summary, // overloaded summary to store list class
		strings.Join(gana.Map(makeListItem, content.List), "\n"))
}

// listNumbered gives us a numbered list html representation
func (e *state) listNumbered(content *yunyun.Content) string {
	// Hijack this type for galleries
	if content.IsGallery() {
		return e.gallery(content)
	}
	return fmt.Sprintf(`
<div class="olist" %s>
<ol class="%s">
%s
</ol>
</div>
`,
		content.CustomHtmlTags,
		content.Summary, // overloaded summary to store list class
		strings.Join(gana.Map(makeListItem, content.List), "\n"))
}

// sourceCode gives us a source code html representation
func (e *state) sourceCode(content *yunyun.Content) string {
	return fmt.Sprintf(`
<div class="coding" %s>
<div class="listingblock">
<pre class="highlight"><code class="language-%s" data-lang="%s">%s</code></pre>
</div>
</div>
`,
		content.CustomHtmlTags,
		narumi.MapSourceCodeLang(content.SourceCodeLang),
		content.SourceCodeLang,
		func() string {
			// Remove the nested parser blockers
			s := strings.ReplaceAll(content.SourceCode, ",#", "#")
			// Escape the whatever HTML that is found in source code
			s = html.EscapeString(s)
			return s
		}(),
	)
}

// rawHTML gives us a raw html representation
func (e *state) rawHtml(content *yunyun.Content) string {
	// If the unsafe flag is enabled, don't even wrap it in `mediablock`
	if content.IsRawHtmlUnsafe() {
		return content.RawHtml
	}
	// If responsive enabled, wrap the inner iframe (*probably*) in it.
	if content.IsRawHtmlResponsive() {
		return fmt.Sprintf(responsiveIFrameHtmlTemplate, content.CustomHtmlTags, content.RawHtml)
	}
	return fmt.Sprintf(rawHtmlTemplate, content.CustomHtmlTags, content.RawHtml, content.Caption)
}

// horizontalLine gives us a horizontal line html representation
func (e *state) horizontalLine(content *yunyun.Content) string {
	return `<center>
<hr>
</center>`
}

// attentionBlock gives us a attention block html representation
func (e *state) attentionBlock(content *yunyun.Content) string {
	return fmt.Sprintf(`
<div class="admonitionblock note">
<table>
<tr>
<td class="icon">
<div class="title">%s</div>
</td>
<td class="content">
%s
</td>
</tr>
</table>
</div>`, content.AttentionTitle, processText(content.AttentionText))
}

// table gives an HTML formatted table
func (e *state) table(content *yunyun.Content) string {
	var headers, rows []string
	numRows := len(content.Table)

	// If there are headers, we need to make sure we don't count them as rows.
	if content.TableHeaders {
		headers = make([]string, len(content.Table[0]))
		numRows--
		for j, header := range content.Table[0] {
			headers[j] = fmt.Sprintf("<th>%s</th>", processTableCell(header))
		}
	}

	// If there are rows, start making them.
	if numRows > 0 {
		rows = make([]string, numRows)
		// Exclude the headers from the rows by reslicing.
		if content.TableHeaders {
			content.Table = content.Table[1:]
		}
		// Make the rows.
		for i, row := range content.Table {
			for j, v := range row {
				content.Table[i][j] = fmt.Sprintf("<td>%s</td>", processTableCell(v))
			}
			rows[i] = fmt.Sprintf("<tr>\n%s</tr>", strings.Join(content.Table[i], "\n"))
		}
	}

	// Make the html table.
	tableHtml := fmt.Sprintf(
		`<table>
<thead>
<tr>
%s
</tr>
</thead>
<tbody>
%s
</tbody>
</table>`,
		strings.Join(headers, "\n"),
		strings.Join(rows, "\n"),
	)
	return fmt.Sprintf(tableTemplate, content.CustomHtmlTags, content.Caption, tableHtml)
}

// processTableCell returns the HTML representation of a table cell given its content.
func processTableCell(what string) string {
	if insideCell, isSpecial := tableSpecialCell(what); isSpecial {
		return insideCell
	}
	return processText(what)
}

const (
	// tableSpecialImagePrefix is the prefix for special cells that contain an image.
	tableSpecialImagePrefix = "file:"
)

// tableSpecialCell checks if the cell is a special cell, and if so, returns the HTML
// representation of it, and a boolean indicating that it was indeed a special cell.
// for example, if the cell is "#+image: [link][text "description"]
// it will return the HTML representation of the image, and true.
func tableSpecialCell(what string) (string, bool) {
	if link := yunyun.ExtractLink(what); link != nil {
		// If the link is an image, return the HTML representation of it.
		if strings.HasPrefix(link.Link, tableSpecialImagePrefix) {
			return fmt.Sprintf(
				`<img class="image" src="%s" title="%s" alt="%s">`,
				strings.TrimPrefix(link.Link, tableSpecialImagePrefix),
				yunyun.RemoveFormatting(link.Description),
				yunyun.RemoveFormatting(link.Text),
			), true
		}
	}
	return what, false
}

// table gives an HTML formatted table
func (e *state) details(content *yunyun.Content) string {
	if content.Type == yunyun.TypeDetails {
		if content.IsDetails() {
			return fmt.Sprintf("<details>\n<summary>%s</summary>\n<hr>", content.Summary)
		}
		return "</details>"
	}
	return ""
}
