package iterator

import "github.com/sidkurella/goption/option"

// An iterator that return elements that are mapped by f.
type mapIterator[T any, U any] struct {
	inner Iterator[T]
	f     func(T) U
}

// Creates an iterator that maps elements from the original iterator.
// Converts an iterator from type T to U.
func Map[T any, U any](iter Iterator[T], f func(T) U) *mapIterator[T, U] {
	return &mapIterator[T, U]{
		inner: iter,
		f:     f,
	}
}

func (m *mapIterator[T, U]) Next() option.Option[U] {
	return option.Match(m.inner.Next(),
		func(t T) option.Option[U] {
			return option.Some(m.f(t))
		},
		func() option.Option[U] {
			return option.Nothing[U]()
		},
	)
}
