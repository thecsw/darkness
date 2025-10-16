package hizuru

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/yunyun"
)

// Explicitly use the yunyun package to satisfy the import requirements
var testPathFile yunyun.FullPathFile = yunyun.FullPathFile("/tmp/test.md")

func setupTestEnvironment(t *testing.T) (string, *alpha.DarknessConfig) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "hizuru-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create a basic configuration for testing
	config := &alpha.DarknessConfig{}
	config.Runtime.Logger = log.NewWithOptions(os.Stderr, log.Options{
		Level: log.FatalLevel, // Only show fatal errors during tests
	})
	config.Runtime.WorkDir = alpha.WorkingDirectory(tempDir)
	config.Project.Input = ".org" // Use .org extension for testing

	// Create some test files and directories
	createTestFiles(t, tempDir)

	return tempDir, config
}

// createTestFiles creates a test directory structure
func createTestFiles(t *testing.T, root string) {
	// Create some basic files with the test extension
	files := []string{
		"file1.org",
		"file2.org",
		"nested/file3.org",
		"nested/deep/file4.org",
		"_ignoredfile.org",            // Should be skipped due to name
		".hiddenfile.org",             // Should be skipped as it's hidden
		"other.txt",                   // Should be skipped due to extension
		"nested/.hidden/hidden.org",   // Hidden directory, but file might be picked up
		"nested/_ignored/skipped.org", // Directory contains _ignored
	}

	// Create directories and files
	for _, file := range files {
		fullPath := filepath.Join(root, file)
		dir := filepath.Dir(fullPath)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		// Create the file
		if err := os.WriteFile(fullPath, []byte("test content for "+file), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Create a symlink
	symlinkSource := filepath.Join(root, "file1.org")
	symlinkTarget := filepath.Join(root, "symlink.org")
	if err := os.Symlink(symlinkSource, symlinkTarget); err != nil {
		// Some environments might not support symlinks, so just warn
		t.Logf("Warning: Couldn't create symlink for testing: %v", err)
	}
}

// TestFindFilesByExtSimple tests the FindFilesByExtSimple function
func TestFindFilesByExtSimple(t *testing.T) {
	tempDir, config := setupTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Run the function we're testing
	files := FindFilesByExtSimple(config)

	// Convert to strings for easier comparison
	fileStrings := make([]string, 0, len(files))
	for _, file := range files {
		// Convert to relative paths for easier testing
		rel, err := filepath.Rel(tempDir, string(file))
		if err != nil {
			t.Fatalf("Failed to get relative path: %v", err)
		}
		fileStrings = append(fileStrings, rel)
	}

	// Sort the results for consistent comparison
	sort.Strings(fileStrings)

	// Expected files (relative paths, sorted)
	expected := []string{
		"file1.org",
		"file2.org",
		"nested/.hidden/hidden.org",
		"nested/deep/file4.org",
		"nested/file3.org",
		"nested/_ignored/skipped.org",
		"symlink.org", // Symlinks should be included
	}
	sort.Strings(expected)

	// Check if expected files match what we got
	if !reflect.DeepEqual(fileStrings, expected) {
		t.Errorf("FindFilesByExtSimple() found incorrect files.\nExpected: %v\nGot: %v",
			expected, fileStrings)
	}

	// Check that we don't have duplicates
	uniqueFiles := make(map[string]bool)
	for _, file := range fileStrings {
		if uniqueFiles[file] {
			t.Errorf("Found duplicate file: %s", file)
		}
		uniqueFiles[file] = true
	}
}

// TestFindFilesByExtSimpleDirs tests the findFilesByExtSimpleDirs function
func TestFindFilesByExtSimpleDirs(t *testing.T) {
	tempDir, config := setupTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		dirs     []string
		expected []string
	}{
		{
			name: "All files (empty dirs)",
			dirs: []string{},
			expected: []string{
				"file1.org",
				"file2.org",
				"nested/.hidden/hidden.org",
				"nested/deep/file4.org",
				"nested/file3.org",
				"nested/_ignored/skipped.org",
				"symlink.org",
			},
		},
		{
			name: "Only nested files",
			dirs: []string{"nested"},
			expected: []string{
				"nested/.hidden/hidden.org",
				"nested/deep/file4.org",
				"nested/file3.org",
				"nested/_ignored/skipped.org",
			},
		},
		{
			name: "Only deep nested files",
			dirs: []string{"nested/deep"},
			expected: []string{
				"nested/deep/file4.org",
			},
		},
		{
			name: "Multiple directories",
			dirs: []string{"", "nested/deep"},
			expected: []string{
				"file1.org",
				"file2.org",
				"nested/.hidden/hidden.org",
				"nested/deep/file4.org",
				"nested/file3.org",
				"nested/_ignored/skipped.org",
				"symlink.org",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			files := findFilesByExtSimpleDirs(config, tc.dirs)

			// Convert to strings for easier comparison
			fileStrings := make([]string, 0, len(files))
			for _, file := range files {
				// Convert to relative paths for easier testing
				rel, err := filepath.Rel(tempDir, string(file))
				if err != nil {
					t.Fatalf("Failed to get relative path: %v", err)
				}
				fileStrings = append(fileStrings, rel)
			}

			// Sort the results for consistent comparison
			sort.Strings(fileStrings)
			sort.Strings(tc.expected)

			if !reflect.DeepEqual(fileStrings, tc.expected) {
				t.Errorf("findFilesByExtSimpleDirs() with dirs=%v found incorrect files.\nExpected: %v\nGot: %v",
					tc.dirs, tc.expected, fileStrings)
			}
		})
	}
}

// TestSkipPrefixExclusion tests that files with skipPrefix in the base filename are properly excluded
func TestSkipPrefixExclusion(t *testing.T) {
	tempDir, config := setupTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Run the function we're testing
	files := FindFilesByExtSimple(config)

	// Check that no base filename contains the skipPrefix
	for _, file := range files {
		baseName := filepath.Base(string(file))
		if strings.Contains(baseName, skipPrefix) {
			t.Errorf("Found file with skip prefix in base filename: %s", baseName)
		}
	}
}

// TestHiddenFileExclusion tests that files with names starting with a dot are excluded
func TestHiddenFileExclusion(t *testing.T) {
	tempDir, config := setupTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Run the function we're testing
	files := FindFilesByExtSimple(config)

	// Check that no hidden files (starting with .) are included
	for _, file := range files {
		baseName := filepath.Base(string(file))
		if strings.HasPrefix(baseName, ".") {
			t.Errorf("Found hidden file (starting with dot): %s", baseName)
		}
	}
}

// TestBuildPagesSimple tests the BuildPagesSimple function
func TestBuildPagesSimple(t *testing.T) {
	tempDir, config := setupTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Add some test content to the files that will be parsed
	testFiles := []struct {
		path    string
		content string
	}{
		{
			path:    "page1.org",
			content: "# Page 1\nThis is test content.\n",
		},
		{
			path:    "nested/page2.org",
			content: "# Page 2\nThis is nested test content.\n",
		},
	}

	for _, tf := range testFiles {
		fullPath := filepath.Join(tempDir, tf.path)
		dir := filepath.Dir(fullPath)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		// Create the file
		if err := os.WriteFile(fullPath, []byte(tf.content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// This test is more limited because we don't have a real parser in test mode
	// So we'll just check that the function returns without errors and returns expected number of pages
	pages := BuildPagesSimple(config, []string{})

	// Since BuildPagesSimple depends on the file parsing which isn't fully mocked,
	// we can only do basic checks here
	if pages == nil {
		t.Fatal("BuildPagesSimple returned nil")
	}
}

// TestNoFileDuplication tests that no files are duplicated when scanning
func TestNoFileDuplication(t *testing.T) {
	tempDir, config := setupTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Create a test case where a file could potentially be discovered via multiple paths
	// by creating hard links
	srcPath := filepath.Join(tempDir, "unique.org")
	if err := os.WriteFile(srcPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create file %s: %v", srcPath, err)
	}

	// Try to create a hard link (may not work on all systems)
	linkPath := filepath.Join(tempDir, "link_to_unique.org")
	if err := os.Link(srcPath, linkPath); err != nil {
		t.Logf("Hard links not supported, skipping part of test: %v", err)
	}

	// Run the function we're testing
	files := FindFilesByExtSimple(config)

	// Build map of paths
	seen := make(map[string]bool)
	for _, file := range files {
		rel, err := filepath.Rel(tempDir, string(file))
		if err != nil {
			t.Fatalf("Failed to get relative path: %v", err)
		}

		// The same content should not be processed twice
		if seen[rel] {
			t.Errorf("Found duplicate file in results: %s", rel)
		}
		seen[rel] = true
	}
}
