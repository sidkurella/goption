package maputil

import (
	"reflect"

	"github.com/sidkurella/goption/option"
)

type mapIterator[K comparable, V any] struct {
	inner *reflect.MapIter
}

// Returns an iterator the entries in the map.
func (m Map[K, V]) Iter() *mapIterator[K, V] {
	return &mapIterator[K, V]{
		inner: reflect.ValueOf(m.m).MapRange(),
	}
}

func (m *mapIterator[K, V]) Next() option.Option[Entry[K, V]] {
	if !m.inner.Next() {
		return option.Nothing[Entry[K, V]]{}
	}
	k := m.inner.Key().Interface().(K)
	v := m.inner.Value().Interface().(V)

	return option.Some[Entry[K, V]]{
		Value: Entry[K, V]{
			Key:   k,
			Value: v,
		},
	}
}
