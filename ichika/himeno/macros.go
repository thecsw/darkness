package himeno

import (
	"os"
	"path/filepath"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/parse/orgmode"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/rei"
)

const (
	globalMacrosFileBasename = `_macros`
)

// RegisterGlobalMacros will read the global macros and inject them into pages.
func RegisterGlobalMacros(conf *alpha.DarknessConfig) {
	// Recall that this is primarily an orgmode feature, so we will lock it to that input ext.
	if conf.Project.Input != puck.ExtensionOrgmode {
		return
	}
	globalMacrosFile := yunyun.RelativePathFile(globalMacrosFileBasename + conf.Project.Input)
	globalMacrosFileFull := string(conf.Runtime.WorkDir.Join(globalMacrosFile))
	if exists, err := rei.FileExists(globalMacrosFileFull); exists {
		file, err := os.ReadFile(filepath.Clean(globalMacrosFileFull))
		if err != nil {
			conf.Runtime.Logger.Warn("Failed reading global macros file", "file", globalMacrosFileFull, "err", err)
		}
		if orgmode.CollectGlobalMacros(conf, globalMacrosFile, string(file)) {
			conf.Runtime.Logger.Info("Loaded global macros")
		}

	} else if err != nil {
		conf.Runtime.Logger.Error("Failed to see if the global macros file even exists", "err", err)
	}
}
