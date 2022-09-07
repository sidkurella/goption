package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestFilter(t *testing.T) {
	expected := []int{2, 4}
	i1 := &fakeIterator{
		elements: []int{1, 2, 3, 4},
	}
	iter := iterator.Filter[int](i1, func(t int) bool {
		return t%2 == 0
	})
	calls := iterator.Collect[int](iter)
	if !reflect.DeepEqual(calls, expected) {
		t.Fail()
	}
}
