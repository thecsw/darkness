package rss

import "encoding/xml"

type TextInput struct {
	XMLName xml.Name `xml:"textInput"`

	// -----------------
	// Required elements
	// -----------------

	// The label of the Submit button in the text input area.
	Title string `xml:"title"`

	// Explains the text input area.
	Description string `xml:"description"`

	// The name of the text object in the text input area.
	Name string `xml:"name"`

	// The Url of the CGI script that processes text input requests.
	Link string `xml:"link"`
}
