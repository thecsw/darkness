package yunyun

import (
	"regexp"
	"testing"
)

// TestRegexInitialization ensures that BuildRegex properly initializes all regex patterns
func TestRegexInitialization(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	// Test that all regexes are properly initialized
	if BoldText == nil {
		t.Error("BoldText was not initialized")
	}
	if ItalicText == nil {
		t.Error("ItalicText was not initialized")
	}
	if BoldItalicText == nil {
		t.Error("BoldItalicText was not initialized")
	}
	if ItalicBoldText == nil {
		t.Error("ItalicBoldText was not initialized")
	}
	if VerbatimText == nil {
		t.Error("VerbatimText was not initialized")
	}
	if StrikethroughText == nil {
		t.Error("StrikethroughText was not initialized")
	}
	if UnderlineText == nil {
		t.Error("UnderlineText was not initialized")
	}
	if SuperscriptText == nil {
		t.Error("SuperscriptText was not initialized")
	}
	if SubscriptText == nil {
		t.Error("SubscriptText was not initialized")
	}
}

// TestBoldRegex tests the bold text regex pattern
func TestBoldRegex(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name     string
		input    string
		expected bool
		text     string
	}{
		{
			name:     "Simple bold text",
			input:    "This is *bold* text",
			expected: true,
			text:     "bold",
		},
		{
			name:     "Bold text at beginning of line",
			input:    "*Bold* text at beginning",
			expected: true,
			text:     "Bold",
		},
		{
			name:     "Bold text at end of line",
			input:    "Text at the end *bold*",
			expected: true,
			text:     "bold",
		},
		{
			name:     "No bold text",
			input:    "This is not bold text",
			expected: false,
			text:     "",
		},
		{
			name:     "Asterisk in the middle of a word",
			input:    "This is not*bold text",
			expected: false,
			text:     "",
		},
		{
			name:     "Bold text with punctuation",
			input:    "This is (*bold*).",
			expected: true,
			text:     "bold",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := BoldText.FindStringSubmatch(tc.input)
			matched := len(matches) > 0

			if matched != tc.expected {
				t.Errorf("Expected match: %v, got: %v for input: %s", tc.expected, matched, tc.input)
			}

			if matched {
				textIdx := BoldText.SubexpIndex("text")
				if textIdx == -1 {
					t.Fatal("No 'text' capturing group found in BoldText regex")
				}

				if matches[textIdx] != tc.text {
					t.Errorf("Expected text: %q, got: %q for input: %s", tc.text, matches[textIdx], tc.input)
				}
			}
		})
	}
}

// TestItalicRegex tests the italic text regex pattern
func TestItalicRegex(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name     string
		input    string
		expected bool
		text     string
	}{
		{
			name:     "Simple italic text",
			input:    "This is /italic/ text",
			expected: true,
			text:     "italic",
		},
		{
			name:     "Italic text at beginning of line",
			input:    "/Italic/ text at beginning",
			expected: true,
			text:     "Italic",
		},
		{
			name:     "Italic text at end of line",
			input:    "Text at the end /italic/",
			expected: true,
			text:     "italic",
		},
		{
			name:     "No italic text",
			input:    "This is not italic text",
			expected: false,
			text:     "",
		},
		{
			name:     "Slash in the middle of a word",
			input:    "This is not/italic text",
			expected: false,
			text:     "",
		},
		{
			name:     "Italic text with punctuation",
			input:    "This is (/italic/).",
			expected: true,
			text:     "italic",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := ItalicText.FindStringSubmatch(tc.input)
			matched := len(matches) > 0

			if matched != tc.expected {
				t.Errorf("Expected match: %v, got: %v for input: %s", tc.expected, matched, tc.input)
			}

			if matched {
				textIdx := ItalicText.SubexpIndex("text")
				if textIdx == -1 {
					t.Fatal("No 'text' capturing group found in ItalicText regex")
				}

				if matches[textIdx] != tc.text {
					t.Errorf("Expected text: %q, got: %q for input: %s", tc.text, matches[textIdx], tc.input)
				}
			}
		})
	}
}

