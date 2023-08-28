package emilia

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/emilia/alpha"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/yunyun"
	g "github.com/thecsw/gana"
	"github.com/thecsw/komi"
)

// FindFilesByExt finds all files with a given extension.
func FindFilesByExt(conf alpha.DarknessConfig, pool komi.PoolConnector[yunyun.FullPathFile]) <-chan struct{} {
	done := make(chan struct{})
	go func(done chan<- struct{}) {
		if err := godirwalk.Walk(string(conf.Runtime.WorkDir), &godirwalk.Options{
			ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
				conf.Runtime.Logger.Errorf("traversing %s: %v", osPathname, err)
				return godirwalk.SkipNode
			},
			Unsorted: true,
			Callback: func(osPathname string, de *godirwalk.Dirent) error {
				if pool.IsClosed() {
					return nil
				}
				if filepath.Ext(osPathname) != conf.Project.Input || strings.HasPrefix(filepath.Base(osPathname), ".") {
					return nil
				}
				if (conf.Project.ExcludeEnabled && conf.Project.ExcludeRegex.MatchString(osPathname)) ||
					(g.First([]rune(de.Name())) == '.' && de.IsDir()) {
					return filepath.SkipDir
				}
				relPath, err := filepath.Rel(string(conf.Runtime.WorkDir), osPathname)
				if err != nil {
					return fmt.Errorf("finding relative path of %s to %s: %v", osPathname, conf.Runtime.WorkDir, err)
				}
				if err := pool.Submit(conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(relPath))); err != nil {
					conf.Runtime.Logger.Errorf("couldn't submit %s: %v", relPath, err)
				}
				return nil
			},
		}); err != nil {
			conf.Runtime.Logger.Errorf("root traversal: %v", err)
		}
		done <- struct{}{}
	}(done)
	return done
}

// FindFilesByExtSimple is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results.
func FindFilesByExtSimple(conf alpha.DarknessConfig) []yunyun.FullPathFile {
	toReturn := make([]yunyun.FullPathFile, 0, 64)
	addFile := func(v yunyun.FullPathFile) { toReturn = append(toReturn, v) }
	filesPool := komi.New(komi.WorkSimple(addFile))
	<-FindFilesByExt(conf, filesPool)
	filesPool.Close()
	return toReturn
}

// FindFilesByExtSimpleDirs is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results and returns only the results
// which are children of the passed dirs.
func FindFilesByExtSimpleDirs(conf alpha.DarknessConfig, dirs []string) []yunyun.FullPathFile {
	files := FindFilesByExtSimple(conf)
	// If no dirs passed, run no filtering.
	if len(dirs) < 1 {
		return files
	}
	// Only return files that have passed dirs as parents.
	return g.Filter(func(path yunyun.FullPathFile) bool {
		return g.Anyf(func(v string) bool {
			return strings.HasPrefix(string(path),
				string(conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(v))))
		}, dirs)
	}, files)
}
