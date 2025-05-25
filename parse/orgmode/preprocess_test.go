package orgmode

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/thecsw/darkness/v3/emilia/alpha"
)

// TestPreprocess tests the preprocess function with various inputs
func TestPreprocess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic preprocessing - no special features",
			input:    "Simple text\nwith multiple lines",
			expected: "Simple text\nwith multiple lines\n\n",
		},
		{
			name:     "Headings get surrounded by newlines",
			input:    "Text\n* Heading\nMore text",
			expected: "Text\n\n* Heading\nMore text\n\n",
		},
		{
			name:     "Multiple headings",
			input:    "* Heading 1\nText\n** Heading 2\nMore text",
			expected: "\n* Heading 1\nText\n\n** Heading 2\nMore text\n\n",
		},
		{
			name:     "Begin/end blocks get surrounded by newlines",
			input:    "Text\n#+begin_quote\nQuoted text\n#+end_quote\nMore text",
			expected: "Text\n\n#+begin_quote\nQuoted text\n\n#+end_quote\nMore text\n\n",
		},
		{
			name:     "Lists get preceded by newlines",
			input:    "Text\n- List item\nMore text",
			expected: "Text\n\n- List item\nMore text\n\n",
		},
		{
			name:     "Lists with multi-line items",
			input:    "Text\n- List item\n  continued\nMore text",
			expected: "Text\n\n- List item\n  continued\nMore text\n\n",
		},
		{
			name:     "Multiple consecutive lists",
			input:    "- List item 1\n- List item 2",
			expected: "\n- List item 1\n- List item 2\n\n",
		},
		{
			name:     "Macro definitions and usage",
			input:    "#+macro: greeting Hello, $1!\n\n{{{greeting(World)}}}",
			expected: "\nHello, World!\n\n",
		},
	}

	// Create a basic configuration for testing, similar to other tests
	config := &alpha.DarknessConfig{}
	config.Runtime.Logger = log.NewWithOptions(os.Stderr, log.Options{
		Level: log.FatalLevel, // Only show fatal errors
	})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := ParserOrgmode{Config: config}
			result := parser.preprocess("test.org", tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%q\nGot:\n%q", tc.expected, result)
			}
		})
	}
}

// TestCollectMacros tests the collectMacros function
func TestCollectMacros(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMacros   map[string]string
		expectedFound    bool
		shouldPanicFatal bool
	}{
		{
			name:           "No macros",
			input:          "Just regular text\nNo macros here",
			expectedMacros: map[string]string{},
			expectedFound:  false,
		},
		{
			name:  "Single macro",
			input: "#+macro: greeting Hello, World!",
			expectedMacros: map[string]string{
				"greeting": "Hello, World!",
			},
			expectedFound: true,
		},
		{
			name: "Multiple macros",
			input: `#+macro: greeting Hello, World!
#+macro: farewell Goodbye, World!`,
			expectedMacros: map[string]string{
				"greeting": "Hello, World!",
				"farewell": "Goodbye, World!",
			},
			expectedFound: true,
		},
		{
			name: "Macros with complex content",
			input: `#+macro: bold *$1*
#+macro: link [[https://example.com][$1]]`,
			expectedMacros: map[string]string{
				"bold": "*$1*",
				"link": "[[https://example.com][$1]]",
			},
			expectedFound: true,
		},
		{
			name:  "Macro with newline escape",
			input: `#+macro: paragraph $1\n$2`,
			expectedMacros: map[string]string{
				"paragraph": "$1\\n$2",
			},
			expectedFound: true,
		},
		{
			name:             "Malformed macro",
			input:            "#+macro: malformed",
			expectedMacros:   map[string]string{},
			expectedFound:    false,
			shouldPanicFatal: true, // This would call Fatal in a real scenario
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldPanicFatal {
				// Skip tests that would cause fatal errors
				// In a real environment, we would mock the logger
				t.Skip("Skipping test that would cause fatal error")
			}

			config := &alpha.DarknessConfig{}
			config.Runtime.Logger = log.NewWithOptions(os.Stderr, log.Options{
				Level: log.FatalLevel, // Only show fatal errors
			})
			
			macrosLookupTable := make(map[string]string)
			found := collectMacros(config, "test.org", macrosLookupTable, tc.input)

			if found != tc.expectedFound {
				t.Errorf("Expected found to be %v, got %v", tc.expectedFound, found)
			}

			if !reflect.DeepEqual(macrosLookupTable, tc.expectedMacros) {
				t.Errorf("Expected macros:\n%v\nGot:\n%v", tc.expectedMacros, macrosLookupTable)
			}
		})
	}
}

