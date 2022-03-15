package emilia

import (
	"github.com/thecsw/darkness/internals"
)

const (
	mathJax = `<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
	<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>`

	katexJs = `
   <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.15.3/dist/katex.min.css" integrity="sha384-KiWOvVjnN8qwAZbuQyWDIbfCLFhLXNETzBQjA/92pIowpC0d2O3nppDGQVgwd2nB" crossorigin="anonymous">
    <!-- The loading of KaTeX is deferred to speed up page rendering -->
    <script defer src="https://cdn.jsdelivr.net/npm/katex@0.15.3/dist/katex.min.js" integrity="sha384-0fdwu/T/EQMsQlrHCCHoH10pkPLlKA1jL5dFyUOvB3lfeT2540/2g6YgSi2BL14p" crossorigin="anonymous"></script>
    <!-- To automatically render math in text elements, include the auto-render extension: -->
    <script defer src="https://cdn.jsdelivr.net/npm/katex@0.15.3/dist/contrib/auto-render.min.js" integrity="sha384-+XBljXPPiv+OzfbB3cVmLHf4hdUFHlWNZN5spNQ7rmHTXpd7WvJum6fIACpNNfIR" crossorigin="anonymous"
        onload="renderMathInElement(document.body);"></script>
<script>
    document.addEventListener("DOMContentLoaded", function() {
        renderMathInElement(document.body, {
          // customised options
          // • auto-render specific keys, e.g.:
          delimiters: [
              {left: '$$', right: '$$', display: true},
              {left: '$', right: '$', display: false},
              {left: '\\(', right: '\\)', display: true},
              {left: '\\[', right: '\\]', display: false},
              {left: "\\begin{equation}", right: "\\end{equation}", display: true},
              {left: "\\begin{equation*}", right: "\\end{equation*}", display: true},
              {left: "\\begin{align}", right: "\\end{align}", display: true},
              {left: "\\begin{align*}", right: "\\end{align*}", display: true},
              {left: "\\begin{aligned}", right: "\\end{aligned}", display: true},
              {left: "\\begin{aligned*}", right: "\\end{aligned*}", display: true},
              {left: "\\begin{alignat}", right: "\\end{alignat}", display: true},
              {left: "\\begin{gather}", right: "\\end{gather}", display: true},
              {left: "\\begin{CD}", right: "\\end{CD}", display: true}
          ],
          // • rendering keys, e.g.:
          throwOnError : false
        });
    });
</script>
`
	mathJs = katexJs
)

func AddMathSupport(page *internals.Page) {
	// Find any match of the math regexp, if found, add the math script
	for _, content := range page.Contents {
		// If it's in our paragraph
		if content.IsParagraph() {
			// If found, add the script and leave
			if internals.MathRegexp.MatchString(content.Paragraph) {
				page.Scripts = append(page.Scripts, mathJs)
				return
			}
		}
		// Or in the heading
		if content.IsHeading() {
			// If found, add the script and leave
			if internals.MathRegexp.MatchString(content.Heading) {
				page.Scripts = append(page.Scripts, mathJs)
				return
			}
		}
		// Or if it's a list
		if content.IsList() {
			for _, item := range content.List {
				// If found, add the script and leave
				if internals.MathRegexp.MatchString(item) {
					page.Scripts = append(page.Scripts, mathJs)
					return
				}

			}
		}
	}
}
