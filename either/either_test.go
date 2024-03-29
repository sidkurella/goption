package either_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/sidkurella/goption/either"
	"github.com/sidkurella/goption/option"
)

func TestEither_IsFirst(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		if !val.IsFirst() {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("second val")
		if val.IsFirst() {
			t.Fail()
		}
	})
}

func TestEither_IsSecond(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		if val.IsSecond() {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("second val")
		if !val.IsSecond() {
			t.Fail()
		}
	})
}

func TestEither_IsFirstAnd(t *testing.T) {
	t.Run("First, passes predicate", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		if !val.IsFirstAnd(func(t *int) bool { return (*t) == 3 }) {
			t.Fail()
		}
	})
	t.Run("First, fails predicate", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		if val.IsFirstAnd(func(t *int) bool { return (*t) == 4 }) {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("second val")
		if val.IsFirstAnd(func(t *int) bool { return (*t) == 4 }) {
			t.Fail()
		}
	})
}

func TestEither_IsSecondAnd(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		if val.IsSecondAnd(func(t *string) bool { return (*t) == "hello" }) {
			t.Fail()
		}
	})
	t.Run("Second, passes predicate", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("hello")
		if !val.IsSecondAnd(func(t *string) bool { return (*t) == "hello" }) {
			t.Fail()
		}
	})
	t.Run("Second, fails predicate", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("world")
		if val.IsSecondAnd(func(t *string) bool { return (*t) == "hello" }) {
			t.Fail()
		}
	})
}

func TestEither_First(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		expected := option.Some(3)
		if val.First() != expected {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("hello")
		expected := option.Nothing[int]()
		if val.First() != expected {
			t.Fail()
		}
	})
}

func TestEither_Second(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		expected := option.Nothing[string]()
		if val.Second() != expected {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("hello")
		expected := option.Some("hello")
		if val.Second() != expected {
			t.Fail()
		}
	})
}

func TestEither_Unwrap(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		expected := 3
		if val.Unwrap() != expected {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic was expected but did not occur")
			}
		}()

		var val either.Either[int, string] = either.Second[int]("hello")
		_ = val.Unwrap()
	})
}

func TestEither_UnwrapOr(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		expected := 3
		if val.UnwrapOr(4) != expected {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("hello")
		expected := 4
		if val.UnwrapOr(4) != expected {
			t.Fail()
		}
	})
}

func TestEither_UnwrapOrElse(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		calls := 0
		var val either.Either[int, string] = either.First[int, string](3)
		expected := 3
		if val.UnwrapOrElse(func(_ string) int {
			calls++
			return 4
		}) != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		calls := 0
		var val either.Either[int, string] = either.Second[int]("hello")
		expected := 5
		if val.UnwrapOrElse(func(s string) int {
			calls++
			return len(s)
		}) != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestEither_UnwrapSecond(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic was expected but did not occur")
			}
		}()

		var val either.Either[int, string] = either.First[int, string](3)
		_ = val.UnwrapSecond()
	})
	t.Run("Second", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("hello")
		expected := "hello"
		if val.UnwrapSecond() != expected {
			t.Fail()
		}
	})
}

func TestEither_UnwrapSecondOr(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		expected := "world"
		if val.UnwrapSecondOr("world") != expected {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("hello")
		expected := "hello"
		if val.UnwrapSecondOr("world") != expected {
			t.Fail()
		}
	})
}

func TestEither_UnwrapSecondOrElse(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		calls := 0
		var val either.Either[int, string] = either.First[int, string](3)
		expected := "3"
		if val.UnwrapSecondOrElse(func(i int) string {
			calls++
			return strconv.Itoa(i)
		}) != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		calls := 0
		var val either.Either[int, string] = either.Second[int]("hello")
		expected := "hello"
		if val.UnwrapSecondOrElse(func(i int) string {
			calls++
			return strconv.Itoa(i)
		}) != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestEither_Expect(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		var val either.Either[int, string] = either.First[int, string](3)
		expected := 3
		if val.Expect("don't panic") != expected {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		msg := "panic expected"
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("Panic was expected but did not occur")
			} else if r != msg {
				t.Errorf("Panic was expected but message was %v (did not match expected msg %v)", r, msg)
			}
		}()

		var val either.Either[int, string] = either.Second[int]("hello")
		_ = val.Expect(msg)
	})
}

