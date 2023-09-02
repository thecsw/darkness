package alpha

import (
	"github.com/thecsw/darkness/emilia/puck"
)

// setupProjectExtensions sets up the input/output extensions for the project.
func (conf *DarknessConfig) setupProjectExtensions(options Options) {
	// If input/output formats are empty, default to .org/.html respectively.
	if isUnset(conf.Project.Input) {
		conf.Runtime.Logger.Warn("Input format not found, using a default", "ext", puck.ExtensionOrgmode)
		conf.Project.Input = puck.ExtensionOrgmode
	}

	// Output section.
	if isUnset(conf.Project.Output) {
		conf.Runtime.Logger.Warn("Input format not found, using a default", "ext", puck.ExtensionHtml)
		conf.Project.Output = puck.ExtensionHtml
	}

	// If the output extension is not the same as the one in the config,
	// then overwrite the config.
	if !isUnset(options.OutputExtension) && conf.Project.Output != options.OutputExtension {
		conf.Runtime.Logger.Warn("Output extension was overwritten", "ext", options.OutputExtension)
		conf.Project.Output = options.OutputExtension
	}

}