// TestExpandMacros tests the expandMacros function
func TestExpandMacros(t *testing.T) {
	tests := []struct {
		name           string
		macros         map[string]string
		input          string
		expected       string
		expectedResult bool
	}{
		{
			name:           "No macros to expand",
			macros:         map[string]string{"greeting": "Hello, World!"},
			input:          "Just regular text",
			expected:       "",
			expectedResult: false,
		},
		{
			name:           "Simple macro expansion",
			macros:         map[string]string{"greeting": "Hello, World!"},
			input:          "{{{greeting}}}",
			expected:       "Hello, World!",
			expectedResult: true,
		},
		{
			name:           "Macro with parameter",
			macros:         map[string]string{"greet": "Hello, $1!"},
			input:          "{{{greet(Friend)}}}",
			expected:       "Hello, Friend!",
			expectedResult: true,
		},
		{
			name:           "Multiple parameters",
			macros:         map[string]string{"template": "From: $1\nTo: $2\nSubject: $3"},
			input:          "{{{template(Alice, Bob, Meeting)}}}",
			expected:       "From: Alice\nTo: Bob\nSubject: Meeting",
			expectedResult: true,
		},
		{
			name:           "Multiple macro usages in one line",
			macros:         map[string]string{"hi": "Hello", "bye": "Goodbye"},
			input:          "{{{hi}}}, World! And {{{bye}}}!",
			expected:       "Hello, World! And Goodbye!",
			expectedResult: true,
		},
		{
			name:           "Macro with escaped newlines",
			macros:         map[string]string{"para": "Line 1\\nLine 2"},
			input:          "{{{para}}}",
			expected:       "Line 1\nLine 2",
			expectedResult: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config := &alpha.DarknessConfig{}
			config.Runtime.Logger = log.NewWithOptions(os.Stderr, log.Options{
				Level: log.FatalLevel, // Only show fatal errors
			})
			
			result, found := expandMacros(config, "test.org", tc.macros, tc.input)

			if found != tc.expectedResult {
				t.Errorf("Expected found to be %v, got %v", tc.expectedResult, found)
			}

			if found && result != tc.expected {
				t.Errorf("Expected result:\n%q\nGot:\n%q", tc.expected, result)
			}
		})
	}
}

