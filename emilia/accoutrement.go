package emilia

import (
	"strings"

	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

const (
	enableOption  = `t`
	disableOption = `nil`

	optionTomb        = `tomb`
	optionAuthorImage = `author-image`
)

var (
	accotrementActions = map[string]func(string, *yunyun.Accoutrement){
		optionTomb:        accoutrementTomb,
		optionAuthorImage: accoutrementAuthorImage,
	}
)

// InitializeAccoutrement fills accoutrement according to the config
// and default values.
func InitializeAccoutrement(page *yunyun.Page) {
	// Better to use a trie for matching multiple prefixes.
	for _, tombPage := range Config.Website.Tombs {
		if strings.HasPrefix(string(page.Location), string(tombPage)) {
			accoutrementTomb(enableOption, page.Accoutrement)
			break
		}
	}

	// Always enable author image by default.
	page.Accoutrement.AuthorImage = true
}

// FillAccoutrement parses `options` and fills the `target`.
func FillAccoutrement(options *string, target *yunyun.Accoutrement) {
	// Exit immediately if it's an empty string.
	if len(*options) < 1 {
		return
	}
	for _, option := range strings.Split(*options, " ") {
		elements := strings.SplitN(option, ":", 2)
		key := gana.First(elements)
		value := gana.Last(elements)
		// If the value was not passed, assume it's true ("t").
		if len(elements) < 2 {
			value = enableOption
		}
		// If action is found, then execute it.
		if action, ok := accotrementActions[key]; ok {
			action(value, target)
		}
	}
}

func accoutrementTomb(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.Tomb)
}

func accoutrementAuthorImage(what string, target *yunyun.Accoutrement) {
	accoutrementBool(what, &target.AuthorImage)
}

func accoutrementBool(what string, target *bool) {
	switch strings.TrimSpace(what) {
	case enableOption:
		*target = true
	case disableOption:
		*target = false
	}
}
