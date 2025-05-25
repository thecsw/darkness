package orgmode

import (
	"strings"
	"testing"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/yunyun"
)

// TestParser tests the basic functionality of the parser
func TestParser(t *testing.T) {
	// Create a basic configuration for testing
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	// Test parsing an empty string
	page := parser.Do("test.org", "")
	if page == nil {
		t.Fatalf("Parser returned nil page for empty input")
	}

	if len(page.Contents) != 0 {
		t.Errorf("Parser returned non-empty contents for empty input: %v", page.Contents)
	}
}

// TestParsingTitle tests parsing a document with a title (level 1 heading)
func TestParsingTitle(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := "* Title"
	page := parser.Do("test.org", input)

	if page.Title != "Title" {
		t.Errorf("Expected title to be 'Title', got '%s'", page.Title)
	}
	
	// No content should be added to the page when we just have a title
	if len(page.Contents) != 0 {
		t.Errorf("Expected no content elements, got %d", len(page.Contents))
	}
}

// TestParsingHeadings tests parsing a document with headings of different levels
func TestParsingHeadings(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `* Title
** Heading Level 2
*** Heading Level 3
**** Heading Level 4
***** Heading Level 5`

	page := parser.Do("test.org", input)

	if page.Title != "Title" {
		t.Errorf("Expected title to be 'Title', got '%s'", page.Title)
	}

	expectedHeadings := []struct {
		text  string
		level uint32
	}{
		{"Heading Level 2", 2},
		{"Heading Level 3", 3},
		{"Heading Level 4", 4},
		{"Heading Level 5", 5},
	}

	if len(page.Contents) != len(expectedHeadings) {
		t.Fatalf("Expected %d heading elements, got %d", len(expectedHeadings), len(page.Contents))
	}

	for i, expected := range expectedHeadings {
		content := page.Contents[i]
		if !content.IsHeading() {
			t.Errorf("Content at index %d is not a heading", i)
		}
		if content.Heading != expected.text {
			t.Errorf("Expected heading text at index %d to be '%s', got '%s'", i, expected.text, content.Heading)
		}
		if content.HeadingLevel != expected.level {
			t.Errorf("Expected heading level at index %d to be %d, got %d", i, expected.level, content.HeadingLevel)
		}
	}
}

// TestParsingParagraphs tests parsing a document with paragraphs
func TestParsingParagraphs(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `First paragraph with some text.

Second paragraph with some more text.

Third paragraph with even more text.`

	page := parser.Do("test.org", input)

	expectedParagraphs := []string{
		"First paragraph with some text.",
		"Second paragraph with some more text.",
		"Third paragraph with even more text.",
	}

	if len(page.Contents) != len(expectedParagraphs) {
		t.Fatalf("Expected %d paragraph elements, got %d", len(expectedParagraphs), len(page.Contents))
	}

	for i, expected := range expectedParagraphs {
		content := page.Contents[i]
		if !content.IsParagraph() {
			t.Errorf("Content at index %d is not a paragraph", i)
		}
		if content.Paragraph != expected {
			t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, expected, content.Paragraph)
		}
	}
}

// TestParsingHeadingsAndParagraphs tests parsing a document with a mix of headings and paragraphs
func TestParsingHeadingsAndParagraphs(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `* Title
This is a paragraph below the title.

** Section 1
Section 1 paragraph.

** Section 2
Section 2 first paragraph.

Section 2 second paragraph.`

	page := parser.Do("test.org", input)

	if page.Title != "Title" {
		t.Errorf("Expected title to be 'Title', got '%s'", page.Title)
	}

	expected := []struct {
		isHeading  bool
		text       string
		level      uint32
	}{
		{false, "This is a paragraph below the title.", 0},
		{true, "Section 1", 2},
		{false, "Section 1 paragraph.", 0},
		{true, "Section 2", 2},
		{false, "Section 2 first paragraph.", 0},
		{false, "Section 2 second paragraph.", 0},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if exp.isHeading {
			if !content.IsHeading() {
				t.Errorf("Content at index %d should be a heading", i)
			}
			if content.Heading != exp.text {
				t.Errorf("Expected heading text at index %d to be '%s', got '%s'", i, exp.text, content.Heading)
			}
			if content.HeadingLevel != exp.level {
				t.Errorf("Expected heading level at index %d to be %d, got %d", i, exp.level, content.HeadingLevel)
			}
		} else {
			if !content.IsParagraph() {
				t.Errorf("Content at index %d should be a paragraph", i)
			}
			if content.Paragraph != exp.text {
				t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.text, content.Paragraph)
			}
		}
	}
}

