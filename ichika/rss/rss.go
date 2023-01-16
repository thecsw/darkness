package rss

import (
	"encoding/xml"
	"time"
)

const (
	// RSS spec version implemented.
	RSSVersion = "2.0"

	// The date format used in RSS spec.
	RSSFormat = time.RFC822

	// The RSS spec implemented.
	RSSDocs = "https://www.rssboard.org/rss-specification"
)

// At the top level, a RSS document is a <rss> element, with a
// mandatory attribute called version, that specifies the version of
// RSS that the document conforms to. If it conforms to this specification,
// the version attribute must be 2.0.
//
// Subordinate to the <rss> element is a single <channel> element, which
// contains information about the channel (metadata) and its contents.
type RSS struct {
	XMLName xml.Name `xml:"rss"`

	// Must be "2.0"
	Version string `xml:"version,attr"`

	// Subordinate to the <rss> element is a single <channel> element, which
	// contains information about the channel (metadata) and its contents.
	Channel *Channel `xml:"channel"`
}
