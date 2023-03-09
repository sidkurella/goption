package sliceutil

import (
	"github.com/sidkurella/goption/option"
)

type sliceIter[T any] struct {
	i    int // Represents the next element to return from the slice.
	data []T // The underlying data for this iterator.
}

func (s *sliceIter[T]) Next() option.Option[T] {
	if s.i >= len(s.data) {
		return option.Nothing[T]()
	}
	ret := option.Some(s.data[s.i])
	s.i++
	return ret
}

// Returns an iterator for the given slice.
// The iterator iterates over elements in the slice; it will shallow-copy them.
func Iter[T any](data []T) *sliceIter[T] {
	return &sliceIter[T]{
		i:    0,
		data: data,
	}
}

type sliceReverseIter[T any] struct {
	i    int // Represents the next element to return from the slice.
	data []T // The underlying data for this iterator.
}

func (s *sliceReverseIter[T]) Next() option.Option[T] {
	if s.i < 0 {
		return option.Nothing[T]()
	}
	ret := option.Some(s.data[s.i])
	s.i--
	return ret
}

// Returns an iterator ranging backwards over the slice.
// The iterator iterates over elements in the slice; it will shallow-copy them.
func ReverseIter[T any](data []T) *sliceReverseIter[T] {
	return &sliceReverseIter[T]{
		i:    len(data) - 1,
		data: data,
	}
}

type slicePointerIter[T any] struct {
	i    int // Represents the next element to return from the slice.
	data []T // The underlying data for this iterator.
}

func (s *slicePointerIter[T]) Next() option.Option[*T] {
	if s.i >= len(s.data) {
		return option.Nothing[*T]()
	}
	ret := option.Some(&(s.data[s.i]))
	s.i++
	return ret
}

// Returns an iterator for the given slice.
// The iterator iterates over pointers to elements in the slice.
func PointerIter[T any](data []T) *slicePointerIter[T] {
	return &slicePointerIter[T]{
		i:    0,
		data: data,
	}
}

type sliceReversePointerIter[T any] struct {
	i    int // Represents the next element to return from the slice.
	data []T // The underlying data for this iterator.
}

func (s *sliceReversePointerIter[T]) Next() option.Option[*T] {
	if s.i < 0 {
		return option.Nothing[*T]()
	}
	ret := option.Some(&(s.data[s.i]))
	s.i--
	return ret
}

// Returns an iterator ranging backwards over the slice.
// The iterator iterates over pointers to elements in the slice.
func ReversePointerIter[T any](data []T) *sliceReversePointerIter[T] {
	return &sliceReversePointerIter[T]{
		i:    len(data) - 1,
		data: data,
	}
}
