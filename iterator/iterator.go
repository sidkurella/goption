package iterator

import (
	"github.com/sidkurella/goption/either"
	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/pair"
	"golang.org/x/exp/constraints"
)

// Iterator returns items via successive Next calls until it has run out.
// It signals that the iterator is now empty by returning Nothing.
// NOTE: After an iterator returns Nothing for the first time, there is no guarantee that successive calls to
// NOTE: Iterator.Next() will continue to return Nothing. If you require this, use Fuse().
type Iterator[T any] interface {
	Next() option.Option[T]
}

// Advances the iterator by n elements.
// This method will eagerly skip n elements by calling next up to n times until Nothing is encountered.
// Returns First[struct{}{}] if successful.
// Returns Second[k] if Nothing is encountered, where k is the number of elements advanced before hitting the end.
func AdvanceBy[T any](iter Iterator[T], n uint64) either.Either[struct{}, uint64] {
	for i := uint64(0); i < n; i++ {
		obj := iter.Next()
		if obj.IsNothing() {
			return either.Second[struct{}, uint64]{Value: i}
		}
	}
	return either.First[struct{}, uint64]{}
}

// Tests if every element of the iterator matches the predicate.
// Applies pred to each element of the iterator, until the iterator returns Nothing.
// Once the first element returns false, all() will short-circuit and exit.
// The empty iterator returns true.
func All[T any](iter Iterator[T], pred func(T) bool) bool {
	return TryFold(iter, true,
		func(_ bool, t T) either.Either[bool, struct{}] {
			// Accumulator must be true at any point here.
			if pred(t) {
				return either.First[bool, struct{}]{Value: true}
			}
			// Signal break from fold since the predicate is now false.
			return either.Second[bool, struct{}]{}
		},
	).UnwrapOr(false)
}

// Tests if any element of the iterator matches the predicate.
// Applies pred to each element of the iterator, until the iterator returns Nothing.
// Once the first element returns true, any() will short-circuit and exit.
// The empty iterator returns false.
func Any[T any](iter Iterator[T], pred func(T) bool) bool {
	return TryFold(iter, false,
		func(b bool, t T) either.Either[bool, struct{}] {
			// Accumulator must be false at any point here.
			if !pred(t) {
				return either.First[bool, struct{}]{Value: false}
			}
			// Signal break from fold since the predicate is now true.
			return either.Second[bool, struct{}]{}
		},
	).UnwrapOr(true)
}

// Consumes the iterator, counting the number of elements until it was exhausted.
// Next() will be called at least once even if the iterator has no elements.
func Count[T any](iter Iterator[T]) uint64 {
	return Fold(iter, uint64(0), func(a uint64, _ T) uint64 {
		return a + 1
	})
}

// Find searches for the first element of the iterator that satisfies the predicate.
// Returns Some[T] for the first element that returns true. Short-circuits upon finding the first true element.
// If no element satisfies the predicate, returns Nothing.
func Find[T any](iter Iterator[T], pred func(T) bool) option.Option[T] {
	return TryFold(iter, struct{}{},
		func(_ struct{}, t T) either.Either[struct{}, option.Option[T]] {
			// Haven't found it yet.
			if pred(t) {
				// Return Second to short-circuit out of here.
				return either.Second[struct{}, option.Option[T]]{
					Value: option.Some[T]{Value: t},
				}
			}
			// Still haven't found it.
			return either.First[struct{}, option.Option[T]]{}
		},
	).UnwrapSecondOr(option.Nothing[T]{})
}

// Folds every element into an accumulator by applying an operation, returning the final either.
// The entire iterator will be consumed by this.
func Fold[T any, A any](iter Iterator[T], a A, f func(A, T) A) A {
	for item := iter.Next(); item.IsSome(); item = iter.Next() {
		a = f(a, item.Unwrap())
	}
	return a
}

// Runs the given closure on each element of the iterator.
func ForEach[T any](iter Iterator[T], f func(T)) {
	Fold(iter, struct{}{},
		func(_ struct{}, t T) struct{} {
			f(t)
			return struct{}{}
		},
	)
}

// Tries to fold every element into an accumulator by applying an operation, returning the final either.
// Short-circuits if the function returns Second, returning the Either.
func TryFold[T any, A any, E any](
	iter Iterator[T], a A, f func(A, T) either.Either[A, E],
) either.Either[A, E] {
	for item := iter.Next(); item.IsSome(); item = iter.Next() {
		res := f(a, item.Unwrap())
		if res.IsSecond() {
			return res
		}
		a = res.Unwrap()
	}
	return either.First[A, E]{Value: a}
}

