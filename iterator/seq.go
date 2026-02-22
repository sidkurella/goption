package iterator

import (
	"iter"

	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/pair"
)

// Closer is an interface for iterators that hold resources and need explicit cleanup.
// Close should be called when the iterator is no longer needed and has not been
// fully consumed. Calling Close on an already-closed or exhausted iterator is safe.
type Closer interface {
	Close()
}

// CloseableIterator is an iterator that holds resources requiring cleanup.
// When the iterator is not consumed to exhaustion, Close must be called to
// release the underlying resources. It is safe to call Close multiple times
// or after the iterator is exhausted.
//
// After Close is called, the iterator will return Nothing for all subsequent
// calls to Next. Any adapters (such as Map, Filter, Take, etc.) that wrap this
// iterator will also return Nothing, since they delegate to the underlying
// iterator's Next method.
type CloseableIterator[T any] interface {
	Iterator[T]
	Closer
}

// ToSeq converts an Iterator[T] to an iter.Seq[T].
// The returned sequence can be used in a for-range loop.
// The iterator will be consumed as the sequence is iterated.
func ToSeq[T any](it Iterator[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			next := it.Next()
			if next.IsNothing() {
				return
			}
			if !yield(next.Unwrap()) {
				return
			}
		}
	}
}

// FromSeq converts an iter.Seq[T] to a CloseableIterator[T].
// The sequence will be pulled from lazily as Next() is called.
//
// The returned iterator holds resources from iter.Pull. It is intended that
// Close be called via defer immediately after creating the iterator. This
// ensures resources are released even if the iterator is not fully consumed.
// Close is safe to call multiple times or after the iterator is exhausted.
//
// Important: You should defer Close even if the iterator is chained into
// other adapters like Map, Filter, Take, etc. The adapters do not propagate
// Close to the underlying iterator, so the original CloseableIterator must
// be closed explicitly.
//
// Example:
//
//	it := iterator.FromSeq(slices.Values(mySlice))
//	defer it.Close()
//	// Safe to chain into adapters after deferring Close
//	mapped := iterator.Map(it, func(x int) int { return x * 2 })
//	result := iterator.Collect(mapped)
func FromSeq[T any](seq iter.Seq[T]) CloseableIterator[T] {
	next, stop := iter.Pull(seq)
	return &seqIterator[T]{
		next: next,
		stop: stop,
	}
}

// seqIterator wraps an iter.Seq pulled iterator.
type seqIterator[T any] struct {
	next func() (T, bool)
	stop func()
	done bool
}

func (s *seqIterator[T]) Next() option.Option[T] {
	if s.done {
		return option.Nothing[T]()
	}
	val, ok := s.next()
	if !ok {
		s.done = true
		s.stop()
		return option.Nothing[T]()
	}
	return option.Some(val)
}

// Close releases the resources held by the iterator.
// After Close is called, Next will return Nothing. Any adapters wrapping this
// iterator will also return Nothing, since they delegate to this iterator.
// It is safe to call Close multiple times or after the iterator is exhausted.
func (s *seqIterator[T]) Close() {
	if !s.done {
		s.done = true
		s.stop()
	}
}

// ToSeq2 converts an Iterator[pair.Pair[K, V]] to an iter.Seq2[K, V].
// This is useful for iterators that yield key-value pairs, such as those
// created by Zip or Enumerate.
func ToSeq2[K any, V any](it Iterator[pair.Pair[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			next := it.Next()
			if next.IsNothing() {
				return
			}
			p := next.Unwrap()
			if !yield(p.First, p.Second) {
				return
			}
		}
	}
}

// FromSeq2 converts an iter.Seq2[K, V] to a CloseableIterator[pair.Pair[K, V]].
// The sequence will be pulled from lazily as Next() is called.
//
// The returned iterator holds resources from iter.Pull2. It is intended that
// Close be called via defer immediately after creating the iterator. This
// ensures resources are released even if the iterator is not fully consumed.
// Close is safe to call multiple times or after the iterator is exhausted.
//
// Important: You should defer Close even if the iterator is chained into
// other adapters like Map, Filter, Take, etc. The adapters do not propagate
// Close to the underlying iterator, so the original CloseableIterator must
// be closed explicitly.
//
// Example:
//
//	it := iterator.FromSeq2(maps.All(myMap))
//	defer it.Close()
//	// Safe to chain into adapters after deferring Close
//	filtered := iterator.Filter(it, func(p pair.Pair[K, V]) bool { return p.First != "" })
//	result := iterator.Collect(filtered)
func FromSeq2[K any, V any](seq iter.Seq2[K, V]) CloseableIterator[pair.Pair[K, V]] {
	next, stop := iter.Pull2(seq)
	return &seq2Iterator[K, V]{
		next: next,
		stop: stop,
	}
}

// seq2Iterator wraps an iter.Seq2 pulled iterator.
type seq2Iterator[K any, V any] struct {
	next func() (K, V, bool)
	stop func()
	done bool
}

func (s *seq2Iterator[K, V]) Next() option.Option[pair.Pair[K, V]] {
	if s.done {
		return option.Nothing[pair.Pair[K, V]]()
	}
	k, v, ok := s.next()
	if !ok {
		s.done = true
		s.stop()
		return option.Nothing[pair.Pair[K, V]]()
	}
	return option.Some(pair.From(k, v))
}

// Close releases the resources held by the iterator.
// After Close is called, Next will return Nothing. Any adapters wrapping this
// iterator will also return Nothing, since they delegate to this iterator.
// It is safe to call Close multiple times or after the iterator is exhausted.
func (s *seq2Iterator[K, V]) Close() {
	if !s.done {
		s.done = true
		s.stop()
	}
}
