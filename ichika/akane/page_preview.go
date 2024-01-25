package akane

import (
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/emilia/reze"
	"github.com/thecsw/darkness/ichika/kuroko"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/komi"
	"github.com/thecsw/rei"
)

// pagePreviewRequest is a request to generate a page preview.
type pagePreviewRequest struct {
	Location yunyun.RelativePathDir
	Title    string
	Time     string
}

// pagePreviewsToGenerate is a set of page previews to generate.
var (
	pagePreviewsToGenerate      = sync.Map{}
	pagePreviewsToGenerateCount = atomic.Uint32{}
)

// RequestPagePreview requests a page preview to be generated.
func RequestPagePreview(location yunyun.RelativePathDir, title string, time string) {
	pagePreviewsToGenerate.Store(location, pagePreviewRequest{
		Location: location,
		Title:    title,
		Time:     time,
	})
	pagePreviewsToGenerateCount.Add(1)
}

const (
	pagePreviewTitleFont = "styles/fonts/EB_Garamond/EBGaramond-VariableFont_wght.ttf"
	pagePreviewNameFont  = "styles/fonts/EB_Garamond/static/EBGaramond-Medium.ttf"
	pagePreviewTimeFont  = "styles/fonts/EB_Garamond/static/EBGaramond-Italic.ttf"
	pagePreviewWidth     = puck.PagePreviewWidth
	pagePreviewHeight    = puck.PagePreviewHeight

	pagePreviewFilename = "preview.jpg"
)

// doPagePreviews generates page previews.
func doPagePreviews(conf *alpha.DarknessConfig) {
	// Let's initialize the page preview generator.
	generator := reze.InitPreviewGenerator(
		pagePreviewTitleFont,
		pagePreviewNameFont,
		pagePreviewTimeFont,
		pagePreviewWidth,
		pagePreviewHeight,
		conf.Website.Color,
		string(conf.Author.Image),
	)
	// Let's make sure we close the generator when we're done.
	defer func(generator reze.PreviewGenerator) {
		err := generator.Close()
		if err != nil {
			logger.Error("Closing reze page preview generator", "err", err)
		}
	}(generator)

	waiting := sync.WaitGroup{}
	waiting.Add(int(pagePreviewsToGenerateCount.Load()))
	skipped := atomic.Int32{}

	processPagePreviewRequest := func(pagePreview pagePreviewRequest) {
		start := time.Now()
		// Find the path to save the preview to.
		relativeTarget := yunyun.RelativePathFile(filepath.Join(string(pagePreview.Location), pagePreviewFilename))
		// Skip if exists, unless forced.
		if !kuroko.Force {
			if exists, _ := rei.FileExists(string(relativeTarget)); exists {
				skipped.Add(1)
				waiting.Done()
				return
			}
		}
		// Get the reader for the generated preview.
		reader := rei.Must(generator.Generate(removeNonPrintables(pagePreview.Title, conf.Title, pagePreview.Time)))
		// Get it with the work directory.
		target := conf.Runtime.WorkDir.Join(relativeTarget)
		// Save the preview as a jpg.
		if err := reze.SaveJpg(reader, string(target)); err != nil {
			logger.Error(
				"Saving page preview",
				"loc", target,
				"err", err,
			)
			return
		}
		logger.Info(
			"Generated page preview",
			"loc", conf.Runtime.WorkDir.Rel(target),
			"elapsed", time.Since(start),
		)
		waiting.Done()
	}

	pageGeneratorPool := komi.NewWithSettings(komi.WorkSimple(processPagePreviewRequest), &komi.Settings{
		Name:     "Komi Page Preview 🫦 ",
		Laborers: runtime.NumCPU(),
	})

	pagePreviewsToGenerate.Range(func(key, value any) bool {
		rei.Try(pageGeneratorPool.Submit(value.(pagePreviewRequest)))
		return true
	})

	waiting.Wait()

	// Block until all work is complete.
	pageGeneratorPool.Close()

	// Write a notice if we skipped any preview generations.
	if numSkipped := skipped.Load(); numSkipped > 0 {
		logger.Warn("Some previews already existed, use -force to overwrite", "skipped", numSkipped)
	}
}

func removeNonPrintables(title, name, time string) (string, string, string) {
	return onlyKeepPrint(title), onlyKeepPrint(name), onlyKeepPrint(time)
}

func onlyKeepPrint(k string) string {
	result := ""
	for _, r := range k {
		if unicode.IsLetter(r) || unicode.IsSpace(r) || unicode.IsDigit(r) || unicode.IsPunct(r) {
			result += string(r)
		}
	}
	return result
}
