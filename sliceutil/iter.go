package sliceutil

import (
	"github.com/sidkurella/goption/option"
)

type sliceIter[T any] struct {
	i    int // Represents the next element to return from the slice.
	data []T // The underlying data for this iterator.
}

func (s *sliceIter[T]) Next() option.Option[*T] {
	if s.i >= len(s.data) {
		return option.Nothing[*T]{}
	}
	ret := option.Some[*T]{Value: &(s.data[s.i])}
	s.i++
	return ret
}

// Returns an iterator for the given slice.
// The iterator iterates over pointers to elements in the slice.
func SliceIterator[T any](data []T) *sliceIter[T] {
	return &sliceIter[T]{
		i:    0,
		data: data,
	}
}
