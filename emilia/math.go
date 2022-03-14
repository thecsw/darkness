package emilia

import (
	"github.com/thecsw/darkness/internals"
)

const (
	mathJax = `<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>`
)

func AddMathSupport(page *internals.Page) {
	// Find any match of the math regexp, if found, add the math script
	for _, content := range page.Contents {
		// If it's in our paragraph
		if content.IsParagraph() {
			// If found, add the script and leave
			if internals.MathRegexp.MatchString(content.Paragraph) {
				page.Scripts = append(page.Scripts, mathJax)
				return
			}
		}
		// Or in the heading
		if content.IsHeading() {
			// If found, add the script and leave
			if internals.MathRegexp.MatchString(content.Heading) {
				page.Scripts = append(page.Scripts, mathJax)
				return
			}
		}
		// Or if it's a list
		if content.IsList() {
			for _, item := range content.List {
				// If found, add the script and leave
				if internals.MathRegexp.MatchString(item) {
					page.Scripts = append(page.Scripts, mathJax)
					return
				}

			}
		}
	}
}
