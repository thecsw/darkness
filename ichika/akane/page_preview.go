package akane

import (
	"github.com/thecsw/darkness/emilia/alpha"
	"github.com/thecsw/darkness/emilia/puck"
	"github.com/thecsw/darkness/emilia/reze"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/rei"
	"path/filepath"
	"unicode"
)

type pagePreview struct {
	Location yunyun.RelativePathDir
	Title    string
	Time     string
}

var (
	pagePreviewsToGenerate = make([]pagePreview, 0, 16)
)

func RequestPagePreview(location yunyun.RelativePathDir, title string, time string) {
	pagePreviewsToGenerate = append(pagePreviewsToGenerate, pagePreview{
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

func doPagePreviews(conf *alpha.DarknessConfig) {
	generator := reze.InitPreviewGenerator(
		pagePreviewTitleFont,
		pagePreviewNameFont,
		pagePreviewTimeFont,
		pagePreviewWidth,
		pagePreviewHeight,
		conf.Website.Color,
		string(conf.Author.Image),
	)
	defer func(generator reze.PreviewGenerator) {
		err := generator.Close()
		if err != nil {
			logger.Error("closing reze page preview generator", "err", err)
		}
	}(generator)

	// Let's start going through the page preview requests.
	for _, pagePreview := range pagePreviewsToGenerate {
		reader := rei.Must(generator.Generate(removeNonPrintables(pagePreview.Title, conf.Title, pagePreview.Time)))
		relativeTarget := yunyun.RelativePathFile(filepath.Join(string(pagePreview.Location), pagePreviewFilename))
		target := conf.Runtime.WorkDir.Join(relativeTarget)
		if err := reze.SaveJpg(reader, string(target)); err != nil {
			logger.Error("saving jpg preview", "loc", target, "err", err)
			continue
		}
		logger.Info("generated preview", "loc", target)
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
