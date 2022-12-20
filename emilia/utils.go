package emilia

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// InputFilenameToOutput converts input filename to the filename to write.
func InputFilenameToOutput(file yunyun.FullPathFile) string {
	outputFilename := strings.Replace(filepath.Base(string(file)),
		Config.Project.Input, Config.Project.Output, 1)
	return filepath.Join(filepath.Dir(string(file)), outputFilename)
}

// InputToOutput converts a single input file to its output.
func InputToOutput(file yunyun.FullPathFile) string {
	data, err := ioutil.ReadFile(filepath.Clean(string(file)))
	if err != nil {
		panic(err)
	}
	page := ParserBuilder.BuildParser(Pack(file, string(data))).Parse()
	return EnrichAndExportPage(EnrichPage(page))
}

// Pack cleans the filename from absolute workspace prefix.
func Pack(filename yunyun.FullPathFile, data string) (yunyun.RelativePathFile, string) {
	return RelPathToWorkdir(filename), data
}

// PackRef cleans the filename from absolute workspace prefix by refs.
func PackRef(filename *yunyun.FullPathFile, data *string) (yunyun.RelativePathFile, string) {
	return RelPathToWorkdir(*filename), *data
}

// relPathToWorkdir returns path trimmed by the workspace
func RelPathToWorkdir(filename yunyun.FullPathFile) yunyun.RelativePathFile {
	return yunyun.RelativePathFile(strings.TrimPrefix(string(filename), Config.WorkDir+`/`))
}

// EnrichAndExportPage automatically applies all the emilia enhancements
// and converts Page into an html document.
func EnrichAndExportPage(page *yunyun.Page) string {
	result := AddHolosceneTitles(
		ExporterBuilder.BuildExporter(EnrichPage(page)).Export(),
		func() int {
			if strings.HasSuffix(string(page.Location), "quotes") {
				return -1
			}
			return 1
		}())
	return result
}

// EnrichPage applies common emilia enhancements.
func EnrichPage(page *yunyun.Page) *yunyun.Page {
	return page.Options(
		WithResolvedComments(),
		WithEnrichedHeadings(),
		WithFootnotes(),
		WithMathSupport(),
		WithSourceCodeTrimmedLeftWhitespace(),
		WithSyntaxHighlighting(),
		WithLazyGalleries(),
	)
}

// GalleryItem is a struct that holds the gallery item path
// and a flag whether it is external (URL regexp matches).
type GalleryItem struct {
	// Item is the link that was provided.
	Item yunyun.RelativePathFile
	// Path is the path of the local gallery source file.
	Path yunyun.RelativePathDir
	// IsExternal runs a URL regexp check.
	IsExternal bool
	// Text found through the link regexp.
	Text string
	// Description found through the link regexp.
	Description string
	// OriginalLine is the original line that include org options.
	OriginalLine string
	// Link is an optional parameter that the gallery item should
	// also link to something.
	Link string
}

// NewGalleryItem creates a new helper `GalleryItem` and
// decides whether the passed item is an external link or not.
func NewGalleryItem(page *yunyun.Page, content *yunyun.Content, wholeLine string) *GalleryItem {
	extractedLinks := yunyun.ExtractLinks(wholeLine)
	// If image wasn't found, then the whole line should be counted as the image path.
	image := wholeLine
	text := ""
	description := ""
	if len(extractedLinks) > 0 {
		image = extractedLinks[0].Link
		text = extractedLinks[0].Text
		description = extractedLinks[0].Description
	}
	optionalLink := ""
	if len(extractedLinks) > 1 {
		optionalLink = extractedLinks[1].Link
	}
	return &GalleryItem{
		Item:         yunyun.RelativePathFile(image),
		Path:         yunyun.JoinPaths(page.Location, content.GalleryPath),
		IsExternal:   yunyun.URLRegexp.MatchString(image),
		Text:         text,
		Description:  description,
		OriginalLine: wholeLine,
		Link:         optionalLink,
	}
}

func GalleryImage(item *GalleryItem) yunyun.FullPathFile {
	if item.IsExternal {
		return yunyun.FullPathFile(item.Item)
	}
	return JoinPath(yunyun.JoinRelativePaths(item.Path, item.Item))
}

// GalleryPreview takes an original image's path and returns
// the preview path of it.
func GalleryPreview(item *GalleryItem) yunyun.FullPathFile {
	if item.IsExternal {
		return JoinPath(yunyun.JoinRelativePaths(item.Path,
			yunyun.RelativePathFile(md5String(string(item.Item))+"_preview.jpeg"),
		))
	}
	filename := filepath.Base(string(item.Item))
	ext := filepath.Ext(filename)
	return JoinPath(yunyun.JoinRelativePaths(item.Path,
		yunyun.RelativePathFile(strings.TrimSuffix(filename, ext)+"_preview"+ext),
	))
}

func md5String(what string) string {
	ans := md5.Sum([]byte(what))
	return hex.EncodeToString(ans[:])
}
