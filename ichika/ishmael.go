package ichika

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thecsw/darkness/emilia"
)

// NewDarknessCommandFunc creates a default darkness config in the current directory
// if one already exists, nothing will happen, except a notice of that
func NewDarknessCommandFunc() {
	if len(os.Args) != 3 {
		fmt.Println("you forgot to add a directory name after new")
		return
	}
	dirName := os.Args[2]
	f, err := os.Open(filepath.Clean(dirName))
	if err == nil {
		fmt.Println("this directory already exists, bailing")
		if err := f.Close(); err != nil {
			fmt.Printf("failed to close directory: %s\n", err.Error())
		}
		return
	}
	if err := os.Mkdir(dirName, os.FileMode(0777)); err != nil {
		fmt.Println("couldn't create a directory for you:", err.Error())
		return
	}
	// Create a buffer for the tar file so we can start untarring it
	tarBuf := bytes.NewReader(defaultDarknessTemplate)
	if err := emilia.Untar(tarBuf, dirName); err != nil {
		fmt.Println("failed at flushing the template files:", err.Error())
		return
	}
	// Done
	fmt.Printf("Done! Go to %s and start creating!\n(run darkness build in there)\n", dirName)
}
