package emilia

import (
	"strings"

	"github.com/thecsw/darkness/v3/yunyun"
)

const (
	enableOption    = `t`
	disableOption   = `nil`
	delimiterOption = ':'

	optionDate              = `date`
	optionDraft             = `draft`
	optionTomb              = `tomb`
	optionAuthorImage       = `author-image`
	optionMath              = `math`
	optionExcludeHtmlHead   = `exclude-html-head`
	optionPreview           = `preview`
	optionPreviewWidth      = `preview-width`
	optionPreviewHeigh      = `preview-height`
	optionPreviewGenerate   = `preview-generate`
	optionPreviewGenerateBg = `preview-generate-bg`
	optionPreviewGenerateFg = `preview-generate-fg`
	optionToc               = `toc`
	optionRssPrefix         = `rss-prefix`
	optionRssTitle          = `rss-title`
)

var accoutrementActions = map[string]func(string, *yunyun.Accoutrement){
	optionDate:              accoutrementDate,
	optionDraft:             accoutrementDraft,
	optionTomb:              accoutrementTomb,
	optionAuthorImage:       accoutrementAuthorImage,
	optionMath:              accoutrementMath,
	optionExcludeHtmlHead:   accoutrementExcludeHtmlScript,
	optionPreview:           accoutrementPreview,
	optionPreviewWidth:      accoutrementPreviewWidth,
	optionPreviewHeigh:      accoutrementPreviewHeight,
	optionPreviewGenerate:   accoutrementPreviewGenerate,
	optionPreviewGenerateBg: accoutrementPreviewGenerateBg,
	optionPreviewGenerateFg: accoutrementPreviewGenerateFg,
	optionToc:               accoutrementToc,
	optionRssPrefix:         accoutrementRssPrefix,
	optionRssTitle:          accoutrementRssTitle,
}

// InitializeAccoutrement fills accoutrement according to the config
// and default values.
func InitializeAccoutrement(tombs []yunyun.RelativePathDir, page *yunyun.Page) {
	// Better to use a trie for matching multiple prefixes.
	for _, tombPage := range tombs {
		if strings.HasPrefix(string(page.Location), string(tombPage)) {
			page.Accoutrement.Tomb.Enable()
			// Just one prefix is enough to deduce tombs.
			break
		}
	}
}

// FillAccoutrement parses `options` and fills the `target`.
func FillAccoutrement(tombs []yunyun.RelativePathDir, options *string, page *yunyun.Page) {
	// Let's first initialize it before filling.
	InitializeAccoutrement(tombs, page)
	// Exit immediately if it's an empty string.
	if len(*options) < 1 {
		return
	}
	for _, option := range strings.Split(*options, " ") {
		key, value := breakOption(option)
		// If action is found, then execute it.
		if action, ok := accoutrementActions[key]; ok {
			action(value, page.Accoutrement)
		}
	}
}

// breakOption breaks the option into two parts, the first part is the
// key, and the second part is the value. If the option doesn't have
// a value, then the second part is `enableOption` by default.
func breakOption(what string) (string, string) {
	for i := len(what) - 1; i >= 0; i-- {
		if what[i] == delimiterOption {
			return what[:i], what[i+1:]
		}
	}
	// By default return the whole string as the first one,
	// and enable option to the right.
	return what, enableOption
}

// accoutrementDraft sets the date inclusion in the exported document option.
func accoutrementDate(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.Date)
}

// accoutrementDraft sets the draft option of the accoutrement.
func accoutrementDraft(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.Draft)
}

// accoutrementTomb sets the tomb option of the accoutrement.
func accoutrementTomb(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.Tomb)
}

// accoutrementAuthorImage sets the author image option of the accoutrement.
func accoutrementAuthorImage(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.AuthorImage)
}

// accoutrementMath sets the math option of the accoutrement.
func accoutrementMath(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.Math)
}

// accoutrementExcludeHtmlScript sets the exclude html script option of the accoutrement.
func accoutrementExcludeHtmlScript(what string, target *yunyun.Accoutrement) {
	target.ExcludeHtmlHeadContains = append(target.ExcludeHtmlHeadContains, what)
}

// accoutrementPreview sets the preview option of the accoutrement.
func accoutrementPreview(what string, target *yunyun.Accoutrement) {
	target.Preview = what
}

// accoutrementPreviewWidth sets the preview width option of the accoutrement.
func accoutrementPreviewWidth(what string, target *yunyun.Accoutrement) {
	target.PreviewWidth = what
}

// accoutrementPreviewHeight sets the preview height option of the accoutrement.
func accoutrementPreviewHeight(what string, target *yunyun.Accoutrement) {
	target.PreviewHeight = what
}

// accoutrementPreviewGenerate sets the preview generation option of the accoutrement.
func accoutrementPreviewGenerate(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.PreviewGenerate)
}

// accoutrementPreviewGenerateBg sets the preview generation background option of the accoutrement.
func accoutrementPreviewGenerateBg(what string, target *yunyun.Accoutrement) {
	target.PreviewGenerateBg = what
}

// accoutrementPreviewGenerateFg sets the preview generation foreground option of the accoutrement.
func accoutrementPreviewGenerateFg(what string, target *yunyun.Accoutrement) {
	target.PreviewGenerateFg = what
}

// accoutrementToc sets the toc option of the accoutrement.
func accoutrementToc(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.Toc)
}

// accoutrementRssPrefix sets the rss prefix option of the accoutrement.
func accoutrementRssPrefix(what string, target *yunyun.Accoutrement) {
	target.RssPrefix = what
}

// accoutrementRssTitle sets the rss title option of the accoutrement.
func accoutrementRssTitle(what string, target *yunyun.Accoutrement) {
	target.RssTitle = what
}

// accoutrementBool sets the bool value of the target according to the what.
func accoutrementBool(what string, target *yunyun.AccoutrementFlip) {
	switch strings.TrimSpace(what) {
	case enableOption:
		target.Enable()
	case disableOption:
		target.Disable()
	}
}
