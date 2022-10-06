package sliceutil_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/sidkurella/goption/pair"
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

func TestFirst(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		in := []int{1, 2, 3}
		fst := sliceutil.First(in)
		if fst.Unwrap() != 1 {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		in := []int{}
		fst := sliceutil.First(in)
		if !fst.IsNothing() {
			t.Fail()
		}
	})
}

func TestLast(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		in := []int{1, 2, 3}
		lst := sliceutil.Last(in)
		if lst.Unwrap() != 3 {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		in := []int{}
		lst := sliceutil.Last(in)
		if !lst.IsNothing() {
			t.Fail()
		}
	})
}

func TestStartsWith(t *testing.T) {
	t.Run("starts with", func(t *testing.T) {
		in := []int{1, 2, 3}
		if !sliceutil.StartsWith(in, []int{1, 2}) {
			t.Fail()
		}
	})
	t.Run("doesn't start with", func(t *testing.T) {
		in := []int{1, 2, 3}
		if sliceutil.StartsWith(in, []int{2}) {
			t.Fail()
		}
	})
	t.Run("too short", func(t *testing.T) {
		in := []int{1, 2, 3}
		if sliceutil.StartsWith(in, []int{1, 2, 3, 4, 5}) {
			t.Fail()
		}
	})
}

func TestStartsWithFunc(t *testing.T) {
	t.Run("starts with", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		startsWith := sliceutil.StartsWithFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "3"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		if !startsWith {
			t.Fail()
		}
	})
	t.Run("doesn't start with", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		startsWith := sliceutil.StartsWithFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 3, Second: "3"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		if startsWith {
			t.Fail()
		}
	})
	t.Run("too short", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		startsWith := sliceutil.StartsWithFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
			{First: 4, Second: "4"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		if startsWith {
			t.Fail()
		}
	})
}

// EndsWith
// EndsWithFunc

func TestEndsWith(t *testing.T) {
	t.Run("ends with", func(t *testing.T) {
		in := []int{1, 2, 3}
		if !sliceutil.EndsWith(in, []int{2, 3}) {
			t.Fail()
		}
	})
	t.Run("doesn't end with", func(t *testing.T) {
		in := []int{1, 2, 3}
		if sliceutil.EndsWith(in, []int{2, 4}) {
			t.Fail()
		}
	})
	t.Run("too short", func(t *testing.T) {
		in := []int{1, 2, 3}
		if sliceutil.EndsWith(in, []int{1, 2, 3, 4, 5}) {
			t.Fail()
		}
	})
}

func TestEndsWithFunc(t *testing.T) {
	t.Run("ends with", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		endsWith := sliceutil.EndsWithFunc(in, []pair.Pair[int, string]{
			{First: 2, Second: "1"},
			{First: 3, Second: "3"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		if !endsWith {
			t.Fail()
		}
	})
	t.Run("doesn't end with", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		endsWith := sliceutil.EndsWithFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 3, Second: "3"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		if endsWith {
			t.Fail()
		}
	})
	t.Run("too short", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		endsWith := sliceutil.EndsWithFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
			{First: 4, Second: "4"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		if endsWith {
			t.Fail()
		}
	})
}

func TestStripPrefix(t *testing.T) {
	t.Run("starts with", func(t *testing.T) {
		in := []int{1, 2, 3, 4}
		ret := sliceutil.StripPrefix(in, []int{1, 2})
		expected := []int{3, 4}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
	t.Run("doesn't start with", func(t *testing.T) {
		in := []int{1, 2, 3, 4}
		ret := sliceutil.StripPrefix(in, []int{1, 3})
		expected := []int{1, 2, 3, 4}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
	t.Run("too short", func(t *testing.T) {
		in := []int{1, 2, 3, 4}
		ret := sliceutil.StripPrefix(in, []int{1, 2, 3, 4, 5})
		expected := []int{1, 2, 3, 4}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
}

func TestStripPrefixFunc(t *testing.T) {
	t.Run("starts with", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		ret := sliceutil.StripPrefixFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "3"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		expected := []pair.Pair[int, string]{
			{First: 3, Second: "3"},
		}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
	t.Run("doesn't start with", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		ret := sliceutil.StripPrefixFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 3, Second: "3"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		expected := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
	t.Run("too short", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		ret := sliceutil.StripPrefixFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
			{First: 4, Second: "2"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		expected := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
}

func TestStripSuffix(t *testing.T) {
	t.Run("ends with", func(t *testing.T) {
		in := []int{1, 2, 3, 4}
		ret := sliceutil.StripSuffix(in, []int{3, 4})
		expected := []int{1, 2}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
	t.Run("doesn't end with", func(t *testing.T) {
		in := []int{1, 2, 3, 4}
		ret := sliceutil.StripSuffix(in, []int{1, 3})
		expected := []int{1, 2, 3, 4}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
	t.Run("too short", func(t *testing.T) {
		in := []int{1, 2, 3, 4}
		ret := sliceutil.StripSuffix(in, []int{1, 2, 3, 4, 5})
		expected := []int{1, 2, 3, 4}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
}

func TestStripSuffixFunc(t *testing.T) {
	t.Run("ends with", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		ret := sliceutil.StripSuffixFunc(in, []pair.Pair[int, string]{
			{First: 2, Second: "1"},
			{First: 3, Second: "3"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		expected := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
		}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
	t.Run("doesn't end with", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		ret := sliceutil.StripSuffixFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 3, Second: "3"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		expected := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
	t.Run("too short", func(t *testing.T) {
		in := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		ret := sliceutil.StripSuffixFunc(in, []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
			{First: 4, Second: "2"},
		}, func(t1, t2 pair.Pair[int, string]) bool {
			return t1.First == t2.First
		})
		expected := []pair.Pair[int, string]{
			{First: 1, Second: "1"},
			{First: 2, Second: "2"},
			{First: 3, Second: "3"},
		}
		if !reflect.DeepEqual(ret, expected) {
			t.Fail()
		}
	})
}
