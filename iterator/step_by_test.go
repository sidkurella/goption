package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestStepBy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := []int{1, 3, 5}
		i1 := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		iter := iterator.StepBy[int](i1, 2)
		calls := iterator.Collect[int](iter)
		if !reflect.DeepEqual(calls, expected) {
			t.Fail()
		}
	})
	t.Run("not long enough", func(t *testing.T) {
		i1 := &fakeIterator{
			elements: []int{1, 2, 3, 4, 5},
		}
		iter := iterator.StepBy[int](i1, 600)
		if iter.Next().Unwrap() != 1 {
			t.Fail()
		}
		if !iter.Next().IsNothing() {
			t.Fail()
		}
	})
}
