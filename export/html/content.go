package html

import (
	"fmt"
	"html"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

const (
	// linkWasNotSpecialFlag is used internally to mark that the read link
	// content was not an embed and should be treated as simple text.
	linkWasNotSpecialFlag yunyun.Bits = yunyun.YunYunStartCustomFlags << iota
	// thisContentOpensWritingFlag shows that we are entering the `.writing` tag.
	thisContentOpensWritingFlag
	// ThisContentClosesWriting shows that we have to close the `.writing` tag.
	thisContentClosesWritingFlag
	// ThisContentClosesDivSection shows that we have to close the previously
	// opened `.sectionbody`
	thisContentClosesDivSectionFlag
)

// setContentFlags sets the content flags.
func (e *ExporterHTML) setContentFlags(v *yunyun.Content) {
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

// buildContent runs the givent `*Content` against known protocols/policies
// and does some funky logic to balance div openings and closures.
func (e *ExporterHTML) buildContent() string {
	built := e.contentFunctions[e.currentContent.Type](e.currentContent)
	e.setContentFlags(e.currentContent)
	return e.resolveDivTags(built)
}

// resolveDivTags applies results from `setContentFlags` by modifying the DOM.
func (e *ExporterHTML) resolveDivTags(built string) string {
	if yunyun.HasFlag(&e.currentContent.Options, thisContentOpensWritingFlag) {
		built = `<div class="writing">` + "\n" + built
	}
	if yunyun.HasFlag(&e.currentContent.Options, thisContentClosesWritingFlag) {
		built = "</div>\n" + built
	}
	if e.inWriting && e.currentContentIndex == e.contentsNum-1 {
		built = built + "\n</div>\n"
	}
	return built
}

// headings gives us a heading html representation
func (e *ExporterHTML) Heading(content *yunyun.Content) string {
	toReturn := fmt.Sprintf(`
<h%d id="%s" class="section-%d">%s</h%d>`,
		content.HeadingLevelAdjusted, // HTML open tag
		extractID(content.Heading),   // ID
		content.HeadingLevel,         // section class
		processText(content.Heading), // Actual title
		content.HeadingLevelAdjusted, // HTML close tag
	)
	e.inHeading = true
	return toReturn
}

// paragraph gives us a paragraph html representation
func (e ExporterHTML) Paragraph(content *yunyun.Content) string {
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

// makeListItem makes an html item
func makeListItem(s string) string {
	return fmt.Sprintf(`
<li>
<p>
%s
</p>
</li>`, processText(s))
}

// list gives us a list html representation
func (e ExporterHTML) List(content *yunyun.Content) string {
	// Hijack this type for galleries
	if content.IsGallery() {
		return e.gallery(content)
	}
	return fmt.Sprintf(`
<div class="ulist">
<ul>
%s
</ul>
</div>
`, strings.Join(gana.Map(makeListItem, content.List), "\n"))
}

var (
	flexOptionPattern = `:flex ([12345])`
	flexOptionRegexp  = regexp.MustCompile(flexOptionPattern)
)

// extractCustomFlex extract custom flex class `:flex [1,5]`
func extractCustomFlex(s string) uint {
	matches := flexOptionRegexp.FindAllStringSubmatch(s, -1)
	if len(matches) < 1 {
		return 0
	}
	if len(matches[0]) < 1 {
		return 0
	}
	ret, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return 0
	}
	return uint(ret)
}

// makeFlexItem will make an item of the flexbox .gallery with 1/3 width
func makeFlexItem(s string, folder string, width uint) string {
	matchLen, link, text, desc := yunyun.ExtractLink(s)
	// Maybe they just didn't use a proper link pattern? Stub the value in instead then.
	if matchLen < 0 {
		link = s
	}
	fullImage := filepath.Join(folder, url.PathEscape(strings.TrimSpace(link)))
	ext := filepath.Ext(fullImage)
	previewImage := strings.TrimSuffix(fullImage, ext) + "_preview" + ext
	if customFlex := extractCustomFlex(s); customFlex != 0 {
		width = customFlex
	}
	return fmt.Sprintf(`<img class="item lazyload flex-%d" src="%s" data-src="%s" title="%s" alt="%s">`,
		width, previewImage, fullImage, desc, text)
}

// gallery will create a flexbox gallery as defined in .gallery css class
func (e ExporterHTML) gallery(content *yunyun.Content) string {
	galleryFolder := ""
	if len(content.GalleryPath) > 0 {
		galleryFolder = content.GalleryPath
	}
	makeFlexItemWithFolder := func(s string) string {
		return makeFlexItem(s, galleryFolder, content.GalleryImagesPerRow)
	}
	return fmt.Sprintf(`
<div class="gallery-container">
<center>
<div class="gallery">
%s
</div>
</center>
</div>
`, strings.Join(gana.Map(makeFlexItemWithFolder, content.List), "\n"))
}

// listNumbered gives us a numbered list html representation
func (e ExporterHTML) ListNumbered(content *yunyun.Content) string {
	// TODO
	return ""
}

// sourceCode gives us a source code html representation
func (e ExporterHTML) SourceCode(content *yunyun.Content) string {
	return fmt.Sprintf(`
<div class="listingblock">
<pre class="highlight"><code class="language-%s" data-lang="%s">%s</code></pre>
</div>
`, emilia.MapSourceCodeLang(content.SourceCodeLang), content.SourceCodeLang, func() string {
		// Remove the nested parser blockers
		s := strings.ReplaceAll(content.SourceCode, ",#", "#")
		// Escape the whatever HTML that is found in source code
		s = html.EscapeString(s)
		return s
	}())
}

// rawHTML gives us a raw html representation
func (e ExporterHTML) RawHTML(content *yunyun.Content) string {
	// If the unsafe flag is enabled, don't even wrap it in `mediablock`
	if content.IsRawHTMLUnsafe() {
		return content.RawHTML
	}
	return fmt.Sprintf(rawHTMLTemplate, content.RawHTML, content.Caption)

}

// horizontalLine gives us a horizontal line html representation
func (e ExporterHTML) HorizontalLine(content *yunyun.Content) string {
	return `<center>
<hr>
</center>`
}

// attentionBlock gives us a attention block html representation
func (e ExporterHTML) AttentionBlock(content *yunyun.Content) string {
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
func (e ExporterHTML) Table(content *yunyun.Content) string {
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
func (e ExporterHTML) Details(content *yunyun.Content) string {
	if content.IsDetails() {
		return fmt.Sprintf("<details>\n<summary>%s</summary>\n<hr>", content.Summary)
	}
	return "</details>"
}