// TestParsingUnorderedLists tests parsing a document with unordered lists
func TestParsingUnorderedLists(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `- Item 1
- Item 2
- Item 3

Paragraph after list.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isList    bool
		paragraph string
		listItems []string
	}{
		{true, "", []string{"Item 1", "Item 2", "Item 3"}},
		{false, "Paragraph after list.", nil},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if exp.isList {
			if !content.IsList() {
				t.Errorf("Content at index %d should be a list", i)
			}
			if len(content.List) != len(exp.listItems) {
				t.Errorf("Expected %d list items at index %d, got %d", len(exp.listItems), i, len(content.List))
			}
			for j, item := range exp.listItems {
				if j < len(content.List) && content.List[j].Text != item {
					t.Errorf("Expected list item at index %d,%d to be '%s', got '%s'", i, j, item, content.List[j].Text)
				}
			}
		} else {
			if !content.IsParagraph() {
				t.Errorf("Content at index %d should be a paragraph", i)
			}
			if content.Paragraph != exp.paragraph {
				t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.paragraph, content.Paragraph)
			}
		}
	}
}

// TestParsingNestedUnorderedLists tests parsing a document with nested unordered lists
func TestParsingNestedUnorderedLists(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `- Item 1
  - Subitem 1.1
  - Subitem 1.2
- Item 2
  - Subitem 2.1
    - Subsubitem 2.1.1
- Item 3

Paragraph after list.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isList    bool
		paragraph string
		listItems []struct {
			text  string
			level uint8
		}
	}{
		{
			isList:    true,
			paragraph: "",
			listItems: []struct {
				text  string
				level uint8
			}{
				{"Item 1", 1},
				{"Subitem 1.1", 2},
				{"Subitem 1.2", 2},
				{"Item 2", 1},
				{"Subitem 2.1", 2},
				{"Subsubitem 2.1.1", 3},
				{"Item 3", 1},
			},
		},
		{
			isList:    false,
			paragraph: "Paragraph after list.",
			listItems: nil,
		},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if exp.isList {
			if !content.IsList() {
				t.Errorf("Content at index %d should be a list", i)
			}
			if len(content.List) != len(exp.listItems) {
				t.Errorf("Expected %d list items at index %d, got %d", len(exp.listItems), i, len(content.List))
			}
			for j, item := range exp.listItems {
				if j < len(content.List) {
					if content.List[j].Text != item.text {
						t.Errorf("Expected list item at index %d,%d to be '%s', got '%s'", i, j, item.text, content.List[j].Text)
					}
					if content.List[j].Level != item.level {
						t.Errorf("Expected list item level at index %d,%d to be %d, got %d", i, j, item.level, content.List[j].Level)
					}
				}
			}
		} else {
			if !content.IsParagraph() {
				t.Errorf("Content at index %d should be a paragraph", i)
			}
			if content.Paragraph != exp.paragraph {
				t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.paragraph, content.Paragraph)
			}
		}
	}
}

