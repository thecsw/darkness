package ichika

// GetDarknessFunc returns Darkness function to run by the
// command supplied, `nil` otherwise.
func GetDarknessFunc(command string) func() {
	if commandFunc, ok := CommandFuncs[DarknessCommand(command)]; ok {
		return commandFunc
	}
	return nil
}
