package html

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/internals"
)

// fancyQuotes replaces boring single and double quotes with fancier Unicode versions
func fancyQuotes(text string) string {
	text = strings.ReplaceAll(text, "'s", "’s")
	text = strings.ReplaceAll(text, "'m", "’m")
	text = strings.ReplaceAll(text, "n't", "n’t")
	text = strings.ReplaceAll(text, "'re", "’re")
	text = strings.ReplaceAll(text, "'ll", "’ll")
	//text = strings.ReplaceAll(text, "`", "‘")
	text = strings.ReplaceAll(text, "``", "“")
	text = strings.ReplaceAll(text, "''", "”")
	text = strings.ReplaceAll(text, "--", "—")
	return text
}

// markupHTMLMapping maps the regex markup to html replacements
var markupHTMLMapping = map[*regexp.Regexp]string{
	internals.ItalicText:        `$1<em>$2</em>$3`,
	internals.BoldText:          `$1<strong>$2</strong>$3`,
	internals.VerbatimText:      `$1<code>$2</code>$3`,
	internals.StrikethroughText: `$1<s>$2</s>$3`,
	internals.UnderlineText:     `$1<u>$2</u>$3`,
}

// markupHTML replaces the markup regexes defined in internal with HTML tags
func markupHTML(text string) string {
	// To make bold italics, it has to be wrapped in /*...*/
	// instead of */.../*
	text = internals.BoldItalicTextBegin.ReplaceAllString(text, `$1/*`)
	text = internals.BoldItalicTextEnd.ReplaceAllString(text, `*/$1`)
	for source, replacement := range markupHTMLMapping {
		text = source.ReplaceAllString(text, replacement)
	}
	text = internals.KeyboardRegexp.ReplaceAllString(text, `<kbd>$1</kbd>`)
	text = internals.NewLineRegexp.ReplaceAllString(text, `$1<br>`)
	return text
}

// processText returns a properly formatted HTML of a text
func processText(text string) string {
	text = html.EscapeString(fancyQuotes(text))
	text = markupHTML(text)
	text = strings.ReplaceAll(text, "◼", `<b style="color:#ba3925">◼︎</b>`)
	text = internals.LinkRegexp.ReplaceAllString(text, `<a href="$1">$2</a>`)

	//text = internals.MathRegexp.ReplaceAllString(text, `\($1\)`)

	text = internals.FootnotePostProcessingRegexp.ReplaceAllStringFunc(text, func(what string) string {
		num, _ := strconv.Atoi(strings.ReplaceAll(what, "!", ""))
		// get the footnote HTML body
		footnote := fmt.Sprintf(
			`<a id="_footnoteref_%d" class="footnote" href="#_footnotedef_%d" title="View footnote.">%s</a>`,
			num, num, footnoteLabel(num))
		// Decide if we need to wrap the footnote in square brackets
		if emilia.Config.Website.FootnoteBrackets {
			footnote = "[" + footnote + "]"
		}
		return `
<sup class="footnote">` + footnote + `</sup>
`
	})

	return strings.TrimSpace(text)
}

// processTitle returns a properly formatted HTML of a title
func processTitle(title string) string {
	title = fancyQuotes(title)
	title = markupHTML(title)
	title = internals.MathRegexp.ReplaceAllString(title, `\($1\)`)
	return title
}

// flattenFormatting returns a plain-text to be fit into the description
func flattenFormatting(what string) string {
	what = fancyQuotes(what)
	// To make bold italics, it has to be wrapped in /*...*/
	// instead of */.../*
	what = internals.BoldItalicTextBegin.ReplaceAllString(what, `$1/*`)
	what = internals.BoldItalicTextEnd.ReplaceAllString(what, `*/$1`)
	for source := range markupHTMLMapping {
		what = source.ReplaceAllString(what, `$1$2$3`)
	}
	what = internals.KeyboardRegexp.ReplaceAllString(what, `$1`)
	what = internals.NewLineRegexp.ReplaceAllString(what, `$1`)
	return what
}

// extractID returns a properly formatted ID for a heading title
func extractID(heading string) string {
	// Check if heading is a link
	match := internals.LinkRegexp.FindStringSubmatch(heading)
	if len(match) > 0 {
		heading = match[2] // 0 is whole match, 1 is link, 2 is title
	}
	res := "_"
	for _, c := range heading {
		if unicode.IsSpace(c) || unicode.IsPunct(c) || unicode.IsSymbol(c) {
			res += "_"
			continue
		}
		if c <= unicode.MaxASCII {
			res += string(unicode.ToLower(c))
		}
	}
	return res
}

var (
	numToRomanMap = map[int]string{
		1:    "I",
		4:    "IV",
		5:    "V",
		9:    "IX",
		10:   "X",
		40:   "XL",
		50:   "L",
		90:   "XC",
		100:  "C",
		400:  "CD",
		500:  "D",
		900:  "CM",
		1000: "M",
	}
	romanNumOrder   = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	numToRomanSmall = map[int]string{
		1:  "I",
		2:  "II",
		3:  "III",
		4:  "IV",
		5:  "V",
		6:  "VI",
		7:  "VII",
		8:  "VIII",
		9:  "IX",
		10: "X",
		11: "XI",
		12: "XII",
		13: "XIII",
		14: "XIV",
	}
)

// numberToRoman converts an integer to a roman numeral
// adapted from https://pencilprogrammer.com/python-programs/convert-integer-to-roman-numerals/
func numberToRoman(num int) string {
	if num <= 14 {
		return numToRomanSmall[num]
	}
	res := ""
	for _, v := range romanNumOrder {
		if num != 0 {
			quot := num / v
			if quot != 0 {
				for x := 0; x < quot; x++ {
					res += numToRomanMap[v]
				}
			}
			num %= v
		}
	}
	return res
}

func footnoteLabel(num int) string {
	if emilia.Config.Website.RomanFootnotes {
		return numberToRoman(num)
	}
	return strconv.Itoa(num)
}
