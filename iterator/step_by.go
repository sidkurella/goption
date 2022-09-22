package iterator

import "github.com/sidkurella/goption/option"

type stepByIterator[T any] struct {
	inner    Iterator[T]
	gotFirst bool
	n        uint64
}

// Creates an iterator starting with the same element, but stepping by n each time.
// Will always return the first element of the iterator, regardless of the step provided.
func StepBy[T any](inner Iterator[T], n uint64) *stepByIterator[T] {
	return &stepByIterator[T]{
		inner:    inner,
		n:        n,
		gotFirst: false,
	}
}

func (s *stepByIterator[T]) Next() option.Option[T] {
	if !s.gotFirst {
		s.gotFirst = true
		return s.inner.Next()
	}
	return Nth(s.inner, s.n-1)
}
