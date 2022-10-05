package set_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/set"
)

func TestSet_Iter(t *testing.T) {
	expected := []int{1, 3, 5}
	m := set.From(map[int]struct{}{
		1: {},
		3: {},
		5: {},
	})
	actual := iterator.Collect(m.Iter())
	sort.Slice(actual, func(i, j int) bool {
		e1 := actual[i]
		e2 := actual[j]
		return e1 < e2
	})
	if !reflect.DeepEqual(actual, expected) {
		t.Fail()
	}
}
