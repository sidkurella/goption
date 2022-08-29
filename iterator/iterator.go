package iterator

import (
	"github.com/sidkurella/goption/option"
)

// Iterator returns items via successive Next calls until it has run out.
// It signals that the iterator is now empty by returning None.
type Iterator[T any] interface {
	Next() option.Option[T]
}

// IntoIterator is an interface representing something that can turn into an Iterator.
type IntoIterator[T any] interface {
	IntoIter() Iterator[T]
}
