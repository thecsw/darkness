package akane

import (
	"path/filepath"
	"runtime"
	"sync"
	"unicode"

	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/emilia/reze"
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

// pagePreviewsToGenerate is a list of page previews to generate.
var pagePreviewsToGenerate = make([]pagePreviewRequest, 0, 16)

// RequestPagePreview requests a page preview to be generated.
func RequestPagePreview(location yunyun.RelativePathDir, title string, time string) {
	pagePreviewsToGenerate = append(pagePreviewsToGenerate, pagePreviewRequest{
		Location: location,
		Title:    title,
		Time:     time,
	})
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
	// Clear the pagePreviewsToGenerate slice when we're done.
	defer func() {
		pagePreviewsToGenerate = pagePreviewsToGenerate[:0]
	}()

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
	waiting.Add(len(pagePreviewsToGenerate))

	processPagePreviewRequest := func(pagePreview pagePreviewRequest) {
		// Get the reader for the generated preview.
		reader := rei.Must(generator.Generate(removeNonPrintables(pagePreview.Title, conf.Title, pagePreview.Time)))
		// Find the path to save the preview to.
		relativeTarget := yunyun.RelativePathFile(filepath.Join(string(pagePreview.Location), pagePreviewFilename))
		// Get it with the work directory.
		target := conf.Runtime.WorkDir.Join(relativeTarget)
		// Save the preview as a jpg.
		if err := reze.SaveJpg(reader, string(target)); err != nil {
			logger.Error("Saving page preview", "loc", target, "err", err)
			return
		}
		logger.Info("Generated page preview", "loc", conf.Runtime.WorkDir.Rel(target))
		waiting.Done()
	}

	pageGeneratorPool := komi.NewWithSettings(komi.WorkSimple(processPagePreviewRequest), &komi.Settings{
		Name:     "Komi Page Preview 🫦 ",
		Laborers: runtime.NumCPU(),
	})

	// Let's start going through the page preview requests.
	for _, pagePreview := range pagePreviewsToGenerate {
		pageGeneratorPool.Submit(pagePreview)
	}

	waiting.Wait()

	// Block until all work is complete.
	pageGeneratorPool.Close()
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
