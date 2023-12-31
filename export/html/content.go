package html

import (
	"fmt"
	"html"
	"strings"

	"github.com/thecsw/darkness/emilia/narumi"
	"github.com/thecsw/darkness/yunyun"
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
<div class="paragraph%s">
<p>
%s
</p>
</div>`,
		// div class
		paragraphClass(content), processText(content.Paragraph),
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
<div class="ulist">
<ul class="%s">
%s
</ul>
</div>
`,
		content.Summary, // overloaded summary to store list class
		strings.Join(gana.Map(makeListItem, content.List), "\n"))
}

// listNumbered gives us a numbered list html representation
func (e *state) listNumbered(content *yunyun.Content) string {
	// TODO
	return ""
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
	rows := make([]string, len(content.Table))
	for i := range content.Table {
		for j, v := range content.Table[i] {
			topTag := "td"
			if i == 0 && content.TableHeaders {
				topTag = "th"
			}

			// Check if the cell is a special cell
			toWrite := ""
			if insideCell, isSpecial := tableSpecialCell(v); isSpecial {
				toWrite = insideCell
			} else {
				toWrite = processText(v)
			}

			content.Table[i][j] = fmt.Sprintf("<%s>%s</%s>", topTag, toWrite, topTag)
		}
		rows[i] = fmt.Sprintf("<tr>\n%s</tr>", strings.Join(content.Table[i], "\n"))
	}
	tableHtml := fmt.Sprintf("<table>%s</table>", strings.Join(rows, "\n"))
	return fmt.Sprintf(tableTemplate, content.CustomHtmlTags, content.Caption, tableHtml)
}

const (
	// tableSpecialImagePrefix is the prefix for special cells that contain an image.
	tableSpecialImagePrefix = "#+image:"
)

// tableSpecialCell checks if the cell is a special cell, and if so, returns the HTML
// representation of it, and a boolean indicating that it was indeed a special cell.
// for example, if the cell is "#+image: [link][text "description"]
// it will return the HTML representation of the image, and true.
func tableSpecialCell(what string) (string, bool) {
	if link := yunyun.ExtractLink(what); link != nil {
		// If the link is an image, return the HTML representation of it.
		if strings.HasPrefix(link.Text, tableSpecialImagePrefix) {
			return fmt.Sprintf(
				`<img class="image" src="%s" title="%s" alt="%s">`,
				link.Link,
				yunyun.RemoveFormatting(strings.TrimPrefix(link.Description, tableSpecialImagePrefix)),
				yunyun.RemoveFormatting(strings.TrimPrefix(link.Text, tableSpecialImagePrefix)),
			), true
		}
	}
	return what, false

}

// table gives an HTML formatted table
func (e *state) details(content *yunyun.Content) string {
	if content.IsDetails() {
		return fmt.Sprintf("<details>\n<summary>%s</summary>\n<hr>", content.Summary)
	}
	return "</details>"
}
