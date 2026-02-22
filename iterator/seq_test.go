package iterator_test

import (
	"iter"
	"reflect"
	"slices"
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/pair"
)

func TestToSeq(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		it := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		seq := iterator.ToSeq(it)
		var result []int
		for v := range seq {
			result = append(result, v)
		}
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}
	})

	t.Run("empty", func(t *testing.T) {
		it := &fakeIterator{
			elements: []int{},
		}
		seq := iterator.ToSeq(it)
		var result []int
		for v := range seq {
			result = append(result, v)
		}
		if len(result) != 0 {
			t.Errorf("got %v, expected empty slice", result)
		}
	})

	t.Run("early break", func(t *testing.T) {
		it := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		seq := iterator.ToSeq(it)
		var result []int
		for v := range seq {
			result = append(result, v)
			if v == 3 {
				break
			}
		}
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}
	})
}

func TestFromSeq(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		it := iterator.FromSeq(seq)
		result := iterator.Collect(it)
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}
	})

	t.Run("empty", func(t *testing.T) {
		seq := slices.Values([]int{})
		it := iterator.FromSeq(seq)
		result := iterator.Collect(it)
		if len(result) != 0 {
			t.Errorf("got %v, expected empty slice", result)
		}
	})

	t.Run("partial consumption", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		it := iterator.FromSeq(seq)

		// Consume first two elements
		first := it.Next()
		if first.Unwrap() != 1 {
			t.Errorf("first element: got %v, expected 1", first.Unwrap())
		}
		second := it.Next()
		if second.Unwrap() != 2 {
			t.Errorf("second element: got %v, expected 2", second.Unwrap())
		}

		// Collect the rest
		rest := iterator.Collect(it)
		expected := []int{3, 4, 5}
		if !reflect.DeepEqual(rest, expected) {
			t.Errorf("rest: got %v, expected %v", rest, expected)
		}
	})

	t.Run("nothing after exhaustion", func(t *testing.T) {
		seq := slices.Values([]int{1})
		it := iterator.FromSeq(seq)

		first := it.Next()
		if first.Unwrap() != 1 {
			t.Errorf("first element: got %v, expected 1", first.Unwrap())
		}

		// Should return Nothing repeatedly
		for i := 0; i < 3; i++ {
			next := it.Next()
			if !next.IsNothing() {
				t.Errorf("iteration %d: expected Nothing, got %v", i, next)
			}
		}
	})
}

