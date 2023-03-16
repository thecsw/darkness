package emilia

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/yunyun"
	g "github.com/thecsw/gana"
)

// FindFilesByExt finds all files with a given extension.
func FindFilesByExt(pool g.PoolConnector[yunyun.FullPathFile], ext string) <-chan struct{} {
	done := make(chan struct{}, 1)
	go func(done chan<- struct{}) {
		if err := godirwalk.Walk(Config.WorkDir, &godirwalk.Options{
			ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
				fmt.Printf("Encountered an error while traversing %s: %s\n", osPathname, err.Error())
				return godirwalk.SkipNode
			},
			Unsorted: true,
			Callback: func(osPathname string, de *godirwalk.Dirent) error {
				if pool.IsClosed() {
					return nil
				}
				if filepath.Ext(osPathname) != ext || strings.HasPrefix(filepath.Base(osPathname), ".") {
					return nil
				}
				if (Config.Project.ExcludeEnabled && Config.Project.ExcludeRegex.MatchString(osPathname)) ||
					(g.First([]rune(de.Name())) == rune('.') && de.IsDir()) {
					return filepath.SkipDir
				}
				relPath, err := filepath.Rel(Config.WorkDir, osPathname)
				pool.Submit(JoinWorkdir(yunyun.RelativePathFile(relPath)))
				return err
			},
		}); err != nil {
			fmt.Printf("File traversal returned an error: %s\n", err.Error())
		}
		done <- struct{}{}
	}(done)
	return done
}

// FindFilesByExitSimple is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results.
func FindFilesByExtSimple(ext string) []yunyun.FullPathFile {
	toReturn := make([]yunyun.FullPathFile, 0, 64)
	addFile := func(v yunyun.FullPathFile) bool { toReturn = append(toReturn, v); return true }
	filesPool := g.NewPool(addFile, 1, 0)
	filesPool.DisableOutput()
	<-FindFilesByExt(filesPool, ext)
	filesPool.Wait()
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
