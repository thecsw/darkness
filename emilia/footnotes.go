package emilia

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/thecsw/darkness/yunyun"
)

// WithFootnotes resolves footnotes and cleans up the page if necessary
func WithFootnotes() yunyun.PageOption {
	return func(page *yunyun.Page) {
		footnotes := make([]string, 0, 4)
		for i := range page.Contents {
			c := &page.Contents[i]
			// Replace footnotes in paragraphs
			if c.IsParagraph() {
				c.Paragraph = findFootnotes(c.Paragraph, &footnotes)
			}
			// Footnotes can also appear in lists
			if c.IsList() {
				for i := 0; i < len(c.List); i++ {
					c.List[i] = findFootnotes(c.List[i], &footnotes)
				}
			}
		}
		page.Footnotes = footnotes
	}
}

// findFootnotes finds footnotes in a paragraph and replaces them with a footnote reference
func findFootnotes(text string, footnotes *[]string) string {
	matches := yunyun.FootnoteRegexp.FindAllStringSubmatch(text, -1)
	// no footnotes found
	if len(matches) < 1 {
		return text
	}
	newText := text
	for _, match := range matches {
		*footnotes = append(*footnotes, match[1])
		newText = strings.Replace(newText, match[0], fmt.Sprintf("!%d!", len(*footnotes)), 1)
	}
	return newText
}

var (
	// FootnoteLabeler will take an integer and return string representation as
	// defined in the darkness config, with either Roman or Arabic numerals.
	FootnoteLabeler = strconv.Itoa

	// numToRomanMap is used by `numberToRoman` to build Roman numerals.
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
	// romanNumOrder is used by `numberToRoman` to build Roman numerals.
	romanNumOrder = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	// numToRomanSmall is used by `numberToRoman` for fast Roman footnote generation.
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
