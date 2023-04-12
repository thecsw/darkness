package rss

import "encoding/xml"

// <cloud> is an optional sub-element of <channel>.
type Cloud struct {
	XMLName xml.Name `xml:"cloud"`

	// A list of URLs of RSS documents that the client seeks to monitor
	Domain string `xml:"domain,attr"`

	// The client's remote procedure call path.
	Path string `xml:"path,attr"`

	// The name of the remote procedure the cloud should
	// call on the client upon an update.
	RegisterProcedure string `xml:"registerProcedure,attr"`

	// The string "xml-rpc" if the client employs XML-RPC,
	// "soap" for SOAP and "http-post" for REST.
	Protocol string `xml:"protocol,attr"`

	// The client's TCP port.
	Port int `xml:"port,attr"`
}
