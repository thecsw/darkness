package narumi

import (
	"strings"

	"github.com/thecsw/darkness/v3/yunyun"
	"github.com/thecsw/gana"
)

const (
	katexLocalCSS        yunyun.RelativePathFile = `scripts/katex/katex.min.css`
	katexLocalJS         yunyun.RelativePathFile = `scripts/katex/katex.min.js`
	katexLocalAutoRender yunyun.RelativePathFile = `scripts/katex/auto-render.min.js`

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

// WithMathSupport adds math support to the page using javascript injection
func WithMathSupport() yunyun.PageOption {
	return func(page *yunyun.Page) {
		if page == nil || page.Contents == nil || page.Accoutrement == nil {
			return
		}
		// If we found math-related tags or forced by user
		if hasMathEquations(page) && !page.Accoutrement.Math.IsDisabled() {
			page.Scripts = append(page.Scripts, mathJs)
		}
	}
}

// hasMathEquations returns true if the page has any math equations and
// returns false otherwise.
func hasMathEquations(page *yunyun.Page) bool {
	return gana.Anyf(hasEquationInContent, page.Contents)
}

// hasEquationInContent returns true if the content has math equations in it.
func hasEquationInContent(content *yunyun.Content) bool {
	return hasEquationInParagraph(content) ||
		hasEquationInList(content) ||
		hasEquationsInHeading(content)
}

// hasEquationInParagraph returns true if the content is a paragraph
// AND there is some math in there.
func hasEquationInParagraph(content *yunyun.Content) bool {
	if content.IsParagraph() && (strings.Contains(content.Paragraph, `\begin`) ||
		yunyun.MathRegexp.MatchString(content.Paragraph)) {
		return true
	}
	// If none of the above worked, give up on this paragraph.
	return false
}

// hasEquationInList returns true if the list has math equations.
func hasEquationInList(content *yunyun.Content) bool {
	if !content.IsList() {
		return false
	}
	return gana.Anyf(
		yunyun.MathRegexp.MatchString,
		gana.Map(func(t yunyun.ListItem) string { return t.Text }, content.List),
	)
}

// hasEquationsInHeading returns true if the heading has an equation.
func hasEquationsInHeading(content *yunyun.Content) bool {
	if !content.IsHeading() {
		return false
	}
	return yunyun.MathRegexp.MatchString(content.Heading)
}
