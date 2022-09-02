package sliceutil_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/sidkurella/goption/sliceutil"
)

func TestMap(t *testing.T) {
	in := []int{1, 2, 3}
	expected := []string{"2", "3", "4"}
	out := sliceutil.Map(in, func(t int) string {
		return strconv.Itoa(t + 1)
	})
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestIndexMap(t *testing.T) {
	in := []int64{1, 2, 3}
	expected := []string{"1", "3", "5"}
	out := sliceutil.IndexMap(in, func(i int, f int64) string {
		return strconv.FormatInt(int64(i)+f, 10)
	})
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestReverse(t *testing.T) {
	in := []int64{1, 2, 3}
	expected := []int64{3, 2, 1}
	out := sliceutil.Reverse(in)
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
	if !reflect.DeepEqual(in, expected) {
		t.Fail() // Reverses in-place.
	}
}

func TestReversed(t *testing.T) {
	in := []int64{1, 2, 3}
	original := []int64{1, 2, 3}
	expected := []int64{3, 2, 1}
	out := sliceutil.Reversed(in)
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
	if !reflect.DeepEqual(in, original) {
		t.Fail() // Does not reverse in place.
	}
}

func TestFoldLeft(t *testing.T) {
	in := []int{1, 2, 3}
	out := sliceutil.FoldLeft(in, "",
		func(a string, i int) string {
			return a + strconv.Itoa(i*i)
		},
	)
	expected := "149"
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestFoldRight(t *testing.T) {
	in := []int{1, 2, 3}
	out := sliceutil.FoldRight(in, "",
		func(a string, i int) string {
			return a + strconv.Itoa(i*i)
		},
	)
	expected := "941"
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}
