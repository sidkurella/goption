package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestIntersperseWith(t *testing.T) {
	start := 50
	expected := []int{1, 100, 2, 200, 3, 400, 4, 800}
	i1 := &fakeIterator{
		elements: []int{1, 2, 3, 4},
	}
	iter := iterator.IntersperseWith[int](i1, func() int {
		start = start * 2
		return start
	})
	calls := iterator.Collect[int](iter)
	if !reflect.DeepEqual(calls, expected) {
		t.Fail()
	}
}
