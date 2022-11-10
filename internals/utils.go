package internals

import "golang.org/x/exp/constraints"

// Min returns the minimum of two numbers.
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two numbers.
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// First returns the first elemnt of the given array, zero value otherwise.
func First[T any](x []T) T {
	if len(x) > 0 {
		return x[0]
	}
	return ZeroValue[T]()
}

// Last returns the last element of the given array, zero value otherwise.
func Last[T any](x []T) T {
	if len(x) > 0 {
		return x[len(x)-1]
	}
	return ZeroValue[T]()
}

// ZeroValue returns the zero value of any type.
func ZeroValue[T any]() T {
	var t T
	return t
}

// Map calls the function on every array element and returns results in list.
func Map[T, S any](f func(T) S, arr []T) []S {
	what := make([]S, 0, len(arr))
	for _, v := range arr {
		what = append(what, f(v))
	}
	return what
}

// Filter calls the function and returns list of values that returned true.
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

// Take returns up to the first `num` elements.
func Take[T any](num int, arr []T) []T {
	if len(arr) < num {
		return arr
	}
	return arr[:num]
}

// Tail returns up to the last `num` elements.
func Tail[T any](num int, arr []T) []T {
	if len(arr) < num {
		return arr
	}
	return arr[len(arr)-num : len(arr)-1]
}

// Drop drops up to the first `num` elements.
func Drop[T any](num int, arr []T) []T {
	if len(arr) < num {
		return []T{}
	}
	return arr[num:]
}

// DropString drops up to the first `num` runes of a string.
func DropString(num int, what string) string {
	if len(what) < num {
		return ""
	}
	return what[num:]
}

// Any returns true if all any element in the list matches the given value.
func Any[T comparable](val T, arr []T) bool {
	for _, v := range arr {
		if val == v {
			return true
		}
	}
	return false
}

// All returns true if all elements in the list match the given value.
func All[T comparable](val T, arr []T) bool {
	for _, v := range arr {
		if val != v {
			return false
		}
	}
	return true
}

// Repeat creates a list of given size consisting of the same given value.
func Repeat[T any](val T, size int) []T {
	arr := make([]T, 0, size)
	for i := 0; i < size; i++ {
		arr = append(arr, val)
	}
	return arr
}
