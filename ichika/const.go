package ichika

// DarknessCommand is a type to enforce input typing.
type DarknessCommand string

const (
	newDarknessCommand DarknessCommand = `new`
	buildCommand       DarknessCommand = `build`
	serveCommand       DarknessCommand = `serve`
	cleanCommand       DarknessCommand = `clean`
	meguminCommand     DarknessCommand = `megumin`
	misaCommand        DarknessCommand = `misa`
	lalatinaCommand    DarknessCommand = `lalatina`
	aquaCommand        DarknessCommand = `aqua`
)

// CommandFuncs maps supplied darkness command to the function
// that needs to get executed.
var CommandFuncs = map[DarknessCommand]func(){
	newDarknessCommand: NewDarknessCommandFunc,
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
