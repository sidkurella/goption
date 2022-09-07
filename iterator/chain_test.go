package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestChain(t *testing.T) {
	elements := []int{}
	expected := []int{1, 2, 3, 10, 20, 30}
	i1 := &fakeIterator{
		elements: []int{1, 2, 3},
	}
	i2 := &fakeIterator{
		elements: []int{10, 20, 30},
	}
	iter := iterator.Chain[int](i1, i2)
	iterator.ForEach[int](iter, func(i int) {
		elements = append(elements, i)
	})
	if !reflect.DeepEqual(elements, expected) {
		t.Fail()
	}
}

func TestChain_FirstEmpty(t *testing.T) {
	elements := []int{}
	expected := []int{10, 20, 30}
	i1 := &fakeIterator{
		elements: []int{},
	}
	i2 := &fakeIterator{
		elements: []int{10, 20, 30},
	}
	iter := iterator.Chain[int](i1, i2)
	iterator.ForEach[int](iter, func(i int) {
		elements = append(elements, i)
	})
	if !reflect.DeepEqual(elements, expected) {
		t.Fail()
	}
}
