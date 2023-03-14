package emilia

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
	g "github.com/thecsw/gana"
)

var (
	NumFoundFiles int32 = 0
)

// FindFilesByExt finds all files with a given extension.
func FindFilesByExt(pool gana.PoolConnector[yunyun.FullPathFile], ext string, wg *sync.WaitGroup) {
	NumFoundFiles = 0
	if err := godirwalk.Walk(Config.WorkDir, &godirwalk.Options{
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			fmt.Printf("Encountered an error while traversing %s: %s\n", osPathname, err.Error())
			return godirwalk.SkipNode
		},
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if filepath.Ext(osPathname) != ext || strings.HasPrefix(filepath.Base(osPathname), ".") {
				return nil
			}
			if (Config.Project.ExcludeEnabled && Config.Project.ExcludeRegex.MatchString(osPathname)) ||
				(g.First([]rune(de.Name())) == rune('.') && de.IsDir()) {
				return filepath.SkipDir
			}
			relPath, err := filepath.Rel(Config.WorkDir, osPathname)
			pool.Submit(JoinWorkdir(yunyun.RelativePathFile(relPath)))
			atomic.AddInt32(&NumFoundFiles, 1)
			return err
		},
	}); err != nil {
		fmt.Printf("File traversal returned an error: %s\n", err.Error())
	}
	wg.Done()
}

// FindFilesByExitSimple is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results.
func FindFilesByExtSimple(ext string) []yunyun.FullPathFile {
	toReturn := make([]yunyun.FullPathFile, 0, 64)
	wg := &sync.WaitGroup{}
	addFile := func(v yunyun.FullPathFile) bool { toReturn = append(toReturn, v); wg.Done(); return true }
	filesPool := gana.NewPool(addFile, 0, 0, "exp")
	go FindFilesByExt(filesPool, ext, wg)
	filesPool.Close()
	return toReturn
}

// FindFilesByExitSimple is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results and returns only the results
// which are children of the passed dirs.
func FindFilesByExtSimpleDirs(ext string, dirs []string) []yunyun.FullPathFile {
	files := FindFilesByExtSimple(ext)
	// If no dirs passed, run no filtering.
	if len(dirs) < 1 {
		return files
	}
	// Only return files that have passed dirs as parents.
	return g.Filter(func(path yunyun.FullPathFile) bool {
		return g.Anyf(func(v string) bool {
			return strings.HasPrefix(string(path),
				string(JoinWorkdir(yunyun.RelativePathFile(v))))
		}, dirs)
	}, files)
}
