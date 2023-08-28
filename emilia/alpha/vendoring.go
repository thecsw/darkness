package alpha

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/thecsw/rei"
)

func (conf *DarknessConfig) SetupGalleryVendoring(options Options) {
	// Work through the vendored galleries.
	conf.Runtime.VendorGalleries = options.VendorGalleries

	// Show a warning that we're going to vendor the galleries and that
	// the user should add the vendor directory to their .gitignore.
	if conf.Runtime.VendorGalleries {
		cmdColor := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff50a2"))
		yellow := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffff00"))
		fmt.Println("I'm going to vendor all gallery paths!")
		fmt.Println("If this is the first time, it will take a while... otherwise,",
			yellow.Render("an instant"))
		fmt.Printf("Please add %s to your .gitignore, so you don't pollute your git objects.\n",
			cmdColor.Render(string(conf.Project.DarknessVendorDirectory)))
		fmt.Println()
		if err := rei.Mkdir(filepath.Join(string(conf.Runtime.WorkDir),
			string(conf.Project.DarknessVendorDirectory))); err != nil {
			conf.Runtime.Logger.Warnf("creating vendor directory %s: %v", conf.Project.DarknessVendorDirectory, err)
			conf.Runtime.Logger.Warn("disabling vendoring by force")
			conf.Runtime.VendorGalleries = false
		}
	}

}