// TestBoldItalicRegex tests the bold-italic text regex pattern (*/text/*)
func TestBoldItalicRegex(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name     string
		input    string
		expected bool
		text     string
	}{
		{
			name:     "Simple bold-italic text",
			input:    "This is */bold-italic/* text",
			expected: true,
			text:     "bold-italic",
		},
		{
			name:     "Bold-italic text at beginning of line",
			input:    "*/Bold-italic/* text at beginning",
			expected: true,
			text:     "Bold-italic",
		},
		{
			name:     "Bold-italic text at end of line",
			input:    "Text at the end */bold-italic/*",
			expected: true,
			text:     "bold-italic",
		},
		{
			name:     "No bold-italic text",
			input:    "This is *bold* and /italic/ but not combined",
			expected: false,
			text:     "",
		},
		{
			name:     "Bold-italic text with punctuation",
			input:    "This is (*/bold-italic/*).",
			expected: true,
			text:     "bold-italic",
		},
		{
			// This test case is specifically testing the issue mentioned in the problem statement
			name:     "Bold-italic with descriptive content",
			input:    "Look at this: */It is violently brilliant/* text",
			expected: true,
			text:     "It is violently brilliant",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Run the test multiple times to check for deterministic behavior
			for i := range 10 {
				matches := BoldItalicText.FindStringSubmatch(tc.input)
				matched := len(matches) > 0

				if matched != tc.expected {
					t.Errorf("Iteration %d: Expected match: %v, got: %v for input: %s", i, tc.expected, matched, tc.input)
				}

				if matched {
					textIdx := BoldItalicText.SubexpIndex("text")
					if textIdx == -1 {
						t.Fatal("No 'text' capturing group found in BoldItalicText regex")
					}

					if matches[textIdx] != tc.text {
						t.Errorf("Iteration %d: Expected text: %q, got: %q for input: %s", i, tc.text, matches[textIdx], tc.input)
					}
				}
			}
		})
	}
}

// TestItalicBoldRegex tests the italic-bold text regex pattern (/*text*/)
func TestItalicBoldRegex(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name     string
		input    string
		expected bool
		text     string
	}{
		{
			name:     "Simple italic-bold text",
			input:    "This is /*italic-bold*/ text",
			expected: true,
			text:     "italic-bold",
		},
		{
			name:     "Italic-bold text at beginning of line",
			input:    "/*Italic-bold*/ text at beginning",
			expected: true,
			text:     "Italic-bold",
		},
		{
			name:     "Italic-bold text at end of line",
			input:    "Text at the end /*italic-bold*/",
			expected: true,
			text:     "italic-bold",
		},
		{
			name:     "No italic-bold text",
			input:    "This is *bold* and /italic/ but not combined",
			expected: false,
			text:     "",
		},
		{
			name:     "Italic-bold text with punctuation",
			input:    "This is (/*italic-bold*/).",
			expected: true,
			text:     "italic-bold",
		},
		{
			// This test case is specifically testing the issue mentioned in the problem statement
			name:     "Italic-bold with descriptive content",
			input:    "Look at this: /*It is violently brilliant*/ text",
			expected: true,
			text:     "It is violently brilliant",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Run the test multiple times to check for deterministic behavior
			for i := range 10 {
				matches := ItalicBoldText.FindStringSubmatch(tc.input)
				matched := len(matches) > 0

				if matched != tc.expected {
					t.Errorf("Iteration %d: Expected match: %v, got: %v for input: %s", i, tc.expected, matched, tc.input)
				}

				if matched {
					textIdx := ItalicBoldText.SubexpIndex("text")
					if textIdx == -1 {
						t.Fatal("No 'text' capturing group found in ItalicBoldText regex")
					}

					if matches[textIdx] != tc.text {
						t.Errorf("Iteration %d: Expected text: %q, got: %q for input: %s", i, tc.text, matches[textIdx], tc.input)
					}
				}
			}
		})
	}
}

