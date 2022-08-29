package iterator

import (
	"github.com/sidkurella/goption/option"
)

// Iterator returns items via successive Next calls until it has run out.
// It signals that the iterator is now empty by returning None.
type Iterator[T any] interface {
	Next() option.Option[T]
}

// TODO: Implement iterator convenience methods.

// IntoIterator is an interface representing something that can turn into an Iterator.
type IntoIterator[T any] interface {
	IntoIter() Iterator[T]
}

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

func SliceIterator[T any](data []T) Iterator[*T] {
	return &sliceIter[T]{
		i:    0,
		data: data,
	}
}
