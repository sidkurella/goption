package sliceutil

import (
	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/pair"
)

// Map applies the given function to each element of T, collecting the results in a new slice.
func Map[T any, U any](s []T, f func(T) U) []U {
	if s == nil {
		return nil
	}

	ret := make([]U, 0, len(s))
	for _, v := range s {
		ret = append(ret, f(v))
	}
	return ret
}

// IndexMap applies the given function to each element of T, collecting the results in a new slice.
// The function may access the index of the element as well.
func IndexMap[T any, U any](s []T, f func(int, T) U) []U {
	if s == nil {
		return nil
	}

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
	if s == nil {
		return nil
	}

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

// Returns the first element of the slice. Returns Nothing if it is empty.
func First[T any](s []T) option.Option[T] {
	if len(s) > 0 {
		return option.Some(s[0])
	}
	return option.Nothing[T]()
}

// Returns the last element of the slice. Returns Nothing if it is empty.
func Last[T any](s []T) option.Option[T] {
	if len(s) > 0 {
		return option.Some(s[len(s)-1])
	}
	return option.Nothing[T]()
}

// Returns if the slice `haystack` starts with the prefix `needle`.
func StartsWith[T comparable](haystack []T, needle []T) bool {
	return StartsWithFunc(haystack, needle, func(t1 T, t2 T) bool {
		return t1 == t2
	})
}

// Returns if the slice `haystack` starts with the prefix `needle`.
// Equality is tested with the given function.
func StartsWithFunc[T any](haystack []T, needle []T, isEqual func(T, T) bool) bool {
	if len(haystack) < len(needle) {
		return false
	}
	return iterator.All[pair.Pair[T, T]](
		iterator.Zip[T, T](Iter(haystack), Iter(needle)),
		func(t pair.Pair[T, T]) bool {
			return isEqual(t.First, t.Second)
		},
	)
}

// Returns if the slice `haystack` ends with the suffix `needle`.
func EndsWith[T comparable](haystack []T, needle []T) bool {
	return EndsWithFunc(haystack, needle, func(t1 T, t2 T) bool {
		return t1 == t2
	})
}

// Returns if the slice `haystack` ends with with the suffix `needle`.
// Equality is tested with the given function.
func EndsWithFunc[T any](haystack []T, needle []T, isEqual func(T, T) bool) bool {
	if len(haystack) < len(needle) {
		return false
	}
	return iterator.All[pair.Pair[T, T]](
		iterator.Zip[T, T](ReverseIter(haystack), ReverseIter(needle)),
		func(t pair.Pair[T, T]) bool {
			return isEqual(t.First, t.Second)
		},
	)
}

// Returns the slice `haystack` with the prefix `needle` removed.
// If it doesn't start with the prefix, then the original `haystack` is returned.
func StripPrefix[T comparable](haystack []T, needle []T) []T {
	if StartsWith(haystack, needle) {
		return haystack[len(needle):]
	}
	return haystack
}

// Returns the slice `haystack` with the prefix `needle` removed.
// If it doesn't start with the prefix, then the original `haystack` is returned.
// Equality is tested with the given function.
func StripPrefixFunc[T any](haystack []T, needle []T, isEqual func(T, T) bool) []T {
	if StartsWithFunc(haystack, needle, isEqual) {
		return haystack[len(needle):]
	}
	return haystack
}

// Returns the slice `haystack` with the suffix `needle` removed.
// If it doesn't end with the suffix, then the original `haystack` is returned.
func StripSuffix[T comparable](haystack []T, needle []T) []T {
	if EndsWith(haystack, needle) {
		return haystack[:len(haystack)-len(needle)]
	}
	return haystack
}

// Returns the slice `haystack` with the suffix `needle` removed.
// If it doesn't end with the suffix, then the original `haystack` is returned.
// Equality is tested with the given function.
func StripSuffixFunc[T any](haystack []T, needle []T, isEqual func(T, T) bool) []T {
	if EndsWithFunc(haystack, needle, isEqual) {
		return haystack[:len(haystack)-len(needle)]
	}
	return haystack
}

// Truncates the slice to the given max length.
func Truncate[T any](s []T, maxLength int) []T {
	if len(s) < maxLength {
		return s
	}
	return s[:maxLength]
}
