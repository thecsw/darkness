package html

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/rem"
	"github.com/thecsw/darkness/v3/ichika/akane"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
)

// safeIntToUint safely converts an int to uint by handling negative values
func safeIntToUint(val int) uint {
	if val < 0 {
		return 0
	}
	return uint(val)
}

var (
	flexOptionRegexp  = regexp.MustCompile(`:flex (\d+)`)
	vFlexOptionRegexp = regexp.MustCompile(`:v-flex (\d+)`)
	flexBreakRegexp   = regexp.MustCompile(`:flex-break(?:\s|$)`)
)

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
	return safeIntToUint(ret)
}

// extractCustomVFlex extracts an item's requested share of the remaining
// vertical space in a gallery column. A value outside 1-100 is ignored.
func extractCustomVFlex(s string) uint {
	matches := vFlexOptionRegexp.FindAllStringSubmatch(s, -1)
	if len(matches) < 1 || len(matches[0]) < 2 {
		return 0
	}
	ret, err := strconv.Atoi(matches[0][1])
	if err != nil || ret < 1 || ret > 100 {
		return 0
	}
	return uint(ret)
}

// hasFlexBreak reports whether an item should start a new gallery column.
func hasFlexBreak(s string) bool {
	return flexBreakRegexp.MatchString(s)
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

// makeFlexItemContent renders a gallery image without its horizontal layout
// wrapper. Vertical gallery columns add semantic classes for the stylesheet.
func makeFlexItemContent(conf *alpha.DarknessConfig, item rem.GalleryItem, vertical bool) string {
	galleryItemClass := "gallery-item"
	imageClasses := ""
	// The preview is also set as the anchor's background so it stays visible
	// while lazysizes swaps the img's src (the browser blanks the img during
	// the full image download).
	backgroundSize := "cover"
	if vertical {
		galleryItemClass += " gallery-row-link"
		imageClasses += " gallery-row-image"
		backgroundSize = "contain"
	}
	imageClasses += resolveCustomFlexItemClasses(item.OriginalLine)
	return fmt.Sprintf(`<a%s class="%s" style="background-image: url('%s'); background-size: %s;">
<img class="item lazyload%s" src="%s" data-src="%s" title="%s" alt="%s">
</a>`,
		// Optionally link the gallery image to something.
		hrefGalleryTagIfLinkGiven(item),
		galleryItemClass,
		// Path to the gallery image's preview (underlay during the swap).
		rem.GalleryPreview(conf, item),
		// cover/contain to match how the full image is displayed.
		backgroundSize,
		// Additionally-enabled options, like no-zoom.
		imageClasses,
		// Path to the gallery image's preview.
		rem.GalleryPreview(conf, item),
		// Path to the image (either external, local, or vendored).
		processGalleryItem(conf, item),
		// The text to show on the image hover.
		processTitle(item.Description),
		// The alt description of the image.
		processTitle(item.Text), // using processTitle for lighter markup
	)
}

// makeFlexItem will make an item of the flexbox .gallery with 1/3 width.
func makeFlexItem(conf *alpha.DarknessConfig, item rem.GalleryItem, width uint) string {
	// See if there is a custom flex width requested for the item.
	if customFlex := extractCustomFlex(item.OriginalLine); customFlex != 0 {
		width = customFlex
	}
	return fmt.Sprintf("<div class=\"flex-%d hide-overflow ease-transition\">\n%s\n</div>",
		// The percentage (or flex class) of the page's width to occupy.
		width,
		makeFlexItemContent(conf, item, false),
	)
}

type galleryColumnItem struct {
	item          rem.GalleryItem
	verticalShare float64
}

type galleryColumn struct {
	items []galleryColumnItem
	width uint
}

// makeGalleryColumns groups consecutive :v-flex items into a flex column.
// Each value consumes that percentage of the vertical space remaining in the
// column. The first unannotated item fills the rest and closes the column.
func makeGalleryColumns(items []rem.GalleryItem, defaultWidth uint) []galleryColumn {
	columns := make([]galleryColumn, 0, len(items))
	for itemIndex := 0; itemIndex < len(items); {
		item := items[itemIndex]
		verticalFlex := extractCustomVFlex(item.OriginalLine)
		if verticalFlex == 0 {
			columns = append(columns, galleryColumn{
				items: []galleryColumnItem{{item: item, verticalShare: 100}},
				width: galleryColumnWidth(item, defaultWidth),
			})
			itemIndex++
			continue
		}

		column := galleryColumn{width: galleryColumnWidth(item, defaultWidth)}
		remaining := 100.0
		for itemIndex < len(items) {
			item = items[itemIndex]
			if len(column.items) > 0 && hasFlexBreak(item.OriginalLine) {
				break
			}
			verticalFlex = extractCustomVFlex(item.OriginalLine)
			if verticalFlex == 0 {
				column.items = append(column.items, galleryColumnItem{item: item, verticalShare: remaining})
				itemIndex++
				break
			}

			share := remaining * float64(verticalFlex) / 100
			column.items = append(column.items, galleryColumnItem{item: item, verticalShare: share})
			remaining -= share
			itemIndex++
			if remaining == 0 {
				break
			}
		}
		columns = append(columns, column)
	}
	return columns
}

func galleryColumnWidth(item rem.GalleryItem, defaultWidth uint) uint {
	if customFlex := extractCustomFlex(item.OriginalLine); customFlex != 0 {
		return customFlex
	}
	return defaultWidth
}

func makeVerticalFlexColumn(conf *alpha.DarknessConfig, column galleryColumn) string {
	rows := make([]string, 0, len(column.items))
	for _, row := range column.items {
		rows = append(rows, fmt.Sprintf(`<div class="v-flex gallery-row ease-transition" style="--v-flex: %s%%;">
%s
</div>`,
			strconv.FormatFloat(row.verticalShare, 'f', -1, 64),
			makeFlexItemContent(conf, row.item, true),
		))
	}
	return fmt.Sprintf(`<div class="flex-%d gallery-column">
%s
</div>`, column.width, strings.Join(rows, "\n"))
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
	items := make([]rem.GalleryItem, 0, len(content.List))
	for _, listItem := range content.List {
		items = append(items, rem.NewGalleryItem(e.conf, e.page, content, listItem.Text))
	}
	columns := makeGalleryColumns(items, content.GalleryImagesPerRow)
	makeGalleryColumn := func(column galleryColumn) string {
		if len(column.items) == 1 && extractCustomVFlex(column.items[0].item.OriginalLine) == 0 {
			return makeFlexItem(e.conf, column.items[0].item, column.width)
		}
		return makeVerticalFlexColumn(e.conf, column)
	}
	return fmt.Sprintf(`
<div class="gallery-container">
<center>
<div class="gallery">
%s
</div>
</center>
</div>
`, strings.Join(gana.Map(makeGalleryColumn, columns), "\n"))
}