func TestToSeq2(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		it := &fakePairIterator{
			elements: []pair.Pair[int, string]{
				{First: 1, Second: "a"},
				{First: 2, Second: "b"},
				{First: 3, Second: "c"},
			},
		}
		seq := iterator.ToSeq2(it)
		var keys []int
		var values []string
		for k, v := range seq {
			keys = append(keys, k)
			values = append(values, v)
		}
		expectedKeys := []int{1, 2, 3}
		expectedValues := []string{"a", "b", "c"}
		if !reflect.DeepEqual(keys, expectedKeys) {
			t.Errorf("keys: got %v, expected %v", keys, expectedKeys)
		}
		if !reflect.DeepEqual(values, expectedValues) {
			t.Errorf("values: got %v, expected %v", values, expectedValues)
		}
	})

	t.Run("empty", func(t *testing.T) {
		it := &fakePairIterator{
			elements: []pair.Pair[int, string]{},
		}
		seq := iterator.ToSeq2(it)
		count := 0
		for range seq {
			count++
		}
		if count != 0 {
			t.Errorf("got %d iterations, expected 0", count)
		}
	})

	t.Run("early break", func(t *testing.T) {
		it := &fakePairIterator{
			elements: []pair.Pair[int, string]{
				{First: 1, Second: "a"},
				{First: 2, Second: "b"},
				{First: 3, Second: "c"},
			},
		}
		seq := iterator.ToSeq2(it)
		var keys []int
		for k, _ := range seq {
			keys = append(keys, k)
			if k == 2 {
				break
			}
		}
		expected := []int{1, 2}
		if !reflect.DeepEqual(keys, expected) {
			t.Errorf("got %v, expected %v", keys, expected)
		}
	})

	t.Run("with enumerate", func(t *testing.T) {
		it := &fakeStringIterator{
			elements: []string{"a", "b", "c"},
		}
		enumIter := iterator.Enumerate[string](it)
		seq := iterator.ToSeq2(enumIter)
		var indices []int
		var values []string
		for idx, val := range seq {
			indices = append(indices, idx)
			values = append(values, val)
		}
		expectedIndices := []int{0, 1, 2}
		expectedValues := []string{"a", "b", "c"}
		if !reflect.DeepEqual(indices, expectedIndices) {
			t.Errorf("indices: got %v, expected %v", indices, expectedIndices)
		}
		if !reflect.DeepEqual(values, expectedValues) {
			t.Errorf("values: got %v, expected %v", values, expectedValues)
		}
	})

	t.Run("with zip", func(t *testing.T) {
		it1 := &fakeIterator{
			elements: []int{1, 2, 3},
		}
		it2 := &fakeStringIterator{
			elements: []string{"a", "b", "c"},
		}
		zipIter := iterator.Zip[int, string](it1, it2)
		seq := iterator.ToSeq2(zipIter)
		var keys []int
		var values []string
		for k, v := range seq {
			keys = append(keys, k)
			values = append(values, v)
		}
		expectedKeys := []int{1, 2, 3}
		expectedValues := []string{"a", "b", "c"}
		if !reflect.DeepEqual(keys, expectedKeys) {
			t.Errorf("keys: got %v, expected %v", keys, expectedKeys)
		}
		if !reflect.DeepEqual(values, expectedValues) {
			t.Errorf("values: got %v, expected %v", values, expectedValues)
		}
	})
}

func TestFromSeq2(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		// Use a deterministic sequence for testing
		seq := func(yield func(string, int) bool) {
			for _, k := range []string{"a", "b", "c"} {
				if !yield(k, m[k]) {
					return
				}
			}
		}
		it := iterator.FromSeq2(seq)
		result := iterator.Collect(it)
		expected := []pair.Pair[string, int]{
			{First: "a", Second: 1},
			{First: "b", Second: 2},
			{First: "c", Second: 3},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}
	})

	t.Run("empty", func(t *testing.T) {
		seq := func(yield func(string, int) bool) {}
		it := iterator.FromSeq2[string, int](seq)
		result := iterator.Collect(it)
		if len(result) != 0 {
			t.Errorf("got %v, expected empty slice", result)
		}
	})

	t.Run("partial consumption", func(t *testing.T) {
		seq := func(yield func(int, string) bool) {
			pairs := []struct {
				k int
				v string
			}{{1, "a"}, {2, "b"}, {3, "c"}}
			for _, p := range pairs {
				if !yield(p.k, p.v) {
					return
				}
			}
		}
		it := iterator.FromSeq2(seq)

		// Consume first element
		first := it.Next()
		p := first.Unwrap()
		if p.First != 1 || p.Second != "a" {
			t.Errorf("first element: got %v, expected {1, a}", p)
		}

		// Collect the rest
		rest := iterator.Collect(it)
		expected := []pair.Pair[int, string]{
			{First: 2, Second: "b"},
			{First: 3, Second: "c"},
		}
		if !reflect.DeepEqual(rest, expected) {
			t.Errorf("rest: got %v, expected %v", rest, expected)
		}
	})

	t.Run("nothing after exhaustion", func(t *testing.T) {
		seq := func(yield func(int, string) bool) {
			yield(1, "a")
		}
		it := iterator.FromSeq2(seq)

		first := it.Next()
		p := first.Unwrap()
		if p.First != 1 || p.Second != "a" {
			t.Errorf("first element: got %v, expected {1, a}", p)
		}

		// Should return Nothing repeatedly
		for i := 0; i < 3; i++ {
			next := it.Next()
			if !next.IsNothing() {
				t.Errorf("iteration %d: expected Nothing, got %v", i, next)
			}
		}
	})
}

