package either_test

import (
	"strconv"
	"testing"

	"github.com/sidkurella/goption/either"
	"github.com/sidkurella/goption/option"
)

func TestEither_IsRight(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		if !val.IsRight() {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"err val"}
		if val.IsRight() {
			t.Fail()
		}
	})
}

func TestEither_IsLeft(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		if val.IsLeft() {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"err val"}
		if !val.IsLeft() {
			t.Fail()
		}
	})
}

func TestEither_IsRightAnd(t *testing.T) {
	t.Run("Right, passes predicate", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		if !val.IsRightAnd(func(t *int) bool { return (*t) == 3 }) {
			t.Fail()
		}
	})
	t.Run("Right, fails predicate", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		if val.IsRightAnd(func(t *int) bool { return (*t) == 4 }) {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"err val"}
		if val.IsRightAnd(func(t *int) bool { return (*t) == 4 }) {
			t.Fail()
		}
	})
}

func TestEither_IsLeftAnd(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		if val.IsLeftAnd(func(t *string) bool { return (*t) == "hello" }) {
			t.Fail()
		}
	})
	t.Run("Left, passes predicate", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		if !val.IsLeftAnd(func(t *string) bool { return (*t) == "hello" }) {
			t.Fail()
		}
	})
	t.Run("Left, fails predicate", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"world"}
		if val.IsLeftAnd(func(t *string) bool { return (*t) == "hello" }) {
			t.Fail()
		}
	})
}

func TestEither_Right(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		expected := option.Some[int]{Value: 3}
		if val.Right() != expected {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		expected := option.Nothing[int]{}
		if val.Right() != expected {
			t.Fail()
		}
	})
}

func TestEither_Left(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		expected := option.Nothing[string]{}
		if val.Left() != expected {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		expected := option.Some[string]{Value: "hello"}
		if val.Left() != expected {
			t.Fail()
		}
	})
}

func TestEither_Unwrap(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		expected := 3
		if val.Unwrap() != expected {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic was expected but did not occur")
			}
		}()

		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		_ = val.Unwrap()
	})
}

func TestEither_UnwrapOr(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		expected := 3
		if val.UnwrapOr(4) != expected {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		expected := 4
		if val.UnwrapOr(4) != expected {
			t.Fail()
		}
	})
}

func TestEither_UnwrapOrElse(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		calls := 0
		var val either.Either[string, int] = either.Right[string, int]{3}
		expected := 3
		if val.UnwrapOrElse(func(_ string) int {
			calls++
			return 4
		}) != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		calls := 0
		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		expected := 5
		if val.UnwrapOrElse(func(s string) int {
			calls++
			return len(s)
		}) != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestEither_UnwrapLeft(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic was expected but did not occur")
			}
		}()

		var val either.Either[string, int] = either.Right[string, int]{3}
		_ = val.UnwrapLeft()
	})
	t.Run("Left", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		expected := "hello"
		if val.UnwrapLeft() != expected {
			t.Fail()
		}
	})
}

func TestEither_UnwrapLeftOr(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		expected := "world"
		if val.UnwrapLeftOr("world") != expected {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		expected := "hello"
		if val.UnwrapLeftOr("world") != expected {
			t.Fail()
		}
	})
}

func TestEither_UnwrapLeftOrElse(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		calls := 0
		var val either.Either[string, int] = either.Right[string, int]{3}
		expected := "3"
		if val.UnwrapLeftOrElse(func(i int) string {
			calls++
			return strconv.Itoa(i)
		}) != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		calls := 0
		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		expected := "hello"
		if val.UnwrapLeftOrElse(func(i int) string {
			calls++
			return strconv.Itoa(i)
		}) != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestEither_Expect(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		var val either.Either[string, int] = either.Right[string, int]{3}
		expected := 3
		if val.Expect("don't panic") != expected {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		msg := "panic expected"
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("Panic was expected but did not occur")
			} else if r != msg {
				t.Errorf("Panic was expected but message was %v (did not match expected msg %v)", r, msg)
			}
		}()

		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		_ = val.Expect(msg)
	})
}

