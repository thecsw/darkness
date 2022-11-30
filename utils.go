package main

import (
	"strings"

	_ "github.com/thecsw/darkness/export/html"
	_ "github.com/thecsw/darkness/export/template"
	_ "github.com/thecsw/darkness/parse/markdown"
	_ "github.com/thecsw/darkness/parse/orgmode"
	_ "github.com/thecsw/darkness/parse/template"
)

// fdb cleans the filename from absolute workspace prefix.
func fdb(filename, data *string) (string, string) {
	return strings.TrimPrefix(*filename, workDir+"/"), *data
}
