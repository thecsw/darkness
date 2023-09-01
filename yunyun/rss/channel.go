package rss

import "encoding/xml"

// Channel
// Here's a list of the required channel elements, each with a brief
// description, an example, and where available, a pointer to a more
// complete description.
type Channel struct {
	// Allows processes to register with a cloud to be notified of updates
	// to the channel, implementing a lightweight publish-subscribe protocol
	// for RSS feeds.
	//
	// Example: "<cloud domain="rpc.sys.com" port="80" path="/RPC2" registerProcedure="pingMe" protocol="soap"/>"
	//
	// See more: https://www.rssboard.org/rsscloud-interface
	Cloud *Cloud `xml:"cloud,omitempty"`

	// A hint for aggregators telling them which days they can skip.
	// This element contains up to seven <day> sub-elements whose value
	// is Monday, Tuesday, Wednesday, Thursday, Friday, Saturday or Sunday.
	// Aggregators may not read the channel during days listed in the <skipDays>
	// element.
	SkipDays *SkipDays `xml:"skipDays,omitempty"`

	// A hint for aggregators telling them which hours they can skip.
	// This element contains up to 24 <hour> sub-elements whose value is
	// a number between 0 and 23, representing a time in GMT, when aggregators,
	// if they support the feature, may not read the channel on hours listed
	// in the <skipHours> element. The hour beginning at midnight is hour zero.
	SkipHours *SkipHours `xml:"skipHours,omitempty"`

	// Specifies a text input box that can be displayed with the channel.
	TextInput *TextInput `xml:"textInput,omitempty"`

	// Specifies a GIF, JPEG or PNG image that can be displayed with the channel.
	Image *Image `xml:"image,omitempty"`

	XMLName xml.Name `xml:"channel"`

	// Copyright notice for content in the channel.
	//
	// Example: "Copyright 2002, Spartanburg Herald-Journal"
	Copyright string `xml:"copyright,omitempty"`

	// Email address for person responsible for editorial content.
	//
	// Example: "geo@herald.com (George Matesky)"
	ManagingEditor string `xml:"managingEditor,omitempty"`

	// The publication date for the content in the channel. For example,
	// the New York Times publishes on a daily basis, the publication date
	// flips once every 24 hours. That's when the pubDate of the channel
	// changes. All date-times in RSS conform to the Date and Time Specification of
	// RFC 822, with the exception that the year may be expressed with two characters
	// or four characters (four preferred).
	//
	// Example: "Sat, 07 Sep 2002 00:00:01 GMT"
	PubDate string `xml:"pubDate,omitempty"`

	// The last time the content of the channel changed.
	//
	// Example: "Sat, 07 Sep 2002 09:42:31 GMT"
	LastBuildDate string `xml:"lastBuildDate,omitempty"`

	// Specify one or more categories that the channel belongs to.
	// Follows the same rules as the <item>-level category element.
	//
	// Example: "<category>Newspapers</category>"
	Category string `xml:"category,omitempty"`

	// A string indicating the program used to generate the channel.
	//
	// Example: "MightyInHouse Content System v2.3"
	Generator string `xml:"generator,omitempty"`

	// A Url that points to the documentation for the format used in the
	// RSS file. It's probably a pointer to this page. It's for people who
	// might stumble across an RSS file on a Web server 25 years from now
	// and wonder what it is.
	//
	// Example: "https://www.rssboard.org/rss-specification"
	Docs string `xml:"docs,omitempty"`

	// Email address for person responsible for technical issues
	// relating to channel.
	//
	// Example: "betty@herald.com (Betty Guernsey)"
	WebMaster string `xml:"webMaster,omitempty"`

	// -----------------
	// Required elements
	// -----------------

	// The name of the channel. It's how people refer to your
	// service. If you have an HTML website that contains the same
	// information as your RSS file, the title of your channel should
	// be the same as the title of your website.
	//
	// Example: "GoUpstate.com News Headlines"
	Title string `xml:"title"`

	// -----------------
	// Optional elements
	// -----------------

	// The language the channel is written in. This allows
	// aggregators to group all Italian language sites, for
	// example, on a single page. A list of allowable values
	// for this element, as provided by Netscape, is here.
	// You may also use values defined by the W3C.
	//
	// Example: "en-us"
	Language string `xml:"language,omitempty"`

	// The PICS rating for the channel.
	//
	// See https://www.w3.org/PICS/
	//
	// Deprecated in favor of POWDER (https://www.w3.org/2007/powder/)?
	Rating string `xml:"rating,omitempty"`

	// Phrase or sentence describing the channel.
	//
	// The latest news from GoUpstate.com, a Spartanburg
	// Herald-Journal Web site.
	Description string `xml:"description"`

	// The Url to the HTML website corresponding to the channel.
	//
	// Example: "http://www.goupstate.com/"
	Link string `xml:"link"`

	// A channel may contain any number of <item>s. An item may represent a "story" --
	// much like a story in a newspaper or magazine; if so its description is a synopsis
	// of the story, and the link points to the full story. An item may also be complete
	// in itself, if so, the description contains the text (entity-encoded HTML is allowed;
	// see examples), and the link and title may be omitted. All elements of an item are
	// optional, however at least one of title or description must be present.
	Items []Item `xml:"item,omitempty"`

	// ttl stands for time to live. It's a number of minutes that indicates
	// how long a channel can be cached before refreshing from the source.
	//
	// Example: "<ttl>60</ttl>"
	TTL int `xml:"ttl,omitempty"`
}
