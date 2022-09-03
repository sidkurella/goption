package sliceutil_test

import (
	"testing"

	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/sliceutil"
)

func TestIter(t *testing.T) {
	t.Run("list has elements", func(t *testing.T) {
		list := []int{1, 2, 3}
		iter := sliceutil.Iter(list)
		first := option.Some[*int]{Value: &list[0]}
		second := option.Some[*int]{Value: &list[1]}
		third := option.Some[*int]{Value: &list[2]}

		item := iter.Next()
		if item != first || *item.Unwrap() != 1 {
			t.Fail()
		}

		item = iter.Next()
		if item != second || *item.Unwrap() != 2 {
			t.Fail()
		}

		item = iter.Next()
		if item != third || *item.Unwrap() != 3 {
			t.Fail()
		}

		end := option.Nothing[*int]{}
		if iter.Next() != end {
			t.Fail()
		}
	})
	t.Run("empty list", func(t *testing.T) {
		list := []int{}
		iter := sliceutil.Iter(list)

		end := option.Nothing[*int]{}
		if iter.Next() != end {
			t.Fail()
		}
	})
}

func TestReverseIter(t *testing.T) {
	t.Run("list has elements", func(t *testing.T) {
		list := []int{1, 2, 3}
		iter := sliceutil.ReverseIter(list)
		first := option.Some[*int]{Value: &list[2]}
		second := option.Some[*int]{Value: &list[1]}
		third := option.Some[*int]{Value: &list[0]}

		item := iter.Next()
		if item != first || *item.Unwrap() != 3 {
			t.Fail()
		}

		item = iter.Next()
		if item != second || *item.Unwrap() != 2 {
			t.Fail()
		}

		item = iter.Next()
		if item != third || *item.Unwrap() != 1 {
			t.Fail()
		}

		end := option.Nothing[*int]{}
		if iter.Next() != end {
			t.Fail()
		}
	})
	t.Run("empty list", func(t *testing.T) {
		list := []int{}
		iter := sliceutil.ReverseIter(list)

		end := option.Nothing[*int]{}
		if iter.Next() != end {
			t.Fail()
		}
	})
}
