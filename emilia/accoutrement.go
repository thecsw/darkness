package emilia

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

const (
	enableOption    = `t`
	disableOption   = `nil`
	delimiterOption = ':'

	optionTomb            = `tomb`
	optionAuthorImage     = `author-image`
	optionMath            = `math`
	optionExcludeHtmlHead = `exclude-html-head`
)

var (
	accotrementActions = map[string]func(string, *yunyun.Accoutrement){
		optionTomb:            accoutrementTomb,
		optionAuthorImage:     accoutrementAuthorImage,
		optionMath:            accoutrementMath,
		optionExcludeHtmlHead: accoutrementExcludeHtmlScript,
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

func accoutrementTomb(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.Tomb)
}

func accoutrementAuthorImage(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.AuthorImage)
}

func accoutrementMath(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.Math)
}

func accoutrementExcludeHtmlScript(what string, target *yunyun.Accoutrement) {
	target.ExcludeHtmlHeadContains = append(target.ExcludeHtmlHeadContains, what)
}

func accoutrementBool(what string, target *yunyun.AccoutrementFlip) {
	cleaned := strings.TrimSpace(what)
	if cleaned == enableOption {
		target.Enable()
	}
	if cleaned == disableOption {
		target.Disable()
	}
}