func TestRoundTrip(t *testing.T) {
	t.Run("ToSeq then FromSeq", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5}
		it := &fakeIterator{
			elements: original,
		}
		seq := iterator.ToSeq(it)
		it2 := iterator.FromSeq(seq)
		result := iterator.Collect(it2)
		if !reflect.DeepEqual(result, original) {
			t.Errorf("got %v, expected %v", result, original)
		}
	})

	t.Run("FromSeq then ToSeq", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5}
		seq := slices.Values(original)
		it := iterator.FromSeq(seq)
		seq2 := iterator.ToSeq(it)
		var result []int
		for v := range seq2 {
			result = append(result, v)
		}
		if !reflect.DeepEqual(result, original) {
			t.Errorf("got %v, expected %v", result, original)
		}
	})

	t.Run("ToSeq2 then FromSeq2", func(t *testing.T) {
		original := []pair.Pair[int, string]{
			{First: 1, Second: "a"},
			{First: 2, Second: "b"},
			{First: 3, Second: "c"},
		}
		it := &fakePairIterator{
			elements: original,
		}
		seq := iterator.ToSeq2(it)
		it2 := iterator.FromSeq2(seq)
		result := iterator.Collect(it2)
		if !reflect.DeepEqual(result, original) {
			t.Errorf("got %v, expected %v", result, original)
		}
	})

	t.Run("FromSeq2 then ToSeq2", func(t *testing.T) {
		original := []pair.Pair[string, int]{
			{First: "a", Second: 1},
			{First: "b", Second: 2},
			{First: "c", Second: 3},
		}
		seq := func(yield func(string, int) bool) {
			for _, p := range original {
				if !yield(p.First, p.Second) {
					return
				}
			}
		}
		it := iterator.FromSeq2(seq)
		defer it.Close()
		seq2 := iterator.ToSeq2(it)
		var result []pair.Pair[string, int]
		for k, v := range seq2 {
			result = append(result, pair.From(k, v))
		}
		if !reflect.DeepEqual(result, original) {
			t.Errorf("got %v, expected %v", result, original)
		}
	})
}

func TestFromSeqWithStdlibIter(t *testing.T) {
	t.Run("slices.All", func(t *testing.T) {
		original := []string{"a", "b", "c"}
		seq := slices.All(original)
		it := iterator.FromSeq2(seq)
		result := iterator.Collect(it)
		expected := []pair.Pair[int, string]{
			{First: 0, Second: "a"},
			{First: 1, Second: "b"},
			{First: 2, Second: "c"},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}
	})
}

// Test that the conversion works with iterator adapters.
func TestSeqWithAdapters(t *testing.T) {
	t.Run("FromSeq with Map", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		it := iterator.FromSeq(seq)
		mapped := iterator.Map(it, func(x int) int { return x * 2 })
		result := iterator.Collect(mapped)
		expected := []int{2, 4, 6, 8, 10}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}
	})

	t.Run("FromSeq with Filter", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		it := iterator.FromSeq(seq)
		filtered := iterator.Filter(it, func(x int) bool { return x%2 == 0 })
		result := iterator.Collect(filtered)
		expected := []int{2, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}
	})

	t.Run("ToSeq from filtered iterator", func(t *testing.T) {
		it := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		filtered := iterator.Filter[int](it, func(x int) bool { return x%2 == 1 })
		seq := iterator.ToSeq(filtered)
		var result []int
		for v := range seq {
			result = append(result, v)
		}
		expected := []int{1, 3, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}
	})
}

