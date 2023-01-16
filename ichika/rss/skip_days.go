package rss

import "encoding/xml"

// A hint for aggregators telling them which days they can skip.
// This element contains up to seven <day> sub-elements whose value
// is Monday, Tuesday, Wednesday, Thursday, Friday, Saturday or Sunday.
// Aggregators may not read the channel during days listed in the <skipDays>
// element.
type SkipDays struct {
	XMLName xml.Name `xml:"skipDays"`

	// Monday, Tuesday, Wednesday, Thursday, Friday, Saturday or Sunday.
	Day []string `xml:"day"`
}
