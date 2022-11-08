package internals

type PageOption func(*Page)

const (
	defaultPageFilename   = "unknown"
	defaultPageTitle      = "no title"
	defaulteDateHoloscene = true
	defaultDate           = "0; 12000 H.E."
	defaultURL            = "unknown"
)

func NewPage(options ...PageOption) *Page {
	p := &Page{
		File:          defaultPageFilename,
		Title:         defaultPageTitle,
		Date:          defaultDate,
		DateHoloscene: defaulteDateHoloscene,
		URL:           defaultURL,
		Contents:      []Content{},
		Footnotes:     []string{},
		Scripts:       []string{},
		Stylesheets:   []string{},
	}
	return p.Options(options...)
}

func (p *Page) Options(options ...PageOption) *Page {
	for _, option := range options {
		option(p)
	}
	return p
}

func WithFilename(filename string) PageOption {
	return func(p *Page) {
		p.File = filename
	}
}

func WithURL(url string) PageOption {
	return func(p *Page) {
		p.URL = url
	}
}

func WithContents(contents []Content) PageOption {
	return func(p *Page) {
		p.Contents = contents
	}
}
