package rss

import "encoding/xml"

type Category struct {
	XMLName xml.Name `xml:"category"`

	Value  string `xml:",chardata"`
	Domain string `xml:"domain,attr"`
}
