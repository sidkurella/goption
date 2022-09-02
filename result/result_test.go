package result_test

import (
	"strconv"
	"testing"

	"github.com/sidkurella/goption/option"
	"github.com/sidkurella/goption/result"
)

func TestResult_IsOk(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		if !val.IsOk() {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"err val"}
		if val.IsOk() {
			t.Fail()
		}
	})
}

func TestResult_IsErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		if val.IsErr() {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"err val"}
		if !val.IsErr() {
			t.Fail()
		}
	})
}

func TestResult_IsOkAnd(t *testing.T) {
	t.Run("Ok, passes predicate", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		if !val.IsOkAnd(func(t *int) bool { return (*t) == 3 }) {
			t.Fail()
		}
	})
	t.Run("Ok, fails predicate", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		if val.IsOkAnd(func(t *int) bool { return (*t) == 4 }) {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"err val"}
		if val.IsOkAnd(func(t *int) bool { return (*t) == 4 }) {
			t.Fail()
		}
	})
}

func TestResult_IsErrAnd(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		if val.IsErrAnd(func(t *string) bool { return (*t) == "hello" }) {
			t.Fail()
		}
	})
	t.Run("Err, passes predicate", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		if !val.IsErrAnd(func(t *string) bool { return (*t) == "hello" }) {
			t.Fail()
		}
	})
	t.Run("Err, fails predicate", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"world"}
		if val.IsErrAnd(func(t *string) bool { return (*t) == "hello" }) {
			t.Fail()
		}
	})
}

func TestResult_Ok(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		expected := option.Some[int]{Value: 3}
		if val.Ok() != expected {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		expected := option.Nothing[int]{}
		if val.Ok() != expected {
			t.Fail()
		}
	})
}

func TestResult_Err(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		expected := option.Nothing[string]{}
		if val.Err() != expected {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		expected := option.Some[string]{Value: "hello"}
		if val.Err() != expected {
			t.Fail()
		}
	})
}

func TestResult_Unwrap(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		expected := 3
		if val.Unwrap() != expected {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic was expected but did not occur")
			}
		}()

		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		_ = val.Unwrap()
	})
}

func TestResult_UnwrapOr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		expected := 3
		if val.UnwrapOr(4) != expected {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		expected := 4
		if val.UnwrapOr(4) != expected {
			t.Fail()
		}
	})
}

func TestResult_UnwrapOrElse(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		calls := 0
		var val result.Result[int, string] = result.Ok[int, string]{3}
		expected := 3
		if val.UnwrapOrElse(func(_ string) int {
			calls++
			return 4
		}) != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		calls := 0
		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		expected := 5
		if val.UnwrapOrElse(func(s string) int {
			calls++
			return len(s)
		}) != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestResult_UnwrapErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Panic was expected but did not occur")
			}
		}()

		var val result.Result[int, string] = result.Ok[int, string]{3}
		_ = val.UnwrapErr()
	})
	t.Run("Err", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		expected := "hello"
		if val.UnwrapErr() != expected {
			t.Fail()
		}
	})
}

func TestResult_UnwrapErrOr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		expected := "world"
		if val.UnwrapErrOr("world") != expected {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		expected := "hello"
		if val.UnwrapErrOr("world") != expected {
			t.Fail()
		}
	})
}

func TestResult_UnwrapErrOrElse(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		calls := 0
		var val result.Result[int, string] = result.Ok[int, string]{3}
		expected := "3"
		if val.UnwrapErrOrElse(func(i int) string {
			calls++
			return strconv.Itoa(i)
		}) != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		calls := 0
		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		expected := "hello"
		if val.UnwrapErrOrElse(func(i int) string {
			calls++
			return strconv.Itoa(i)
		}) != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestResult_Expect(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		var val result.Result[int, string] = result.Ok[int, string]{3}
		expected := 3
		if val.Expect("don't panic") != expected {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		msg := "panic expected"
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("Panic was expected but did not occur")
			} else if r != msg {
				t.Errorf("Panic was expected but message was %v (did not match expected msg %v)", r, msg)
			}
		}()

		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		_ = val.Expect(msg)
	})
}

func TestResult_ExpectErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		msg := "panic expected"
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("Panic was expected but did not occur")
			} else if r != msg {
				t.Errorf("Panic was expected but message was %v (did not match expected msg %v)", r, msg)
			}
		}()

		var val result.Result[int, string] = result.Ok[int, string]{3}
		_ = val.ExpectErr(msg)
	})
	t.Run("Err", func(t *testing.T) {
		var val result.Result[int, string] = result.Err[int, string]{"hello"}
		expected := "hello"
		if val.ExpectErr("don't panic") != expected {
			t.Fail()
		}
	})
}

