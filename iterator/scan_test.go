package iterator_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/option"
)

func TestScan(t *testing.T) {
	expected := []string{"-1", "-2", "-6"}
	i1 := &fakeIterator{
		elements: []int{1, 2, 3},
	}
	iter := iterator.Scan[int](i1, 1, func(state *int, t int) option.Option[string] {
		*state = *state * t
		return option.Some(strconv.Itoa(-*state))
	})
	calls := iterator.Collect[string](iter)
	if !reflect.DeepEqual(calls, expected) {
		t.Fail()
	}
}
