package iterator

import (
	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/pair"
	"github.com/sidkurella/goption/result"
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
// Returns Ok[struct{}{}] if successful.
// Returns Err[k] if Nothing is encountered, where k is the number of elements advanced before hitting the end.
func AdvanceBy[T any](iter Iterator[T], n uint64) result.Result[struct{}, uint64] {
	for i := uint64(0); i < n; i++ {
		obj := iter.Next()
		if obj.IsNothing() {
			return result.Err[struct{}](i)
		}
	}
	return result.Ok[struct{}, uint64](struct{}{})
}

// Tests if every element of the iterator matches the predicate.
// Applies pred to each element of the iterator, until the iterator returns Nothing.
// Once the first element returns false, all() will short-circuit and exit.
// The empty iterator returns true.
func All[T any](iter Iterator[T], pred func(T) bool) bool {
	return TryFold(iter, true,
		func(_ bool, t T) result.Result[bool, struct{}] {
			// Accumulator must be true at any point here.
			if pred(t) {
				return result.Ok[bool, struct{}](true)
			}
			// Signal break from fold since the predicate is now false.
			return result.Err[bool](struct{}{})
		},
	).UnwrapOr(false)
}

// Tests if any element of the iterator matches the predicate.
// Applies pred to each element of the iterator, until the iterator returns Nothing.
// Once the first element returns true, any() will short-circuit and exit.
// The empty iterator returns false.
func Any[T any](iter Iterator[T], pred func(T) bool) bool {
	return TryFold(iter, false,
		func(b bool, t T) result.Result[bool, struct{}] {
			// Accumulator must be false at any point here.
			if !pred(t) {
				return result.Ok[bool, struct{}](false)
			}
			// Signal break from fold since the predicate is now true.
			return result.Err[bool](struct{}{})
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
		func(_ struct{}, t T) result.Result[struct{}, option.Option[T]] {
			// Haven't found it yet.
			if pred(t) {
				// Return Err to short-circuit out of here.
				return result.Err[struct{}](
					option.Some(t),
				)
			}
			// Still haven't found it.
			return result.Ok[struct{}, option.Option[T]](struct{}{})
		},
	).UnwrapErrOr(option.Nothing[T]())
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
// Short-circuits if the function returns Err, returning the Result.
func TryFold[T any, A any, E any](
	iter Iterator[T], a A, f func(A, T) result.Result[A, E],
) result.Result[A, E] {
	for item := iter.Next(); item.IsSome(); item = iter.Next() {
		res := f(a, item.Unwrap())
		if res.IsErr() {
			return res
		}
		a = res.Unwrap()
	}
	return result.Ok[A, E](a)
}

// Advances the iterator by n and returns the nth next item.
// Count starts from 0, so Nth(I, 0) returns the current element.
// The iterator is not rewinded, so preceding elements will be discarded.
// Subsequent calls (even to Nth(I, 0)) will return different values.
// Returns Nothing if n is greater or equal to the length of the iterator.
func Nth[T any](iter Iterator[T], n uint64) option.Option[T] {
	return option.AndThen(
		AdvanceBy(iter, n).Ok(),
		func(_ struct{}) option.Option[T] {
			return iter.Next()
		},
	)
}

// Last returns the final element of the iterator, before it returns Nothing.
// Returns Nothing if the iterator is empty.
func Last[T any](iter Iterator[T]) option.Option[T] {
	return Fold(iter, option.Nothing[T](),
		func(_ option.Option[T], t T) option.Option[T] {
			return option.Some(t)
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
	return Fold(iter, option.Nothing[T](),
		func(o option.Option[T], t T) option.Option[T] {
			val, first := o.Get()
			if !first || !less(t, val) { // If there is no current value, or the new is >= than the current, update.
				val = t
			}
			return option.Some(val)
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
	return Fold(iter, option.Nothing[T](),
		func(o option.Option[T], t T) option.Option[T] {
			val, first := o.Get()
			if !first || less(t, val) { // If there is no current value, or the new is < than the current, update.
				val = t
			}
			return option.Some(val)
		},
	)
}

// Collect returns all the elements of the iterator into a slice.
func Collect[T any](iter Iterator[T]) []T {
	return Fold(iter, []T{}, func(a []T, t T) []T {
		return append(a, t)
	})
}

// CollectInto collects the iterator into the given collection.
// The provided collection is modified to hold the elements in iter.
func CollectInto[T any, C Collection[T]](iter Iterator[T], collection C) C {
	collection.Append(Collect(iter)...)
	return collection
}

// TryCollect attempts to collect an iterator of Result[T, E] into []T.
// It short-circuits upon reaching the first Err variant, instead returning Err[E].
func TryCollect[T any, E any](iter Iterator[result.Result[T, E]]) result.Result[[]T, E] {
	return TryFold(iter, []T{}, func(a []T, res result.Result[T, E]) result.Result[[]T, E] {
		return result.AndThen(res, func(t T) result.Result[[]T, E] {
			return result.Ok[[]T, E](append(a, t))
		})
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
			return option.Some(i)
		}
		i++
	}
	return option.Nothing[uint64]()
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

// Collection represents a collection of elements of type T.
type Collection[T any] interface {
	// Inserts the provided elements into the collection.
	Append(...T)
}
