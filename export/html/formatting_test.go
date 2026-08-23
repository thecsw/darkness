package html

import "testing"

func TestPlainTitleRemovesMathSyntax(t *testing.T) {
	got := plainTitle(`$\sqrt{2}$ is irrational ◻`)
	want := `√2 is irrational ◻`
	if got != want {
		t.Fatalf("plainTitle() = %q, want %q", got, want)
	}
}

func TestPlainTitleRemovesFormatting(t *testing.T) {
	got := plainTitle(`*The* $x^2$ and \(y\)`)
	want := `The x^2 and y`
	if got != want {
		t.Fatalf("plainTitle() = %q, want %q", got, want)
	}
}
