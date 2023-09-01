package alpha

import (
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// JoinGeneric joins the target path with the final root path (url or local).
func (conf RuntimeConfig) JoinGeneric(what ...string) string {
	if !conf.isUrlLocal {
		return conf.UrlPath.JoinPath(what...).String()
	}
	return filepath.Join(append(conf.urlSlice, what...)...)
}

// Join joins the target path with the final root path (url or local).
func (conf RuntimeConfig) Join(relative ...yunyun.RelativePathFile) yunyun.FullPathFile {
	return yunyun.FullPathFile(conf.JoinGeneric(yunyun.AnyPathsToStrings(relative)...))
}

// JoinGeneric joins target path with the working directory.
func (workDir WorkingDirectory) JoinGeneric(target string) string {
	return filepath.Join(string(workDir), target)
}

// Join joins target path with the working directory.
func (workDir WorkingDirectory) Join(target yunyun.RelativePathFile) yunyun.FullPathFile {
	return yunyun.FullPathFile(workDir.JoinGeneric(string(target)))
}

// Rel returns path trimmed by the workspace
func (workDir WorkingDirectory) Rel(filename yunyun.FullPathFile) yunyun.RelativePathFile {
	return yunyun.RelativePathFile(strings.TrimPrefix(string(filename), string(workDir+`/`)))
}

// PackRel cleans the filename from absolute workspace prefix.
func (workDir WorkingDirectory) PackRel(filename yunyun.FullPathFile, data string) (yunyun.RelativePathFile, string) {
	return workDir.Rel(filename), data
}

// PackRelRef cleans the filename from absolute workspace prefix by refs.
func (workDir WorkingDirectory) PackRelRef(filename *yunyun.FullPathFile, data *string) (yunyun.RelativePathFile, string) {
	return workDir.Rel(*filename), *data
}