// Advances the iterator by n and returns the nth next item.
// Count starts from 0, so Nth(I, 0) returns the current element.
// The iterator is not rewinded, so preceding elements will be discarded.
// Subsequent calls (even to Nth(I, 0)) will return different values.
// Returns Nothing if n is greater or equal to the length of the iterator.
func Nth[T any](iter Iterator[T], n uint64) option.Option[T] {
	return option.AndThen(
		AdvanceBy(iter, n).First(),
		func(_ struct{}) option.Option[T] {
			return iter.Next()
		},
	)
}

// Last returns the final element of the iterator, before it returns Nothing.
// Returns Nothing if the iterator is empty.
func Last[T any](iter Iterator[T]) option.Option[T] {
	return Fold[T, option.Option[T]](iter, option.Nothing[T]{},
		func(_ option.Option[T], t T) option.Option[T] {
			return option.Some[T]{Value: t}
		},
	)
}

// Max returns the maximum element of the iterator.
// Returns the last element if multiple elements are equally maximal.
// Returns Nothing if the iterator is empty.
func Max[T constraints.Ordered](iter Iterator[T]) option.Option[T] {
	return MaxBy(iter, func(t1 T, t2 T) bool {
		return t1 < t2
	})
}

// MaxBy returns the maximum element of the iterator with respect to the specified less function.
// less(a, b) should return true if a is less than b, and false otherwise.
// Returns the last element if multiple elements are equally maximal.
// Returns Nothing if the iterator is empty.
func MaxBy[T any](iter Iterator[T], less func(T, T) bool) option.Option[T] {
	return Fold[T, option.Option[T]](iter, option.Nothing[T]{},
		func(o option.Option[T], t T) option.Option[T] {
			val, first := o.Get()
			if !first || !less(t, val) { // If there is no current value, or the new is >= than the current, update.
				val = t
			}
			return option.Some[T]{Value: val}
		},
	)
}

// Min returns the minimum element of the iterator.
// Returns the first element if multiple elements are equally minimal.
// Returns Nothing if the iterator is empty.
func Min[T constraints.Ordered](iter Iterator[T]) option.Option[T] {
	return MinBy(iter, func(t1 T, t2 T) bool {
		return t1 < t2
	})
}

// MinBy returns the minimum element of the iterator with respect to the specified less function.
// less(a, b) should return true if a is less than b, and false otherwise.
// Returns the first element if multiple elements are equally minimal.
// Returns Nothing if the iterator is empty.
func MinBy[T any](iter Iterator[T], less func(T, T) bool) option.Option[T] {
	return Fold[T, option.Option[T]](iter, option.Nothing[T]{},
		func(o option.Option[T], t T) option.Option[T] {
			val, first := o.Get()
			if !first || less(t, val) { // If there is no current value, or the new is < than the current, update.
				val = t
			}
			return option.Some[T]{Value: val}
		},
	)
}

// Collect returns all the elements of the iterator into a slice.
func Collect[T any](iter Iterator[T]) []T {
	return Fold(iter, []T{}, func(a []T, t T) []T {
		return append(a, t)
	})
}

// Consumes an iterator, producing two lists from it.
// The first contains all the elements the predicate returned true for, and the second, false.
func Partition[T any](iter Iterator[T], f func(T) bool) ([]T, []T) {
	trueList := []T{}
	falseList := []T{}
	ForEach(iter, func(t T) {
		if f(t) {
			trueList = append(trueList, t)
		} else {
			falseList = append(falseList, t)
		}
	})
	return trueList, falseList
}

// Searches for an element in an iterator, returning its index.
// Returns Nothing if it was not found.
// Consumes the iterator up to the item that returned true.
func Position[T any](iter Iterator[T], pred func(T) bool) option.Option[uint64] {
	i := uint64(0)
	for item := iter.Next(); item.IsSome(); item = iter.Next() {
		if pred(item.Unwrap()) {
			return option.Some[uint64]{Value: i}
		}
		i++
	}
	return option.Nothing[uint64]{}
}

// Consumes an entire iterator of pairs, producing two collections, for the first and second elements respectively.
func Unzip[T any, U any](iter Iterator[pair.Pair[T, U]]) ([]T, []U) {
	firstList := []T{}
	secondList := []U{}
	ForEach(iter, func(t pair.Pair[T, U]) {
		firstList = append(firstList, t.First)
		secondList = append(secondList, t.Second)
	})
	return firstList, secondList
}

// IntoIterator is an interface representing something that can turn into an Iterator.
type IntoIterator[T any] interface {
	IntoIter() Iterator[T]
}
