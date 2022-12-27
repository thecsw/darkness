package emilia

import "os"

// FileExists returns true if file exists, false otherwise (in error too).
func FileExists(path string) bool {
	info, err := os.Stat(string(path))
	return info != nil && !os.IsNotExist(err)
}

// Mkdir creates a directory and reports fatal errors.
func Mkdir(path string) error {
	// Make sure that the vendor directory exists.
	err := os.Mkdir(string(path), 0755)
	// If we couldn't create the vendor directory and it doesn't
	// exist, then turn off the vendor option.
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
