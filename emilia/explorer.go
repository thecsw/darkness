package emilia

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/yunyun"
	g "github.com/thecsw/gana"
	"github.com/thecsw/komi"
)

// FindFilesByExt finds all files with a given extension.
func FindFilesByExt(pool komi.PoolConnector[yunyun.FullPathFile], ext string) <-chan struct{} {
	done := make(chan struct{})
	go func(done chan<- struct{}) {
		if err := godirwalk.Walk(Config.WorkDir, &godirwalk.Options{
			ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
				Logger.Errorf("traversing %s: %v", osPathname, err)
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
				if err != nil {
					return fmt.Errorf("finding relative path of %s to %s: %v", osPathname, Config.WorkDir, err)
				}
				pool.Submit(JoinWorkdir(yunyun.RelativePathFile(relPath)))
				return nil

			},
		}); err != nil {
			Logger.Errorf("root traversal: %v", err)
		}
		done <- struct{}{}
	}(done)
	return done
}

// FindFilesByExitSimple is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results.
func FindFilesByExtSimple(ext string) []yunyun.FullPathFile {
	toReturn := make([]yunyun.FullPathFile, 0, 64)
	addFile := func(v yunyun.FullPathFile) { toReturn = append(toReturn, v) }
	filesPool := komi.New(komi.WorkSimple(addFile))
	<-FindFilesByExt(filesPool, ext)
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