// Compile-time check that FromSeq2 works with the standard library map iteration.
func TestFromSeq2WithMaps(t *testing.T) {
	t.Run("maps.All", func(t *testing.T) {
		// Create a deterministic test using a custom Seq2
		seq := func(yield func(string, int) bool) {
			data := []struct {
				k string
				v int
			}{{"x", 10}, {"y", 20}, {"z", 30}}
			for _, d := range data {
				if !yield(d.k, d.v) {
					return
				}
			}
		}
		it := iterator.FromSeq2(seq)

		// Use the iterator with standard operations
		count := iterator.Count(it)
		if count != 3 {
			t.Errorf("got count %d, expected 3", count)
		}
	})
}

// Verify that iter package types are used correctly.
var (
	_ iter.Seq[int]          = iterator.ToSeq[int](nil)
	_ iter.Seq2[int, string] = iterator.ToSeq2[int, string](nil)
)

// Verify that FromSeq and FromSeq2 return CloseableIterator.
var (
	_ iterator.CloseableIterator[int]                 = iterator.FromSeq(slices.Values([]int{}))
	_ iterator.CloseableIterator[pair.Pair[int, int]] = iterator.FromSeq2(slices.All([]int{}))
)

func TestCloseableIterator(t *testing.T) {
	t.Run("FromSeq Close before exhaustion", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		it := iterator.FromSeq(seq)

		// Consume partially
		first := it.Next()
		if first.Unwrap() != 1 {
			t.Errorf("first element: got %v, expected 1", first.Unwrap())
		}

		// Close early
		it.Close()

		// Should return Nothing after close
		next := it.Next()
		if !next.IsNothing() {
			t.Errorf("expected Nothing after Close, got %v", next)
		}
	})

	t.Run("FromSeq Close after exhaustion is safe", func(t *testing.T) {
		seq := slices.Values([]int{1})
		it := iterator.FromSeq(seq)

		// Fully consume
		_ = it.Next()
		exhausted := it.Next()
		if !exhausted.IsNothing() {
			t.Errorf("expected Nothing after exhaustion, got %v", exhausted)
		}

		// Close after exhaustion should not panic
		it.Close()
		it.Close() // Multiple closes should be safe
	})

	t.Run("FromSeq Close multiple times is safe", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		it := iterator.FromSeq(seq)

		// Close multiple times should not panic
		it.Close()
		it.Close()
		it.Close()

		// Should still return Nothing
		next := it.Next()
		if !next.IsNothing() {
			t.Errorf("expected Nothing after Close, got %v", next)
		}
	})

	t.Run("FromSeq2 Close before exhaustion", func(t *testing.T) {
		seq := slices.All([]string{"a", "b", "c"})
		it := iterator.FromSeq2(seq)

		// Consume partially
		first := it.Next()
		p := first.Unwrap()
		if p.First != 0 || p.Second != "a" {
			t.Errorf("first element: got %v, expected {0, a}", p)
		}

		// Close early
		it.Close()

		// Should return Nothing after close
		next := it.Next()
		if !next.IsNothing() {
			t.Errorf("expected Nothing after Close, got %v", next)
		}
	})

	t.Run("FromSeq2 Close after exhaustion is safe", func(t *testing.T) {
		seq := func(yield func(int, string) bool) {
			yield(1, "a")
		}
		it := iterator.FromSeq2(seq)

		// Fully consume
		_ = it.Next()
		exhausted := it.Next()
		if !exhausted.IsNothing() {
			t.Errorf("expected Nothing after exhaustion, got %v", exhausted)
		}

		// Close after exhaustion should not panic
		it.Close()
		it.Close() // Multiple closes should be safe
	})

	t.Run("FromSeq2 Close multiple times is safe", func(t *testing.T) {
		seq := slices.All([]int{1, 2, 3})
		it := iterator.FromSeq2(seq)

		// Close multiple times should not panic
		it.Close()
		it.Close()
		it.Close()

		// Should still return Nothing
		next := it.Next()
		if !next.IsNothing() {
			t.Errorf("expected Nothing after Close, got %v", next)
		}
	})

	t.Run("defer Close pattern", func(t *testing.T) {
		result := func() []int {
			seq := slices.Values([]int{1, 2, 3, 4, 5})
			it := iterator.FromSeq(seq)
			defer it.Close()

			// Only take first 2 elements
			var out []int
			for i := 0; i < 2; i++ {
				next := it.Next()
				if next.IsSome() {
					out = append(out, next.Unwrap())
				}
			}
			return out
		}()

		expected := []int{1, 2}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}
	})
}

