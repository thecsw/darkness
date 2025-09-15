package yunyun

import "strings"

// quotesReplace is the map to replace
var quotesReplace = map[string]string{
	"'s":  "’s",
	"'d":  "’d",
	"s'":  "s’",
	"'m":  "’m",
	"n't": "n’t",
	"'re": "’re",
	"'ve": "’ve",
	"'ll": "’ll",
	"L'":  "L’",
	"``":  "“",
	"''":  "”",
}

// FancyText replaces boring single and double quotes with fancier Unicode versions
func FancyText(text string) string {
	for k, v := range quotesReplace {
		text = strings.ReplaceAll(text, k, v)
	}
	text = strings.ReplaceAll(text, "---", "—") // em-dash
	text = strings.ReplaceAll(text, "--", "–")  // en dash
	return text
}
