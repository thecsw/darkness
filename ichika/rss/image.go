package rss

// <image> is an optional sub-element of <channel>, which
// contains three required and three optional sub-elements.
type Image struct {
	// -----------------
	// Required elements
	// -----------------

	// The URL of a GIF, JPEG or PNG image that represents the channel.
	Url string `xml:"url"`

	// Describes the image, it's used in the ALT attribute of the HTML
	// <img> tag when the channel is rendered in HTML.
	Title string `xml:"title"`

	// URL of the site, when the channel is rendered, the image is a
	// link to the site. (Note, in practice the image <title> and <link>
	// should have the same value as the channel's <title> and <link>.
	Link string `xml:"link"`

	// Contains text that is included in the TITLE attribute of the
	// link formed around the image in the HTML rendering.
	Description string `xml:"description,omitempty"`

	// -----------------
	// Optional elements
	// -----------------

	// Optional elements include <width> and <height>, numbers,
	// indicating the width and height of the image in pixels.
	Width int `xml:"width,omitempty"`

	// Optional elements include <width> and <height>, numbers,
	// indicating the width and height of the image in pixels.
	Height int `xml:"height,omitempty"`
}
