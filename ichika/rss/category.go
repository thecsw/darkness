package rss

import "encoding/xml"

// Category It has one optional attribute, domain, a string that identifies
// a categorization taxonomy.
//
// The value of the element is a forward-slash-separated string that
// identifies a hierarchic location in the indicated taxonomy. Processors
// may establish conventions for the interpretation of categories. Two
// examples are provided below:
//
//   - `<category>Grateful Dead</category>`
//   - `<category domain="http://www.fool.com/cusips">MSFT</category>`
//
// You may include as many category elements as you need to, for different
// domains, and to have an item cross-referenced in different parts of the
// same domain.
type Category struct {
	XMLName xml.Name `xml:"category"`

	// Name of the category.
	Value string `xml:",chardata"`

	// URL to the category.
	Domain string `xml:"domain,attr"`
}
