package rss

import "encoding/xml"

// Source is the RSS channel that the item came from.
//
// Its value is the name of the RSS channel that the item came from,
// derived from its <title>. It has one required attribute, url,
// which links to the XMLization of the source.
//
// The purpose of this element is to propagate credit for links, to
// publicize the sources of news items. It can be used in the Post
// command of an aggregator. It should be generated automatically when
// forwarding an item from an aggregator to a weblog authoring tool.
type Source struct {
	XMLName xml.Name `xml:"source"`

	// Title of the source.
	Value string `xml:",chardata"`

	// Url of the source.
	Url string `xml:"url,attr"`
}
