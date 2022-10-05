package option_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/sidkurella/goption/option"
)

func TestOption_IsSome(t *testing.T) {
	t.Run("true for Some", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		if !opt.IsSome() {
			t.Fail()
		}
	})
	t.Run("false for Nothing", func(t *testing.T) {
		var opt option.Option[int] = option.Nothing[int]{}
		if opt.IsSome() {
			t.Fail()
		}
	})
}

func TestOption_IsNothing(t *testing.T) {
	t.Run("false for Some", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		if opt.IsNothing() {
			t.Fail()
		}
	})
	t.Run("true for Nothing", func(t *testing.T) {
		var opt option.Option[int] = option.Nothing[int]{}
		if !opt.IsNothing() {
			t.Fail()
		}
	})
}

func TestOption_Unwrap(t *testing.T) {
	t.Run("unwrap successful for Some", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		if opt.Unwrap() != 4 {
			t.Fail()
		}
	})
	t.Run("panics for Nothing", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic was expected but did not occur")
			}
		}()

		var opt option.Option[int] = option.Nothing[int]{}
		_ = opt.Unwrap()
	})
}

func TestOption_UnwrapOr(t *testing.T) {
	t.Run("returns internal value for Some", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		if opt.UnwrapOr(3) != 4 {
			t.Fail()
		}
	})
	t.Run("returns default for Nothing", func(t *testing.T) {
		var opt option.Option[int] = option.Nothing[int]{}
		if opt.UnwrapOr(3) != 3 {
			t.Fail()
		}
	})
}

func TestOption_Get(t *testing.T) {
	t.Run("get succeeds for Some", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		val, ok := opt.Get()
		if !ok || val != 4 {
			t.Fail()
		}
	})
	t.Run("returns default for Nothing", func(t *testing.T) {
		var opt option.Option[int] = option.Nothing[int]{}
		val, ok := opt.Get()
		if ok || val != 0 {
			t.Fail()
		}
	})
}

func TestOption_Expect(t *testing.T) {
	t.Run("successful for Some", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		if opt.Expect("don't panic") != 4 {
			t.Fail()
		}
	})
	t.Run("panics on Nothing", func(t *testing.T) {
		msg := "panic expected"
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("Panic was expected but did not occur")
			} else if r != msg {
				t.Errorf("Panic was expected but message was %v (did not match expected msg %v)", r, msg)
			}
		}()

		var opt option.Option[int] = option.Nothing[int]{}
		_ = opt.Expect(msg)
	})
}

func TestOption_IsSomeAnd(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		if !opt.IsSomeAnd(func(t *int) bool { return (*t) == 4 }) {
			t.Fail()
		}
	})
	t.Run("failed", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		if opt.IsSomeAnd(func(t *int) bool { return (*t) == 3 }) {
			t.Fail()
		}
	})
	t.Run("fails for Nothing", func(t *testing.T) {
		var opt option.Option[int] = option.Nothing[int]{}
		if opt.IsSomeAnd(func(t *int) bool { return (*t) == 3 }) {
			t.Fail()
		}
	})
}

func TestOption_Filter(t *testing.T) {
	t.Run("filter passes", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		if opt.Filter(func(t *int) bool { return (*t) == 4 }) != opt {
			t.Fail()
		}
	})
	t.Run("filter fails", func(t *testing.T) {
		var opt option.Option[int] = option.Some[int]{Value: 4}
		var target option.Option[int] = option.Nothing[int]{}
		if opt.Filter(func(t *int) bool { return (*t) == 3 }) != target {
			t.Fail()
		}
	})
	t.Run("Nothing returns Nothing", func(t *testing.T) {
		var opt option.Option[int] = option.Nothing[int]{}
		if opt.Filter(func(t *int) bool { return (*t) == 3 }) != opt {
			t.Fail()
		}
	})
}

func TestOption_And(t *testing.T) {
	t.Run("Some, Some", func(t *testing.T) {
		res := option.And[int, string](option.Some[int]{3}, option.Some[string]{"test"})
		expected := option.Some[string]{"test"}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Some, Nothing", func(t *testing.T) {
		res := option.And[int, string](option.Some[int]{3}, option.Nothing[string]{})
		expected := option.Nothing[string]{}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing, Some", func(t *testing.T) {
		res := option.And[int, string](option.Nothing[int]{}, option.Some[string]{"test"})
		expected := option.Nothing[string]{}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing, Nothing", func(t *testing.T) {
		res := option.And[int, string](option.Nothing[int]{}, option.Nothing[string]{})
		expected := option.Nothing[string]{}
		if res != expected {
			t.Fail()
		}
	})
}

