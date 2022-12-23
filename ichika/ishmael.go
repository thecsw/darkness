package ichika

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/thecsw/darkness/emilia"
)

var (
	//go:embed ishmael/ishmael.tar.gz
	defaultDarknessTemplate []byte
)

// NewDarknessCommandFunc creates a default darkness config in the current directory
// if one already exists, nothing will happen, except a notice of that
func NewDarknessCommandFunc() {
	if len(os.Args) != 3 {
		fmt.Println("you forgot to add a directory name after new")
		os.Exit(1)
	}
	dirName := os.Args[2]
	f, err := os.Open(filepath.Clean(dirName))
	if err == nil {
		fmt.Println("this directory already exists, bailing")
		if err := f.Close(); err != nil {
			fmt.Printf("failed to close directory: %s\n", err.Error())
		}
		os.Exit(1)
	}
	if err := os.Mkdir(dirName, os.FileMode(0777)); err != nil {
		fmt.Println("couldn't create a directory for you:", err.Error())
		os.Exit(1)
	}

	// Create the darkness template reader.
	defaultDarknessTemplateReader := bytes.NewReader(defaultDarknessTemplate)

	// Uncompress the gzip file.
	gzipBuf, err := gzip.NewReader(defaultDarknessTemplateReader)
	if err != nil {
		fmt.Println("couldn't read the default template, fatal: " + err.Error())
		os.Exit(1)
	}

	// Create a buffer for the tar file so we can start untarring it.
	if err := emilia.Untar(gzipBuf, dirName); err != nil {
		fmt.Println("failed at flushing the template files:", err.Error())
		return
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
