package maputil

import (
	"fmt"

	"github.com/sidkurella/goption/either"
	"github.com/sidkurella/goption/option"
)

// An error indicating that the given key is already in the map, with the given value.
type OccupiedError[K comparable, V any] struct {
	Key   K
	Value V
}

// String representation of this OccupiedError.
func (o OccupiedError[K, V]) Error() string {
	return fmt.Sprintf("key %v is occupied (value: %v)", o.Key, o.Value)
}

// Represents an individual entry (key-value pair) in the map.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// A HashMap type backed by an underlying Go map.
type Map[K comparable, V any] struct {
	m map[K]V
}

// Returns a new map, initialized for use.
func New[K comparable, V any]() Map[K, V] {
	return Map[K, V]{
		m: map[K]V{},
	}
}

// From creates a map from an existing built-in map.
// The provided map is not cloned. It is the same map by reference.
func From[K comparable, V any](m map[K]V) Map[K, V] {
	return Map[K, V]{
		m: m,
	}
}

// Into returns a Go map from this map.
// The provided map is not cloned. It is the same map by reference.
func (m Map[K, V]) Into() map[K]V {
	return m.m
}

// Clear deletes all key-value pairs from the map.
func (m Map[K, V]) Clear() {
	for k := range m.m {
		delete(m.m, k)
	}
}

// ContainsKey returns true if the map contains a value for the specified key.
func (m Map[K, V]) ContainsKey(k K) bool {
	_, ok := m.m[k]
	return ok
}

// IsEmpty returns if the map has no key-value pairs in it.
func (m Map[K, V]) IsEmpty() bool {
	return len(m.m) == 0
}

// Insert inserts a key-value pair into the map.
// If a value V already existed for this key, Some[V] is returned. Otherwise, Nothing is returned.
func (m Map[K, V]) Insert(k K, v V) option.Option[V] {
	val, ok := m.m[k]
	m.m[k] = v
	if ok {
		return option.Some[V]{Value: val}
	}
	return option.Nothing[V]{}
}

// Len returns the number of key-value pairs in the map.
func (m Map[K, V]) Len() int {
	return len(m.m)
}

// TryInsert tries to insert a key-value pair into the map.
// If a value V already exists for this key, OccupiedError is returned.
// Otherwise, the new value is returned.
func (m Map[K, V]) TryInsert(k K, v V) either.Either[OccupiedError[K, V], V] {
	val, ok := m.m[k]
	if !ok {
		m.m[k] = v
		return either.Right[OccupiedError[K, V], V]{Value: v}
	}
	return either.Left[OccupiedError[K, V], V]{
		Value: OccupiedError[K, V]{Key: k, Value: val},
	}
}

// Remove removes a key from the map. If there was a value present at the key, it is returned.
func (m Map[K, V]) Remove(k K) option.Option[V] {
	val, ok := m.m[k]
	if ok {
		delete(m.m, k)
		return option.Some[V]{Value: val}
	}
	return option.Nothing[V]{}
}

// Get returns the value for a given key. Returns Nothing if it does not exist.
func (m Map[K, V]) Get(k K) option.Option[V] {
	val, ok := m.m[k]
	if ok {
		return option.Some[V]{Value: val}
	}
	return option.Nothing[V]{}
}

// Filter removes all entries from the map that do not satisfy the given predicate.
// This modifies the map in-place.
func (m Map[K, V]) Filter(pred func(K, V) bool) {
	for k, v := range m.m {
		if !pred(k, v) {
			delete(m.m, k)
		}
	}
}

// Collect collects a map into a slice of key-value pairs.
// NOTE: Since there is no defined iteration order on maps, the order of elements in the slice will be random.
func (m Map[K, V]) Collect() []Entry[K, V] {
	ret := make([]Entry[K, V], 0, len(m.m))
	for k, v := range m.m {
		ret = append(ret, Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return ret
}

// CollectKeys collects a map into a slice of keys.
// NOTE: Since there is no defined iteration order on maps, the order of elements in the slice will be random.
func (m Map[K, V]) CollectKeys() []K {
	ret := make([]K, 0, len(m.m))
	for k := range m.m {
		ret = append(ret, k)
	}
	return ret
}

// CollectValues collects a map into a slice of values.
// NOTE: Since there is no defined iteration order on maps, the order of elements in the slice will be random.
func (m Map[K, V]) CollectValues() []V {
	ret := make([]V, 0, len(m.m))
	for _, v := range m.m {
		ret = append(ret, v)
	}
	return ret
}

// ForEach calls the given closure for every key-value pair in the map.
// NOTE: Since there is no defined iteration order on maps, the order of elements traversed will be random.
func (m Map[K, V]) ForEach(f func(K, V)) {
	for k, v := range m.m {
		f(k, v)
	}
}

// Extend inserts every element of m2 into the map. Returns a reference to the original map (for chaining).
// m2 is not modified; elements are copied.
// If any key is already present in m, its value will be overwritten with the value from m2.
func (m Map[K, V]) Extend(m2 Map[K, V]) Map[K, V] {
	m2.ForEach(func(k K, v V) {
		_ = m.Insert(k, v)
	})
	return m
}

// Apply applies f to each key-value pair in m, returning a new map with the resultant key-value pairs.
func Apply[K1 comparable, V1 any, K2 comparable, V2 any](m Map[K1, V1], f func(K1, V1) (K2, V2)) Map[K2, V2] {
	ret := make(map[K2]V2, len(m.m))
	for k, v := range m.m {
		k2, v2 := f(k, v)
		ret[k2] = v2
	}
	return From(ret)
}

// FilterMap applies f to each key-value pair in m.
// If f returns Some[Entry] then the entry is kept in the new map.
// If it returns Nothing, the entry is discarded.
func FilterMap[K1 comparable, V1 any, K2 comparable, V2 any](
	m Map[K1, V1], f func(K1, V1) option.Option[Entry[K2, V2]],
) Map[K2, V2] {
	ret := make(map[K2]V2, len(m.m))
	for k, v := range m.m {
		entry, ok := f(k, v).Get()
		if ok {
			ret[entry.Key] = entry.Value
		}
	}
	return From(ret)
}

// Fold calls f successively with each key-value pair in the map and the current accumulator.
// The accumulator value is then updated with the new return value of f.
// NOTE: Map has no defined iteration order. It may differ in successive runs.
// Do not use Fold with any function that is order-dependent.
func Fold[K comparable, V any, A any](m Map[K, V], a A, f func(A, K, V) A) A {
	for k, v := range m.m {
		a = f(a, k, v)
	}
	return a
}
