package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/either"
	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/pair"
)

type fakeIterator struct {
	elements []int
	i        int
}

func (f *fakeIterator) Next() option.Option[int] {
	if f.i < len(f.elements) {
		ret := option.Some[int]{Value: f.elements[f.i]}
		f.i++
		return ret
	}
	return option.Nothing[int]{}
}

type fakeStringIterator struct {
	elements []string
	i        int
}

func (f *fakeStringIterator) Next() option.Option[string] {
	if f.i < len(f.elements) {
		ret := option.Some[string]{Value: f.elements[f.i]}
		f.i++
		return ret
	}
	return option.Nothing[string]{}
}

type fakePairIterator struct {
	elements []pair.Pair[int, string]
	i        int
}

func (f *fakePairIterator) Next() option.Option[pair.Pair[int, string]] {
	if f.i < len(f.elements) {
		ret := option.Some[pair.Pair[int, string]]{Value: f.elements[f.i]}
		f.i++
		return ret
	}
	return option.Nothing[pair.Pair[int, string]]{}
}

func TestAdvanceBy(t *testing.T) {
	iter := &fakeIterator{
		elements: []int{1, 2, 3, 4, 5},
	}
	if iter.Next().Unwrap() != 1 {
		t.Fail()
	}
	res := iterator.AdvanceBy[int](iter, 2)
	if res.Unwrap() != struct{}{} {
		t.Fail()
	}
	if iter.Next().Unwrap() != 4 {
		t.Fail()
	}
	res = iterator.AdvanceBy[int](iter, 3)
	if res.UnwrapSecond() != uint64(1) {
		t.Fail()
	}
}

func TestAll(t *testing.T) {
	t.Run("all succeeds", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		all := iterator.All[int](iter, func(t int) bool {
			return t < 6
		})
		if !all {
			t.Fail()
		}
	})
	t.Run("all fails", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		all := iterator.All[int](iter, func(t int) bool {
			return t < 4
		})
		if all {
			t.Fail()
		}
		if iter.Next().Unwrap() != 5 { // Iterator is not fully consumed; should be one after failing element.
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{},
		}
		all := iterator.All[int](iter, func(t int) bool {
			return t < 4
		})
		if !all {
			t.Fail()
		}
	})
}

func TestAny(t *testing.T) {
	t.Run("any succeeds", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		any := iterator.Any[int](iter, func(t int) bool {
			return t > 3
		})
		if !any {
			t.Fail()
		}
		if iter.Next().Unwrap() != 5 { // Iterator is not fully consumed; should be one after succeeding element.
			t.Fail()
		}
	})
	t.Run("any fails", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		any := iterator.Any[int](iter, func(t int) bool {
			return t > 10
		})
		if any {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{},
		}
		any := iterator.Any[int](iter, func(t int) bool {
			return t < 4
		})
		if any {
			t.Fail()
		}
	})
}

func TestCount(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		if iterator.Count[int](iter) != 5 {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{},
		}
		if iterator.Count[int](iter) != 0 {
			t.Fail()
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		ret := iterator.Find[int](iter, func(t int) bool {
			return t == 3
		})
		if ret.Unwrap() != 3 {
			t.Fail()
		}
		if iter.Next().Unwrap() != 4 { // Iterator should continue after first found item.
			t.Fail()
		}
	})
}

func TestFold(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		ret := iterator.Fold[int](iter, 1,
			func(a int, t int) int {
				return a * t
			},
		)
		if ret != 120 {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{},
		}
		ret := iterator.Fold[int](iter, 0,
			func(a int, t int) int {
				return a + t
			},
		)
		if ret != 0 {
			t.Fail()
		}
	})
}

func TestForEach(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		calls := map[int]struct{}{}
		expected := map[int]struct{}{
			1: {},
			2: {},
			3: {},
			4: {},
			5: {},
		}
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		iterator.ForEach[int](iter,
			func(t int) {
				calls[t] = struct{}{}
			},
		)
		if !reflect.DeepEqual(calls, expected) {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		calls := map[int]struct{}{}
		expected := map[int]struct{}{}
		iter := &fakeIterator{
			elements: []int{},
		}
		iterator.ForEach[int](iter,
			func(t int) {
				calls[t] = struct{}{}
			},
		)
		if !reflect.DeepEqual(calls, expected) {
			t.Fail()
		}
	})
}

func TestTryFold(t *testing.T) {
	t.Run("early exit", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		ret := iterator.TryFold[int](iter, 1,
			func(a int, t int) either.Either[int, int] {
				if t < 4 {
					return either.First[int, int]{Value: a * t}
				}
				return either.Second[int, int]{Value: a}
			},
		)
		if ret.UnwrapSecond() != 6 {
			t.Fail()
		}
		if iter.Next().Unwrap() != 5 { // Iterator should be at the next non-failed element.
			t.Fail()
		}
	})
	t.Run("no early exit", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		ret := iterator.TryFold[int](iter, 1,
			func(a int, t int) either.Either[int, int] {
				return either.First[int, int]{Value: a * t}
			},
		)
		if ret.Unwrap() != 120 {
			t.Fail()
		}
	})
}

