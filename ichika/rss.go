package ichika

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/ichika/rss"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

const (
	rssXMLFilename = "feed.xml"
	rssGenerator   = "Darkness (sandyuraz.com/darkness)"
)

func rssf(rssFilename string, rssDirectories []string, dryRun bool) {
	// Get all all the pages we can build out.
	allPages := buildPagesSimple(rssDirectories)

	// Try to retrieve the top root page to get channel description. If not found, use the
	// website's title as the description.
	topPage := gana.First(gana.Filter(func(page *yunyun.Page) bool { return page.Location == "." }, allPages))
	rootDescription := emilia.Config.RSS.Description
	if topPage != nil {
		rootDescription = emilia.GetDescription(topPage, emilia.Config.Website.DescriptionLength*4)
	}
	// If both the top page and RSS config have no description, default to the title.
	if len(rootDescription) < 1 {
		rootDescription = emilia.Config.Title
	}

	// Get all pages that have dates defined, we only use those to be included in the rss feed.
	pages := Pages(gana.Filter(func(page *yunyun.Page) bool { return getDate(page) != nil }, allPages))

	// Sort the pages in descending order of dates.
	sort.Sort(pages)

	// Create RSS items.
	items := make([]*rss.Item, 0, len(pages))

	for _, page := range pages {
		categoryName, categoryLocation := page.Title, page.Location
		if categoryPage := getCategory(page, allPages); categoryPage != nil {
			categoryName = categoryPage.Title
			categoryLocation = categoryPage.Location
		}

		items = append(items, &rss.Item{
			Title:       yunyun.RemoveFormatting(page.Title),
			Link:        emilia.Config.URL + string(page.Location),
			Description: emilia.GetDescription(page, emilia.Config.Website.DescriptionLength*4) + " [ Continue reading... ]",
			Category: &rss.Category{
				Value:  categoryName,
				Domain: emilia.Config.URL + string(categoryLocation),
			},
			Guid:    &rss.Guid{Value: emilia.Config.URL + string(page.Location), IsPermaLink: true},
			PubDate: getDate(page).Format(rss.RSSFormat),
			Source:  &rss.Source{Value: emilia.Config.Title, URL: emilia.Config.URL},
		})
	}

	// Create the final feed.
	feed := &rss.RSS{
		Version: rss.RSSVersion,
		Channel: &rss.Channel{
			XMLName:        xml.Name{},
			Title:          emilia.Config.Title,
			Link:           emilia.Config.URL,
			Description:    rootDescription,
			Language:       emilia.Config.RSS.Language,
			Copyright:      emilia.Config.RSS.Copyright,
			ManagingEditor: emilia.Config.RSS.ManagingEditor,
			WebMaster:      emilia.Config.RSS.WebMaster,
			PubDate:        getDate(gana.First(pages)).Format(rss.RSSFormat),
			LastBuildDate:  time.Now().Format(rss.RSSFormat),
			Category:       emilia.Config.RSS.Category,
			Generator:      rssGenerator,
			Docs:           rss.RSSDocs,
			TTL:            60,
			Items:          items,
		},
	}

	xmlTarget := string(emilia.JoinWorkdir(rssXMLFilename))
	feedXml, err := os.Create(filepath.Clean(xmlTarget))
	if err != nil {
		fmt.Printf("couldn't create %s: %s\n", xmlTarget, err)
		os.Exit(1)
	}
	encoder := xml.NewEncoder(feedXml)
	encoder.Indent("", "  ")
	if err := encoder.Encode(feed); err != nil {
		fmt.Printf("failed to encode %s: %s\n", xmlTarget, err)
		os.Exit(1)
	}
	if err := feedXml.Close(); err != nil {
		fmt.Printf("failed to close %s: %s", xmlTarget, err)
		os.Exit(1)
	}
	fmt.Printf("Created rss file in %s\n", xmlTarget)
}

// getDate takes a page and returns its date if any found.
func getDate(page *yunyun.Page) *time.Time {
	parsed := emilia.ConvertHoloscene(page.Date)
	if parsed != nil && parsed.Day() != 31 && parsed.Year() != 2000 {
		return parsed
	}
	return nil
}

var (
	categoryCache = make(map[string]*yunyun.Page)
)

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
func (p Pages) Less(i, j int) bool { return getDate(p[i]).Unix() > getDate(p[j]).Unix() }

// Will return a slice of built pages that have dirs as parents (empty dirs will return everything).
func buildPagesSimple(dirs []string) Pages {
	inputs := emilia.FindFilesByExtSimpleDirs(emilia.Config.Project.Input, dirs)
	pages := make([]*yunyun.Page, 0, len(inputs))
	for _, input := range inputs {
		bundle := openFile(input)
		data, err := io.ReadAll(bundle.Second)
		if err != nil {
			fmt.Printf("failed to read %s: %s", input, err)
			continue
		}
		page := emilia.ParserBuilder.BuildParser(
			emilia.FullPathToWorkDirRel(bundle.First),
			string(data),
		).Parse()
		pages = append(pages, page)
	}
	return pages
}
