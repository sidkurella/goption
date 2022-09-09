package iterator_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestMap(t *testing.T) {
	expected := []string{"10", "20", "30", "40"}
	i1 := &fakeIterator{
		elements: []int{1, 2, 3, 4},
	}
	iter := iterator.Map[int](i1, func(t int) string {
		return strconv.Itoa(t * 10)
	})
	calls := iterator.Collect[string](iter)
	if !reflect.DeepEqual(calls, expected) {
		t.Fail()
	}
}
