package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestTake(t *testing.T) {
	t.Run("take success", func(t *testing.T) {
		expected := []int{1, 2, 3}
		i1 := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		iter := iterator.Take[int](i1, 3)
		calls := iterator.Collect[int](iter)
		if !reflect.DeepEqual(calls, expected) {
			t.Fail()
		}
	})
	t.Run("not long enough", func(t *testing.T) {
		expected := []int{1, 2, 3, 4, 5}
		i1 := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		iter := iterator.Take[int](i1, 100)
		calls := iterator.Collect[int](iter)
		if !reflect.DeepEqual(calls, expected) {
			t.Fail()
		}
	})
}
