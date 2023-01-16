package rss

import "encoding/xml"

// A hint for aggregators telling them which hours they can skip.
// This element contains up to 24 <hour> sub-elements whose value is
// a number between 0 and 23, representing a time in GMT, when aggregators,
// if they support the feature, may not read the channel on hours listed
// in the <skipHours> element. The hour beginning at midnight is hour zero.
type SkipHours struct {
	XMLName xml.Name `xml:"skipHours"`

	// A number between 0 and 23, representing a time in GMT.
	Hour []int `xml:"hour"`
}
