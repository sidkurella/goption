package iterator

import "github.com/sidkurella/goption/option"

type takeWhileIterator[T any] struct {
	inner Iterator[T]
	done  bool
	f     func(T) bool
}

// Creates an iterator that takes the elements for which pred is true.
// As soon as pred returns false, the predicate is no longer evaluated.
// The rest of the elements will be ignored.
// The element for which false is returned will be removed from the iterator in order to check it.
func TakeWhile[T any](inner Iterator[T], f func(T) bool) *takeWhileIterator[T] {
	return &takeWhileIterator[T]{
		inner: inner,
		f:     f,
	}
}

func (s *takeWhileIterator[T]) Next() option.Option[T] {
	if s.done {
		return option.Nothing[T]{}
	}
	val := s.inner.Next()
	if val.IsNothing() {
		s.done = true
		return val
	}
	if s.f(val.Unwrap()) {
		return val
	}
	s.done = true
	return option.Nothing[T]{}
}
