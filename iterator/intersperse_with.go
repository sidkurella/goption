package iterator

import "github.com/sidkurella/goption/option"

type intersperseWithIterator[T any] struct {
	inner   Iterator[T]
	f       func() T
	nextIsF bool
}

// Creates a new iterator which calls f and places a copy of its return value between items of the original iterator.
// Before the final Nothing, f will be called and its result will be returned.
func IntersperseWith[T any](inner Iterator[T], f func() T) *intersperseWithIterator[T] {
	return &intersperseWithIterator[T]{
		inner:   inner,
		f:       f,
		nextIsF: false,
	}
}

// Returns the next item from the interspersed iterator.
func (i *intersperseWithIterator[T]) Next() option.Option[T] {
	var ret option.Option[T]
	if i.nextIsF {
		ret = option.Some(i.f())
	} else {
		ret = i.inner.Next()
	}
	i.nextIsF = !i.nextIsF
	return ret
}
