package narumi

import (
	"strings"
	"testing"
)

func TestKatexRenderConfigUsesCorrectDisplayModes(t *testing.T) {
	if !strings.Contains(katexRenderConfig, `{left: '\\(', right: '\\)', display: false}`) {
		t.Error(`inline \\(…\\) delimiters must render inline`)
	}
	if !strings.Contains(katexRenderConfig, `{left: '\\[', right: '\\]', display: true}`) {
		t.Error(`display \\[…\\] delimiters must render as a block`)
	}
}
