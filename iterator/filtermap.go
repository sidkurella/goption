package iterator

import "github.com/sidkurella/goption/option"

// An iterator that return elements that are both filtered and mapped by pred.
type filterMapIterator[T any, U any] struct {
	inner Iterator[T]
	pred  func(T) option.Option[U]
}

// Creates an iterator that both filters and maps.
// The returned iterator yields only the values for which the supplied closure returns Some(value).
func FilterMap[T any, U any](iter Iterator[T], pred func(T) option.Option[U]) *filterMapIterator[T, U] {
	return &filterMapIterator[T, U]{
		inner: iter,
		pred:  pred,
	}
}

func (f *filterMapIterator[T, U]) Next() option.Option[U] {
	for item := f.inner.Next(); item.IsSome(); item = f.inner.Next() {
		res := f.pred(item.Unwrap())
		if res.IsSome() {
			return res
		}
	}
	return option.Nothing[U]()
}
