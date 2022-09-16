package iterator

import "github.com/sidkurella/goption/option"

type inspectIterator[T any] struct {
	inner Iterator[T]
	f     func(T)
}

// Does something with each element of an iterator, passing the value on.
func Inspect[T any](inner Iterator[T], f func(T)) *inspectIterator[T] {
	return &inspectIterator[T]{
		inner: inner,
		f:     f,
	}
}

func (i *inspectIterator[T]) Next() option.Option[T] {
	ret := i.inner.Next()
	val, ok := ret.Get()
	if ok {
		i.f(val)
	}
	return ret
}