// TestParsingOrderedLists tests parsing a document with ordered lists
func TestParsingOrderedLists(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `1. First item
2. Second item
3. Third item

Paragraph after list.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isOrderedList bool
		paragraph     string
		listItems     []string
	}{
		{true, "", []string{"First item", "Second item", "Third item"}},
		{false, "Paragraph after list.", nil},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if exp.isOrderedList {
			if !content.IsListNumbered() {
				t.Errorf("Content at index %d should be an ordered list", i)
			}
			if len(content.List) != len(exp.listItems) {
				t.Errorf("Expected %d list items at index %d, got %d", len(exp.listItems), i, len(content.List))
			}
			for j, item := range exp.listItems {
				if j < len(content.List) && content.List[j].Text != item {
					t.Errorf("Expected list item at index %d,%d to be '%s', got '%s'", i, j, item, content.List[j].Text)
				}
			}
		} else {
			if !content.IsParagraph() {
				t.Errorf("Content at index %d should be a paragraph", i)
			}
			if content.Paragraph != exp.paragraph {
				t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.paragraph, content.Paragraph)
			}
		}
	}
}

// TestParsingTables tests parsing a document with tables
func TestParsingTables(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `| Header 1 | Header 2 | Header 3 |
|-
| Cell 1,1 | Cell 1,2 | Cell 1,3 |
| Cell 2,1 | Cell 2,2 | Cell 2,3 |

Paragraph after table.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isTable      bool
		paragraph    string
		tableHeaders bool
		table        [][]string
	}{
		{
			isTable:      true,
			paragraph:    "",
			tableHeaders: true,
			table: [][]string{
				{"Header 1", "Header 2", "Header 3"},
				{"Cell 1,1", "Cell 1,2", "Cell 1,3"},
				{"Cell 2,1", "Cell 2,2", "Cell 2,3"},
			},
		},
		{
			isTable:      false,
			paragraph:    "Paragraph after table.",
			tableHeaders: false,
			table:        nil,
		},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if exp.isTable {
			if !content.IsTable() {
				t.Errorf("Content at index %d should be a table", i)
			}
			if content.TableHeaders != exp.tableHeaders {
				t.Errorf("Expected TableHeaders at index %d to be %v, got %v", i, exp.tableHeaders, content.TableHeaders)
			}
			if len(content.Table) != len(exp.table) {
				t.Errorf("Expected %d table rows at index %d, got %d", len(exp.table), i, len(content.Table))
				continue
			}
			for j, row := range exp.table {
				if len(content.Table[j]) != len(row) {
					t.Errorf("Expected %d cells in row %d at index %d, got %d", len(row), j, i, len(content.Table[j]))
					continue
				}
				for k, cell := range row {
					if content.Table[j][k] != cell {
						t.Errorf("Expected table cell at [%d][%d][%d] to be '%s', got '%s'", i, j, k, cell, content.Table[j][k])
					}
				}
			}
		} else {
			if !content.IsParagraph() {
				t.Errorf("Content at index %d should be a paragraph", i)
			}
			if content.Paragraph != exp.paragraph {
				t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.paragraph, content.Paragraph)
			}
		}
	}
}

// TestParsingSourceCode tests parsing a document with source code blocks
func TestParsingSourceCode(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+begin_src go
func main() {
	fmt.Println("Hello, World!")
}
#+end_src

Paragraph after source code.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isSourceCode     bool
		paragraph        string
		sourceCode       string
		sourceCodeLang   string
	}{
		{
			isSourceCode:     true,
			paragraph:        "",
			sourceCode:       "func main() {\n\tfmt.Println(\"Hello, World!\")\n}",
			sourceCodeLang:   "go",
		},
		{
			isSourceCode:     false,
			paragraph:        "Paragraph after source code.",
			sourceCode:       "",
			sourceCodeLang:   "",
		},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if exp.isSourceCode {
			if !content.IsSourceCode() {
				t.Errorf("Content at index %d should be source code", i)
			}
			if content.SourceCode != exp.sourceCode {
				t.Errorf("Expected source code at index %d to be '%s', got '%s'", i, exp.sourceCode, content.SourceCode)
			}
			if content.SourceCodeLang != exp.sourceCodeLang {
				t.Errorf("Expected source code language at index %d to be '%s', got '%s'", i, exp.sourceCodeLang, content.SourceCodeLang)
			}
		} else {
			if !content.IsParagraph() {
				t.Errorf("Content at index %d should be a paragraph", i)
			}
			if content.Paragraph != exp.paragraph {
				t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.paragraph, content.Paragraph)
			}
		}
	}
}

// TestParsingRawHtml tests parsing a document with HTML export blocks
func TestParsingRawHtml(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+begin_export html
<div class="custom">
  <p>Custom HTML</p>
</div>
#+end_export

Paragraph after HTML.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isRawHtml  bool
		paragraph  string
		rawHtml    string
	}{
		{
			isRawHtml:  true,
			paragraph:  "",
			rawHtml:    "<div class=\"custom\">\n  <p>Custom HTML</p>\n</div>\n",
		},
		{
			isRawHtml:  false,
			paragraph:  "Paragraph after HTML.",
			rawHtml:    "",
		},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if exp.isRawHtml {
			if !content.IsRawHtml() {
				t.Errorf("Content at index %d should be raw HTML", i)
			}
			if content.RawHtml != exp.rawHtml {
				t.Errorf("Expected raw HTML at index %d to be '%s', got '%s'", i, exp.rawHtml, content.RawHtml)
			}
		} else {
			if !content.IsParagraph() {
				t.Errorf("Content at index %d should be a paragraph", i)
			}
			if content.Paragraph != exp.paragraph {
				t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.paragraph, content.Paragraph)
			}
		}
	}
}

