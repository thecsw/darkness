package emilia

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
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
		// If it's vendored, then retrieve a local copy (if doesn't already
		// exist) and stub it in as the full path
		if Config.VendorGalleries {
			return galleryVendorItem(item)
		}
		return yunyun.FullPathFile(item.Item)
	}
	return JoinPath(yunyun.JoinRelativePaths(item.Path, item.Item))
}

// GalleryPreview takes an original image's path and returns
// the preview path of it. Previews are always .jpg
func GalleryPreview(item *GalleryItem) yunyun.FullPathFile {
	if item.IsExternal {
		return JoinPath(yunyun.JoinRelativePaths(item.Path, galleryItemHash(item)))
	}
	filename := filepath.Base(string(item.Item))
	ext := filepath.Ext(filename)
	return JoinPath(yunyun.JoinRelativePaths(item.Path,
		yunyun.RelativePathFile(strings.TrimSuffix(filename, ext)+"_small.jpg"),
	))
}

const (
	// vendorDirectory is the name of the dir where vendor images are stored.
	vendorDirectory yunyun.RelativePathDir = "darkness_vendor"
)

var (
	vendorClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			MaxConnsPerHost:     1,
		},
		Timeout: 10 * time.Second,
	}
	vendorLock = &sync.Mutex{}
)

// galleryVendorItem vendors given item and returns the full path of the file.
//
// Only call this function on remote images, it's up to the user to make the
// .IsExternal check before calling this. SLOW function because of network calls.
//
// If the vendoring fails at any point, fallback to the remote image path.
func galleryVendorItem(item *GalleryItem) yunyun.FullPathFile {
	// Process only one vendor request at a time.
	vendorLock.Lock()
	// Unlock so the next vendor request can get processed.
	defer vendorLock.Unlock()

	hashedName := galleryItemHash(item)
	vendoredPath := yunyun.JoinRelativePaths(vendorDirectory, hashedName)
	localVendoredPath := filepath.Join(Config.WorkDir, string(vendoredPath))

	// Create the two types of return.
	fallbackReturn := yunyun.FullPathFile(item.Item)
	expectedReturn := JoinPath(vendoredPath)

	// Check if the image was already vendored, if it was, return it immediately.
	if FileExists(localVendoredPath) {
		return expectedReturn
	}

	start := time.Now()
	fmt.Printf("Vendoring %s... ", hashedName)

	req, err := http.NewRequest(http.MethodGet, string(item.Item), nil)
	if err != nil {
		fmt.Printf("Failed to create a request: %s", err.Error())
		return fallbackReturn
	}

	resp, err := vendorClient.Do(req)
	// resp, err := http.Get(string(item.Item))
	if err != nil {
		fmt.Printf("Failed to vendor %s: %s\n", hashedName, err.Error())
		return fallbackReturn
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		fmt.Printf("Got status %d: %s\n", resp.StatusCode, resp.Status)
		return fallbackReturn
	}
	defer resp.Body.Close()

	// Read the image with imaging and convert it by force to jpeg.
	//
	// Respect EXIF flags with AutoOrientation turned on.
	img, err := imaging.Decode(resp.Body, imaging.AutoOrientation(true))
	if err != nil {
		fmt.Printf("Failed to decode %s: %s\n", hashedName, err.Error())
		return fallbackReturn
	}

	// Open the file writer and encode the image there.
	imgFile, err := os.Create(localVendoredPath)
	if err != nil {
		fmt.Printf("Failed to create file %s: %s\n", localVendoredPath, err.Error())
		return fallbackReturn
	}
	defer imgFile.Close()

	// Decode the image into the file.
	if err := imaging.Encode(imgFile, img, imaging.JPEG); err != nil {
		fmt.Printf("Failde to encode %s: %s\n", hashedName, err.Error())
		return fallbackReturn
	}

	finish := time.Now()
	fmt.Printf("done in %d ms\n", finish.Sub(start).Milliseconds())

	// Finally.
	return expectedReturn
}

// galleryItemHash returns a hashed name of a gallery item link.
func galleryItemHash(item *GalleryItem) yunyun.RelativePathFile {
	return yunyun.RelativePathFile(sha256String(string(item.Item))[:7] + ".jpg")
}

// sha256String hashes given string to sha256.
func sha256String(what string) string {
	ans := sha256.Sum256([]byte(what))
	return hex.EncodeToString(ans[:])
}

// FileExists returns true if file exists, false otherwise (in error too).
func FileExists(path string) bool {
	info, err := os.Stat(string(path))
	return info != nil && !os.IsNotExist(err)
}
