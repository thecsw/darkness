package internals

// Number is a type for numbers
type Number interface {
	int | float64
}

// Min returns the minimum of two numbers
func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two numbers
func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}
