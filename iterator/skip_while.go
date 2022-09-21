package iterator

import "github.com/sidkurella/goption/option"

type skipWhileIterator[T any] struct {
	inner   Iterator[T]
	skipped bool
	f       func(T) bool
}

// Creates an iterator that skips the elements for which pred is true.
// As soon as pred returns false, the predicate is no longer evaluated.
// Even if further elements would pass the predicate, they will still be yielded.
func SkipWhile[T any](inner Iterator[T], f func(T) bool) *skipWhileIterator[T] {
	return &skipWhileIterator[T]{
		inner: inner,
		f:     f,
	}
}

func (s *skipWhileIterator[T]) Next() option.Option[T] {
	if !s.skipped {
		for item := s.inner.Next(); item.IsSome(); item = s.inner.Next() {
			val := item.Unwrap()
			if !s.f(val) {
				s.skipped = true
				return option.Some[T]{Value: val}
			}
		}
		return option.Nothing[T]{}
	}
	return s.inner.Next()
}
