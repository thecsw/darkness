package misa

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/narumi"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/ichika/hizuru"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/darkness/yunyun/rss"
	"github.com/thecsw/gana"
)

const (
	// rssGenerator is the generator string used in the RSS feed.
	rssGenerator = "Darkness (sandyuraz.com/darkness)"
)

// GenerateRssFeed generates an RSS feed based on the given config and directories.
func GenerateRssFeed(conf *alpha.DarknessConfig, rssFilename string, rssDirectories []string, dryRun bool) {
	// Get all all the pages we can build out.
	allPages := hizuru.BuildPagesSimple(conf, rssDirectories)
	// Try to retrieve the top root page to get channel description. If not found, use the
	// website's title as the description.
	topPage := gana.First(gana.Filter(func(page *yunyun.Page) bool { return page.Location == "." }, allPages))
	rootDescription := conf.RSS.Description
	if topPage != nil {
		rootDescription = getDescription(topPage, conf.Website.DescriptionLength*4)
	}
	// If both the top page and RSS config have no description, default to the title.
	if len(rootDescription) < 1 {
		rootDescription = conf.Title
	}

	sort.Slice(allPages, func(i, j int) bool { return allPages[i].Title < allPages[j].Title })

	// Get all pages that have dates defined, we only use those to be included in the rss feed.
	pages := Pages(gana.Filter(func(page *yunyun.Page) bool {
		_, dateFound := getDate(page)
		return dateFound
	}, allPages))

	// Sort the pages in descending order of dates.
	sort.Sort(pages)

	// Create RSS items.
	items := make([]rss.Item, 0, len(pages))

	func() {
		defer puck.Stopwatch("Built RSS pages", "num", len(pages)).Record()
		for _, page := range pages {
			// Skip drafts.
			if page.Accoutrement.Draft.IsEnabled() {
				continue
			}
			// Create the category name and location.
			categoryName, categoryLocation := page.Title, page.Location
			if categoryPage := getCategory(page, allPages); categoryPage != nil {
				categoryName = categoryPage.Title
				categoryLocation = categoryPage.Location
			}

			// Create the RSS item.
			items = append(items, rss.Item{
				XMLName: xml.Name{},
				Title:   yunyun.RemoveFormatting(yunyun.FancyText(page.Title)),
				Link:    conf.Url + string(page.Location),
				Description: yunyun.FancyText(getDescription(page, conf.Website.DescriptionLength*4)) +
					" [ Continue reading... ]",
				Author:    page.Author,
				Category:  &rss.Category{Value: categoryName, Domain: conf.Url + string(categoryLocation)},
				Enclosure: &rss.Enclosure{},
				Guid:      &rss.Guid{Value: conf.Url + string(page.Location), IsPermaLink: true},
				PubDate:   mustDate(page).Format(rss.RSSFormat),
				Source:    &rss.Source{Value: conf.Title, Url: conf.Url},
			})
		}
	}()

	// Create the final feed.
	feed := &rss.RSS{
		Version: rss.RSSVersion,
		Channel: &rss.Channel{
			XMLName:        xml.Name{},
			Title:          yunyun.FancyText(conf.Title),
			Link:           conf.Url,
			Description:    yunyun.FancyText(rootDescription),
			Language:       conf.RSS.Language,
			Copyright:      conf.RSS.Copyright,
			ManagingEditor: conf.RSS.ManagingEditor,
			WebMaster:      conf.RSS.WebMaster,
			PubDate:        mustDate(gana.First(pages)).Format(rss.RSSFormat),
			LastBuildDate:  time.Now().Format(rss.RSSFormat),
			Category:       conf.RSS.Category,
			Generator:      rssGenerator,
			Docs:           rss.RSSDocs,
			TTL:            60,
			Items:          items,
		},
	}

	xmlTarget := "stdout"
	xmlFile := os.Stdout
	var err error
	if !dryRun {
		xmlTarget = string(conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(rssFilename)))
		xmlFile, err = os.Create(filepath.Clean(xmlTarget))
		if err != nil {
			puck.Logger.Error("Creating file", "path", xmlTarget, "err", err)
			os.Exit(1)
		}
	}

	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "  ")
	if err := encoder.Encode(feed); err != nil {
		puck.Logger.Error("Encoding to xml", "path", xmlTarget, "err", err)
		os.Exit(1)
	}
	if err := xmlFile.Close(); err != nil {
		puck.Logger.Error("Closing file", "path", xmlTarget, "err", err)
		os.Exit(1)
	}
	puck.Logger.Print("Created rss file", "path", xmlTarget)
}

// getDate takes a page and returns its date if any found.
func getDate(page *yunyun.Page) (time.Time, bool) {
	parsed := narumi.ConvertHoloscene(page.Date)
	return parsed, parsed.Unix() != 0 && parsed.Day() != 31 && parsed.Year() != 2000
}

var categoryCache = make(map[string]*yunyun.Page)

func getCategory(page *yunyun.Page, pages Pages) *yunyun.Page {
	categoryName := strings.TrimSuffix(string(page.Location), "/"+filepath.Base(string(page.Location)))
	if v, ok := categoryCache[categoryName]; ok {
		return v
	}
	for _, allPage := range pages {
		if allPage.Location == yunyun.RelativePathDir(categoryName) {
			categoryCache[categoryName] = allPage
			return allPage
		}
	}
	return nil
}

// Pages is custom type of slice of pages to enable sorting.
type Pages []*yunyun.Page

// Len returns the number of pages.
func (p Pages) Len() int { return len(p) }

// Swap swaps lol.
func (p Pages) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Less sorts the array in descending order.
func (p Pages) Less(i, j int) bool { return mustDate(p[i]).Unix() > mustDate(p[j]).Unix() }

func mustDate(v *yunyun.Page) time.Time {
	t, f := getDate(v)
	if !f {
		panic("must be date")
	}
	return t
}

const (
	// Minimum length of the description
	descriptionMinLength = 14
)

// getDescription returns the description of the page
// It will return the first paragraph that is not empty and not a holoscene time
// If no such paragraph is found, it will return an empty string
// If the description is less than 14 characters, it will return an empty string
func getDescription(page *yunyun.Page, length int) string {
	// Find the first paragraph for description
	description := ""
	for _, content := range page.Contents {
		// We are only looking for paragraphs
		if !content.IsParagraph() {
			continue
		}
		// Skip holoscene times
		paragraph := strings.TrimSpace(content.Paragraph)
		if paragraph == "" || puck.HEregex.MatchString(paragraph) {
			continue
		}

		cleanText := yunyun.RemoveFormatting(paragraph[:gana.Min(len(paragraph), length+10)])
		description = cleanText[:gana.Max(len(cleanText)-10, 0)] + "..."
		if len(description) < descriptionMinLength {
			continue
		}
		break
	}
	return description
}
