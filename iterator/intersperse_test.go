package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestIntersperse(t *testing.T) {
	expected := []int{1, 100, 2, 100, 3, 100, 4, 100}
	i1 := &fakeIterator{
		elements: []int{1, 2, 3, 4},
	}
	iter := iterator.Intersperse[int](i1, 100)
	calls := iterator.Collect[int](iter)
	if !reflect.DeepEqual(calls, expected) {
		t.Fail()
	}
}
