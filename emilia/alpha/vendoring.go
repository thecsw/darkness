package alpha

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/thecsw/rei"
)

// setupGalleryVendoring sets up the vendoring of galleries.
func (conf *DarknessConfig) setupGalleryVendoring(options Options) {
	// Work through the vendored galleries.
	conf.Runtime.VendorGalleries = options.VendorGalleries

	// If we're not vendoring, then we're done.
	if !conf.Runtime.VendorGalleries {
		return
	}

	// Show a warning that we're going to vendor the galleries and that
	// the user should add the vendor directory to their .gitignore.
	cmdColor := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff50a2"))
	yellow := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffff00"))
	fmt.Println("I'm going to vendor all gallery paths!")
	fmt.Println("If this is the first time, it will take a while... otherwise,",
		yellow.Render("an instant"))
	fmt.Printf("Please add %s to your .gitignore, so you don't pollute your git objects.\n",
		cmdColor.Render(string(conf.Project.DarknessVendorDirectory)))
	fmt.Println()

	// Create the vendor directory.
	err := rei.Mkdir(filepath.Join(string(conf.Runtime.WorkDir), string(conf.Project.DarknessVendorDirectory)))

	// If we can create the vendor directory, then we're done.
	if err == nil {
		return
	}

	// Otherwise, we're going to disable vendoring.
	conf.Runtime.Logger.Warnf("creating vendor directory %s: %v", conf.Project.DarknessVendorDirectory, err)
	conf.Runtime.Logger.Warn("disabling vendoring by force")
	conf.Runtime.VendorGalleries = false
}
