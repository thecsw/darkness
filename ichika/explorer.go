package ichika

import (
	"fmt"
	"github.com/thecsw/rei"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/yunyun"
	g "github.com/thecsw/gana"
)

// FindFilesByExt finds all files with a given extension.
func FindFilesByExt(conf *alpha.DarknessConfig, inputFiles chan<- yunyun.FullPathFile) {
	if err := godirwalk.Walk(string(conf.Runtime.WorkDir), &godirwalk.Options{
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			conf.Runtime.Logger.Errorf("traversing %s: %v", osPathname, err)
			return godirwalk.SkipNode
		},
		Unsorted: true,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
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
			inputFiles <- conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(relPath))
			return nil
		},
	}); err != nil {
		conf.Runtime.Logger.Errorf("root traversal: %v", err)
	}
	close(inputFiles)
}

// FindFilesByExtSimple is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results.
func FindFilesByExtSimple(conf *alpha.DarknessConfig) []yunyun.FullPathFile {
	c := make(chan yunyun.FullPathFile)
	go FindFilesByExt(conf, c)
	return rei.Collect(c)
}

// FindFilesByExtSimpleDirs is the same as `FindFilesByExt` but it simply blocks the
// parent goroutine until it processes all the results and returns only the results
// which are children of the passed dirs.
func FindFilesByExtSimpleDirs(conf *alpha.DarknessConfig, dirs []string) []yunyun.FullPathFile {
	files := FindFilesByExtSimple(conf)
	// If no dirs passed, run no filtering.
	if len(dirs) < 1 {
		return files
	}
	// Only return files that have passed dirs as parents.
	return g.Filter(func(inputFilename yunyun.FullPathFile) bool {
		return g.Anyf(func(v string) bool {
			return strings.HasPrefix(string(inputFilename),
				string(conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(v))))
		}, dirs)
	}, files)
}
