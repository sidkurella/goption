package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestInspect(t *testing.T) {
	expectedCalls := []int{1, 2, 3}
	actualCalls := []int{}
	i := &fakeIterator{
		elements: []int{1, 2, 3},
	}
	iter := iterator.Inspect[int](i, func(i int) {
		actualCalls = append(actualCalls, i)
	})
	vals := iterator.Collect[int](iter)
	if !reflect.DeepEqual(vals, expectedCalls) {
		t.Fail()
	}
	if !reflect.DeepEqual(actualCalls, expectedCalls) {
		t.Fail()
	}
}