func TestEither_ExpectLeft(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		msg := "panic expected"
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("Panic was expected but did not occur")
			} else if r != msg {
				t.Errorf("Panic was expected but message was %v (did not match expected msg %v)", r, msg)
			}
		}()

		var val either.Either[string, int] = either.Right[string, int]{3}
		_ = val.ExpectLeft(msg)
	})
	t.Run("Left", func(t *testing.T) {
		var val either.Either[string, int] = either.Left[string, int]{"hello"}
		expected := "hello"
		if val.ExpectLeft("don't panic") != expected {
			t.Fail()
		}
	})
}

func TestEither_And(t *testing.T) {
	t.Run("Right, Right", func(t *testing.T) {
		res := either.And[string, int, float64](
			either.Right[string, int]{3},
			either.Right[string, float64]{2.0},
		)
		expected := either.Right[string, float64]{2.0}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Right, Left", func(t *testing.T) {
		res := either.And[string, int, float64](
			either.Right[string, int]{3},
			either.Left[string, float64]{"err"},
		)
		expected := either.Left[string, float64]{"err"}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Left, Right", func(t *testing.T) {
		res := either.And[string, int, float64](
			either.Left[string, int]{"err"},
			either.Right[string, float64]{2.0},
		)
		expected := either.Left[string, float64]{"err"}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Left, Left", func(t *testing.T) {
		res := either.And[string, int, float64](
			either.Left[string, int]{"some err"},
			either.Left[string, float64]{"other err"},
		)
		expected := either.Left[string, float64]{"some err"}
		if res != expected {
			t.Fail()
		}
	})
}

