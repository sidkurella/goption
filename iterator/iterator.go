package iterator

import (
	"github.com/sidkurella/goption/option"
)

// Iterator returns items via successive Next calls until it has run out.
// It signals that the iterator is now empty by returning None.
type Iterator[T any] interface {
	Next() option.Option[T]
}

// Advances the iterator by n and returns the nth next item.
func Take[T any](iter Iterator[T], n int) option.Option[T] {
	var ret option.Option[T]
	for i := 0; i < n; i++ {
		ret = iter.Next()
		if ret.IsNothing() {
			return ret
		}
	}
	return ret
}

// IntoIterator is an interface representing something that can turn into an Iterator.
type IntoIterator[T any] interface {
	IntoIter() Iterator[T]
}
