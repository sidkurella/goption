package sliceutil

// Map applies the given function to each element of T, collecting the results in a new slice.
func Map[T any, U any](s []T, f func(T) U) []U {
	ret := make([]U, 0, len(s))
	for _, v := range s {
		ret = append(ret, f(v))
	}
	return ret
}

// IndexMap applies the given function to each element of T, collecting the results in a new slice.
// The function may access the index of the element as well.
func IndexMap[T any, U any](s []T, f func(int, T) U) []U {
	ret := make([]U, 0, len(s))
	for i, v := range s {
		ret = append(ret, f(i, v))
	}
	return ret
}

// Reverse reverses a slice in-place.
func Reverse[T any](s []T) []T {
	max := len(s)
	start := 0
	end := max - start - 1
	for start < end {
		s[start], s[end] = s[end], s[start]
		start++
		end--
	}
	return s
}

// Reversed returns a new slice that is the reverse of the provided one.
// The old one is not modified.
func Reversed[T any](s []T) []T {
	ret := make([]T, 0, len(s))
	max := len(s)
	for i := 0; i < max; i++ {
		ret = append(ret, s[max-i-1])
	}
	return ret
}

// FoldLeft calls f successively with each value in the list and the current accumulator.
// The accumulator value is then updated with the new return value of f.
// FoldLeft proceeds from the left, calling f(a, f[0]) first.
func FoldLeft[T any, A any](s []T, a A, f func(A, T) A) A {
	for _, v := range s {
		a = f(a, v)
	}
	return a
}

// FoldRight calls f successively with each value in the list and the current accumulator.
// The accumulator value is then updated with the new return value of f.
// FoldRight proceeds from the right, calling f(a, f[len(s)-1]) first.
func FoldRight[T any, A any](s []T, a A, f func(A, T) A) A {
	max := len(s)
	for i := 0; i < max; i++ {
		a = f(a, s[max-i-1])
	}
	return a
}
