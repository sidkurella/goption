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
