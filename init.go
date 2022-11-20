package main

// Naked imports so the parsers and converts run their inits.
import (
	_ "github.com/thecsw/darkness/export/html"
	_ "github.com/thecsw/darkness/export/template"
	_ "github.com/thecsw/darkness/parse/orgmode"
	_ "github.com/thecsw/darkness/parse/template"
)