func TestOption_AndThen(t *testing.T) {
	t.Run("Some, f returns Some", func(t *testing.T) {
		calls := 0
		res := option.AndThen[int](option.Some[int]{3}, func(t int) option.Option[string] {
			calls++
			return option.Some[string]{"test"}
		})
		expected := option.Some[string]{"test"}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Some, f returns Nothing", func(t *testing.T) {
		calls := 0
		res := option.AndThen[int](option.Some[int]{3}, func(t int) option.Option[string] {
			calls++
			return option.Nothing[string]{}
		})
		expected := option.Nothing[string]{}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Nothing, f not called", func(t *testing.T) {
		calls := 0
		res := option.AndThen[int](option.Nothing[int]{}, func(t int) option.Option[string] {
			calls++
			return option.Nothing[string]{}
		})
		expected := option.Nothing[string]{}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestOption_Or(t *testing.T) {
	t.Run("Some, Some", func(t *testing.T) {
		res := option.Some[int]{3}.Or(option.Some[int]{4})
		expected := option.Some[int]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Some, Nothing", func(t *testing.T) {
		res := option.Some[int]{3}.Or(option.Nothing[int]{})
		expected := option.Some[int]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing, Some", func(t *testing.T) {
		res := option.Nothing[int]{}.Or(option.Some[int]{4})
		expected := option.Some[int]{4}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing, Nothing", func(t *testing.T) {
		res := option.Nothing[int]{}.Or(option.Nothing[int]{})
		expected := option.Nothing[int]{}
		if res != expected {
			t.Fail()
		}
	})
}

func TestOption_OrElse(t *testing.T) {
	t.Run("Some, f not called", func(t *testing.T) {
		calls := 0
		res := option.Some[int]{3}.OrElse(func() option.Option[int] {
			calls++
			return option.Some[int]{4}
		})
		expected := option.Some[int]{3}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Nothing, f returned Some", func(t *testing.T) {
		calls := 0
		res := option.Nothing[int]{}.OrElse(func() option.Option[int] {
			calls++
			return option.Some[int]{4}
		})
		expected := option.Some[int]{4}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Nothing, f returned Nothing", func(t *testing.T) {
		calls := 0
		res := option.Nothing[int]{}.OrElse(func() option.Option[int] {
			calls++
			return option.Nothing[int]{}
		})
		expected := option.Nothing[int]{}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestOption_Xor(t *testing.T) {
	t.Run("Some, Some", func(t *testing.T) {
		res := option.Some[int]{3}.Xor(option.Some[int]{4})
		expected := option.Nothing[int]{}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Some, Nothing", func(t *testing.T) {
		res := option.Some[int]{3}.Xor(option.Nothing[int]{})
		expected := option.Some[int]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing, Some", func(t *testing.T) {
		res := option.Nothing[int]{}.Xor(option.Some[int]{4})
		expected := option.Some[int]{4}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing, Nothing", func(t *testing.T) {
		res := option.Nothing[int]{}.Xor(option.Nothing[int]{})
		expected := option.Nothing[int]{}
		if res != expected {
			t.Fail()
		}
	})
}

func TestOption_Flatten(t *testing.T) {
	t.Run("Some[Some]", func(t *testing.T) {
		res := option.Flatten[int](
			option.Some[option.Option[int]]{option.Some[int]{4}},
		)
		expected := option.Some[int]{4}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Some[Nothing]", func(t *testing.T) {
		res := option.Flatten[int](
			option.Some[option.Option[int]]{option.Nothing[int]{}},
		)
		expected := option.Nothing[int]{}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		res := option.Flatten[int](
			option.Nothing[option.Option[int]]{},
		)
		expected := option.Nothing[int]{}
		if res != expected {
			t.Fail()
		}
	})
}

func TestOption_Map(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		res := option.Map[int](option.Some[int]{4},
			func(t int) string {
				return strconv.Itoa(t)
			},
		)
		expected := option.Some[string]{"4"}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		res := option.Map[int](option.Nothing[int]{},
			func(t int) string {
				return strconv.Itoa(t)
			},
		)
		expected := option.Nothing[string]{}
		if res != expected {
			t.Fail()
		}
	})
}

func TestOption_MapOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		res := option.MapOr[int](option.Some[int]{4}, "default",
			func(t int) string {
				return strconv.Itoa(t)
			},
		)
		expected := option.Some[string]{"4"}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		res := option.MapOr[int](option.Nothing[int]{}, "default",
			func(t int) string {
				return strconv.Itoa(t)
			},
		)
		expected := option.Some[string]{"default"}
		if res != expected {
			t.Fail()
		}
	})
}

func TestOption_MapOrElse(t *testing.T) {
	t.Run("Some, default not called", func(t *testing.T) {
		calls := 0
		res := option.MapOrElse[int](option.Some[int]{4},
			func() string {
				calls++
				return "default"
			},
			func(t int) string {
				return strconv.Itoa(t)
			},
		)
		expected := option.Some[string]{"4"}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Nothing, default called", func(t *testing.T) {
		calls := 0
		res := option.MapOrElse[int](option.Nothing[int]{},
			func() string {
				calls++
				return "default"
			},
			func(t int) string {
				return strconv.Itoa(t)
			},
		)
		expected := option.Some[string]{"default"}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestOption_Match(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		someArmCalls := 0
		nothingArmCalls := 0
		res := option.Match[int](option.Some[int]{4},
			func(t option.Some[int]) string {
				someArmCalls++
				return strconv.Itoa(t.Value)
			},
			func(t option.Nothing[int]) string {
				nothingArmCalls++
				return "nothing"
			},
		)
		expected := "4"
		if res != expected || someArmCalls != 1 || nothingArmCalls != 0 {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		someArmCalls := 0
		nothingArmCalls := 0
		res := option.Match[int](option.Nothing[int]{},
			func(t option.Some[int]) string {
				someArmCalls++
				return strconv.Itoa(t.Value)
			},
			func(t option.Nothing[int]) string {
				nothingArmCalls++
				return "nothing"
			},
		)
		expected := "nothing"
		if res != expected || someArmCalls != 0 || nothingArmCalls != 1 {
			t.Fail()
		}
	})
}

func TestOption_From(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		res := option.From(3, true)
		expected := option.Some[int]{Value: 3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("not ok", func(t *testing.T) {
		res := option.From(3, false)
		expected := option.Nothing[int]{}
		if res != expected {
			t.Fail()
		}
	})
}

func TestOption_FromError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		res := option.FromError(3, nil)
		expected := option.Some[int]{Value: 3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("not nil", func(t *testing.T) {
		res := option.FromError(3, fmt.Errorf("error"))
		expected := option.Nothing[int]{}
		if res != expected {
			t.Fail()
		}
	})
}
