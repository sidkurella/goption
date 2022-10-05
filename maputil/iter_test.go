package maputil_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/maputil"
)

func TestMap_Iter(t *testing.T) {
	expected := []maputil.Entry[int, string]{
		{
			Key:   1,
			Value: "2",
		},
		{
			Key:   3,
			Value: "4",
		},
		{
			Key:   5,
			Value: "6",
		},
	}
	m := maputil.From(map[int]string{
		1: "2",
		3: "4",
		5: "6",
	})
	actual := iterator.Collect[maputil.Entry[int, string]](m.Iter())
	sort.Slice(actual, func(i, j int) bool {
		e1 := actual[i]
		e2 := actual[j]
		return e1.Key < e2.Key
	})
	if !reflect.DeepEqual(actual, expected) {
		t.Fail()
	}
}
