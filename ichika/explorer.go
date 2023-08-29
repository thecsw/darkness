package ichika

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/ichika/makima"
	"github.com/thecsw/darkness/yunyun"
	g "github.com/thecsw/gana"
	"github.com/thecsw/komi"
)

// FindFilesByExt finds all files with a given extension.
func FindFilesByExt(
	conf *alpha.DarknessConfig,
	pool komi.PoolConnector[*makima.Control],
	freshContext2 makima.Control,
) <-chan struct{} {
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

				freshContext := &makima.Control{
					Conf:     freshContext2.Conf,
					Parser:   freshContext2.Parser,
					Exporter: freshContext2.Exporter,
				}
				freshContext.InputFilename = conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(relPath))

				if err := pool.Submit(freshContext); err != nil {
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
func FindFilesByExtSimple(conf *alpha.DarknessConfig) []*makima.Control {
	toReturn := make([]*makima.Control, 0, 64)
	addFile := func(c *makima.Control) { toReturn = append(toReturn, c) }
	filesPool := komi.New(komi.WorkSimple(addFile))
	freshConfig := makima.Control{Conf: conf}
	<-FindFilesByExt(conf, filesPool, freshConfig)
	filesPool.Close()
	return toReturn
}

// FindFilesByExtSimpleDirs is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results and returns only the results
// which are children of the passed dirs.
func FindFilesByExtSimpleDirs(conf *alpha.DarknessConfig, dirs []string) []*makima.Control {
	files := FindFilesByExtSimple(conf)
	// If no dirs passed, run no filtering.
	if len(dirs) < 1 {
		return files
	}
	// Only return files that have passed dirs as parents.
	return g.Filter(func(c *makima.Control) bool {
		return g.Anyf(func(v string) bool {
			return strings.HasPrefix(string(c.InputFilename),
				string(conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(v))))
		}, dirs)
	}, files)
}
