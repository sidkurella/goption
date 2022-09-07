package iterator

import (
	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/pair"
)

// Returns a Pair of the index and the value provided by the wrapped iterator.
type enumerateIterator[T any] struct {
	inner Iterator[T]
	i     int
}

// Returns a new iterator that delivers Pairs of indexes and values from the wrapped iterator.
func Enumerate[T any](iter Iterator[T]) *enumerateIterator[T] {
	return &enumerateIterator[T]{
		inner: iter,
		i:     0,
	}
}

// Returns the next item from the iterator.
// This item is a Pair of the index and the value from the wrapped iterator.
func (e *enumerateIterator[T]) Next() option.Option[pair.Pair[int, T]] {
	val, ok := e.inner.Next().Get()
	if !ok {
		return option.Nothing[pair.Pair[int, T]]{}
	}
	ret := option.Some[pair.Pair[int, T]]{
		Value: pair.Pair[int, T]{
			First:  e.i,
			Second: val,
		},
	}
	e.i++
	return ret
}