// TestParsingHorizontalLine tests parsing a document with horizontal lines
func TestParsingHorizontalLine(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `Paragraph before.

-----

Paragraph after.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isHorizontalLine bool
		paragraph        string
	}{
		{
			isHorizontalLine: false,
			paragraph:        "Paragraph before.",
		},
		{
			isHorizontalLine: true,
			paragraph:        "",
		},
		{
			isHorizontalLine: false,
			paragraph:        "Paragraph after.",
		},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if exp.isHorizontalLine {
			if !content.IsHorizontalLine() {
				t.Errorf("Content at index %d should be a horizontal line", i)
			}
		} else {
			if !content.IsParagraph() {
				t.Errorf("Content at index %d should be a paragraph", i)
			}
			if content.Paragraph != exp.paragraph {
				t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.paragraph, content.Paragraph)
			}
		}
	}
}

// TestParsingAttentionBlocks tests parsing a document with attention blocks
func TestParsingAttentionBlocks(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `NOTE: This is a note.

WARNING: This is a warning.

IMPORTANT: This is important.

CAUTION: Be careful here.

TIP: This is a helpful tip.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isAttentionBlock bool
		attentionTitle   string
		attentionText    string
	}{
		{
			isAttentionBlock: true,
			attentionTitle:   "NOTE",
			attentionText:    "This is a note. ",
		},
		{
			isAttentionBlock: true,
			attentionTitle:   "WARNING",
			attentionText:    "This is a warning. ",
		},
		{
			isAttentionBlock: true,
			attentionTitle:   "IMPORTANT",
			attentionText:    "This is important. ",
		},
		{
			isAttentionBlock: true,
			attentionTitle:   "CAUTION",
			attentionText:    "Be careful here. ",
		},
		{
			isAttentionBlock: true,
			attentionTitle:   "TIP",
			attentionText:    "This is a helpful tip. ",
		},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if exp.isAttentionBlock {
			if !content.IsAttentionBlock() {
				t.Errorf("Content at index %d should be an attention block", i)
			}
			if content.AttentionTitle != exp.attentionTitle {
				t.Errorf("Expected attention title at index %d to be '%s', got '%s'", i, exp.attentionTitle, content.AttentionTitle)
			}
			if content.AttentionText != exp.attentionText {
				t.Errorf("Expected attention text at index %d to be '%s', got '%s'", i, exp.attentionText, content.AttentionText)
			}
		}
	}
}

// TestParsingOptions tests parsing a document with various options
func TestParsingOptions(t *testing.T) {
	config := &alpha.DarknessConfig{
		RSS: alpha.RSSConfig{
			DefaultAuthor: "Default Author",
		},
	}
	parser := ParserOrgmode{Config: config}

	input := `#+title: Document Title
#+author: John Doe
#+date: 2023-05-15
#+caption: An image caption
#+html_head: <style>body { font-family: Arial; }</style>

This is a paragraph.`

	page := parser.Do("test.org", input)

	if page.Author != "John Doe" {
		t.Errorf("Expected author to be 'John Doe', got '%s'", page.Author)
	}

	if page.Date != "2023-05-15" {
		t.Errorf("Expected date to be '2023-05-15', got '%s'", page.Date)
	}

	if len(page.HtmlHead) != 1 || page.HtmlHead[0] != "<style>body { font-family: Arial; }</style>" {
		t.Errorf("Expected 1 HTML head with specific content, got %v", page.HtmlHead)
	}

	if len(page.Contents) != 1 || !page.Contents[0].IsParagraph() {
		t.Errorf("Expected 1 paragraph content element, got %d elements", len(page.Contents))
	}
}

// TestParsingQuotes tests parsing a document with quote blocks
func TestParsingQuotes(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+begin_quote
This is a quoted text.
It spans multiple lines.
#+end_quote

This is a normal paragraph.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isParagraph bool
		isQuote     bool
		text        string
	}{
		{
			isParagraph: true,
			isQuote:     true,
			text:        "This is a quoted text. It spans multiple lines.",
		},
		{
			isParagraph: true,
			isQuote:     false,
			text:        "This is a normal paragraph.",
		},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if content.IsParagraph() != exp.isParagraph {
			t.Errorf("Content at index %d paragraph status mismatch", i)
		}

		if content.IsQuote() != exp.isQuote {
			t.Errorf("Content at index %d quote status mismatch, expected %v, got %v", i, exp.isQuote, content.IsQuote())
		}

		if exp.isParagraph && content.Paragraph != exp.text {
			t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.text, content.Paragraph)
		}
	}
}

