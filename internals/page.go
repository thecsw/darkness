package internals

// Page is a struct for holding the page contents
type Page struct {
	// Title is the title of the page
	Title string
	// Date is the date of the page
	Date string
	// DateHoloscene tells us whether the first paragraph
	// on the page is given as holoscene date stamp
	DateHoloscene bool
	// URL is the URL of the page
	URL string
	// Contents is the contents of the page
	Contents []Content
	// Footnotes is the footnotes of the page
	Footnotes []string
	// Scripts is the scripts of the page
	Scripts []string
	// Stylesheets is the list of css of the page
	Stylesheets []string
}

// MetaTag is a struct for holding the meta tag
type MetaTag struct {
	// Name is the name of the meta tag
	Name string
	// Content is the content of the meta tag
	Content string
	// Propery is the property of the meta tag
	Property string
}

// Link is a struct for holding the link tag
type Link struct {
	// Rel is the rel of the link tag
	Rel string
	// Type is the type of the link tag
	Type string
	// Href is the href of the link tag
	Href string
}
