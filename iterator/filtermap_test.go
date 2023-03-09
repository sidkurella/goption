package iterator_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/option"
)

func TestFilterMap(t *testing.T) {
	expected := []string{"20", "40"}
	i1 := &fakeIterator{
		elements: []int{1, 2, 3, 4},
	}
	iter := iterator.FilterMap[int](i1, func(t int) option.Option[string] {
		if t%2 == 0 {
			return option.Some(strconv.Itoa(t * 10))
		}
		return option.Nothing[string]()
	})
	calls := iterator.Collect[string](iter)
	if !reflect.DeepEqual(calls, expected) {
		t.Fail()
	}
}