// TestParsingCenters tests parsing a document with center blocks
func TestParsingCenters(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+begin_center
This text should be centered.
It spans multiple lines.
#+end_center

This is a normal paragraph.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isParagraph bool
		isCentered  bool
		text        string
	}{
		{
			isParagraph: true,
			isCentered:  true,
			text:        "This text should be centered. It spans multiple lines.",
		},
		{
			isParagraph: true,
			isCentered:  false,
			text:        "This is a normal paragraph.",
		},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if content.IsParagraph() != exp.isParagraph {
			t.Errorf("Content at index %d paragraph status mismatch", i)
		}

		if content.IsCentered() != exp.isCentered {
			t.Errorf("Content at index %d centered status mismatch, expected %v, got %v", i, exp.isCentered, content.IsCentered())
		}

		if exp.isParagraph && content.Paragraph != exp.text {
			t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.text, content.Paragraph)
		}
	}
}

// TestParsingDetails tests parsing a document with details blocks
func TestParsingDetails(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+begin_details Custom Summary
This is hidden content inside a details block.
Multiple lines of content.
#+end_details

This is a normal paragraph.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isDetails   bool
		isParagraph bool
		isStart     bool
		summary     string
		text        string
	}{
		{isDetails: true, isParagraph: false, isStart: true, summary: "Custom Summary", text: ""},
		{isDetails: true, isParagraph: true, isStart: false, summary: "", text: "This is hidden content inside a details block. Multiple lines of content."},
		{isDetails: false, isParagraph: false, isStart: false, summary: "", text: ""},
		{isDetails: false, isParagraph: true, isStart: false, summary: "", text: "This is a normal paragraph."},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if content.IsDetails() != exp.isDetails {
			t.Errorf("Content at index %d details status mismatch, expected %v, got %v", i, exp.isDetails, content.IsDetails())
		}

		if content.IsParagraph() != exp.isParagraph {
			t.Errorf("Content at index %d paragraph status mismatch, expected %v, got %v", i, exp.isParagraph, content.IsParagraph())
		}

		if content.Type == yunyun.TypeDetails && i == 0 && content.Summary != "Custom Summary" {
			t.Errorf("Expected summary for details block to be 'Custom Summary', got '%s'", content.Summary)
		}

		if exp.isParagraph && content.Paragraph != exp.text {
			t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.text, content.Paragraph)
		}
	}
}

// TestParsingLinks tests parsing a document with links
func TestParsingLinks(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	yunyun.ActiveMarkings.BuildRegex()
	linkRegexp = yunyun.LinkRegexp

	input := `[[https://example.com][Example Website]]

[[https://example.org]]

This is a paragraph with a [[https://example.net][link]] inside.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isLink       bool
		isParagraph  bool
		link         string
		linkTitle    string
		paragraphHas string
	}{
		{
			isLink:       true,
			isParagraph:  false,
			link:         "https://example.com",
			linkTitle:    "Example Website",
			paragraphHas: "",
		},
		{
			isLink:       true,
			isParagraph:  false,
			link:         "https://example.org",
			linkTitle:    "",
			paragraphHas: "",
		},
		{
			isLink:       false,
			isParagraph:  true,
			link:         "",
			linkTitle:    "",
			paragraphHas: "This is a paragraph with a [[https://example.net][link]] inside.",
		},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		if content.IsLink() != exp.isLink {
			t.Errorf("Content at index %d link status mismatch, expected %v, got %v", i, exp.isLink, content.IsLink())
		}

		if content.IsParagraph() != exp.isParagraph {
			t.Errorf("Content at index %d paragraph status mismatch, expected %v, got %v", i, exp.isParagraph, content.IsParagraph())
		}

		if exp.isLink {
			if content.Link != exp.link {
				t.Errorf("Expected link at index %d to be '%s', got '%s'", i, exp.link, content.Link)
			}

			if content.LinkTitle != exp.linkTitle {
				t.Errorf("Expected link title at index %d to be '%s', got '%s'", i, exp.linkTitle, content.LinkTitle)
			}
		}

		if exp.isParagraph && content.Paragraph != exp.paragraphHas {
			t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.paragraphHas, content.Paragraph)
		}
	}
}

// TestParsingGallery tests parsing a document with a gallery block
func TestParsingGallery(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+begin_gallery :path images/gallery :num 4
This text is inside the gallery.
#+end_gallery

This is a normal paragraph.`

	page := parser.Do("test.org", input)

	if len(page.Contents) < 2 {
		t.Fatalf("Expected at least 2 content elements, got %d", len(page.Contents))
	}

	// First content is the gallery paragraph
	if !page.Contents[0].IsParagraph() {
		t.Errorf("Expected first content to be a paragraph")
	}

	if !page.Contents[0].IsGallery() {
		t.Errorf("Expected first content to be marked as gallery")
	}

	if string(page.Contents[0].GalleryPath) != "images/gallery" {
		t.Errorf("Expected gallery path to be 'images/gallery', got '%s'", page.Contents[0].GalleryPath)
	}

	if page.Contents[0].GalleryImagesPerRow != 4 {
		t.Errorf("Expected gallery images per row to be 4, got %d", page.Contents[0].GalleryImagesPerRow)
	}

	// Second content is a normal paragraph
	if !page.Contents[1].IsParagraph() || page.Contents[1].IsGallery() {
		t.Errorf("Expected second content to be a normal paragraph")
	}

	if page.Contents[1].Paragraph != "This is a normal paragraph." {
		t.Errorf("Expected second paragraph to be 'This is a normal paragraph.', got '%s'", page.Contents[1].Paragraph)
	}
}

