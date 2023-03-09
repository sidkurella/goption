package iterator

import "github.com/sidkurella/goption/option"

type mapWhileIterator[T any, U any] struct {
	inner Iterator[T]
	pred  func(T) option.Option[U]
}

// Returns an iterator yielding elements from inner, mapped by pred.
// The iterator stops once pred returns Nothing. Yields elements while pred continues to return Some.
func MapWhile[T any, U any](inner Iterator[T], pred func(T) option.Option[U]) *mapWhileIterator[T, U] {
	return &mapWhileIterator[T, U]{
		inner: inner,
		pred:  pred,
	}
}

func (m *mapWhileIterator[T, U]) Next() option.Option[U] {
	return option.Match(m.inner.Next(),
		func(t T) option.Option[U] {
			return m.pred(t)
		},
		func() option.Option[U] {
			return option.Nothing[U]()
		},
	)
}
