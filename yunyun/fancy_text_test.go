package yunyun

import "testing"

func TestFancyText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"no replacements", "hello world", "hello world"},
		{"contraction s", "it's", "it’s"},
		{"contraction d", "he'd", "he’d"},
		{"possessive", "dogs' toy", "dogs’ toy"},
		{"contraction m", "i'm", "i’m"},
		{"negation", "don't", "don’t"},
		{"contraction re", "they're", "they’re"},
		{"contraction ve", "we've", "we’ve"},
		{"contraction ll", "she'll", "she’ll"},
		{"double backticks", "``quote``", "“quote“"},
		{"double single quotes", "''quote''", "”quote”"},
		{"em dash", "a---b", "a—b"},
		{"en dash", "a--b", "a–b"},
		{"multiple replacements", "it's ``done'' -- yay", "it’s “done” – yay"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := FancyText(tc.input); got != tc.want {
				t.Fatalf("FancyText(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
