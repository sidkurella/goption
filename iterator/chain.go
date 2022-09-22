package iterator

import "github.com/sidkurella/goption/option"

// Implements an iterator that delivers all elements from first, then second once first is exhausted.
type chainIterator[T any] struct {
	first          Iterator[T]
	second         Iterator[T]
	firstExhausted bool
}

// Returns a new iterator that delivers all elements from first, then second once first is exhausted.
func Chain[T any](first Iterator[T], second Iterator[T]) *chainIterator[T] {
	return &chainIterator[T]{
		first:          first,
		second:         second,
		firstExhausted: false,
	}
}

// Returns the next item from the iterator. If the first iterator is empty, returns the next element from the second.
func (c *chainIterator[T]) Next() option.Option[T] {
	if !c.firstExhausted {
		firstValue := c.first.Next()
		if firstValue.IsNothing() {
			c.firstExhausted = true
		} else {
			return firstValue
		}
	}
	return c.second.Next()
}
