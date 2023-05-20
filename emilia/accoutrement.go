package emilia

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

const (
	enableOption    = `t`
	disableOption   = `nil`
	delimiterOption = ':'

	optionDraft           = `draft`
	optionTomb            = `tomb`
	optionAuthorImage     = `author-image`
	optionMath            = `math`
	optionExcludeHtmlHead = `exclude-html-head`
	optionPreview         = `preview`
)

var (
	accotrementActions = map[string]func(string, *yunyun.Accoutrement){
		optionDraft:           accoutrementDraft,
		optionTomb:            accoutrementTomb,
		optionAuthorImage:     accoutrementAuthorImage,
		optionMath:            accoutrementMath,
		optionExcludeHtmlHead: accoutrementExcludeHtmlScript,
		optionPreview:         accoutrementPreview,
	}
)

// InitializeAccoutrement fills accoutrement according to the config
// and default values.
func InitializeAccoutrement(page *yunyun.Page) {
	// Better to use a trie for matching multiple prefixes.
	for _, tombPage := range Config.Website.Tombs {
		if strings.HasPrefix(string(page.Location), string(tombPage)) {
			page.Accoutrement.Tomb.Enable()
			// Just one prefix is enough to deduce tombs.
			break
		}
	}
}

// FillAccoutrement parses `options` and fills the `target`.
func FillAccoutrement(options *string, page *yunyun.Page) {
	// Let's first initialize it before filling.
	InitializeAccoutrement(page)
	// Exit immediately if it's an empty string.
	if len(*options) < 1 {
		return
	}
	for _, option := range strings.Split(*options, " ") {
		key, value := breakOption(option)
		// If action is found, then execute it.
		if action, ok := accotrementActions[key]; ok {
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

// accoutrementBool sets the bool value of the target according to the what.
func accoutrementBool(what string, target *yunyun.AccoutrementFlip) {
	switch strings.TrimSpace(what) {
	case enableOption:
		target.Enable()
	case disableOption:
		target.Disable()
	}
}
