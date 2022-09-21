package iterator

import "github.com/sidkurella/goption/option"

type skipIterator[T any] struct {
	inner   Iterator[T]
	skipped bool
	n       uint64
}

// Creates an iterator that skips the first n elements.
// If the iterator is too short, an empty iterator is returned.
func Skip[T any](inner Iterator[T], n uint64) *skipIterator[T] {
	return &skipIterator[T]{
		inner: inner,
		n:     n,
	}
}

func (s *skipIterator[T]) Next() option.Option[T] {
	if !s.skipped {
		ret := Nth(s.inner, s.n)
		s.skipped = true
		return ret
	}
	return s.inner.Next()
}
