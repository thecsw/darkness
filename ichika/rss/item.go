package rss

import "encoding/xml"

// A channel may contain any number of `<item>`s. An item may represent a "story" --
// much like a story in a newspaper or magazine; if so its description is a synopsis
// of the story, and the link points to the full story. An item may also be complete
// in itself, if so, the description contains the text (entity-encoded HTML is allowed;
// see examples), and the link and title may be omitted. All elements of an item are
// optional, however at least one of title or description must be present.
type Item struct {
	XMLName xml.Name `xml:"item"`

	// --------------------------
	// Kind of required elements?
	// --------------------------

	// The title of the item.
	//
	// Example: "Venice Film Festival Tries to Quit Sinking"
	Title string `xml:"title"`

	// The URL of the item.
	//
	// Example: "http://nytimes.com/2004/12/07FEST.html"
	Link string `xml:"link"`

	// The item synopsis.
	//
	// Example: "<description>Some of the most heated chatter
	// at the Venice Film Festival this week was about the way
	// that the arrival of the stars at the Palazzo del Cinema
	// was being staged.</description>"
	Description string `xml:"description"`

	// Email address of the author of the item.
	//
	// It's the email address of the author of the item. For newspapers
	// and magazines syndicating via RSS, the author is the person who wrote
	// the article that the <item> describes. For collaborative weblogs, th
	// author of the item might be different from the managing editor or webmaster.
	// For a weblog authored by a single individual it would make sense to omit the
	// <author> element.
	Author string `xml:"author"`

	// It has one optional attribute, domain, a string that identifies
	// a categorization taxonomy.
	//
	// The value of the element is a forward-slash-separated string that
	// identifies a hierarchic location in the indicated taxonomy. Processors
	// may establish conventions for the interpretation of categories. Two
	// examples are provided below:
	//
	//  - `<category>Grateful Dead</category>`
	//  - `<category domain="http://www.fool.com/cusips">MSFT</category>`
	//
	// You may include as many category elements as you need to, for different
	// domains, and to have an item cross-referenced in different parts of the
	// same domain.
	Category *Category `xml:"category"`

	// URL of a page for comments relating to the item. If present, it is the
	// url of the comments page for the item.
	//
	// Example: "<comments>http://ekzemplo.com/entry/4403/comments</comments>"
	//
	// See more: https://www.rssboard.org/rss-weblog-comments-use-case
	Comments string `xml:"comments"`

	// Describes a media object that is attached to the item.
	Enclosure *Enclosure `xml:"enclosure"`

	// A string that uniquely identifies the item.
	//
	// guid stands for globally unique identifier. It's a string that
	// uniquely identifies the item. When present, an aggregator may
	// choose to use this string to determine if an item is new.
	//
	// Example: "<guid>http://some.server.com/weblogItem3207</guid>"
	Guid *Guid `xml:"guid"`

	// Its value is a date, indicating when the item was published. If
	// it's a date in the future, aggregators may choose to not display
	// the item until that date.
	PubDate string `xml:"pubDate"`

	// The RSS channel that the item came from.
	//
	// Its value is the name of the RSS channel that the item came from,
	// derived from its <title>. It has one required attribute, url,
	// which links to the XMLization of the source.
	//
	// The purpose of this element is to propagate credit for links, to
	// publicize the sources of news items. It can be used in the Post
	// command of an aggregator. It should be generated automatically when
	// forwarding an item from an aggregator to a weblog authoring tool.
	Source *Source `xml:"source"`
}
