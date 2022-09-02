package maputil_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/sidkurella/goption/maputil"
	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/result"
)

func TestContainsKey(t *testing.T) {
	m := maputil.From(map[string]int{
		"key 1": 1,
		"key 2": 2,
		"key 3": 3,
	})
	if !m.ContainsKey("key 1") || m.ContainsKey("key 4") {
		t.Fail()
	}
}

func TestLen(t *testing.T) {
	m := maputil.From(map[string]int{
		"key 1": 1,
		"key 2": 2,
		"key 3": 3,
	})
	empty := maputil.New[int, string]()
	if m.Len() != 3 || empty.Len() != 0 {
		t.Fail()
	}
}

func TestIsEmpty(t *testing.T) {
	m := maputil.From(map[string]int{
		"key 1": 1,
		"key 2": 2,
		"key 3": 3,
	})
	empty := maputil.New[int, string]()
	if m.IsEmpty() || !empty.IsEmpty() {
		t.Fail()
	}
}

func TestClear(t *testing.T) {
	m := maputil.From(map[string]int{
		"key 1": 1,
		"key 2": 2,
		"key 3": 3,
	})
	empty := maputil.New[string, int]()
	m.Clear()
	if !m.IsEmpty() || !reflect.DeepEqual(m, empty) {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		m := maputil.From(map[string]int{
			"key 1": 1,
			"key 2": 2,
			"key 3": 3,
		})
		res := m.Get("key 3")
		expected := option.Some[int]{Value: 3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("key does not", func(t *testing.T) {
		m := maputil.From(map[string]int{
			"key 1": 1,
			"key 2": 2,
			"key 3": 3,
		})
		res := m.Get("key 4")
		expected := option.Nothing[int]{}
		if res != expected {
			t.Fail()
		}
	})
}

func TestInsert(t *testing.T) {
	t.Run("key doesn't already exist", func(t *testing.T) {
		m := maputil.From(map[string]int{
			"key 1": 1,
			"key 2": 2,
			"key 3": 3,
		})
		res := m.Insert("key 4", 4)
		getVal := m.Get("key 4")
		expectedGet := option.Some[int]{Value: 4}
		if !res.IsNothing() || !m.ContainsKey("key 4") || getVal != expectedGet {
			t.Fail()
		}
	})
	t.Run("key already exists", func(t *testing.T) {
		m := maputil.From(map[string]int{
			"key 1": 1,
			"key 2": 2,
			"key 3": 3,
		})
		res := m.Insert("key 3", 4)
		getVal := m.Get("key 3")
		expected := option.Some[int]{Value: 3}
		expectedGet := option.Some[int]{Value: 4}
		if res != expected || !m.ContainsKey("key 3") || getVal != expectedGet {
			t.Fail()
		}
	})
}

func TestTryInsert(t *testing.T) {
	t.Run("key doesn't already exist", func(t *testing.T) {
		m := maputil.From(map[string]int{
			"key 1": 1,
			"key 2": 2,
			"key 3": 3,
		})
		res := m.TryInsert("key 4", 4)
		expectedRes := result.Ok[int, maputil.OccupiedError[string, int]]{Value: 4}
		getVal := m.Get("key 4")
		expectedGet := option.Some[int]{Value: 4}
		if res != expectedRes || !m.ContainsKey("key 4") || getVal != expectedGet {
			t.Fail()
		}
	})
	t.Run("key already exists", func(t *testing.T) {
		m := maputil.From(map[string]int{
			"key 1": 1,
			"key 2": 2,
			"key 3": 3,
		})
		res := m.TryInsert("key 3", 4)
		expectedRes := result.Err[int, maputil.OccupiedError[string, int]]{
			Value: maputil.OccupiedError[string, int]{
				Key:   "key 3",
				Value: 3,
			},
		}
		getVal := m.Get("key 3")
		expectedGet := option.Some[int]{Value: 3}
		if res != expectedRes || !m.ContainsKey("key 3") || getVal != expectedGet {
			t.Fail()
		}
	})
}

func TestRemove(t *testing.T) {
	t.Run("key doesn't already exist", func(t *testing.T) {
		m := maputil.From(map[string]int{
			"key 1": 1,
			"key 2": 2,
			"key 3": 3,
		})
		res := m.Remove("key 4")
		getVal := m.Get("key 4")
		expectedGet := option.Nothing[int]{}
		if !res.IsNothing() || m.ContainsKey("key 4") || getVal != expectedGet {
			t.Fail()
		}
	})
	t.Run("key already exists", func(t *testing.T) {
		m := maputil.From(map[string]int{
			"key 1": 1,
			"key 2": 2,
			"key 3": 3,
		})
		res := m.Remove("key 3")
		expectedRes := option.Some[int]{Value: 3}
		getVal := m.Get("key 3")
		expectedGet := option.Nothing[int]{}
		if res != expectedRes || m.ContainsKey("key 3") || getVal != expectedGet {
			t.Fail()
		}
	})
}

func TestFilter(t *testing.T) {
	m := maputil.From(map[int]string{
		1: "one",
		2: "two",
		3: "three",
	})
	m.Filter(func(i int, _ string) bool {
		return i%2 == 0
	})
	expected := map[int]string{
		2: "two",
	}
	if !reflect.DeepEqual(m.Into(), expected) {
		t.Fail()
	}
}

func TestCollect(t *testing.T) {
	m := maputil.From(map[int]string{
		1: "one",
		2: "two",
		3: "three",
	})
	out := m.Collect()
	sort.Slice(out, func(i, j int) bool {
		l := out[i]
		r := out[j]
		return l.Key < r.Key
	})
	expected := []maputil.Entry[int, string]{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	}
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestCollectKeys(t *testing.T) {
	m := maputil.From(map[int]string{
		1: "one",
		2: "two",
		3: "three",
	})
	out := m.CollectKeys()
	sort.Slice(out, func(i, j int) bool {
		l := out[i]
		r := out[j]
		return l < r
	})
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestCollectValues(t *testing.T) {
	m := maputil.From(map[int]string{
		1: "one",
		2: "two",
		3: "three",
	})
	out := m.CollectValues()
	sort.Slice(out, func(i, j int) bool {
		l := out[i]
		r := out[j]
		return l < r
	})
	expected := []string{"one", "three", "two"}
	if !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestForEach(t *testing.T) {
	calls := map[maputil.Entry[int, string]]struct{}{}
	m := maputil.From(map[int]string{
		1: "one",
		2: "two",
		3: "three",
	})
	m.ForEach(func(i int, s string) {
		calls[maputil.Entry[int, string]{
			Key:   i,
			Value: s,
		}] = struct{}{}
	})
	expected := map[maputil.Entry[int, string]]struct{}{
		{1, "one"}:   {},
		{2, "two"}:   {},
		{3, "three"}: {},
	}
	if !reflect.DeepEqual(calls, expected) {
		t.Fail()
	}
}

func TestApply(t *testing.T) {
	calls := map[maputil.Entry[int, string]]struct{}{}
	m := maputil.From(map[int]string{
		1: "one",
		2: "two",
		3: "three",
	})
	out := maputil.Apply(m,
		func(i int, s string) (float64, int64) {
			calls[maputil.Entry[int, string]{
				Key:   i,
				Value: s,
			}] = struct{}{}

			return float64(i + 1), int64(len(s) * 2)
		},
	)
	expected := maputil.From(map[float64]int64{
		2: 6,
		3: 6,
		4: 10,
	})
	expectedCalls := map[maputil.Entry[int, string]]struct{}{
		{1, "one"}:   {},
		{2, "two"}:   {},
		{3, "three"}: {},
	}
	if !reflect.DeepEqual(calls, expectedCalls) || !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestFilterMap(t *testing.T) {
	calls := map[maputil.Entry[int, string]]struct{}{}
	m := maputil.From(map[int]string{
		1: "one",
		2: "two",
		3: "three",
	})
	out := maputil.FilterMap(m,
		func(i int, s string) option.Option[maputil.Entry[float64, int64]] {
			calls[maputil.Entry[int, string]{
				Key:   i,
				Value: s,
			}] = struct{}{}

			if i%2 == 0 {
				return option.Nothing[maputil.Entry[float64, int64]]{}
			}

			return option.Some[maputil.Entry[float64, int64]]{
				Value: maputil.Entry[float64, int64]{
					Key:   float64(i + 1),
					Value: int64(len(s) * 2),
				},
			}
		},
	)
	expected := maputil.From(map[float64]int64{
		2: 6,
		4: 10,
	})
	expectedCalls := map[maputil.Entry[int, string]]struct{}{
		{1, "one"}:   {},
		{2, "two"}:   {},
		{3, "three"}: {},
	}
	if !reflect.DeepEqual(calls, expectedCalls) || !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}

func TestFold(t *testing.T) {
	calls := map[maputil.Entry[int, string]]struct{}{}
	m := maputil.From(map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
	})
	out := maputil.Fold(m, float64(1),
		func(a float64, i int, s string) float64 {
			calls[maputil.Entry[int, string]{
				Key:   i,
				Value: s,
			}] = struct{}{}

			return a * float64(i)
		},
	)
	expected := float64(24)
	expectedCalls := map[maputil.Entry[int, string]]struct{}{
		{1, "one"}:   {},
		{2, "two"}:   {},
		{3, "three"}: {},
		{4, "four"}:  {},
	}
	if !reflect.DeepEqual(calls, expectedCalls) || !reflect.DeepEqual(out, expected) {
		t.Fail()
	}
}
