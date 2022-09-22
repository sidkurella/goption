package iterator_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/option"
)

func TestMapWhile(t *testing.T) {
	expected := []string{"10", "20"}
	i1 := &fakeIterator{
		elements: []int{1, 2, 3, 4},
	}
	iter := iterator.MapWhile[int](i1, func(t int) option.Option[string] {
		if t >= 3 {
			return option.Nothing[string]{}
		}
		return option.Some[string]{Value: strconv.Itoa(t * 10)}
	})
	calls := iterator.Collect[string](iter)
	if !reflect.DeepEqual(calls, expected) {
		t.Fail()
	}
}
