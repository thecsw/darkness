package misa

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/ichika/hizuru"
	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
	"github.com/thecsw/haruhi"
	"github.com/thecsw/rei"
)

const (
	// indexNowKeyPattern should be an md53
	indexNowKeyPattern = "^[a-fA-F0-9]{32}[.]txt$"
)

var (
	indexNowKeyRegex = regexp.MustCompile(indexNowKeyPattern)
)

// NotifySearchEngines notifies search engines of the updated URLs through indexnow.org.
func NotifySearchEngines(conf *alpha.DarknessConfig, indexNowKey yunyun.RelativePathFile, dryRun bool) {
	if !indexNowKeyRegex.MatchString(string(indexNowKey)) {
		logger.Fatalf("indexnow file text should match the pattern %s", indexNowKeyPattern)
	}

	indexNowKeyContents := rei.Must(os.ReadFile(filepath.Clean(string(indexNowKey))))
	indexNowKeyContentsString := strings.TrimSpace(string(indexNowKeyContents))

	// We need to verify that the file's name matches the contents, strip the file extension.
	indexNowKeyFilename := strings.TrimSuffix(string(indexNowKey), ".txt")
	if indexNowKeyFilename != indexNowKeyContentsString {
		logger.Fatalf("indexnow file name %s should match its contents %s",
			indexNowKeyFilename, indexNowKeyContentsString)
	}

	// Let's get all the URLs to update.
	allPages := gana.Map(
		func(p *yunyun.Page) string { return string(conf.Runtime.JoinDir(p.Location)) },
		hizuru.BuildPagesSimple(conf, nil))
	logger.Infof("notifying search engines of %d URLs", len(allPages))

	// Notify search engines.
	if !dryRun {
		for _, searchEngine := range conf.External.SearchEngines {
			if err := notifySearchEngineMultiple(searchEngine, indexNowKeyContentsString, conf.Url, allPages); err != nil {
				logger.Warnf("failed to notify search engine %s: %s", searchEngine, err)
				continue
			}
			logger.Infof("notified search engine %s", searchEngine)
		}
	}
}

type indexNowRequestMultipleUrls struct {
	Host    string   `json:"host"`
	Key     string   `json:"key"`
	UrlList []string `json:"urlList"`
}

func notifySearchEngineMultiple(
	searchEngineUrl string,
	indexNowKey string,
	host string,
	urlsToUpdate []string,
) error {
	path := fmt.Sprintf("https://%s", searchEngineUrl)

	// From, https://www.indexnow.org/documentation
	// > You can submit up to 10,000 URLs per post, mixing
	// > http and https URLs if needed.
	payload := indexNowRequestMultipleUrls{
		Host:    host,
		Key:     indexNowKey,
		UrlList: urlsToUpdate,
	}
	resp, cancel, err := haruhi.URL(path).
		Path("/indexnow").
		Method(http.MethodPost).
		Header("Content-Type", "application/json; charset=utf-8").
		BodyJson(payload).
		Response()
	defer cancel()
	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("failed to notify search engine %s: %w", searchEngineUrl, err)
	}

	if resp.StatusCode != 200 {
		errorBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to notify search engine %s: %s with error: %s", searchEngineUrl, resp.Status, errorBody)
	}
	return nil
}
