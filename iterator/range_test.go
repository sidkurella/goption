package iterator_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/iterator"
)

func TestRange(t *testing.T) {
	t.Run("range", func(t *testing.T) {
		i := iterator.Range(1, 5)
		expected := []int{1, 2, 3, 4}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
	t.Run("range empty", func(t *testing.T) {
		i := iterator.Range(3, 1)
		expected := []int{}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
}

func TestRangeBy(t *testing.T) {
	t.Run("range forwards", func(t *testing.T) {
		i := iterator.RangeBy(1, 6, 2)
		expected := []int{1, 3, 5}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
	t.Run("range backwards", func(t *testing.T) {
		i := iterator.RangeBy(6, 1, -1)
		expected := []int{6, 5, 4, 3, 2}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
	t.Run("range empty", func(t *testing.T) {
		i := iterator.RangeBy(3, 1, 1)
		expected := []int{}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
}

func TestRangeInclusive(t *testing.T) {
	t.Run("range", func(t *testing.T) {
		i := iterator.RangeInclusive(1, 5)
		expected := []int{1, 2, 3, 4, 5}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
	t.Run("range empty", func(t *testing.T) {
		i := iterator.RangeInclusive(3, 1)
		expected := []int{}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
}

func TestRangeInclusiveBy(t *testing.T) {
	t.Run("range forwards", func(t *testing.T) {
		i := iterator.RangeInclusiveBy(1, 7, 2)
		expected := []int{1, 3, 5, 7}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
	t.Run("range backwards", func(t *testing.T) {
		i := iterator.RangeInclusiveBy(6, 1, -1)
		expected := []int{6, 5, 4, 3, 2, 1}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
	t.Run("range empty", func(t *testing.T) {
		i := iterator.RangeInclusiveBy(3, 1, 1)
		expected := []int{}
		actual := iterator.Collect[int](i)
		if !reflect.DeepEqual(expected, actual) {
			t.Fail()
		}
	})
}
