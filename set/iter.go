package set

import (
	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/maputil"
)

// Returns an iterator over the keys in the set.
func (s Set[K]) Iter() iterator.Iterator[K] {
	return iterator.Map[maputil.Entry[K, struct{}]](
		s.m.Iter(), func(e maputil.Entry[K, struct{}]) K {
			return e.Key
		},
	)
}
