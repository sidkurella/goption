package iterator

import "github.com/sidkurella/goption/option"

// Returns an iterator that only returns elements that pass pred.
type filterIterator[T any] struct {
	inner Iterator[T]
	pred  func(T) bool
}

// Returns an iterator that only returns elements from iter that pass pred.
func Filter[T any](iter Iterator[T], pred func(T) bool) *filterIterator[T] {
	return &filterIterator[T]{
		inner: iter,
		pred:  pred,
	}
}

// Returns the next element from the inner iterator that passes the filter predicate.
func (f *filterIterator[T]) Next() option.Option[T] {
	return Find(f.inner, f.pred)
}