func TestEither_ExpectSecond(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		msg := "panic expected"
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("Panic was expected but did not occur")
			} else if r != msg {
				t.Errorf("Panic was expected but message was %v (did not match expected msg %v)", r, msg)
			}
		}()

		var val either.Either[int, string] = either.First[int, string](3)
		_ = val.ExpectSecond(msg)
	})
	t.Run("Second", func(t *testing.T) {
		var val either.Either[int, string] = either.Second[int]("hello")
		expected := "hello"
		if val.ExpectSecond("don't panic") != expected {
			t.Fail()
		}
	})
}

func TestEither_And(t *testing.T) {
	t.Run("First, First", func(t *testing.T) {
		res := either.And(
			either.First[int, string](3),
			either.First[float64, string](2.0),
		)
		expected := either.First[float64, string](2.0)
		if res != expected {
			t.Fail()
		}
	})
	t.Run("First, Second", func(t *testing.T) {
		res := either.And(
			either.First[int, string](3),
			either.Second[float64]("second"),
		)
		expected := either.Second[float64]("second")
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Second, First", func(t *testing.T) {
		res := either.And(
			either.Second[int]("second"),
			either.First[float64, string](2.0),
		)
		expected := either.Second[float64]("second")
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Second, Second", func(t *testing.T) {
		res := either.And(
			either.Second[int]("some second"),
			either.Second[float64]("other second"),
		)
		expected := either.Second[float64]("some second")
		if res != expected {
			t.Fail()
		}
	})
}

