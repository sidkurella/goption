package set

import (
	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/maputil"
	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/sliceutil"
)

// A HashSet type backed by an HashMap.
type Set[K comparable] struct {
	m maputil.Map[K, struct{}]
}

// Returns a new set, initialized for use.
func New[K comparable]() Set[K] {
	return Set[K]{
		m: maputil.New[K, struct{}](),
	}
}

// From creates a set from an existing built-in map.
// The provided map is not cloned. The same map by reference is used for underlying storage.
func From[K comparable](m map[K]struct{}) Set[K] {
	return Set[K]{
		m: maputil.From(m),
	}
}

// FromSlice creates a set from an existing slice.
func FromSlice[K comparable](s []K) Set[K] {
	ret := New[K]()
	iterator.ForEach[K](sliceutil.Iter(s), func(k K) {
		ret.Insert(k)
	})
	return ret
}

// Into returns the underlying Go map backing this set.
// The provided map is not cloned. It is the same map by reference.
func (s Set[K]) Into() map[K]struct{} {
	return s.m.Into()
}

// Returns if the set is empty (no values) or not.
func (s Set[K]) IsEmpty() bool {
	return s.m.IsEmpty()
}

// Removes all elements from the set.
func (s Set[K]) Clear() {
	s.m.Clear()
}

// Returns the number of values in the set.
func (s Set[K]) Len() int {
	return s.m.Len()
}

// Returns if the given value is in the set.
func (s Set[K]) Contains(k K) bool {
	return s.m.ContainsKey(k)
}

// Inserts the given value into the set. Returns if the value was already in the set or not.
func (s Set[K]) Insert(k K) bool {
	return s.m.Insert(k, struct{}{}).IsSome()
}

// Append inserts the given values into the set.
func (s Set[K]) Append(keys ...K) {
	for _, k := range keys {
		s.Insert(k)
	}
}

// Deletes the value from the set. Returns if the value was already in the set or not.
func (s Set[K]) Remove(k K) bool {
	return s.m.Remove(k).IsSome()
}

// Filter removes all entries from the set that do not satisfy the given predicate.
// This modifies the set in-place.
func (s Set[K]) Filter(pred func(K) bool) {
	s.m.Filter(func(k K, _ struct{}) bool {
		return pred(k)
	})
}

// Collects all the values in the set into a slice.
// NOTE: Since there is no defined iteration order on sets, the order of elements in the slice will be random.
func (s Set[K]) Collect() []K {
	return s.m.CollectKeys()
}

// Calls the given closure for every value in the set.
// NOTE: Since there is no defined iteration order on sets, the order of elements traversed will be random.
func (s Set[K]) ForEach(f func(K)) {
	s.m.ForEach(func(k K, _ struct{}) {
		f(k)
	})
}

// Returns a new set representing the union of the two sets.
// The original sets are not modified.
func (s Set[K]) Union(s2 Set[K]) Set[K] {
	return Set[K]{
		m: maputil.New[K, struct{}]().Extend(s.m).Extend(s2.m),
	}
}

// Intersection returns a new set representing the intersection of the two sets; that is, all elements in both s and s2.
// The original sets are not modified.
func (s Set[K]) Intersection(s2 Set[K]) Set[K] {
	ret := New[K]()
	s.ForEach(func(k K) {
		if s2.Contains(k) {
			ret.Insert(k)
		}
	})
	return ret
}

// Difference returns a new set with all elements in s that are not in s2.
// The original sets are not modified.
func (s Set[K]) Difference(s2 Set[K]) Set[K] {
	ret := New[K]()
	s.ForEach(func(k K) {
		if !s2.Contains(k) {
			ret.Insert(k)
		}
	})
	return ret
}

// SymmetricDifference returns a new set with all elements in only one set.
// The original sets are not modified.
func (s Set[K]) SymmetricDifference(s2 Set[K]) Set[K] {
	ret := s.Union(s2)
	ret.ForEach(func(k K) {
		if s.Contains(k) && s2.Contains(k) {
			ret.Remove(k)
		}
	})
	return ret
}

// IsDisjoint returns if the two sets have no common elements; that is, their intersection is empty.
func (s Set[K]) IsDisjoint(s2 Set[K]) bool {
	return s.Intersection(s2).IsEmpty()
}

// IsSubset returns if the set s is a subset of s2; that is, all elements in s are contained in s2.
func (s Set[K]) IsSubset(s2 Set[K]) bool {
	return s.Difference(s2).IsEmpty()
}

// IsSuperset returns if the set s is a superset of s2; that is, all elements in s2 are contained in s.
func (s Set[K]) IsSuperset(s2 Set[K]) bool {
	return s2.Difference(s).IsEmpty()
}

// Applies f to each value in the set, returning a new set with the resultant values.
func Map[K1 comparable, K2 comparable](s Set[K1], f func(K1) K2) Set[K2] {
	return Set[K2]{
		m: maputil.Apply(s.m,
			func(k K1, _ struct{}) (K2, struct{}) {
				return f(k), struct{}{}
			},
		),
	}
}

// FilterMap applies f to each value in s.
// If f returns Some, then the entry is kept in the new set (with the new value).
// If it returns Nothing, the entry is discarded.
func FilterMap[K1 comparable, K2 comparable](s Set[K1], f func(K1) option.Option[K2]) Set[K2] {
	return Set[K2]{
		m: maputil.FilterMap(s.m,
			func(k K1, _ struct{}) option.Option[maputil.Entry[K2, struct{}]] {
				return option.Match(f(k),
					func(s option.Some[K2]) option.Option[maputil.Entry[K2, struct{}]] {
						return option.Some[maputil.Entry[K2, struct{}]]{
							Value: maputil.Entry[K2, struct{}]{
								Key:   s.Value,
								Value: struct{}{},
							},
						}
					},
					func(_ option.Nothing[K2]) option.Option[maputil.Entry[K2, struct{}]] {
						return option.Nothing[maputil.Entry[K2, struct{}]]{}
					},
				)
			},
		),
	}
}

// Fold calls f successively with each value in the  map and the current accumulator.
// The accumulator value is then updated with the new return value of f.
// NOTE: Set has no defined iteration order. It may differ in successive runs.
// Do not use Fold with any function that is order-dependent.
func Fold[K comparable, A any](s Set[K], a A, f func(A, K) A) A {
	return maputil.Fold(s.m, a,
		func(a A, k K, _ struct{}) A {
			return f(a, k)
		},
	)
}
