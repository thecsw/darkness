package misa

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/puck"
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
	initLog()
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
	allPagesRelative := gana.Map(
		func(p *yunyun.Page) yunyun.RelativePathDir { return p.Location },
		hizuru.BuildPagesSimple(conf, nil))
	logger.Infof("starting to track recent changes for %d URLs", len(allPagesRelative))

	// Let's filter the pages only if their most recent modified
	// date is after the published online (if found).
	lastBuilt, err := getLastBuilt(conf)
	if err != nil {
		logger.Warn("getting last built, sending all pages", "err", err)
	}
	logger.Info("Found the previous remote build", "last_built", lastBuilt.Local().Format(time.RFC850))

	if err == nil {
		filteredPages := make([]yunyun.RelativePathDir, 0, len(allPagesRelative))
		indexPath := "/index" + yunyun.FullPathFile(conf.Project.Input)
		for _, allPageRelative := range allPagesRelative {
			fsPath := conf.Runtime.WorkDir.Join(yunyun.RelativePathFile(allPageRelative)) + indexPath
			lastModTime, err := alpha.ExtractGitLastModified(conf, yunyun.RelativePathFile(fsPath))
			if err != nil {
				logger.Warnf("couldn't get git's last modified time for %s: %v", fsPath, err)
			}
			logger.Debug("Extracted git last modified",
				"path", conf.Runtime.WorkDir.Rel(fsPath), "git_last_mod", lastModTime.Local().Format(time.RFC850))
			if lastModTime.After(*lastBuilt) {
				logger.Info("Found page modified since the last build",
					"path", conf.Runtime.WorkDir.Rel(fsPath), "git_last_mod", lastModTime.Local().Format(time.RFC850))
				filteredPages = append(filteredPages, allPageRelative)
			}
		}
		allPagesRelative = filteredPages
	}

	if len(allPagesRelative) == 0 {
		logger.Warn("There are no new updates to recrawl or reindex, bailing out")
	} else {
		logger.Infof("Found %d pages that need a recrawl and reindex", len(allPagesRelative))
	}

	// Notify search engines.
	if !dryRun {
		for _, searchEngine := range conf.External.SearchEngines {
			if err := notifySearchEngineMultiple(conf, searchEngine, indexNowKeyContentsString, conf.Url, allPagesRelative); err != nil {
				logger.Warnf("failed to notify search engine %s: %s", searchEngine, err)
				continue
			}
			logger.Infof("notified search engine %s", searchEngine)
		}
	} else {
		logger.Warn("Dryrun mode active, I am not notifying search engines")
	}
}

func getLastBuilt(conf *alpha.DarknessConfig) (*time.Time, error) {
	remotePath := conf.Runtime.UrlPath.JoinPath(puck.LastBuildTimestampFile).String()
	lastBuilt, err := haruhi.URL(remotePath).ResponseString()
	if err != nil {
		return nil, fmt.Errorf("collecting last_built.txt from %s: %v", remotePath, err)
	}
	lastBuiltTime, err := time.Parse(time.RFC3339, lastBuilt)
	if err != nil {
		return nil, fmt.Errorf("parsing last_built.txt from %s: %v", remotePath, err)
	}
	return &lastBuiltTime, nil
}

type indexNowRequestMultipleUrls struct {
	Host    string   `json:"host"`
	Key     string   `json:"key"`
	UrlList []string `json:"urlList"`
}

func notifySearchEngineMultiple(
	conf *alpha.DarknessConfig,
	searchEngineUrl string,
	indexNowKey string,
	host string,
	relPathsToUpdate []yunyun.RelativePathDir,
) error {
	path := fmt.Sprintf("https://%s", searchEngineUrl)

	// convert them to URLs
	urlsToUpdate := gana.Map(func(rel yunyun.RelativePathDir) string {
		return string(conf.Runtime.JoinDir(rel))
	}, relPathsToUpdate)

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