// TestAllFormatTogether tests all formatting options together to check for interference
func TestAllFormatTogether(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	// Test each format individually with distinct inputs to avoid overlaps
	inputs := map[*regexp.Regexp]struct {
		text  string
		match string
	}{
		BoldText:          {text: "This is *bold* text", match: "bold"},
		ItalicText:        {text: "This is /italic/ text", match: "italic"},
		BoldItalicText:    {text: "This is */bold-italic/* text", match: "bold-italic"},
		ItalicBoldText:    {text: "This is /*italic-bold*/ text", match: "italic-bold"},
		VerbatimText:      {text: "This is ~verbatim~ text", match: "verbatim"},
		StrikethroughText: {text: "This is +strikethrough+ text", match: "strikethrough"},
		UnderlineText:     {text: "This is _underline_ text", match: "underline"},
	}

	for regex, data := range inputs {
		matches := regex.FindStringSubmatch(data.text)
		if len(matches) == 0 {
			t.Errorf("Format not matched: %v", data.text)
			continue
		}

		textIdx := regex.SubexpIndex("text")
		if matches[textIdx] != data.match {
			t.Errorf("Expected match: %q, got: %q for text: %s", data.match, matches[textIdx], data.text)
		}
	}

	// Test the specific problematic case from the issue
	problematicInput := "/*It is violently brilliant*/"
	matches := ItalicBoldText.FindStringSubmatch(problematicInput)

	if len(matches) == 0 {
		t.Error("Failed to match the problematic case")
	} else {
		textIdx := ItalicBoldText.SubexpIndex("text")
		expectedText := "It is violently brilliant"
		if matches[textIdx] != expectedText {
			t.Errorf("Expected text: %q, got: %q", expectedText, matches[textIdx])
		}
	}
}

// TestVerbatimText tests verbatim text formatting
func TestVerbatimText(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name     string
		input    string
		expected bool
		text     string
	}{
		{
			name:     "Verbatim with tilde",
			input:    "This is ~verbatim~ text",
			expected: true,
			text:     "verbatim",
		},
		{
			name:     "Verbatim with equals",
			input:    "This is =verbatim= text",
			expected: true,
			text:     "verbatim",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := VerbatimText.FindStringSubmatch(tc.input)
			matched := len(matches) > 0

			if matched != tc.expected {
				t.Errorf("Expected match: %v, got: %v for input: %s", tc.expected, matched, tc.input)
			}

			if matched {
				textIdx := VerbatimText.SubexpIndex("text")
				if textIdx == -1 {
					t.Fatal("No 'text' capturing group found in VerbatimText regex")
				}

				if matches[textIdx] != tc.text {
					t.Errorf("Expected text: %q, got: %q for input: %s", tc.text, matches[textIdx], tc.input)
				}
			}
		})
	}
}

// TestStrikethroughText tests strikethrough text formatting
func TestStrikethroughText(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name     string
		input    string
		expected bool
		text     string
	}{
		{
			name:     "Simple strikethrough",
			input:    "This is +strikethrough+ text",
			expected: true,
			text:     "strikethrough",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := StrikethroughText.FindStringSubmatch(tc.input)
			matched := len(matches) > 0

			if matched != tc.expected {
				t.Errorf("Expected match: %v, got: %v for input: %s", tc.expected, matched, tc.input)
			}

			if matched {
				textIdx := StrikethroughText.SubexpIndex("text")
				if textIdx == -1 {
					t.Fatal("No 'text' capturing group found in StrikethroughText regex")
				}

				if matches[textIdx] != tc.text {
					t.Errorf("Expected text: %q, got: %q for input: %s", tc.text, matches[textIdx], tc.input)
				}
			}
		})
	}
}

