package iterator

import "github.com/sidkurella/goption/option"

type intersperseIterator[T any] struct {
	inner      Iterator[T]
	item       T
	nextIsItem bool
}

// Creates a new iterator which places a copy of separator between adjacent items of the original iterator.
// A copy of item will be placed before the final Nothing.
// The item will only be shallow-copied. If you desire deep copying/bespoke behavior, use IntersperseWith.
func Intersperse[T any](inner Iterator[T], item T) *intersperseIterator[T] {
	return &intersperseIterator[T]{
		inner:      inner,
		item:       item,
		nextIsItem: false,
	}
}

// Returns the next item from the interspersed iterator.
func (i *intersperseIterator[T]) Next() option.Option[T] {
	var ret option.Option[T]
	if i.nextIsItem {
		ret = option.Some(i.item)
	} else {
		ret = i.inner.Next()
	}
	i.nextIsItem = !i.nextIsItem
	return ret
}
