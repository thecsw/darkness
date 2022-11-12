package parse

import "github.com/thecsw/darkness/yunyun"

// Parser is an interface used to define import packages,
// which convert source data into a yunyun `Page`.
type Parser interface {
	// Parse takes the `data` and `filename` and returns `*Page`.
	Parse(string, string) *yunyun.Page
}