// TestUnderlineText tests underline text formatting
func TestUnderlineText(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name     string
		input    string
		expected bool
		text     string
	}{
		{
			name:     "Simple underline",
			input:    "This is _underline_ text",
			expected: true,
			text:     "underline",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := UnderlineText.FindStringSubmatch(tc.input)
			matched := len(matches) > 0

			if matched != tc.expected {
				t.Errorf("Expected match: %v, got: %v for input: %s", tc.expected, matched, tc.input)
			}

			if matched {
				textIdx := UnderlineText.SubexpIndex("text")
				if textIdx == -1 {
					t.Fatal("No 'text' capturing group found in UnderlineText regex")
				}

				if matches[textIdx] != tc.text {
					t.Errorf("Expected text: %q, got: %q for input: %s", tc.text, matches[textIdx], tc.input)
				}
			}
		})
	}
}

// TestSuperscriptText tests superscript text formatting
func TestSuperscriptText(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name     string
		input    string
		expected bool
		text     string
	}{
		{
			name:     "Simple superscript",
			input:    "This is ^{{superscript}} text",
			expected: true,
			text:     "superscript",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := SuperscriptText.FindStringSubmatch(tc.input)
			matched := len(matches) > 0

			if matched != tc.expected {
				t.Errorf("Expected match: %v, got: %v for input: %s", tc.expected, matched, tc.input)
			}

			if matched {
				textIdx := SuperscriptText.SubexpIndex("text")
				if textIdx == -1 {
					t.Fatal("No 'text' capturing group found in SuperscriptText regex")
				}

				if matches[textIdx] != tc.text {
					t.Errorf("Expected text: %q, got: %q for input: %s", tc.text, matches[textIdx], tc.input)
				}
			}
		})
	}
}

// TestSubscriptText tests subscript text formatting
func TestSubscriptText(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name     string
		input    string
		expected bool
		text     string
	}{
		{
			name:     "Simple subscript",
			input:    "This is _{{subscript}} text",
			expected: true,
			text:     "subscript",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := SubscriptText.FindStringSubmatch(tc.input)
			matched := len(matches) > 0

			if matched != tc.expected {
				t.Errorf("Expected match: %v, got: %v for input: %s", tc.expected, matched, tc.input)
			}

			if matched {
				textIdx := SubscriptText.SubexpIndex("text")
				if textIdx == -1 {
					t.Fatal("No 'text' capturing group found in SubscriptText regex")
				}

				if matches[textIdx] != tc.text {
					t.Errorf("Expected text: %q, got: %q for input: %s", tc.text, matches[textIdx], tc.input)
				}
			}
		})
	}
}

// TestDeterministicMatching runs multiple iterations of matching to check for consistency
func TestDeterministicMatching(t *testing.T) {
	// Initialize regexes
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name  string
		input string
		regex *regexp.Regexp
		text  string
	}{
		{
			name:  "Bold-italic case from issue",
			input: "/*It is violently brilliant*/",
			regex: ItalicBoldText,
			text:  "It is violently brilliant",
		},
		{
			name:  "Bold-italic alternative syntax",
			input: "*/It is violently brilliant/*",
			regex: BoldItalicText,
			text:  "It is violently brilliant",
		},
		{
			name:  "Bold with special chars",
			input: "*Special-chars!*",
			regex: BoldText,
			text:  "Special-chars!",
		},
		{
			name:  "Italic with special chars",
			input: "/Special-chars!/",
			regex: ItalicText,
			text:  "Special-chars!",
		},
	}

	// Run each test pattern multiple times
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Run 100 iterations to ensure deterministic behavior
			for i := range 100 {
				matches := tc.regex.FindStringSubmatch(tc.input)
				if len(matches) == 0 {
					t.Fatalf("Iteration %d: Failed to match: %s", i, tc.input)
				}

				textIdx := tc.regex.SubexpIndex("text")
				if textIdx == -1 {
					t.Fatal("No 'text' capturing group found in regex")
				}

				matchedText := matches[textIdx]
				if matchedText != tc.text {
					t.Fatalf("Iteration %d: Expected text: %q, got: %q", i, tc.text, matchedText)
				}
			}
		})
	}
}