func TestResult_And(t *testing.T) {
	t.Run("Ok, Ok", func(t *testing.T) {
		res := result.And[int, string, float64](
			result.Ok[int, string]{3},
			result.Ok[float64, string]{2.0},
		)
		expected := result.Ok[float64, string]{2.0}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Ok, Err", func(t *testing.T) {
		res := result.And[int, string, float64](
			result.Ok[int, string]{3},
			result.Err[float64, string]{"err"},
		)
		expected := result.Err[float64, string]{"err"}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Err, Ok", func(t *testing.T) {
		res := result.And[int, string, float64](
			result.Err[int, string]{"err"},
			result.Ok[float64, string]{2.0},
		)
		expected := result.Err[float64, string]{"err"}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Err, Err", func(t *testing.T) {
		res := result.And[int, string, float64](
			result.Err[int, string]{"some err"},
			result.Err[float64, string]{"other err"},
		)
		expected := result.Err[float64, string]{"some err"}
		if res != expected {
			t.Fail()
		}
	})
}

func TestResult_AndThen(t *testing.T) {
	t.Run("Ok, f returns Ok", func(t *testing.T) {
		calls := 0
		res := result.AndThen[int, string](
			result.Ok[int, string]{3},
			func(i int) result.Result[float64, string] {
				calls++
				return result.Ok[float64, string]{2.0}
			},
		)
		expected := result.Ok[float64, string]{2.0}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Ok, f returns Err", func(t *testing.T) {
		calls := 0
		res := result.AndThen[int, string](
			result.Ok[int, string]{3},
			func(i int) result.Result[float64, string] {
				calls++
				return result.Err[float64, string]{"hello"}
			},
		)
		expected := result.Err[float64, string]{"hello"}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Err, f not called", func(t *testing.T) {
		calls := 0
		res := result.AndThen[int, string](
			result.Err[int, string]{"hello"},
			func(i int) result.Result[float64, string] {
				calls++
				return result.Err[float64, string]{"world"}
			},
		)
		expected := result.Err[float64, string]{"hello"}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestResult_Or(t *testing.T) {
	t.Run("Ok, Ok", func(t *testing.T) {
		res := result.Or[int, string, float64](
			result.Ok[int, string]{3},
			result.Ok[int, float64]{4},
		)
		expected := result.Ok[int, float64]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Ok, Err", func(t *testing.T) {
		res := result.Or[int, string, float64](
			result.Ok[int, string]{3},
			result.Err[int, float64]{2.0},
		)
		expected := result.Ok[int, float64]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Err, Ok", func(t *testing.T) {
		res := result.Or[int, string, float64](
			result.Err[int, string]{"err"},
			result.Ok[int, float64]{2},
		)
		expected := result.Ok[int, float64]{2}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Err, Err", func(t *testing.T) {
		res := result.Or[int, string, float64](
			result.Err[int, string]{"some err"},
			result.Err[int, float64]{1.0},
		)
		expected := result.Err[int, float64]{1.0}
		if res != expected {
			t.Fail()
		}
	})
}

func TestResult_OrElse(t *testing.T) {
	t.Run("Ok, f not called", func(t *testing.T) {
		calls := 0
		res := result.OrElse[int, string](
			result.Ok[int, string]{3},
			func(_ string) result.Result[int, float64] {
				calls++
				return result.Err[int, float64]{2.0}
			},
		)
		expected := result.Ok[int, float64]{3}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Err, f returns Ok", func(t *testing.T) {
		calls := 0
		res := result.OrElse[int, string](
			result.Err[int, string]{"hello"},
			func(s string) result.Result[int, float64] {
				calls++
				return result.Ok[int, float64]{len(s)}
			},
		)
		expected := result.Ok[int, float64]{5}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Err, f returns Err", func(t *testing.T) {
		calls := 0
		res := result.OrElse[int, string](
			result.Err[int, string]{"hello"},
			func(s string) result.Result[int, float64] {
				calls++
				return result.Err[int, float64]{float64(len(s))}
			},
		)
		expected := result.Err[int, float64]{5}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestResult_Flatten(t *testing.T) {
	t.Run("Ok[Ok]", func(t *testing.T) {
		res := result.Flatten[int, string](
			result.Ok[result.Result[int, string], string]{
				result.Ok[int, string]{3},
			},
		)
		expected := result.Ok[int, string]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Ok[Err]", func(t *testing.T) {
		res := result.Flatten[int, string](
			result.Ok[result.Result[int, string], string]{
				result.Err[int, string]{"err"},
			},
		)
		expected := result.Err[int, string]{"err"}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		res := result.Flatten[int, string](
			result.Err[result.Result[int, string], string]{"hello"},
		)
		expected := result.Err[int, string]{"hello"}
		if res != expected {
			t.Fail()
		}
	})
}

func TestResult_Map(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		calls := 0
		res := result.Map[int, string](
			result.Ok[int, string]{3},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := result.Ok[float64, string]{4}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		calls := 0
		res := result.Map[int, string](
			result.Err[int, string]{"err"},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := result.Err[float64, string]{"err"}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestResult_MapErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		calls := 0
		res := result.MapErr[int, string](
			result.Ok[int, string]{3},
			func(i string) float64 {
				calls++
				return float64(len(i) + 1)
			},
		)
		expected := result.Ok[int, float64]{3}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		calls := 0
		res := result.MapErr[int, string](
			result.Err[int, string]{"err"},
			func(i string) float64 {
				calls++
				return float64(len(i) + 1)
			},
		)
		expected := result.Err[int, float64]{4}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}

func TestResult_MapOr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		calls := 0
		res := result.MapOr[int, string](
			result.Ok[int, string]{3},
			float64(600),
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := result.Ok[float64, string]{4}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		calls := 0
		res := result.MapOr[int, string](
			result.Err[int, string]{"err"},
			float64(600),
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := result.Ok[float64, string]{600}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
}

func TestResult_MapOrElse(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		calls := 0
		defaultCalls := 0
		res := result.MapOrElse[int, string](
			result.Ok[int, string]{3},
			func(e string) float64 {
				defaultCalls++
				return float64(len(e))
			},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := result.Ok[float64, string]{4}
		if res != expected || calls != 1 || defaultCalls != 0 {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		calls := 0
		defaultCalls := 0
		res := result.MapOrElse[int, string](
			result.Err[int, string]{"hello world"},
			func(e string) float64 {
				defaultCalls++
				return float64(len(e))
			},
			func(i int) float64 {
				calls++
				return float64(i + 1)
			},
		)
		expected := result.Ok[float64, string]{11}
		if res != expected || calls != 0 || defaultCalls != 1 {
			t.Fail()
		}
	})
}

func TestResult_Match(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		okCalls := 0
		errCalls := 0
		res := result.Match[int, string](result.Ok[int, string]{3},
			func(o result.Ok[int, string]) float64 {
				okCalls++
				return float64(o.Value + 1)
			},
			func(e result.Err[int, string]) float64 {
				errCalls++
				return float64(600)
			},
		)
		expected := float64(4)
		if res != expected || okCalls != 1 || errCalls != 0 {
			t.Fail()
		}
	})
	t.Run("Err", func(t *testing.T) {
		okCalls := 0
		errCalls := 0
		res := result.Match[int, string](result.Err[int, string]{"hello world"},
			func(o result.Ok[int, string]) float64 {
				okCalls++
				return float64(600)
			},
			func(e result.Err[int, string]) float64 {
				errCalls++
				return float64(len(e.Value) + 1)
			},
		)
		expected := float64(12)
		if res != expected || okCalls != 0 || errCalls != 1 {
			t.Fail()
		}
	})
}

func TestResult_OkOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		res := result.OkOr[int](option.Some[int]{Value: 3}, "err")
		expected := result.Ok[int, string]{3}
		if res != expected {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		res := result.OkOr[int](option.Nothing[int]{}, "err")
		expected := result.Err[int, string]{"err"}
		if res != expected {
			t.Fail()
		}
	})
}

func TestResult_OkOrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		calls := 0
		res := result.OkOrElse[int](option.Some[int]{Value: 3}, func() string {
			calls++
			return "hello"
		})
		expected := result.Ok[int, string]{3}
		if res != expected || calls != 0 {
			t.Fail()
		}
	})
	t.Run("Nothing", func(t *testing.T) {
		calls := 0
		res := result.OkOrElse[int](option.Nothing[int]{}, func() string {
			calls++
			return "hello"
		})
		expected := result.Err[int, string]{"hello"}
		if res != expected || calls != 1 {
			t.Fail()
		}
	})
}
