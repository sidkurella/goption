package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestSkipWhile(t *testing.T) {
	t.Run("skip success", func(t *testing.T) {
		expected := []int{3, 4, 5}
		i1 := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		iter := iterator.SkipWhile[int](i1, func(i int) bool {
			return i < 3
		})
		calls := iterator.Collect[int](iter)
		if !reflect.DeepEqual(calls, expected) {
			t.Fail()
		}
	})
	t.Run("not long enough", func(t *testing.T) {
		i1 := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		iter := iterator.SkipWhile[int](i1, func(i int) bool {
			return i < 6
		})
		if !iter.Next().IsNothing() {
			t.Fail()
		}
	})
	t.Run("stops after initial false", func(t *testing.T) {
		expected := []int{3, -4, 5}
		i1 := &fakeIterator{
			elements: []int{-1, -2, 3, -4, 5},
		}
		iter := iterator.SkipWhile[int](i1, func(i int) bool {
			return i < 0
		})
		calls := iterator.Collect[int](iter)
		if !reflect.DeepEqual(calls, expected) {
			t.Fail()
		}
	})
}
