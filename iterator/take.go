package iterator

import "github.com/sidkurella/goption/option"

type takeIterator[T any] struct {
	inner Iterator[T]
	left  uint64
}

// Creates an iterator yielding the first n elements, or fewer if the provided iterator ends first.
func Take[T any](inner Iterator[T], n uint64) *takeIterator[T] {
	return &takeIterator[T]{
		inner: inner,
		left:  n,
	}
}

func (t *takeIterator[T]) Next() option.Option[T] {
	if t.left > 0 {
		t.left--
		return t.inner.Next()
	}
	return option.Nothing[T]{}
}
