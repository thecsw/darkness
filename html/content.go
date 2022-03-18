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
		headings, paragraph, list, listNumbered,
		link, sourceCode, rawHTML, horizontalLine, attentionBlock,
	}
)

// headings gives us a heading html representation
func headings(content *internals.Content) string {
	start := ``
	if !content.HeadingChild && !content.HeadingFirst {
		start = `</div>` + "\n" + `</div>`
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
		toReturn += "\n" + `</div>`
	}
	return toReturn
}

// paragraph gives us a paragraph html representation
func paragraph(content *internals.Content) string {
	text := processText(content.Paragraph)
	return fmt.Sprintf(
		`
<div class="paragraph">
<p>%s</p>
</div>`,
		text,
	)
}

// list gives us a list html representation
func list(content *internals.Content) string {
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

// listNumbered gives us a numbered list html representation
func listNumbered(content *internals.Content) string {
	// TODO
	return ""
}

// sourceCode gives us a source code html representation
func sourceCode(content *internals.Content) string {
	return fmt.Sprintf(`
<div class="listingblock">
<div class="content">
<pre class="highlight">
<code class="language-%s" data-lang="%s">%s</code>
</pre>
</div>
</div>
`, content.SourceCodeLang, content.SourceCodeLang, func(sourceCode string) string {
		// Remove the nested parser blockers
		s := strings.ReplaceAll(sourceCode, ",#", "#")
		// Escape the whatever HTML that is found in source code
		s = html.EscapeString(s)
		return s
	}(content.SourceCode))
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