func TestCloserInterface(t *testing.T) {
	t.Run("FromSeq implements Closer", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		it := iterator.FromSeq(seq)

		// Should be assignable to Closer
		var closer iterator.Closer = it
		closer.Close()
	})

	t.Run("FromSeq2 implements Closer", func(t *testing.T) {
		seq := slices.All([]int{1, 2, 3})
		it := iterator.FromSeq2(seq)

		// Should be assignable to Closer
		var closer iterator.Closer = it
		closer.Close()
	})
}

func TestCloseAfterAdapter(t *testing.T) {
	t.Run("FromSeq Close after Map adapter partial consumption", func(t *testing.T) {
		closed := false
		seq := func(yield func(int) bool) {
			defer func() { closed = true }()
			for i := 1; i <= 10; i++ {
				if !yield(i) {
					return
				}
			}
		}
		it := iterator.FromSeq(seq)
		defer it.Close()

		// Chain into Map adapter
		mapped := iterator.Map(it, func(x int) int { return x * 2 })

		// Only consume 2 elements through the adapter
		first := mapped.Next()
		if first.Unwrap() != 2 {
			t.Errorf("first element: got %v, expected 2", first.Unwrap())
		}
		second := mapped.Next()
		if second.Unwrap() != 4 {
			t.Errorf("second element: got %v, expected 4", second.Unwrap())
		}

		// Close the original iterator (deferred, but we can also call explicitly)
		it.Close()

		// Verify resources were released
		if !closed {
			t.Error("expected sequence cleanup to have run after Close")
		}

		// Further calls to the original iterator should return Nothing
		next := it.Next()
		if !next.IsNothing() {
			t.Errorf("expected Nothing after Close, got %v", next)
		}

		// Further calls to the adapted iterator should also return Nothing
		adaptedNext := mapped.Next()
		if !adaptedNext.IsNothing() {
			t.Errorf("expected adapted iterator to return Nothing after Close, got %v", adaptedNext)
		}
	})

	t.Run("FromSeq Close after Filter adapter partial consumption", func(t *testing.T) {
		closed := false
		seq := func(yield func(int) bool) {
			defer func() { closed = true }()
			for i := 1; i <= 10; i++ {
				if !yield(i) {
					return
				}
			}
		}
		it := iterator.FromSeq(seq)
		defer it.Close()

		// Chain into Filter adapter
		filtered := iterator.Filter(it, func(x int) bool { return x%2 == 0 })

		// Only consume 1 element through the adapter
		first := filtered.Next()
		if first.Unwrap() != 2 {
			t.Errorf("first element: got %v, expected 2", first.Unwrap())
		}

		// Close the original iterator
		it.Close()

		// Verify resources were released
		if !closed {
			t.Error("expected sequence cleanup to have run after Close")
		}

		// Further calls to the adapted iterator should return Nothing
		adaptedNext := filtered.Next()
		if !adaptedNext.IsNothing() {
			t.Errorf("expected adapted iterator to return Nothing after Close, got %v", adaptedNext)
		}
	})

	t.Run("FromSeq Close after Take adapter", func(t *testing.T) {
		closed := false
		seq := func(yield func(int) bool) {
			defer func() { closed = true }()
			for i := 1; i <= 100; i++ {
				if !yield(i) {
					return
				}
			}
		}
		it := iterator.FromSeq(seq)
		defer it.Close()

		// Chain into Take adapter - only want first 3
		taken := iterator.Take(it, 3)

		// Fully consume the Take adapter (but not the underlying sequence)
		result := iterator.Collect(taken)
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}

		// Close the original iterator
		it.Close()

		// Verify resources were released
		if !closed {
			t.Error("expected sequence cleanup to have run after Close")
		}
	})

	t.Run("FromSeq Close after chained adapters", func(t *testing.T) {
		closed := false
		seq := func(yield func(int) bool) {
			defer func() { closed = true }()
			for i := 1; i <= 20; i++ {
				if !yield(i) {
					return
				}
			}
		}
		it := iterator.FromSeq(seq)
		defer it.Close()

		// Chain multiple adapters: Filter -> Map -> Take
		filtered := iterator.Filter(it, func(x int) bool { return x%2 == 0 })
		mapped := iterator.Map(filtered, func(x int) int { return x * 10 })
		taken := iterator.Take(mapped, 2)

		// Consume through the chain
		result := iterator.Collect(taken)
		expected := []int{20, 40} // 2*10, 4*10
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}

		// Close the original iterator
		it.Close()

		// Verify resources were released
		if !closed {
			t.Error("expected sequence cleanup to have run after Close")
		}

		// Further calls to all adapters in the chain should return Nothing
		filteredNext := filtered.Next()
		if !filteredNext.IsNothing() {
			t.Errorf("expected filtered iterator to return Nothing after Close, got %v", filteredNext)
		}
		mappedNext := mapped.Next()
		if !mappedNext.IsNothing() {
			t.Errorf("expected mapped iterator to return Nothing after Close, got %v", mappedNext)
		}
		takenNext := taken.Next()
		if !takenNext.IsNothing() {
			t.Errorf("expected taken iterator to return Nothing after Close, got %v", takenNext)
		}
	})

	t.Run("FromSeq2 Close after Map adapter", func(t *testing.T) {
		closed := false
		seq := func(yield func(int, string) bool) {
			defer func() { closed = true }()
			data := []string{"a", "b", "c", "d", "e"}
			for i, v := range data {
				if !yield(i, v) {
					return
				}
			}
		}
		it := iterator.FromSeq2(seq)
		defer it.Close()

		// Chain into Map adapter
		mapped := iterator.Map(it, func(p pair.Pair[int, string]) string {
			return p.Second + p.Second
		})

		// Only consume 2 elements
		first := mapped.Next()
		if first.Unwrap() != "aa" {
			t.Errorf("first element: got %v, expected aa", first.Unwrap())
		}
		second := mapped.Next()
		if second.Unwrap() != "bb" {
			t.Errorf("second element: got %v, expected bb", second.Unwrap())
		}

		// Close the original iterator
		it.Close()

		// Verify resources were released
		if !closed {
			t.Error("expected sequence cleanup to have run after Close")
		}

		// Further calls to the original iterator should return Nothing
		next := it.Next()
		if !next.IsNothing() {
			t.Errorf("expected Nothing after Close, got %v", next)
		}

		// Further calls to the adapted iterator should also return Nothing
		adaptedNext := mapped.Next()
		if !adaptedNext.IsNothing() {
			t.Errorf("expected adapted iterator to return Nothing after Close, got %v", adaptedNext)
		}
	})

	t.Run("FromSeq defer pattern with adapter does cleanup", func(t *testing.T) {
		closed := false
		result := func() []int {
			seq := func(yield func(int) bool) {
				defer func() { closed = true }()
				for i := 1; i <= 100; i++ {
					if !yield(i) {
						return
					}
				}
			}
			it := iterator.FromSeq(seq)
			defer it.Close() // This is the key pattern

			// Use adapters
			mapped := iterator.Map(it, func(x int) int { return x * 2 })
			taken := iterator.Take(mapped, 3)
			return iterator.Collect(taken)
		}()

		expected := []int{2, 4, 6}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, expected %v", result, expected)
		}

		// Verify cleanup happened when function returned
		if !closed {
			t.Error("expected sequence cleanup to have run after function return")
		}
	})
}
