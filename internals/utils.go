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

func First[T any](x []T) T {
	if len(x) > 0 {
		return x[0]
	}
	return ZeroValue[T]()
}

func Last[T any](x []T) T {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return ZeroValue[T]()
}

func ZeroValue[T any]() T {
	var t T
	return t
}