func TestEither_AndThen(t *testing.T) {
	t.Run("First, f returns First", func(t *testing.T) {
		calls := 0
		res := either.AndThen(
			either.First[int, string](3),
			func(i int) either.Either[float64, string] {
				calls++
				return either.First[float64, string](2.0)
			},
		)
		expected := either.First[float64, string](2.0)
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("First, f returns Second", func(t *testing.T) {
		calls := 0
		res := either.AndThen(
			either.First[int, string](3),
			func(i int) either.Either[float64, string] {
				calls++
				return either.Second[float64]("hello")
			},
		)
		expected := either.Second[float64]("hello")
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Second, f not called", func(t *testing.T) {
		calls := 0
		res := either.AndThen(
			either.Second[int]("hello"),
			func(i int) either.Either[float64, string] {
				calls++
				return either.Second[float64]("world")
			},
		)
		expected := either.Second[float64]("hello")
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestEither_Or(t *testing.T) {
	t.Run("First, First", func(t *testing.T) {
		res := either.Or(
			either.First[int, string](3),
			either.First[int, float64](4),
		)
		expected := either.First[int, float64](3)
		if res != expected {
			t.Fail()
		}
	})
	t.Run("First, Second", func(t *testing.T) {
		res := either.Or(
			either.First[int, string](3),
			either.Second[int](2.0),
		)
		expected := either.First[int, float64](3)
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Second, First", func(t *testing.T) {
		res := either.Or(
			either.Second[int]("second"),
			either.First[int, float64](2),
		)
		expected := either.First[int, float64](2)
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Second, Second", func(t *testing.T) {
		res := either.Or(
			either.Second[int]("some second"),
			either.Second[int](1.0),
		)
		expected := either.Second[int](1.0)
		if res != expected {
			t.Fail()
		}
	})
}

func TestEither_OrElse(t *testing.T) {
	t.Run("First, f not called", func(t *testing.T) {
		calls := 0
		res := either.OrElse(
			either.First[int, string](3),
			func(_ string) either.Either[int, float64] {
				calls++
				return either.Second[int](2.0)
			},
		)
		expected := either.First[int, float64](3)
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Second, f returns First", func(t *testing.T) {
		calls := 0
		res := either.OrElse(
			either.Second[int]("hello"),
			func(s string) either.Either[int, float64] {
				calls++
				return either.First[int, float64](len(s))
			},
		)
		expected := either.First[int, float64](5)
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Second, f returns Second", func(t *testing.T) {
		calls := 0
		res := either.OrElse(
			either.Second[int]("hello"),
			func(s string) either.Either[int, float64] {
				calls++
				return either.Second[int](float64(len(s)))
			},
		)
		expected := either.Second[int, float64](5)
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestEither_Flatten(t *testing.T) {
	t.Run("First[First]", func(t *testing.T) {
		res := either.Flatten(
			either.First[either.Either[int, string], string](
				either.First[int, string](3),
			),
		)
		expected := either.First[int, string](3)
		if res != expected {
			t.Fail()
		}
	})
	t.Run("First[Second]", func(t *testing.T) {
		res := either.Flatten(
			either.First[either.Either[int, string], string](
				either.Second[int]("second"),
			),
		)
		expected := either.Second[int]("second")
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		res := either.Flatten(
			either.Second[either.Either[int, string]]("hello"),
		)
		expected := either.Second[int]("hello")
		if res != expected {
			t.Fail()
		}
	})
}

func TestEither_Map(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		calls := 0
		res := either.Map(
			either.First[int, string](3),
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.First[float64, string](4)
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		calls := 0
		res := either.Map(
			either.Second[int]("second"),
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.Second[float64]("second")
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestEither_MapSecond(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		calls := 0
		res := either.MapSecond(
			either.First[int, string](3),
			func(i string) float64 {
				calls++
				return float64(len(i) + 1)
			},
		)
		expected := either.First[int, float64](int(3))
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		calls := 0
		res := either.MapSecond(
			either.Second[int]("second"),
			func(i string) float64 {
				calls++
				return float64(len(i) + 1)
			},
		)
		expected := either.Second[int](float64(7))
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestEither_MapOr(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		calls := 0
		res := either.MapOr(
			either.First[int, string](3),
			float64(600),
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.First[float64, string](4)
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		calls := 0
		res := either.MapOr(
			either.Second[int]("second"),
			float64(600),
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.First[float64, string](600)
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestEither_MapOrElse(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		calls := 0
		defaultCalls := 0
		res := either.MapOrElse(
			either.First[int, string](3),
			func(e string) float64 {
				defaultCalls++
				return float64(len(e))
			},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.First[float64, string](4)
		if res != expected || calls != 1 || defaultCalls != 0 {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		calls := 0
		defaultCalls := 0
		res := either.MapOrElse(
			either.Second[int]("hello world"),
			func(e string) float64 {
				defaultCalls++
				return float64(len(e))
			},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.First[float64, string](11)
		if res != expected || calls != 0 || defaultCalls != 1 {
			t.Fail()
		}
	})
}

func TestEither_Match(t *testing.T) {
	t.Run("First", func(t *testing.T) {
		firstCalls := 0
		secondCalls := 0
		res := either.Match(either.First[int, string](3),
			func(f int) float64 {
				firstCalls++
				return float64(f + 1)
			},
			func(s string) float64 {
				secondCalls++
				return float64(600)
			},
		)
		expected := float64(4)
		if res != expected || firstCalls != 1 || secondCalls != 0 {
			t.Fail()
		}
	})
	t.Run("Second", func(t *testing.T) {
		firstCalls := 0
		secondCalls := 0
		res := either.Match(either.Second[int]("hello world"),
			func(f int) float64 {
				firstCalls++
				return float64(600)
			},
			func(s string) float64 {
				secondCalls++
				return float64(len(s) + 1)
			},
		)
		expected := float64(12)
		if res != expected || firstCalls != 0 || secondCalls != 1 {
			t.Fail()
		}
	})
}

func TestEither_FirstOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		res := either.FirstOr(option.Some(3), "second")
		expected := either.First[int, string](3)
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		res := either.FirstOr(option.Nothing[int](), "second")
		expected := either.Second[int]("second")
		if res != expected {
			t.Fail()
		}
	})
}

func TestEither_FirstOrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		calls := 0
		res := either.FirstOrElse(option.Some(3), func() string {
			calls++
			return "hello"
		})
		expected := either.First[int, string](3)
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		calls := 0
		res := either.FirstOrElse(option.Nothing[int](), func() string {
			calls++
			return "hello"
		})
		expected := either.Second[int]("hello")
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestEither_From(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		res := either.From(3, nil)
		expected := either.First[int, error](3)
		if res != expected {
			t.Fail()
		}
	})
	t.Run("not nil", func(t *testing.T) {
		err := fmt.Errorf("error")
		res := either.From(3, err)
		expected := either.Second[int](err)
		if res != expected {
			t.Fail()
		}
	})
}
