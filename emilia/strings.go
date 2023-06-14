package emilia

import "strings"

// CountRunesLeft counts the number of times a rune appears
// at the beginning of a string.
func CountRunesLeft(s string, r rune) uint8 {
	count := uint8(0)
	for _, rr := range s {
		if rr != r {
			return count
		}
		count++
	}
	return count
}

// CountStringsLeft counts the number of times a substring appears
// at the beginning of a string.
func CountStringsLeft(s, substr string) uint8 {
	count := uint8(0)
	for strings.HasPrefix(s, substr) {
		count++
		s = strings.TrimPrefix(s, substr)
	}
	return count
}
