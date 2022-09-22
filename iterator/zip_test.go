package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/pair"
)

func TestZip(t *testing.T) {
	t.Run("equal lengths", func(t *testing.T) {
		expected := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
			{First: 4, Second: "4"},
		}
		i1 := &fakeIterator{
			elements: []int{1, 2, 3, 4},
		}
		i2 := &fakeStringIterator{
			elements: []string{"1", "2", "3", "4"},
		}
		iter := iterator.Zip[int, string](i1, i2)
		res := iterator.Collect[pair.Pair[int, string]](iter)
		if !reflect.DeepEqual(res, expected) {
			t.Fail()
		}
	})
	t.Run("first iterator shorter", func(t *testing.T) {
		expected := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "30"},
		}
		i1 := &fakeIterator{
			elements: []int{1, 2, 3},
		}
		i2 := &fakeStringIterator{
			elements: []string{"1", "2", "30", "4"},
		}
		iter := iterator.Zip[int, string](i1, i2)
		res := iterator.Collect[pair.Pair[int, string]](iter)
		if !reflect.DeepEqual(res, expected) {
			t.Fail()
		}
	})
	t.Run("second iterator shorter", func(t *testing.T) {
		expected := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 30, Second: "3"},
		}
		i1 := &fakeIterator{
			elements: []int{1, 2, 30, 4},
		}
		i2 := &fakeStringIterator{
			elements: []string{"1", "2", "3"},
		}
		iter := iterator.Zip[int, string](i1, i2)
		res := iterator.Collect[pair.Pair[int, string]](iter)
		if !reflect.DeepEqual(res, expected) {
			t.Fail()
		}
	})
}
