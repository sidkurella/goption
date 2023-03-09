package iterator

import "github.com/sidkurella/goption/option"

type fuseIterator[T any] struct {
	inner         Iterator[T]
	returnNothing bool
}

// Creates an iterator that will always return Nothing after the first Nothing.
func Fuse[T any](inner Iterator[T]) *fuseIterator[T] {
	return &fuseIterator[T]{
		inner:         inner,
		returnNothing: false,
	}
}

func (f *fuseIterator[T]) Next() option.Option[T] {
	if f.returnNothing {
		return option.Nothing[T]()
	}
	ret := f.inner.Next()
	f.returnNothing = ret.IsNothing()
	return ret
}
