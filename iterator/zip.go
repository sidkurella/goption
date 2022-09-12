package iterator

import (
	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/pair"
)

type zipIterator[T any, U any] struct {
	first  Iterator[T]
	second Iterator[U]
}

// ‘Zips up’ two iterators into a single iterator of pairs.
// zip() returns a new iterator that will iterate over two other iterators.
// The iterator returns a tuple where the first element comes from the first iterator,
// and the second element comes from the second iterator.
// If either iterator returns None, next from the zipped iterator will return None.
func Zip[T any, U any](first Iterator[T], second Iterator[U]) *zipIterator[T, U] {
	return &zipIterator[T, U]{
		first:  first,
		second: second,
	}
}

// Returns the next item from the zipped-up iterator.
// If either iterator returns None, next from the zipped iterator will return None.
func (z *zipIterator[T, U]) Next() option.Option[pair.Pair[T, U]] {
	valFirst := z.first.Next()
	if valFirst.IsNothing() {
		return option.Nothing[pair.Pair[T, U]]{}
	}
	valSecond := z.second.Next()
	if valSecond.IsNothing() {
		return option.Nothing[pair.Pair[T, U]]{}
	}
	return option.Some[pair.Pair[T, U]]{
		Value: pair.Pair[T, U]{
			First:  valFirst.Unwrap(),
			Second: valSecond.Unwrap(),
		},
	}
}
