package html

import (
	"fmt"
	"html"
	"strings"

	"github.com/thecsw/darkness/internals"
)

var (
	// contentFunctions is a map of functions that process content
	contentFunctions = []func(*internals.Content) string{
		headings,
		paragraph,
		list,
		listNumbered,
		link,
		sourceCode,
		rawHTML,
		horizontalLine,
		attentionBlock,
		table,
		details,
	}
)

// headings gives us a heading html representation
func headings(content *internals.Content) string {
	start := ``
	if !content.HeadingChild && !content.HeadingFirst {
		start = "</div>\n</div>"
	}
	if content.HeadingChild && !content.HeadingFirst {
		start = `</div>`
	}
	toReturn := fmt.Sprintf(start+`
<div class="sect%d">
<h%d id="%s">%s</h%d>
<div class="sectionbody">`,
		content.HeadingLevel-1,       // CSS class
		content.HeadingLevel,         // HTML open tag
		extractID(content.Heading),   // ID
		processText(content.Heading), // Actual title
		content.HeadingLevel,         // HTML close tag
	)
	if content.HeadingLast {
		toReturn += "</div>\n</div>\n"
	}
	return toReturn
}

// paragraph gives us a paragraph html representation
func paragraph(content *internals.Content) string {
	return fmt.Sprintf(
		`
<div class="paragraph%s">
<p>
%s
</p>
</div>`,
		// div class
		func() string {
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
		}(),
		processText(content.Paragraph),
	)
}

// list gives us a list html representation
func list(content *internals.Content) string {
	// Hijack this type for galleries
	if internals.HasFlag(&content.Options, internals.InGalleryFlag) {
		return gallery(content)
	}

	elements := make([]string, len(content.List))
	for i, item := range content.List {
		elements[i] = fmt.Sprintf(`
<li>
<p>
%s
</p>
</li>`, processText(item))
	}

	return fmt.Sprintf(`
<div class="ulist">
<ul>
%s
</ul>
</div>
`, strings.Join(elements, "\n"))
}

func gallery(content *internals.Content) string {
	flexItems := make([]string, len(content.List))
	for i, item := range content.List {
		flexItems[i] = fmt.Sprintf(`<img class="item" height="33%%" width="33%%" src="%s">`, strings.TrimSpace(item))
	}
	return fmt.Sprintf(`
<center>
<div class="gallery">
%s
</div>
</center>
`, strings.Join(flexItems, "\n"))
}

// listNumbered gives us a numbered list html representation
func listNumbered(content *internals.Content) string {
	// TODO
	return ""
}

var sourceCodeLang = map[string]string{
	"":   "plaintext",
	"sh": "bash",
}

func mapSourceCodeLang(s string) string {
	if v, ok := sourceCodeLang[s]; ok {
		return v
	}
	return s
}

// sourceCode gives us a source code html representation
func sourceCode(content *internals.Content) string {
	lang := mapSourceCodeLang(content.SourceCodeLang)
	return fmt.Sprintf(`
<div class="listingblock">
<div class="content">
<pre class="highlight"><code class="language-%s" data-lang="%s">%s</code></pre>
</div>
</div>
`, lang, lang, func() string {
		// Remove the nested parser blockers
		s := strings.ReplaceAll(content.SourceCode, ",#", "#")
		// Escape the whatever HTML that is found in source code
		s = html.EscapeString(s)
		return s
	}())
}

// rawHTML gives us a raw html representation
func rawHTML(content *internals.Content) string {
	return content.RawHTML
}

// horizontalLine gives us a horizontal line html representation
func horizontalLine(content *internals.Content) string {
	return `<hr>`
}

// attentionBlock gives us a attention block html representation
func attentionBlock(content *internals.Content) string {
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
func table(content *internals.Content) string {
	rows := make([]string, len(content.Table))
	for i := range content.Table {
		for j, v := range content.Table[i] {
			topTag := "td"
			if i == 0 && content.TableHeaders {
				topTag = "th"
			}
			content.Table[i][j] = fmt.Sprintf("<%s>%s</%s>", topTag, processText(v), topTag)
		}
		rows[i] = fmt.Sprintf("<tr>\n%s</tr>", strings.Join(content.Table[i], "\n"))
	}
	tableHTML := fmt.Sprintf("<table>%s</table>", strings.Join(rows, "\n"))
	return fmt.Sprintf(tableTemplate, content.Caption, tableHTML)
}

// table gives an HTML formatted table
func details(content *internals.Content) string {
	if internals.HasFlag(&content.Options, internals.InDetailsFlag) {
		return fmt.Sprintf("<details>\n<summary>%s</summary>\n<hr>", content.Summary)
	}
	return "</details>"
}
