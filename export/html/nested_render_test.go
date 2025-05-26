package html

import (
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/parse/orgmode"
)

func TestNestedSourceCodeInDetailsRendering(t *testing.T) {
	// Create a simple test document with a details block containing source code
	input := `#+begin_details Click to see code
This is some text before the code block.

#+begin_src go
func main() {
	fmt.Println("Hello from inside details!")
}
#+end_src

#+end_details

This is outside the details block.`

	// Parse the document using the orgmode parser
	parser := orgmode.ParserOrgmode{}
	page := parser.Do("test.org", input)
	
	// Export to HTML
	exporter := ExporterHtml{
		Config: &alpha.DarknessConfig{},
	}
	htmlReader := exporter.Do(page)
	
	// Convert to string for testing
	buf := new(strings.Builder)
	_, err := io.Copy(buf, htmlReader)
	if err != nil {
		t.Fatalf("Failed to read HTML output: %v", err)
	}
	
	htmlOutput := buf.String()
	
	// Check for nesting - details should appear before source code which should appear before closing details
	detailsPattern := regexp.MustCompile(`(?s)<details>.*?<summary>Click to see code</summary>.*?<div class="coding".*?</details>`)
	
	if !detailsPattern.MatchString(htmlOutput) {
		t.Error("Source code is not properly nested inside details block in HTML")
		
		// If the test fails, help debug by showing key segment positions
		t.Logf("Details open position: %d", strings.Index(htmlOutput, "<details>"))
		t.Logf("Source code position: %d", strings.Index(htmlOutput, "class=\"language-go\""))
		t.Logf("Details close position: %d", strings.Index(htmlOutput, "</details>"))
	}
}