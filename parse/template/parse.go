package template

import "github.com/thecsw/darkness/yunyun"

// Parse will convert the input format into `*yunyun.Page`.
func (p ParserTemplate) Parse() *yunyun.Page {
	// Define your specific markup here and always
	// ask Yunyun to build the regex rules for you.
	yunyun.ActiveMarkings.BuildRegex()
	panic("implement me")
}
