package ichika

import (
	_ "embed"
)

var (
	//go:embed ishmael/ishmael.tar
	defaultDarknessTemplate []byte
)

// DarknessCommand is a type to enforce input typing.
type DarknessCommand string

const (
	newDarknessCommand DarknessCommand = `new`
	oneFileCommand     DarknessCommand = `file`
	buildCommand       DarknessCommand = `build`
	serveCommand       DarknessCommand = `serve`
	cleanCommand       DarknessCommand = `clean`
	meguminCommand     DarknessCommand = `megumin`
	misaCommand        DarknessCommand = `misa`
	lalatinaCommand    DarknessCommand = `lalatina`
	aquaCommand        DarknessCommand = `aqua`
)

var (
	// CommandFuncs maps supplied darkness command to the function
	// that needs to get executed.
	CommandFuncs = map[DarknessCommand]func(){
		newDarknessCommand: NewDarknessCommandFunc,
		oneFileCommand:     OneFileCommandFunc,
		buildCommand:       BuildCommandFunc,
		serveCommand:       ServeCommandFunc,
		cleanCommand:       CleanCommandFunc,
		meguminCommand:     MeguminCommandFunc,
		misaCommand:        MisaCommandFunc,
		lalatinaCommand:    LalatinaCommandFunc,
		aquaCommand:        AquaCommandFunc,

		// All the help commands
		`-h`:     HelpCommandFunc,
		`help`:   HelpCommandFunc,
		`-help`:  HelpCommandFunc,
		`--help`: HelpCommandFunc,
	}
)