// TestParsingComments tests that comments are ignored
func TestParsingComments(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `# This is a comment and should be ignored
This is a normal paragraph.
# Another comment
Another paragraph.`

	page := parser.Do("test.org", input)

	// The parser combines adjacent paragraphs/text into a single paragraph
	expected := []string{
		"This is a normal paragraph. Another paragraph.",
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		if !page.Contents[i].IsParagraph() {
			t.Errorf("Expected content at index %d to be a paragraph", i)
		}
		if page.Contents[i].Paragraph != exp {
			t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp, page.Contents[i].Paragraph)
		}
	}
}

// TestParsingHolosceneDate tests the holoscene date detection
func TestParsingHolosceneDate(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `127; 12024 H.E.

This is a normal paragraph.`

	page := parser.Do("test.org", input)

	if !page.DateHoloscene {
		t.Errorf("Expected DateHoloscene to be true")
	}

	if page.Date != "127; 12024 H.E." {
		t.Errorf("Expected Date to be '127; 12024 H.E.', got '%s'", page.Date)
	}

	// Two content elements: one for date paragraph and one for the second paragraph
	if len(page.Contents) != 2 {
		t.Errorf("Expected 2 content elements, got %d", len(page.Contents))
	}
}

// TestParsingMixedContent tests parsing a document with a mix of different types of content
func TestParsingMixedContent(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `* Main Title

Introduction paragraph.

** Section 1

#+begin_quote
This is a quoted text.
#+end_quote

- List item 1
- List item 2
  - Subitem 2.1

** Section 2

| Header 1 | Header 2 |
|-
| Value 1  | Value 2  |

#+begin_src python
def hello():
    print("Hello, World!")
#+end_src

WARNING: This is a warning message.

-----

** Conclusion

Final paragraph.`

	page := parser.Do("test.org", input)

	// Check title
	if page.Title != "Main Title" {
		t.Errorf("Expected title to be 'Main Title', got '%s'", page.Title)
	}

	// We should have many content elements (paragraph, headings, quote, list, table, source code, attention block, horizontal line, paragraph)
	if len(page.Contents) < 10 {
		t.Fatalf("Expected at least 10 content elements, got %d", len(page.Contents))
	}

	// Check types in order
	contentTypes := []struct {
		index int
		check func(*yunyun.Content) bool
		desc  string
	}{
		{0, func(c *yunyun.Content) bool { return c.IsParagraph() }, "paragraph"},
		{1, func(c *yunyun.Content) bool { return c.IsHeading() && c.Heading == "Section 1" }, "Section 1 heading"},
		{2, func(c *yunyun.Content) bool { return c.IsParagraph() && c.IsQuote() }, "quote"},
		{3, func(c *yunyun.Content) bool { return c.IsList() }, "list"},
		{4, func(c *yunyun.Content) bool { return c.IsHeading() && c.Heading == "Section 2" }, "Section 2 heading"},
		{5, func(c *yunyun.Content) bool { return c.IsTable() }, "table"},
		{6, func(c *yunyun.Content) bool { return c.IsSourceCode() }, "source code"},
		{7, func(c *yunyun.Content) bool { return c.IsAttentionBlock() }, "attention block"},
		{8, func(c *yunyun.Content) bool { return c.IsHorizontalLine() }, "horizontal line"},
		{9, func(c *yunyun.Content) bool { return c.IsHeading() && c.Heading == "Conclusion" }, "Conclusion heading"},
		{10, func(c *yunyun.Content) bool { return c.IsParagraph() }, "final paragraph"},
	}

	for _, ct := range contentTypes {
		if ct.index < len(page.Contents) {
			if !ct.check(page.Contents[ct.index]) {
				t.Errorf("Content at index %d should be %s", ct.index, ct.desc)
			}
		} else {
			t.Errorf("Content at index %d not found, expected %s", ct.index, ct.desc)
		}
	}
}

