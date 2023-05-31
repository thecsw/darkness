package yunyun

// PageOption representions a function that can be passed
// to a new `Page` instantiation to modify the state.
type PageOption func(*Page)

const (
	defaultPageFilename   = "unknown"
	defaultPageTitle      = "no title"
	defaulteDateHoloscene = true
	defaultDate           = "0; 12000 H.E."
	defaultURL            = "unknown"
	defaultPreviewWidth   = `1280`
	defaultPreviewHeight  = `640`
)

// NewPage creates a new `Page` and runs passed options.
func NewPage(options ...PageOption) *Page {
	p := &Page{
		File:          defaultPageFilename,
		Title:         defaultPageTitle,
		Author:        "",
		Date:          defaultDate,
		DateHoloscene: defaulteDateHoloscene,
		Location:      defaultURL,
		Contents:      nil,
		Footnotes:     make([]string, 0, 2),
		Scripts:       make([]string, 0, 4),
		Stylesheets:   make([]string, 0, 2),
		HtmlHead:      make([]string, 0, 2),
		Accoutrement: &Accoutrement{
			ExcludeHtmlHeadContains: make([]string, 0, 2),
			PreviewWidth:            defaultPreviewWidth,
			PreviewHeight:           defaultPreviewHeight,
		},
	}
	return p.Options(options...)
}

// Options runs the provided options.
func (p *Page) Options(options ...PageOption) *Page {
	for _, option := range options {
		option(p)
	}
	return p
}

// WithFilename sets the filename.
func WithFilename(filename RelativePathFile) PageOption {
	return func(p *Page) {
		p.File = filename
	}
}

// WithLocation sets the URL.
func WithLocation(url RelativePathDir) PageOption {
	return func(p *Page) {
		p.Location = url
	}
}

// WithContents sets the contents.
func WithContents(contents Contents) PageOption {
	return func(p *Page) {
		p.Contents = contents
	}
}
