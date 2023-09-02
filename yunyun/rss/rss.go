package rss

import (
	"encoding/xml"
	"time"
)

const (
	// RSSVersion spec version implemented.
	RSSVersion = "2.0"

	// RSSFormat date format used in RSS spec.
	RSSFormat = time.RFC1123

	// RSSDocs RSS spec implemented.
	RSSDocs = "https://www.rssboard.org/rss-specification"
)

// RSS document is a <rss> element, with a
// mandatory attribute called version, that specifies the version of
// RSS that the document conforms to. If it conforms to this specification,
// the version attribute must be 2.0.
//
// Subordinate to the <rss> element is a single <channel> element, which
// contains information about the channel (metadata) and its contents.
type RSS struct {
	// Subordinate to the <rss> element is a single <channel> element, which
	// contains information about the channel (metadata) and its contents.
	Channel *Channel `xml:"channel"`
	XMLName xml.Name `xml:"rss"`

	// Must be "2.0"
	Version string `xml:"version,attr"`
}
