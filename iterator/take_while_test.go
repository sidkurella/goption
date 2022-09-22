package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestTakeWhile(t *testing.T) {
	t.Run("take success", func(t *testing.T) {
		expected := []int{1, 2}
		left := []int{4, 5}
		i1 := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		iter := iterator.TakeWhile[int](i1, func(i int) bool {
			return i < 3
		})
		calls := iterator.Collect[int](iter)
		actualLeft := iterator.Collect[int](i1)
		if !reflect.DeepEqual(calls, expected) {
			t.Fail()
		}
		if !reflect.DeepEqual(actualLeft, left) {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		i1 := &fakeIterator{
			elements: []int{},
		}
		iter := iterator.TakeWhile[int](i1, func(i int) bool {
			return i < 6
		})
		if !iter.Next().IsNothing() {
			t.Fail()
		}
	})
	t.Run("stops after initial false", func(t *testing.T) {
		expected := []int{-1, -2}
		left := []int{-4, 5}
		i1 := &fakeIterator{
			elements: []int{-1, -2, 3, -4, 5},
		}
		iter := iterator.TakeWhile[int](i1, func(i int) bool {
			return i < 0
		})
		calls := iterator.Collect[int](iter)
		actualLeft := iterator.Collect[int](i1)
		if !reflect.DeepEqual(calls, expected) {
			t.Fail()
		}
		if !reflect.DeepEqual(actualLeft, left) {
			t.Fail()
		}
	})
}