func TestEither_AndThen(t *testing.T) {
	t.Run("Right, f returns Right", func(t *testing.T) {
		calls := 0
		res := either.AndThen[string, int](
			either.Right[string, int]{3},
			func(i int) either.Either[string, float64] {
				calls++
				return either.Right[string, float64]{2.0}
			},
		)
		expected := either.Right[string, float64]{2.0}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Right, f returns Left", func(t *testing.T) {
		calls := 0
		res := either.AndThen[string, int](
			either.Right[string, int]{3},
			func(i int) either.Either[string, float64] {
				calls++
				return either.Left[string, float64]{"hello"}
			},
		)
		expected := either.Left[string, float64]{"hello"}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Left, f not called", func(t *testing.T) {
		calls := 0
		res := either.AndThen[string, int](
			either.Left[string, int]{"hello"},
			func(i int) either.Either[string, float64] {
				calls++
				return either.Left[string, float64]{"world"}
			},
		)
		expected := either.Left[string, float64]{"hello"}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestEither_Or(t *testing.T) {
	t.Run("Right, Right", func(t *testing.T) {
		res := either.Or[string, float64, int](
			either.Right[string, int]{3},
			either.Right[float64, int]{4},
		)
		expected := either.Right[float64, int]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Right, Left", func(t *testing.T) {
		res := either.Or[string, float64, int](
			either.Right[string, int]{3},
			either.Left[float64, int]{2.0},
		)
		expected := either.Right[float64, int]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Left, Right", func(t *testing.T) {
		res := either.Or[string, float64, int](
			either.Left[string, int]{"err"},
			either.Right[float64, int]{2},
		)
		expected := either.Right[float64, int]{2}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Left, Left", func(t *testing.T) {
		res := either.Or[string, float64, int](
			either.Left[string, int]{"some err"},
			either.Left[float64, int]{1.0},
		)
		expected := either.Left[float64, int]{1.0}
		if res != expected {
			t.Fail()
		}
	})
}

func TestEither_OrElse(t *testing.T) {
	t.Run("Right, f not called", func(t *testing.T) {
		calls := 0
		res := either.OrElse[string, float64, int](
			either.Right[string, int]{3},
			func(_ string) either.Either[float64, int] {
				calls++
				return either.Left[float64, int]{2.0}
			},
		)
		expected := either.Right[float64, int]{3}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Left, f returns Right", func(t *testing.T) {
		calls := 0
		res := either.OrElse[string, float64, int](
			either.Left[string, int]{"hello"},
			func(s string) either.Either[float64, int] {
				calls++
				return either.Right[float64, int]{len(s)}
			},
		)
		expected := either.Right[float64, int]{5}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Left, f returns Left", func(t *testing.T) {
		calls := 0
		res := either.OrElse[string, float64, int](
			either.Left[string, int]{"hello"},
			func(s string) either.Either[float64, int] {
				calls++
				return either.Left[float64, int]{float64(len(s))}
			},
		)
		expected := either.Left[float64, int]{5}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestEither_Flatten(t *testing.T) {
	t.Run("Right[Right]", func(t *testing.T) {
		res := either.Flatten[string, int](
			either.Right[string, either.Either[string, int]]{
				either.Right[string, int]{3},
			},
		)
		expected := either.Right[string, int]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Right[Left]", func(t *testing.T) {
		res := either.Flatten[string, int](
			either.Right[string, either.Either[string, int]]{
				either.Left[string, int]{"err"},
			},
		)
		expected := either.Left[string, int]{"err"}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		res := either.Flatten[string, int](
			either.Left[string, either.Either[string, int]]{"hello"},
		)
		expected := either.Left[string, int]{"hello"}
		if res != expected {
			t.Fail()
		}
	})
}

func TestEither_Map(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		calls := 0
		res := either.Map[string, int](
			either.Right[string, int]{3},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.Right[string, float64]{4}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		calls := 0
		res := either.Map[string, int](
			either.Left[string, int]{"err"},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.Left[string, float64]{"err"}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestEither_MapLeft(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		calls := 0
		res := either.MapLeft[string, float64, int](
			either.Right[string, int]{3},
			func(i string) float64 {
				calls++
				return float64(len(i) + 1)
			},
		)
		expected := either.Right[float64, int]{3}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		calls := 0
		res := either.MapLeft[string, float64, int](
			either.Left[string, int]{"err"},
			func(i string) float64 {
				calls++
				return float64(len(i) + 1)
			},
		)
		expected := either.Left[float64, int]{4}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestEither_MapOr(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		calls := 0
		res := either.MapOr[string, int](
			either.Right[string, int]{3},
			float64(600),
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.Right[string, float64]{4}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		calls := 0
		res := either.MapOr[string, int](
			either.Left[string, int]{"err"},
			float64(600),
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.Right[string, float64]{600}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestEither_MapOrElse(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		calls := 0
		defaultCalls := 0
		res := either.MapOrElse[string, int](
			either.Right[string, int]{3},
			func(e string) float64 {
				defaultCalls++
				return float64(len(e))
			},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.Right[string, float64]{4}
		if res != expected || calls != 1 || defaultCalls != 0 {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		calls := 0
		defaultCalls := 0
		res := either.MapOrElse[string, int](
			either.Left[string, int]{"hello world"},
			func(e string) float64 {
				defaultCalls++
				return float64(len(e))
			},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := either.Right[string, float64]{11}
		if res != expected || calls != 0 || defaultCalls != 1 {
			t.Fail()
		}
	})
}

func TestEither_Match(t *testing.T) {
	t.Run("Right", func(t *testing.T) {
		rightCalls := 0
		leftCalls := 0
		res := either.Match[string, int](either.Right[string, int]{3},
			func(l either.Left[string, int]) float64 {
				leftCalls++
				return float64(600)
			},
			func(r either.Right[string, int]) float64 {
				rightCalls++
				return float64(r.Value + 1)
			},
		)
		expected := float64(4)
		if res != expected || rightCalls != 1 || leftCalls != 0 {
			t.Fail()
		}
	})
	t.Run("Left", func(t *testing.T) {
		rightCalls := 0
		leftCalls := 0
		res := either.Match[string, int](either.Left[string, int]{"hello world"},
			func(l either.Left[string, int]) float64 {
				leftCalls++
				return float64(len(l.Value) + 1)
			},
			func(r either.Right[string, int]) float64 {
				rightCalls++
				return float64(600)
			},
		)
		expected := float64(12)
		if res != expected || rightCalls != 0 || leftCalls != 1 {
			t.Fail()
		}
	})
}

func TestEither_RightOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		res := either.RightOr[string, int](option.Some[int]{Value: 3}, "err")
		expected := either.Right[string, int]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		res := either.RightOr[string, int](option.Nothing[int]{}, "err")
		expected := either.Left[string, int]{"err"}
		if res != expected {
			t.Fail()
		}
	})
}

func TestEither_RightOrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		calls := 0
		res := either.RightOrElse[string, int](option.Some[int]{Value: 3}, func() string {
			calls++
			return "hello"
		})
		expected := either.Right[string, int]{3}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		calls := 0
		res := either.RightOrElse[string, int](option.Nothing[int]{}, func() string {
			calls++
			return "hello"
		})
		expected := either.Left[string, int]{"hello"}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}
