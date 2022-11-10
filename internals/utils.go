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

func Map[T, S any](f func(T) S, arr []T) []S {
	what := make([]S, 0, len(arr))
	for _, v := range arr {
		what = append(what, f(v))
	}
	return what
}

func Filter[T any](f func(T) bool, arr []T) []T {
	what := make([]T, 0, len(arr))
	for _, v := range arr {
		if !f(v) {
			continue
		}
		what = append(what, v)
	}
	return what
}

func Take[T any](num int, arr []T) []T {
	if len(arr) < num {
		return arr
	}
	return arr[:num]
}

func Tail[T any](num int, arr []T) []T {
	if len(arr) < num {
		return arr
	}
	return arr[len(arr)-num : len(arr)-1]
}

func Drop[T any](num int, arr []T) []T {
	if len(arr) < num {
		return []T{}
	}
	return arr[num:]
}

func DropString(num int, what string) string {
	if len(what) < num {
		return ""
	}
	return what[num:]
}

func Any[T comparable](val T, arr []T) bool {
	for _, v := range arr {
		if val == v {
			return true
		}
	}
	return false
}

func Repeat[T any](val T, size int) []T {
	arr := make([]T, 0, size)
	for i := 0; i < size; i++ {
		arr = append(arr, val)
	}
	return arr
}
