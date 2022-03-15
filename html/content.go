package html

import (
	"fmt"

	"github.com/thecsw/darkness/internals"
)

var (
	contentFunctions = []func(*internals.Content) string{
		headings, paragraph, list, listNumbered,
		link, sourceCode, rawHTML, horizontalLine, attentionBlock,
	}
)

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

func list(content *internals.Content) string {
	elements := ""
	for _, item := range content.List {
		elements += fmt.Sprintf(`
<li>
<p>
%s
</p>
</li>
`, processText(item))
	}
	return fmt.Sprintf(`
<div class="ulist">
<ul>
%s
</ul>
</div>
`, elements)
}

func listNumbered(content *internals.Content) string {
	// TODO
	return ""
}

func sourceCode(content *internals.Content) string {
	return fmt.Sprintf(`
<div class="listingblock">
<div class="content">
<pre class="highlight">
<code class="language-%s" data-lang="%s">%s</code>
</pre>
</div>
</div>
`, content.SourceCodeLang, content.SourceCodeLang, content.SourceCode)
}

func rawHTML(content *internals.Content) string {
	return content.RawHTML
}

func horizontalLine(content *internals.Content) string {
	return `<hr>`
}

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
