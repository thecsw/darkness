package emilia

import (
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

var (
	NumFoundFiles int32 = 0
)

// FindFilesByExt finds all files with a given extension.
func FindFilesByExt(inputFilenames chan<- yunyun.FullPathFile, ext string, wg *sync.WaitGroup) {
	NumFoundFiles = 0
	if err := godirwalk.Walk(Config.WorkDir, &godirwalk.Options{
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			fmt.Printf("Encountered an error while traversing %s: %s\n", osPathname, err.Error())
			return godirwalk.SkipNode
		},
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if filepath.Ext(osPathname) != ext {
				return nil
			}
			if (Config.Project.ExcludeEnabled && Config.Project.ExcludeRegex.MatchString(osPathname)) ||
				gana.First([]rune(de.Name())) == rune('.') {
				return filepath.SkipDir
			}
			wg.Add(1)
			relPath, err := filepath.Rel(Config.WorkDir, osPathname)
			inputFilenames <- JoinWorkdir(yunyun.RelativePathFile(relPath))
			atomic.AddInt32(&NumFoundFiles, 1)
			return err
		},
	}); err != nil {
		fmt.Printf("File traversal returned an error: %s\n", err.Error())
	}
	close(inputFilenames)
}

// FindFilesByExitSimple is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results.
func FindFilesByExtSimple(ext string) []yunyun.FullPathFile {
	inputFilenames := make(chan yunyun.FullPathFile, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go FindFilesByExt(inputFilenames, ext, wg)
	toReturn := make([]yunyun.FullPathFile, 0, 64)
	for inputFilename := range inputFilenames {
		toReturn = append(toReturn, inputFilename)
		wg.Done()
	}
	wg.Done()
	return toReturn
}
