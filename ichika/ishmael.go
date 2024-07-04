package ichika

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/rei"
)

//go:embed ishmael.tar.gz
var defaultDarknessTemplate []byte

// newDarknessCommandFunc creates a default darkness config in the current directory
// if one already exists, nothing will happen, except a notice of that
func newDarknessCommandFunc() {
	if len(os.Args) != 3 {
		fmt.Println("you forgot to add a directory name after new")
		os.Exit(1)
	}
	dirName := os.Args[2]
	f, err := os.Open(filepath.Clean(dirName))
	if err == nil {
		fmt.Println("this directory already exists, bailing")
		if err := f.Close(); err != nil {
			puck.Logger.Errorf("closing directory %s: %v", dirName, err)
		}
		os.Exit(1)
	}
	if err := os.MkdirAll(dirName, os.FileMode(0o777)); err != nil {
		puck.Logger.Fatalf("creating your directory %s: %v", dirName, err)
	}

	// Create the darkness template reader.
	defaultDarknessTemplateReader := bytes.NewReader(defaultDarknessTemplate)

	// Uncompress the gzip file.
	gzipBuf, err := gzip.NewReader(defaultDarknessTemplateReader)
	if err != nil {
		puck.Logger.Fatalf("reading default template %s: %v", defaultDarknessTemplate, err)
	}

	// Create a buffer for the tar file so we can start untarring it.
	if err := rei.Untar(gzipBuf, dirName); err != nil {
		puck.Logger.Fatalf("flushing tarred template files: %v", err)
	}

	cmdColor := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff50a2"))
	// Done
	fmt.Printf(`Done!
Go to %s and start creating!

You can %s to build your new shiny
website (with %s to use local paths).

Or even better, try running %s and
darkness will serve your files on a local server!
(it will auto build files during editing)
`, dirName,
		cmdColor.Render("darkness build"),
		cmdColor.Render("-dev"),
		cmdColor.Render("darkness serve"),
	)
}