func TestNth(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		nth := iterator.Nth[int](iter, 3)
		if nth.Unwrap() != 4 {
			t.Fail()
		}
		if iter.Next().Unwrap() != 5 { // Iterator is advanced.
			t.Fail()
		}
	})
	t.Run("0", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		nth := iterator.Nth[int](iter, 0)
		if nth.Unwrap() != 1 {
			t.Fail()
		}
		if iter.Next().Unwrap() != 2 { // Iterator is advanced.
			t.Fail()
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		nth := iterator.Nth[int](iter, 7)
		if !nth.IsNothing() {
			t.Fail()
		}
	})
}

func TestLast(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		last := iterator.Last[int](iter)
		if last.Unwrap() != 5 {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{},
		}
		last := iterator.Last[int](iter)
		if !last.IsNothing() {
			t.Fail()
		}
	})
}

func TestMax(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{2, 1, 5, 3, 4},
		}
		max := iterator.Max[int](iter)
		if max.Unwrap() != 5 {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{},
		}
		max := iterator.Max[int](iter)
		if !max.IsNothing() {
			t.Fail()
		}
	})
}

func TestMaxBy(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		iter := &fakeStringIterator{
			elements: []string{
				"one",
				"two",
				"three",
				"four",
				"five",
			},
		}
		max := iterator.MaxBy[string](iter, func(s1 string, s2 string) bool {
			return len(s1) < len(s2)
		})
		if max.Unwrap() != "three" {
			t.Fail()
		}
	})
	t.Run("last if multiple are equal", func(t *testing.T) {
		iter := &fakeStringIterator{
			elements: []string{
				"one",
				"two",
				"a",
				"b",
			},
		}
		max := iterator.MaxBy[string](iter, func(s1 string, s2 string) bool {
			return len(s1) < len(s2)
		})
		if max.Unwrap() != "two" {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeStringIterator{
			elements: []string{},
		}
		max := iterator.MaxBy[string](iter, func(s1 string, s2 string) bool {
			return len(s1) < len(s2)
		})
		if !max.IsNothing() {
			t.Fail()
		}
	})
}

func TestMin(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{2, 1, 5, 3, 4},
		}
		min := iterator.Min[int](iter)
		if min.Unwrap() != 1 {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{},
		}
		min := iterator.Min[int](iter)
		if !min.IsNothing() {
			t.Fail()
		}
	})
}

func TestMinBy(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		iter := &fakeStringIterator{
			elements: []string{
				"bb",
				"a",
				"ccc",
				"dddd",
				"ee",
			},
		}
		min := iterator.MinBy[string](iter, func(s1 string, s2 string) bool {
			return len(s1) < len(s2)
		})
		if min.Unwrap() != "a" {
			t.Fail()
		}
	})
	t.Run("first if multiple are equal", func(t *testing.T) {
		iter := &fakeStringIterator{
			elements: []string{
				"bb",
				"a",
				"ccc",
				"dddd",
				"e",
			},
		}
		min := iterator.MinBy[string](iter, func(s1 string, s2 string) bool {
			return len(s1) < len(s2)
		})
		if min.Unwrap() != "a" {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeStringIterator{
			elements: []string{},
		}
		min := iterator.MinBy[string](iter, func(s1 string, s2 string) bool {
			return len(s1) < len(s2)
		})
		if !min.IsNothing() {
			t.Fail()
		}
	})
}

func TestCollect(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{2, 1, 5, 3, 4},
		}
		res := iterator.Collect[int](iter)
		if !reflect.DeepEqual(res, iter.elements) {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := &fakeIterator{
			elements: []int{},
		}
		res := iterator.Collect[int](iter)
		if !reflect.DeepEqual(res, []int{}) {
			t.Fail()
		}
	})
}

func TestPartition(t *testing.T) {
	expectedTrue := []int{2, 4}
	expectedFalse := []int{1, 5, 3}
	iter := &fakeIterator{
		elements: []int{2, 1, 5, 3, 4},
	}
	trueList, falseList := iterator.Partition[int](iter, func(t int) bool {
		return t%2 == 0
	})
	if !reflect.DeepEqual(trueList, expectedTrue) {
		t.Fail()
	}
	if !reflect.DeepEqual(falseList, expectedFalse) {
		t.Fail()
	}
}

func TestUnzip(t *testing.T) {
	iter := &fakePairIterator{
		elements: []pair.Pair[int, string]{
			{First: 1, Second: "3"},
			{First: 2, Second: "2"},
			{First: 3, Second: "1"},
		},
	}
	expectedFirst := []int{1, 2, 3}
	expectedSecond := []string{"3", "2", "1"}
	firstList, secondList := iterator.Unzip[int, string](iter)
	if !reflect.DeepEqual(firstList, expectedFirst) {
		t.Fail()
	}
	if !reflect.DeepEqual(secondList, expectedSecond) {
		t.Fail()
	}
}
