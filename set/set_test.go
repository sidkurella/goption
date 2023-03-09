package set_test

import (
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/set"
	"github.com/sidkurella/goption/sliceutil"
)

func TestFromSlice(t *testing.T) {
	expected := set.From(map[int]struct{}{
		1: {},
		2: {},
	})
	actual := set.FromSlice([]int{1, 2})
	if !reflect.DeepEqual(actual, expected) {
		t.Fail()
	}
}

func TestIsEmpty(t *testing.T) {
	t.Run("not empty", func(t *testing.T) {
		s := set.From(map[int]struct{}{
			1: {},
			2: {},
		})
		if s.IsEmpty() {
			t.Fail()
		}
	})
	t.Run("empty", func(t *testing.T) {
		s := set.From(map[int]struct{}{})
		if !s.IsEmpty() {
			t.Fail()
		}
	})
}

func TestClear(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
	})
	if s.IsEmpty() {
		t.Fail()
	}
	s.Clear()
	if !s.IsEmpty() {
		t.Fail()
	}
}

func TestLen(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
	})
	if s.Len() != 2 {
		t.Fail()
	}
}

func TestContains(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
	})
	if !s.Contains(2) || s.Contains(3) {
		t.Fail()
	}
}

func TestInsert(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
	})
	ok := s.Insert(3)
	if ok || !s.Contains(3) { // Should not be in the set already.
		t.Fail()
	}
	ok = s.Insert(2)
	if !ok || !s.Contains(2) { // Should be in the set already.
		t.Fail()
	}
}

func TestAppend(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
	})
	s.Append(1, 2, 3, 4, 5, 6)
	expected := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
		6: {},
	})
	if !reflect.DeepEqual(s, expected) {
		t.Fail()
	}
}

func TestRemove(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
	})
	ok := s.Remove(3)
	if ok || s.Contains(3) { // Should not be in the set already.
		t.Fail()
	}
	ok = s.Remove(2)
	if !ok || s.Contains(2) { // Should be in the set already.
		t.Fail()
	}
}

func TestFilter(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	})
	s.Filter(func(i int) bool {
		return i%2 == 0
	})
	expected := set.From(map[int]struct{}{
		2: {},
		4: {},
	})
	if !reflect.DeepEqual(s, expected) {
		t.Fail()
	}
}

func TestCollect(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	})
	out := s.Collect()
	sort.Slice(out, func(i, j int) bool {
		a := out[i]
		b := out[j]
		return a < b
	})
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestForEach(t *testing.T) {
	calls := map[int]struct{}{}
	expectedCalls := map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	}
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	})
	s.ForEach(func(i int) {
		calls[i] = struct{}{}
	})
	if !reflect.DeepEqual(calls, expectedCalls) {
		t.Fail()
	}
}

func TestUnion(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
	})
	s2 := set.From(map[int]struct{}{
		2: {},
		3: {},
		4: {},
	})
	expected := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	})
	if !reflect.DeepEqual(s.Union(s2), expected) {
		t.Fail()
	}
}

func TestIntersection(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
	})
	s2 := set.From(map[int]struct{}{
		2: {},
		3: {},
		4: {},
	})
	expected := set.From(map[int]struct{}{
		2: {},
		3: {},
	})
	if !reflect.DeepEqual(s.Intersection(s2), expected) {
		t.Fail()
	}
}

func TestDifference(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
	})
	s2 := set.From(map[int]struct{}{
		2: {},
		3: {},
		4: {},
	})
	expected := set.From(map[int]struct{}{
		1: {},
	})
	if !reflect.DeepEqual(s.Difference(s2), expected) {
		t.Fail()
	}
}

func TestSymmetricDifference(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
	})
	s2 := set.From(map[int]struct{}{
		2: {},
		3: {},
		4: {},
	})
	expected := set.From(map[int]struct{}{
		1: {},
		4: {},
	})
	if !reflect.DeepEqual(s.SymmetricDifference(s2), expected) {
		t.Fail()
	}
}

func TestIsDisjoint(t *testing.T) {
	t.Run("disjoint", func(t *testing.T) {
		s := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
		})
		s2 := set.From(map[int]struct{}{
			4: {},
			5: {},
			6: {},
		})
		if !s.IsDisjoint(s2) {
			t.Fail()
		}
	})
	t.Run("not disjoint", func(t *testing.T) {
		s := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
		})
		s2 := set.From(map[int]struct{}{
			3: {},
			4: {},
			5: {},
		})
		if s.IsDisjoint(s2) {
			t.Fail()
		}
	})
}

func TestIsSubset(t *testing.T) {
	t.Run("subset", func(t *testing.T) {
		s := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
		})
		s2 := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
			4: {},
		})
		if !s.IsSubset(s2) {
			t.Fail()
		}
	})
	t.Run("not subset", func(t *testing.T) {
		s := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
			5: {},
		})
		s2 := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
			4: {},
		})
		if s.IsSubset(s2) {
			t.Fail()
		}
	})
}

func TestIsSuperset(t *testing.T) {
	t.Run("superset", func(t *testing.T) {
		s := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
			4: {},
		})
		s2 := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
		})
		if !s.IsSuperset(s2) {
			t.Fail()
		}
	})
	t.Run("not superset", func(t *testing.T) {
		s := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
			5: {},
		})
		s2 := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
			4: {},
		})
		if s.IsSuperset(s2) {
			t.Fail()
		}
	})
}

func TestMap(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	})
	out := set.Map(s,
		func(k int) string {
			return strconv.Itoa(k * 2)
		},
	)
	expected := set.From(map[string]struct{}{
		"2": {},
		"4": {},
		"6": {},
		"8": {},
	})
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestFilterMap(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	})
	out := set.FilterMap(s,
		func(k int) option.Option[string] {
			if k%2 == 0 {
				return option.Nothing[string]()
			}
			return option.Some(strconv.Itoa(k * 2))
		},
	)
	expected := set.From(map[string]struct{}{
		"2": {},
		"6": {},
	})
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestFold(t *testing.T) {
	s := set.From(map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
	})
	out := set.Fold(s, float64(1),
		func(a float64, k int) float64 {
			return float64(a * float64(k))
		},
	)
	if out != float64(24) {
		t.Fail()
	}
}

func TestCollectInto(t *testing.T) {
	t.Run("empty set", func(t *testing.T) {
		i := sliceutil.Iter([]int{
			1, 2, 3,
		})
		s := set.New[int]()
		expected := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
		})
		out := iterator.CollectInto[int](i, s)
		if !reflect.DeepEqual(out, expected) {
			t.Fail()
		}
		if !reflect.DeepEqual(s, expected) {
			t.Fail()
		}
	})
	t.Run("non-empty set", func(t *testing.T) {
		i := sliceutil.Iter([]int{
			1, 2, 3,
		})
		s := set.New[int]()
		s.Insert(1)
		expected := set.From(map[int]struct{}{
			1: {},
			2: {},
			3: {},
		})
		out := iterator.CollectInto[int](i, s)
		if !reflect.DeepEqual(out, expected) {
			t.Fail()
		}
		if !reflect.DeepEqual(s, expected) {
			t.Fail()
		}
	})
}
