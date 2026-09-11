package misa

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/thecsw/darkness/v3/emilia/alpha"
	"github.com/thecsw/darkness/v3/emilia/memento"
	"github.com/thecsw/darkness/v3/ichika/hizuru"
	"github.com/thecsw/darkness/v3/ichika/kuroko"
)

type embedRefreshJob struct {
	request      memento.Request
	pageLocation string
	reference    string
}

func (job embedRefreshJob) identity() string {
	return job.reference
}

type embedRefreshResult struct {
	job    embedRefreshJob
	status string
	err    error
}

// refreshVendoredEmbeds refreshes explicitly referenced, page-local TOML
// manifests as the first phase of the unified vendoring workflow.
func refreshVendoredEmbeds(conf *alpha.DarknessConfig, dryRun bool) {
	initLog()
	seen := make(map[string]embedRefreshJob)
	invalidReferences := 0
	for _, page := range hizuru.BuildPagesSimple(conf, nil) {
		for _, content := range page.Contents {
			if !content.IsLink() {
				continue
			}
			if !strings.HasPrefix(content.Link, "embed:") {
				continue
			}
			reference := strings.TrimPrefix(content.Link, "embed:")
			draft, err := memento.ReadReference(conf, string(page.Location), reference)
			if err != nil {
				invalidReferences++
				logger.Warn("Could not read embed manifest", "manifest", reference, "err", err)
				continue
			}
			request, ok := memento.ParseURL(draft.URL)
			if !ok {
				invalidReferences++
				logger.Warn("Embed manifest has no supported url", "manifest", reference)
				continue
			}
			manifestPath, err := memento.ResolveReferencePath(conf, string(page.Location), reference)
			if err != nil {
				invalidReferences++
				logger.Warn("Could not resolve embed manifest", "manifest", reference, "err", err)
				continue
			}
			seen[manifestPath] = embedRefreshJob{
				request: request, pageLocation: string(page.Location), reference: reference,
			}
		}
	}

	jobsList := make([]embedRefreshJob, 0, len(seen))
	for _, job := range seen {
		jobsList = append(jobsList, job)
	}
	sort.Slice(jobsList, func(i, j int) bool {
		return jobsList[i].identity() < jobsList[j].identity()
	})
	if len(jobsList) == 0 {
		logger.Info("No vendored embed references found", "invalid", invalidReferences)
		return
	}

	workerCount := kuroko.CustomNumWorkers
	if workerCount < 1 {
		workerCount = 1
	}
	jobs := make(chan embedRefreshJob)
	results := make(chan embedRefreshResult)
	importer := memento.NewImporter()
	ctx := context.Background()

	var workers sync.WaitGroup
	for range workerCount {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for job := range jobs {
				status, err := importer.RefreshReference(ctx, conf, job.request,
					job.pageLocation, job.reference, kuroko.Force, dryRun)
				results <- embedRefreshResult{job: job, status: status, err: err}
			}
		}()
	}
	go func() {
		for _, job := range jobsList {
			jobs <- job
		}
		close(jobs)
		workers.Wait()
		close(results)
	}()

	failures := invalidReferences
	for result := range results {
		if result.err != nil {
			failures++
			logger.Warn("Could not refresh embed snapshot", "embed", result.job.identity(), "err", result.err)
			continue
		}
		logger.Info("Embed snapshot", "embed", result.job.identity(), "status", result.status)
	}
	logger.Info("Finished refreshing vendored embeds", "found", len(jobsList), "failures", failures)
}
