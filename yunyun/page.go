package yunyun

import (
	"path/filepath"

	"github.com/thecsw/gana"
)

// Page is a struct for holding the page contents.
type Page struct {
	// File is the original filename of the page (optional).
	File RelativePathFile
	// Location is the Location of the page.
	Location RelativePathDir
	// Title is the title of the page.
	Title string
	// Date is the date of the page.
	Date string
	// DateHoloscene tells us whether the first paragraph
	// on the page is given as holoscene date stamp.
	DateHoloscene bool
	// Contents is the contents of the page.
	Contents Contents
	// Footnotes is the footnotes of the page.
	Footnotes []string
	// Scripts is the scripts of the page.
	Scripts []string
	// Stylesheets is the list of css of the page.
	Stylesheets []string
	// HtmlHead is the list of extra HTML declaration to add in the head.
	HtmlHead []string
	// Accoutrement are additional options enabled on a page.
	Accoutrement *Accoutrement
}

// MetaTag is a struct for holding the meta tag.
type MetaTag struct {
	// Name is the name of the meta tag.
	Name string
	// Content is the content of the meta tag.
	Content string
	// Propery is the property of the meta tag.
	Property string
}

// Link is a struct for holding the link tag.
type Link struct {
	// Rel is the rel of the link tag.
	Rel string
	// Type is the type of the link tag.
	Type string
	// Href is the href of the link tag.
	Href string
}

// RelativePath is used in `Page` for `Location` to make sure
// that we are passing correct understanding of that it should
// only have the relative path to the workspace -- this should
// be a directory with NO filename and NO base.
type RelativePathDir string

// RelativePathFile is similar to `RelativePath` but also
// includes the filename in the end as base.
type RelativePathFile string

// RelativePath is a type constraint for `RelativePath`-style paths.
type RelativePath interface {
	RelativePathDir | RelativePathFile
}

// FullPath is the result of joining emilia root with `RelativePath`.
type FullPathDir string

// FullPathFile is the result of joining emilia root with `RelativePathFile`.
type FullPathFile string

// FullPath is a type constraint for `FullPath`-style types.
type FullPath interface {
	FullPathDir | FullPathFile
}

// AnyPath is a full generalization of relative and full paths.
type AnyPath interface {
	FullPath | RelativePath
}

// RelativePathTrim returns the directory of the relative file.
func RelativePathTrim(filename RelativePathFile) RelativePathDir {
	return RelativePathDir(filepath.Dir(string(filename)))
}

// JoinRelativePaths joins relative paths.
func JoinRelativePaths(dir RelativePathDir, file RelativePathFile) RelativePathFile {
	return JoinPaths(RelativePathFile(dir), file)
}

// JoinPaths joins relative paths.
func JoinPaths[T AnyPath](what ...T) T {
	return T(filepath.Join(AnyPathsToStrings(what)...))
}

// PathsToString converts an array of `AnyPath` to `string`.
func AnyPathsToStrings[T AnyPath](what []T) []string {
	return gana.Map(func(t T) string { return string(t) }, what)
}
