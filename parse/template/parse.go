package template

import "github.com/thecsw/darkness/yunyun"

// Parse will convert the input format into `*yunyun.Page`.
func (p ParserTemplate) Parse() *yunyun.Page {
	// Define your specific markup here and always
	// ask Yunyun to build the regex rules for you.
	yunyun.ActiveMarkings.BuildRegex()

	page := yunyun.NewPage(
		yunyun.WithFilename(p.Filename),
		yunyun.WithLocation(yunyun.RelativePathTrim(p.Filename)),
		yunyun.WithContents(make([]*yunyun.Content, 0, 16)),
	)

	// Remove this panic and proceed with implementation.
	panic("implement me " + page.File)
}
