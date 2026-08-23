package yunyun

import "testing"

func TestPlainTextRendersCommonMathReadably(t *testing.T) {
	ActiveMarkings.BuildRegex()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"square root", `$\sqrt{2}$ is irrational ◻`, `√2 is irrational ◻`},
		{"bold math", `$\mathbf{Theorem}$`, `Theorem`},
		{"fraction", `$\frac{a}{b}$`, `a/b`},
		{"parentheses", `\(x\)`, `x`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := PlainText(test.input); got != test.want {
				t.Fatalf("PlainText(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}
