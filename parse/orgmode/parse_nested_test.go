package orgmode

import (
	"strings"
	"testing"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/yunyun"
)

// TestParsingSourceCodeInDetails tests parsing a document with source code blocks inside details blocks
func TestParsingSourceCodeInDetails(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+begin_details Click to see code
This is some text before the code block.

#+begin_src go
func main() {
	fmt.Println("Hello from inside details!")
}
#+end_src

More content within details.

#+begin_src python
def hello():
    print("Python inside details!")
#+end_src
#+end_details

This is a normal paragraph.`

	page := parser.Do("test.org", input)

	expected := []struct {
		isDetails     bool
		isSourceCode  bool
		isParagraph   bool
		isStart       bool
		summary       string
		text          string
		sourceCode    string
		sourceCodeLang string
	}{
		{isDetails: true, isSourceCode: false, isParagraph: false, isStart: true, summary: "Click to see code", text: "", sourceCode: "", sourceCodeLang: ""},
		{isDetails: true, isSourceCode: false, isParagraph: true, isStart: false, summary: "", text: "This is some text before the code block.", sourceCode: "", sourceCodeLang: ""},
		{isDetails: true, isSourceCode: true, isParagraph: false, isStart: false, summary: "", text: "", sourceCode: "func main() {\n\tfmt.Println(\"Hello from inside details!\")\n}", sourceCodeLang: "go"},
		{isDetails: true, isSourceCode: false, isParagraph: true, isStart: false, summary: "", text: "More content within details.", sourceCode: "", sourceCodeLang: ""},
		{isDetails: true, isSourceCode: true, isParagraph: false, isStart: false, summary: "", text: "", sourceCode: "def hello():\n    print(\"Python inside details!\")", sourceCodeLang: "python"},
		{isDetails: false, isSourceCode: false, isParagraph: false, isStart: false, summary: "", text: "", sourceCode: "", sourceCodeLang: ""},
		{isDetails: false, isSourceCode: false, isParagraph: true, isStart: false, summary: "", text: "This is a normal paragraph.", sourceCode: "", sourceCodeLang: ""},
	}

	if len(page.Contents) != len(expected) {
		t.Fatalf("Expected %d content elements, got %d", len(expected), len(page.Contents))
	}

	for i, exp := range expected {
		content := page.Contents[i]
		
		// Check if this content should be part of details block
		if content.IsDetails() != exp.isDetails {
			t.Errorf("Content at index %d details status mismatch, expected %v, got %v", i, exp.isDetails, content.IsDetails())
		}

		// Check if this content is a source code block
		if content.IsSourceCode() != exp.isSourceCode {
			t.Errorf("Content at index %d source code status mismatch, expected %v, got %v", i, exp.isSourceCode, content.IsSourceCode())
		}

		// Check if this content is a paragraph
		if content.IsParagraph() != exp.isParagraph {
			t.Errorf("Content at index %d paragraph status mismatch, expected %v, got %v", i, exp.isParagraph, content.IsParagraph())
		}

		// Check details summary
		if i == 0 && content.Type == yunyun.TypeDetails && content.Summary != exp.summary {
			t.Errorf("Expected summary for details block to be '%s', got '%s'", exp.summary, content.Summary)
		}

		// Check paragraph text
		if exp.isParagraph && content.Paragraph != exp.text {
			t.Errorf("Expected paragraph text at index %d to be '%s', got '%s'", i, exp.text, content.Paragraph)
		}

		// Check source code
		if exp.isSourceCode {
			if content.SourceCode != exp.sourceCode {
				t.Errorf("Expected source code at index %d to be '%s', got '%s'", i, exp.sourceCode, content.SourceCode)
			}

			if content.SourceCodeLang != exp.sourceCodeLang {
				t.Errorf("Expected source code language at index %d to be '%s', got '%s'", i, exp.sourceCodeLang, content.SourceCodeLang)
			}
		}
	}
}

// TestParsingComplexDetailsWithSourceCode tests parsing complex details blocks with nested elements
func TestParsingComplexDetailsWithSourceCode(t *testing.T) {
	config := &alpha.DarknessConfig{}
	parser := ParserOrgmode{Config: config}

	input := `#+begin_details Complex nested example
This is text at the start of details.

#+begin_src go
func example() {
	// Code inside details
}
#+end_src

#+begin_quote
This is a quote inside details.
#+end_quote

#+begin_center
Centered text inside details
#+end_center

#+begin_src python
def nested():
    """Source code after other blocks"""
    return True
#+end_src
#+end_details

Normal paragraph after details.`

	page := parser.Do("test.org", input)

	// We'll check specific aspects rather than the full structure
	foundDetailsStart := false
	foundTextInDetails := false
	foundGoCode := false
	foundQuoteInDetails := false
	foundCenteredInDetails := false
	foundPythonCode := false
	foundDetailsEnd := false
	foundNormalParagraph := false

	for i, content := range page.Contents {
		// Check details flags are preserved for nested elements
		if i > 0 && i < len(page.Contents)-2 {
			if !content.IsDetails() {
				t.Errorf("Content at index %d should be part of details block", i)
			}
		}

		if content.Type == yunyun.TypeDetails && !content.IsParagraph() {
			if i == 0 {
				foundDetailsStart = true
				if content.Summary != "Complex nested example" {
					t.Errorf("Expected details summary to be 'Complex nested example', got '%s'", content.Summary)
				}
			} else {
				foundDetailsEnd = true
			}
		}

		if content.IsParagraph() && content.IsDetails() && content.Paragraph == "This is text at the start of details." {
			foundTextInDetails = true
		}

		if content.IsSourceCode() && content.IsDetails() {
			if content.SourceCodeLang == "go" {
				foundGoCode = true
				if content.SourceCode != "func example() {\n\t// Code inside details\n}" {
					t.Errorf("Go source code doesn't match expected content")
				}
			} else if content.SourceCodeLang == "python" {
				foundPythonCode = true
				if !strings.Contains(content.SourceCode, "Source code after other blocks") {
					t.Errorf("Python source code doesn't match expected content")
				}
			}
		}

		if content.IsQuote() && content.IsDetails() {
			foundQuoteInDetails = true
		}

		if content.IsCentered() && content.IsDetails() {
			foundCenteredInDetails = true
		}

		if content.IsParagraph() && !content.IsDetails() && content.Paragraph == "Normal paragraph after details." {
			foundNormalParagraph = true
		}
	}

	if !foundDetailsStart {
		t.Error("Missing details block start")
	}
	if !foundTextInDetails {
		t.Error("Missing text inside details")
	}
	if !foundGoCode {
		t.Error("Missing Go source code in details")
	}
	if !foundQuoteInDetails {
		t.Error("Missing quote in details")
	}
	if !foundCenteredInDetails {
		t.Error("Missing centered text in details")
	}
	if !foundPythonCode {
		t.Error("Missing Python source code in details")
	}
	if !foundDetailsEnd {
		t.Error("Missing details block end")
	}
	if !foundNormalParagraph {
		t.Error("Missing normal paragraph after details")
	}
}