// TestParsingMacros tests parsing a document with macro definitions and usage
func TestParsingMacros(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+macro: greeting Hello, *World*!
#+macro: signature Best regards, Author

{{{greeting}}}

This is a normal paragraph.

{{{signature}}}`

	page := parser.Do("test.org", input)

	// There should be at least two paragraphs
	if len(page.Contents) < 2 {
		t.Fatalf("Expected at least 2 content elements, got %d", len(page.Contents))
	}

	// Check if we have the expanded greeting in one of the paragraphs
	foundGreeting := false
	foundSignature := false

	for _, content := range page.Contents {
		if content.IsParagraph() {
			if content.Paragraph == "Hello, *World*!" {
				foundGreeting = true
			}
			if strings.Contains(content.Paragraph, "Best regards, Author") {
				foundSignature = true
			}
		}
	}

	if !foundGreeting {
		t.Errorf("Did not find paragraph with expanded greeting 'Hello, *World*!'")
	}

	if !foundSignature {
		t.Errorf("Did not find paragraph with expanded signature 'Best regards, Author'")
	}
}

// TestParsingMacrosWithParameters tests macros with parameters
func TestParsingMacrosWithParameters(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+macro: greet Hello, $1!

{{{greet(World)}}}

{{{greet(Friend)}}}`

	page := parser.Do("test.org", input)

	// There should be two paragraphs - one with each expanded greeting
	if len(page.Contents) != 2 {
		t.Fatalf("Expected 2 content elements, got %d", len(page.Contents))
	}

	// The first paragraph should contain "Hello, World!"
	if !page.Contents[0].IsParagraph() || page.Contents[0].Paragraph != "Hello, World!" {
		t.Errorf("Expected first paragraph to be 'Hello, World!', got '%s'", page.Contents[0].Paragraph)
	}

	// The second paragraph should contain "Hello, Friend!"
	if !page.Contents[1].IsParagraph() || page.Contents[1].Paragraph != "Hello, Friend!" {
		t.Errorf("Expected second paragraph to be 'Hello, Friend!', got '%s'", page.Contents[1].Paragraph)
	}
}

// TestParsingNestedStructures tests nested structures like lists within quotes, etc.
func TestParsingNestedStructures(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+begin_quote
This is a quoted text with a list:

- Item 1
- Item 2
- Item 3
#+end_quote

#+begin_center
This is centered text with a table:

| Header 1 | Header 2 |
|-
| Cell 1,1 | Cell 1,2 |
#+end_center`

	page := parser.Do("test.org", input)

	// The structure is complex but we expect at least 2 content elements
	if len(page.Contents) < 2 {
		t.Fatalf("Expected at least 2 content elements, got %d", len(page.Contents))
	}

	// First section: a quote containing text and a list
	foundQuote := false
	foundListInQuote := false

	// Second section: a centered block with text and a table
	foundCenter := false
	foundTableInCenter := false

	for _, content := range page.Contents {
		if content.IsQuote() && content.IsParagraph() && strings.Contains(content.Paragraph, "quoted text with a list") {
			foundQuote = true
		}
		if content.IsQuote() && content.IsList() {
			foundListInQuote = true
		}
		if content.IsCentered() && content.IsParagraph() && strings.Contains(content.Paragraph, "centered text with a table") {
			foundCenter = true
		}
		if content.IsCentered() && content.IsTable() {
			foundTableInCenter = true
		}
	}

	// Check that we found all the expected elements
	if !foundQuote {
		t.Errorf("Did not find a quoted paragraph")
	}
	if !foundListInQuote {
		t.Errorf("Did not find a list in a quote")
	}
	if !foundCenter {
		t.Errorf("Did not find a centered paragraph")
	}
	if !foundTableInCenter {
		t.Errorf("Did not find a table in a centered block")
	}
}

