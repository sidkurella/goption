package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/pair"
)

func TestEnumerate(t *testing.T) {
	elements := []pair.Pair[int, string]{}
	expected := []pair.Pair[int, string]{
		{First: 0, Second: "one"},
		{First: 1, Second: "two"},
		{First: 2, Second: "three"},
	}
	i1 := &fakeStringIterator{
		elements: []string{"one", "two", "three"},
	}
	iter := iterator.Enumerate[string](i1)
	iterator.ForEach[pair.Pair[int, string]](iter, func(v pair.Pair[int, string]) {
		elements = append(elements, v)
	})
	if !reflect.DeepEqual(elements, expected) {
		t.Fail()
	}
}

func TestEnumerate_Empty(t *testing.T) {
	elements := []pair.Pair[int, string]{}
	expected := []pair.Pair[int, string]{}
	i1 := &fakeStringIterator{
		elements: []string{},
	}
	iter := iterator.Enumerate[string](i1)
	iterator.ForEach[pair.Pair[int, string]](iter, func(v pair.Pair[int, string]) {
		elements = append(elements, v)
	})
	if !reflect.DeepEqual(elements, expected) {
		t.Fail()
	}
	if !iter.Next().IsNothing() {
		t.Fail()
	}
}
