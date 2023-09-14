package html

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/rem"
	"github.com/thecsw/darkness/ichika/akane"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

var flexOptionRegexp = regexp.MustCompile(`:flex (\d+)`)

// extractCustomFlex extract custom flex class `:flex [1,5]`
func extractCustomFlex(s string) uint {
	matches := flexOptionRegexp.FindAllStringSubmatch(s, -1)
	if len(matches) < 1 {
		return 0
	}
	if len(matches[0]) < 1 {
		return 0
	}
	ret, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return 0
	}
	return uint(ret)
}

// hrefGalleryTagIfLinkGiven returns an href tag if gallery link is found,
// an empty string otherwise.
func hrefGalleryTagIfLinkGiven(item rem.GalleryItem) string {
	if item.Link == "" {
		return ""
	}
	return fmt.Sprintf(` href="%s"`, item.Link)
}

// resolveCustomFlexItemClasses searches for custom support flex item classes.
func resolveCustomFlexItemClasses(wholeLine string) string {
	what := ""
	if strings.Contains(wholeLine, ":no-zoom") {
		what += " no-zoom"
	}
	return what
}

// makeFlexItem will make an item of the flexbox .gallery with 1/3 width
func makeFlexItem(conf *alpha.DarknessConfig, item rem.GalleryItem, width uint) string {
	// See if there is a custom flex width requested for the item.
	if customFlex := extractCustomFlex(item.OriginalLine); customFlex != 0 {
		width = customFlex
	}
	// If the image is external AND vendor galleries option is enabled,
	// then get a local copy of the remote image (if it doesn't already exist)
	// and stub it in.
	return fmt.Sprintf(`<div class="flex-%d hide-overflow ease-transition">
<a%s class="gallery-item">
<img class="item lazyload %s" src="%s" data-src="%s" title="%s" alt="%s">
</a>
</div>`,
		// The percentage (or flex class) of the page's width to occupy.
		width,
		// Optionally link the gallery image to something.
		hrefGalleryTagIfLinkGiven(item),
		// Additionally-enabled options, like no-zoom.
		resolveCustomFlexItemClasses(item.OriginalLine),
		// Path to the gallery image's preview.
		rem.GalleryPreview(conf, item),
		// Path to the image (either external, local, or vendored).
		processGalleryItem(conf, item),
		// The text to show on the image hover.
		item.Description,
		// The alt descriptino of the image.
		item.Text,
	)
}

// processGalleryItem takes a gallery item and returns the full path, while also submitting an
// akane request to download the gallery image.
func processGalleryItem(conf *alpha.DarknessConfig, item rem.GalleryItem) yunyun.FullPathFile {
	path, shouldBeVendored := rem.GalleryImage(conf, item)
	if shouldBeVendored {
		akane.RequestGalleryVendor(item)
	}
	return path
}

// gallery will create a flexbox gallery as defined in .gallery css class
func (e *state) gallery(content *yunyun.Content) string {
	makeFlexItemWithFolder := func(s yunyun.ListItem) string {
		return makeFlexItem(e.conf, rem.NewGalleryItem(e.page, content, s.Text), content.GalleryImagesPerRow)
	}
	return fmt.Sprintf(`
<div class="gallery-container">
<center>
<div class="gallery">
%s
</div>
</center>
</div>
`, strings.Join(gana.Map(makeFlexItemWithFolder, content.List), "\n"))
}
