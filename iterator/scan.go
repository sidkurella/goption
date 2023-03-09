package iterator

import "github.com/sidkurella/goption/option"

type scanIterator[T any, S any, U any] struct {
	inner Iterator[T]
	state *S
	f     func(*S, T) option.Option[U]
}

// An iterator adapter similar to fold that holds internal state and produces a new iterator.
// Calls f with the current state and the current item, and yields a new item mapped by the closure.
// NOTE: The closure is responsible for updating the state; it is not automatically assigned the result of the closure.
func Scan[T any, S any, U any](inner Iterator[T], initial S, f func(*S, T) option.Option[U]) *scanIterator[T, S, U] {
	return &scanIterator[T, S, U]{
		inner: inner,
		state: &initial,
		f:     f,
	}
}

func (s *scanIterator[T, S, U]) Next() option.Option[U] {
	return option.Match(s.inner.Next(),
		func(t T) option.Option[U] {
			return s.f(s.state, t)
		},
		func() option.Option[U] {
			return option.Nothing[U]()
		},
	)
}
