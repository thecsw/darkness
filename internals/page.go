package internals

type Page struct {
	Title     string
	URL       string
	MetaTags  []MetaTag
	Links     []Link
	Contents  []Content
	Footnotes []string
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
