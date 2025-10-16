package misa

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/narumi"
	"github.com/thecsw/darkness/v3/emilia/puck"
	"github.com/thecsw/darkness/v3/ichika/hizuru"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/darkness/v3/yunyun/rss"
	"github.com/thecsw/gana"
)

const (
	// rssGenerator is the generator string used in the RSS feed.
	rssGenerator = "Darkness (sandyuraz.com/darkness)"
)

// GenerateRssFeed generates an RSS feed based on the given config and directories.
func GenerateRssFeed(conf *alpha.DarknessConfig, rssFilename string, rssDirectories []string, dryRun bool) {
	initLog()
	recordGlobalMacros(conf)
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
		_, dateFound := narumi.ConvertHoloscene(page.Date)
		if !dateFound {
			logger.Debug("Skipping because no date found", "page", page.Location)
		}
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
				logger.Warn("Skipping draft", "page", page.Location)
				continue
			}
			// Create the category name and location.
			categoryName, categoryLocation := page.Title, page.Location
			if categoryPage := getCategory(page, allPages); categoryPage != nil {
				categoryName = categoryPage.Title
				categoryLocation = categoryPage.Location
			}

			// Override the title if the page has a custom RSS title.
			finalTitle := page.Title
			if len(page.Accoutrement.RssTitle) > 0 {
				finalTitle = page.Accoutrement.RssTitle
			}

			// Add the RSS prefix to the title.
			finalTitle = page.Accoutrement.RssPrefix + " " + finalTitle

			// Let's update the time if needed.
			parsedDate, isValid := narumi.ConvertHoloscene(page.Date)
			if !isValid {
				logger.Warn("Skipping invalid publication date", "page", page.Location)
				continue
			}
			finalLocation, err := time.LoadLocation(conf.RSS.Timezone)
			// Fallback to UTC
			if err != nil {
				finalLocation = time.UTC
			}
			hour, minute := parsedDate.Hour(), parsedDate.Minute()
			if hour == 0 && minute == 0 {
				hour = conf.RSS.DefaultHour
				minute = 0
			}
			finalDate := time.Date(
				parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
				hour, minute, 0, 0, finalLocation)

			// Create the RSS item.
			items = append(items, rss.Item{
				XMLName: xml.Name{},
				Title:   yunyun.RemoveFormatting(yunyun.FancyText(finalTitle)),
				Link:    string(conf.Runtime.JoinDir(page.Location)),
				Description: yunyun.FancyText(getDescription(page, conf.Website.DescriptionLength*4)) +
					" [ Continue reading... ]",
				Author:    page.Author,
				Category:  &rss.Category{Value: categoryName, Domain: conf.Url + string(categoryLocation)},
				Enclosure: &rss.Enclosure{},
				Guid:      &rss.Guid{Value: conf.Url + string(page.Location), IsPermaLink: true},
				PubDate:   finalDate.Format(rss.RSSFormat),
				Source:    &rss.Source{Value: conf.Title, Url: conf.Url},
			})
		}
	}()

	// Try to find the pub date, if none, then reuse the build date
	buildDate := time.Now().Format(rss.RSSFormat)
	pubDate := buildDate

	firstPage := gana.First(pages)
	if firstPage != nil {
		pubDate = mustDate(firstPage).Format(rss.RSSFormat)
	}

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
			PubDate:        pubDate,
			LastBuildDate:  buildDate,
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
			logger.Error("Creating file", "path", xmlTarget, "err", err)
			os.Exit(1)
		}
	}

	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "  ")
	if err := encoder.Encode(feed); err != nil {
		logger.Error("Encoding to xml", "path", xmlTarget, "err", err)
		os.Exit(1)
	}
	if err := xmlFile.Close(); err != nil {
		logger.Error("Closing file", "path", xmlTarget, "err", err)
		os.Exit(1)
	}
	logger.Info("Created rss file", "path", conf.Runtime.WorkDir.Rel(yunyun.FullPathFile(xmlTarget)))
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
	t, f := narumi.ConvertHoloscene(v.Date)
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