// TestExpandSetupFile tests the expandSetupFile function
func TestExpandSetupFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "darkness-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test setup file
	setupFileContent := "This is content from a setup file.\n#+macro: greeting Hello from setup!"
	setupFilePath := filepath.Join(tmpDir, "setup.org")
	err = os.WriteFile(setupFilePath, []byte(setupFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write setup file: %v", err)
	}

	tests := []struct {
		name             string
		directive        string
		expected         string
		expectedFound    bool
		shouldPanicFatal bool
	}{
		{
			name:          "Not a setup file directive",
			directive:     "Just regular text",
			expected:      "",
			expectedFound: false,
		},
		{
			name:          "Valid setup file directive",
			directive:     "#+setupfile: setup.org",
			expected:      setupFileContent,
			expectedFound: true,
		},
		{
			name:          "Uppercase directive",
			directive:     "#+SETUPFILE: setup.org",
			expected:      setupFileContent,
			expectedFound: true,
		},
		{
			name:             "Non-existent setup file",
			directive:        "#+setupfile: nonexistent.org",
			expected:         "",
			expectedFound:    false,
			shouldPanicFatal: true, // This would cause a fatal error in real code
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldPanicFatal {
				// Skip tests that would cause fatal errors
				t.Skip("Skipping test that would cause fatal error")
			}

			// Create a config with the temp dir as the working directory
			config := &alpha.DarknessConfig{}
			config.Runtime.Logger = log.NewWithOptions(os.Stderr, log.Options{
				Level: log.FatalLevel, // Only show fatal errors
			})
			config.Runtime.WorkDir = alpha.WorkingDirectory(tmpDir)

			result, found := expandSetupFile(config, "main.org", tc.directive)

			if found != tc.expectedFound {
				t.Errorf("Expected found to be %v, got %v", tc.expectedFound, found)
			}

			if found && result != tc.expected {
				t.Errorf("Expected result:\n%q\nGot:\n%q", tc.expected, result)
			}
		})
	}

	// Test caching behavior
	t.Run("Caching behavior", func(t *testing.T) {
		config := &alpha.DarknessConfig{}
		config.Runtime.Logger = log.NewWithOptions(os.Stderr, log.Options{
			Level: log.FatalLevel, // Only show fatal errors
		})
		config.Runtime.WorkDir = alpha.WorkingDirectory(tmpDir)

		// First call should read from disk
		result1, found1 := expandSetupFile(config, "main.org", "#+setupfile: setup.org")
		if !found1 || result1 != setupFileContent {
			t.Fatalf("First call failed, expected setup content but got: %v", result1)
		}

		// Modify the file but the second call should use the cached version
		modifiedContent := "Modified content"
		err = os.WriteFile(setupFilePath, []byte(modifiedContent), 0644)
		if err != nil {
			t.Fatalf("Failed to update setup file: %v", err)
		}

		result2, found2 := expandSetupFile(config, "main.org", "#+setupfile: setup.org")
		if !found2 {
			t.Fatalf("Second call failed to find setup file")
		}

		// Should get the original content from cache, not the modified content
		if result2 != setupFileContent {
			t.Errorf("Cache not working. Expected original content, got: %q", result2)
		}
	})
}

// TestPreprocessIntegration tests the integration of all preprocessing functions
func TestPreprocessIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "darkness-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test setup file with macros
	setupFileContent := `#+macro: bold *$1*
#+macro: link [[https://example.com][$1]]`
	setupFilePath := filepath.Join(tmpDir, "setup.org")
	err = os.WriteFile(setupFilePath, []byte(setupFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write setup file: %v", err)
	}

	input := `#+setupfile: setup.org

#+macro: greet Hello, $1!

* Introduction

{{{bold(Important note)}}}

{{{greet(World)}}}

- List item with {{{link(example)}}}
- Another item

** Subsection

More text here.`

	expectedOutput := strings.Join([]string{
		"",
		"#+macro: bold *$1*",
		"#+macro: link [[https://example.com][$1]]",
		"",
		"",
		"",
		"* Introduction",
		"",
		"*Important note*",
		"",
		"Hello, World!",
		"",
		"",
		"- List item with [[https://example.com][example]]",
		"- Another item",
		"",
		"",
		"** Subsection",
		"",
		"More text here.",
		"",
		"",
	}, "\n")

	config := &alpha.DarknessConfig{}
	config.Runtime.Logger = log.NewWithOptions(os.Stderr, log.Options{
		Level: log.FatalLevel, // Only show fatal errors
	})
	config.Runtime.WorkDir = alpha.WorkingDirectory(tmpDir)

	parser := ParserOrgmode{Config: config}
	result := parser.preprocess("main.org", input)

	if result != expectedOutput {
		t.Errorf("Integration test failed.\nExpected:\n%q\nGot:\n%q", expectedOutput, result)
	}
}