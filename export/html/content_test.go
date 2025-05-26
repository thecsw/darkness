package html

import (
	"io"
	"strings"
	"testing"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/yunyun"
)

// TestDetailsWithSourceCodeExport tests that source code blocks inside details
// are properly exported to HTML
func TestDetailsWithSourceCodeExport(t *testing.T) {
	// Create a page with details and source code blocks
	page := &yunyun.Page{
		Title: "Test Page",
		Contents: []*yunyun.Content{
			// Test a details block opening
			{
				Type:    yunyun.TypeDetails,
				Summary: "Click to see code",
				Options: yunyun.InDetailsFlag,
			},
			// Test a source code block inside a details block
			{
				Type:          yunyun.TypeSourceCode,
				SourceCode:    "func main() {\n\tfmt.Println(\"Hello from inside details!\")\n}",
				SourceCodeLang: "go",
				Options:       yunyun.InDetailsFlag,
			},
			// Test a closing details block
			{
				Type: yunyun.TypeDetails,
			},
		},
	}

	// Create exporter and export the page
	exporter := ExporterHtml{
		Config: &alpha.DarknessConfig{},
	}
	
	reader := exporter.Do(page)
	
	// Read the output
	buf := new(strings.Builder)
	_, err := io.Copy(buf, reader)
	if err != nil {
		t.Fatalf("Failed to read exported content: %v", err)
	}
	
	exportedHTML := buf.String()
	
	// Convert to lowercase for case-insensitive checks
	lowerHTML := strings.ToLower(exportedHTML)

	// Verify that:
	// 1. The <details> tag is opened
	if !strings.Contains(lowerHTML, "<details>") {
		t.Error("Missing <details> opening tag")
	}

	// 2. The details summary is included
	if !strings.Contains(lowerHTML, "<summary>click to see code</summary>") {
		t.Error("Missing or incorrect summary tag")
	}

	// 3. The source code is within the details section
	if !strings.Contains(lowerHTML, "class=\"language-go\" data-lang=\"go\"") {
		t.Error("Missing or incorrect source code language")
	}

	// 4. The source code content is properly included
	if !strings.Contains(lowerHTML, "fmt.println(") {
		t.Error("Missing source code content")
	}

	// 5. The </details> tag is closed
	if !strings.Contains(lowerHTML, "</details>") {
		t.Error("Missing </details> closing tag")
	}

	// The problem we're verifying: source code should be properly contained within the details block
	detailsOpenPos := strings.Index(lowerHTML, "<details>")
	sourceCodeStartPos := strings.Index(lowerHTML, "class=\"language-go\"")
	detailsClosePos := strings.Index(lowerHTML, "</details>")

	// Check the structure: <details> should come before source code, which should come before </details>
	if detailsOpenPos == -1 || sourceCodeStartPos == -1 || detailsClosePos == -1 {
		t.Error("Missing essential HTML elements in output")
	} else if !(detailsOpenPos < sourceCodeStartPos && sourceCodeStartPos < detailsClosePos) {
		t.Error("Source code block is not properly nested within details block in the HTML output")
	}
}
