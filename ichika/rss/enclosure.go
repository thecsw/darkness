package rss

import "encoding/xml"

// Describes a media object that is attached to the item.
//
// It has three required attributes. url says where the enclosure
// is located, length says how big it is in bytes, and type says
// what its type is, a standard MIME type.
//
// See more: https://www.rssboard.org/rss-enclosures-use-case
type Enclosure struct {
	XMLName xml.Name `xml:"enclosure"`

	// The url must be an http url.
	URL string `xml:"url,attr"`

	// Content type of the enclosure.
	Type string `xml:"type,attr"`

	// Length of the image in bytes.
	Length int `xml:"length,attr"`
}
