package internals

type TypeContent uint8

const (
	TypeHeader TypeContent = iota
	TypeParagraph
	TypeList
	TypeListNumbered
	TypeLink
	TypeImage
	TypeYoutube
	TypeSpotifyTrack
	TypeSpotifyPlaylist
	TypeSourceCode
	TypeRawHTML
)

type Page struct {
	Title    string
	URL      string
	MetaTags []MetaTag
	Links    []Link
	Contents []Content
}

type MetaTag struct {
	Name     string
	Content  string
	Property string
}

type Link struct {
	Rel  string
	Type string
	Href string
}

type Content struct {
	Type TypeContent

	HeaderLevel     int
	Header          string
	Paragraph       string
	List            []string
	ListNumbered    []string
	Link            string
	LinkTitle       string
	ImageSource     string
	ImageCaption    string
	Youtube         string
	SpotifyTrack    string
	SpotifyPlaylist string
	SourceCode      string
	SourceCodeLang  string
	RawHTML         string
}