// TestParsingEdgeCases tests various edge cases and less common structures
func TestParsingEdgeCases(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	inputs := []struct {
		name     string
		input    string
		validate func(*yunyun.Page)
	}{
		{
			name: "Empty lines between paragraphs",
			input: `Paragraph 1.



Paragraph 2.`,
			validate: func(page *yunyun.Page) {
				if len(page.Contents) != 2 {
					t.Errorf("Expected 2 paragraphs, got %d", len(page.Contents))
				}
			},
		},
		{
			name:  "Empty document",
			input: ``,
			validate: func(page *yunyun.Page) {
				if len(page.Contents) != 0 {
					t.Errorf("Expected 0 elements, got %d", len(page.Contents))
				}
			},
		},
		{
			name:  "Document with only whitespace",
			input: `   
   
  `,
			validate: func(page *yunyun.Page) {
				if len(page.Contents) != 0 {
					t.Errorf("Expected 0 elements, got %d", len(page.Contents))
				}
			},
		},
		{
			name:  "Document with only comments",
			input: `# Comment 1
# Comment 2
# Comment 3`,
			validate: func(page *yunyun.Page) {
				if len(page.Contents) != 0 {
					t.Errorf("Expected 0 elements, got %d", len(page.Contents))
				}
			},
		},
		{
			name:  "List with only one item",
			input: `- Single item list`,
			validate: func(page *yunyun.Page) {
				if len(page.Contents) != 1 || !page.Contents[0].IsList() {
					t.Errorf("Expected 1 list element, got something else")
				}
				if len(page.Contents[0].List) != 1 {
					t.Errorf("Expected list with 1 item, got %d", len(page.Contents[0].List))
				}
			},
		},
		{
			name: "Nested blocks",
			input: `#+begin_quote
#+begin_center
Centered inside a quote
#+end_center
#+end_quote`,
			validate: func(page *yunyun.Page) {
				if len(page.Contents) < 1 {
					t.Errorf("Expected at least 1 element, got %d", len(page.Contents))
				}
				// The specific behavior here depends on how the parser handles nested blocks
				// We just check that something was parsed successfully
			},
		},
	}

	for _, tc := range inputs {
		t.Run(tc.name, func(t *testing.T) {
			page := parser.Do("test.org", tc.input)
			if page == nil {
				t.Fatalf("Parser returned nil page for input: %s", tc.name)
			}
			tc.validate(page)
		})
	}
}

// TestParsingCombinations tests combinations of multiple features together
func TestParsingCombinations(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	inputs := []struct {
		name     string
		input    string
		validate func(*yunyun.Page)
	}{
		{
			name: "Heading with option",
			input: `* Heading
#+caption: This is a caption

Paragraph under heading.`,
			validate: func(page *yunyun.Page) {
				if page.Title != "Heading" {
					t.Errorf("Expected title 'Heading', got '%s'", page.Title)
				}
				if len(page.Contents) != 1 || !page.Contents[0].IsParagraph() {
					t.Errorf("Expected one paragraph content element")
				}
			},
		},
		{
			name: "Quote with source code",
			input: `#+begin_quote
Example code:

#+begin_src go
func main() {
	fmt.Println("Hello")
}
#+end_src
#+end_quote`,
			validate: func(page *yunyun.Page) {
				foundQuote := false
				foundCode := false
				for _, c := range page.Contents {
					if c.IsQuote() && c.IsParagraph() {
						foundQuote = true
					}
					if c.IsSourceCode() {
						foundCode = true
					}
				}
				if !foundQuote {
					t.Errorf("Expected to find a quote")
				}
				if !foundCode {
					t.Errorf("Expected to find source code")
				}
			},
		},
	}

	for _, tc := range inputs {
		t.Run(tc.name, func(t *testing.T) {
			page := parser.Do("test.org", tc.input)
			if page == nil {
				t.Fatalf("Parser returned nil page for input: %s", tc.name)
			}
			tc.validate(page)
		})
	}
}