package result

import (
	"fmt"

	"github.com/sidkurella/goption/option"
)

// Implements Monad[Result[T, E], Result[U, E], T].
type ResultMonad[T any, E any, U any] struct {
}

func (m ResultMonad[T, E, U]) Bind(val Result[T, E], f func(T) Result[U, E]) Result[U, E] {
	return Match(val,
		func(ok T) Result[U, E] {
			return f(ok)
		},
		func(err E) Result[U, E] {
			return Err[U](err)
		},
	)
}

func (m ResultMonad[T1, E, T2]) Return(val T1) Result[T1, E] {
	return Ok[T1, E](val)
}

type resultVariant int

const (
	resultVariantErr resultVariant = iota
	resultVariantOk
)

// Result type. Represents either the success variant Ok(T) containing a success value
// or an Err(E) variant containing an error type. The default value is Err(*new(E))
// (i.e. Err variant containing the zero value of E).
type Result[T any, E any] struct {
	variant resultVariant
	ok      T
	err     E
}

//=====================================================

func Ok[T any, E any](t T) Result[T, E] {
	return Result[T, E]{
		variant: resultVariantOk,
		ok:      t,
	}
}

func Err[T any, E any](e E) Result[T, E] {
	return Result[T, E]{
		variant: resultVariantErr,
		err:     e,
	}
}

// Returns true if the result is Ok.
func (e Result[T, E]) IsOk() bool {
	return e.variant == resultVariantOk
}

// Returns true if the result is Ok and matches the given predicate.
func (e Result[T, E]) IsOkAnd(pred func(*T) bool) bool {
	return e.IsOk() && pred(&e.ok)
}

// Returns true if the result is Err.
func (e Result[T, E]) IsErr() bool {
	return e.variant == resultVariantErr
}

// Returns true if the result is Err and matches the given predicate.
func (e Result[T, E]) IsErrAnd(pred func(*E) bool) bool {
	return e.IsErr() && pred(&e.err)
}

// Converts from Result[T, E] to Option[T].
// Converts self into an Option[T], discarding the Err value, if any.
func (e Result[T, E]) Ok() option.Option[T] {
	if e.IsOk() {
		return option.Some(e.ok)
	}
	return option.Nothing[T]()
}

// Converts from Result[T, E] to Option[E].
// Converts self into an Option[E], and discarding the Ok value, if any.
func (e Result[T, E]) Err() option.Option[E] {
	if e.IsErr() {
		return option.Some(e.err)
	}
	return option.Nothing[E]()
}

// Returns the contained Ok value. Panics with the given message if the result is not Ok.
func (e Result[T, E]) Expect(msg string) T {
	return Match(e,
		func(t T) T {
			return t
		},
		func(_ E) T {
			panic(msg)
		},
	)
}

// Returns the contained Err value. Panics with the given message if the result is not Err.
func (e Result[T, E]) ExpectErr(msg string) E {
	return Match(e,
		func(_ T) E {
			panic(msg)
		},
		func(e E) E {
			return e
		},
	)
}

// Unwrap returns the contained Ok value. Panics if it is Err.
func (e Result[T, E]) Unwrap() T {
	return e.Expect("result was Err")
}

// Returns the contained Ok value. If the result is Err, returns the provided default.
// Default value is eagerly evaluated. Consider using UnwrapOrElse if providing the return value of a function call.
func (e Result[T, E]) UnwrapOr(defaultValue T) T {
	return Match(e,
		func(t T) T {
			return t
		},
		func(_ E) T {
			return defaultValue
		},
	)
}

// Returns the contained Ok value. If the result is Err, computes the default from the provided closure.
func (e Result[T, E]) UnwrapOrElse(defaultFunc func(E) T) T {
	return Match(e,
		func(t T) T {
			return t
		},
		func(e E) T {
			return defaultFunc(e)
		},
	)
}

// Returns the contained Err value. Panics if it is Ok.
func (e Result[T, E]) UnwrapErr() E {
	return e.ExpectErr("result was Ok")
}

// Returns the contained Err value. If the result is Ok, returns the provided default.
func (e Result[T, E]) UnwrapErrOr(defaultValue E) E {
	return Match(e,
		func(_ T) E {
			return defaultValue
		},
		func(e E) E {
			return e
		},
	)
}

// Returns the contained Err value. If the result is Ok, computes the default from the provided closure.
func (e Result[T, E]) UnwrapErrOrElse(defaultFunc func(T) E) E {
	return Match(e,
		func(t T) E {
			return defaultFunc(t)
		},
		func(e E) E {
			return e
		},
	)
}

// Returns a string representation of this Result.
func (e Result[T, E]) String() string {
	return Match(e,
		func(t T) string {
			return fmt.Sprintf("Ok(%v)", t)
		},
		func(e E) string {
			return fmt.Sprintf("Err(%v)", e)
		},
	)
}

//=====================================================

// Returns res2 if res1 is Ok, otherwise returns the Err value of res1.
func And[T any, E any, U any](res1 Result[T, E], res2 Result[U, E]) Result[U, E] {
	return ResultMonad[T, E, U]{}.Bind(
		res1,
		func(_ T) Result[U, E] {
			return res2
		},
	)
}

// Returns f(T) if res1 is Ok[T], otherwise returns the Err value of res1.
func AndThen[T any, E any, U any](res1 Result[T, E], f func(T) Result[U, E]) Result[U, E] {
	return ResultMonad[T, E, U]{}.Bind(res1, f)
}

// Flattens a result of type Result[Result[T, E], E] to just Result[T, E].
func Flatten[T any, E any](res Result[Result[T, E], E]) Result[T, E] {
	return Match(res,
		func(r Result[T, E]) Result[T, E] {
			return r
		},
		func(e E) Result[T, E] {
			return Err[T](e)
		},
	)
}

// Maps a Result[T, E] to Result[U, E] by applying a function to a contained Ok value.
// Leaves an Err value untouched.
func Map[T any, E any, U any](res Result[T, E], f func(T) U) Result[U, E] {
	return Match(res,
		func(ok T) Result[U, E] {
			return Ok[U, E](f(ok))
		},
		func(e E) Result[U, E] {
			return Err[U](e)
		},
	)
}

// Maps a Result[T, E] to Result[T, E2] by applying a function to a contained Err value.
// Leaves an Ok value untouched.
func MapErr[T any, E any, E2 any](res Result[T, E], f func(E) E2) Result[T, E2] {
	return Match(res,
		func(ok T) Result[T, E2] {
			return Ok[T, E2](ok)
		},
		func(e E) Result[T, E2] {
			return Err[T](f(e))
		},
	)
}

// Maps a Result[T, E] to Result[U, E] by applying a function to a contained Ok value.
// Returns the provided default if it is Err.
// Default value is eagerly evaluated. Consider using MapOrElse if you are passing the return value of a function call.
func MapOr[T any, E any, U any](res Result[T, E], defaultValue U, f func(T) U) Result[U, E] {
	return Match(res,
		func(ok T) Result[U, E] {
			return Ok[U, E](f(ok))
		},
		func(_ E) Result[U, E] {
			return Ok[U, E](defaultValue)
		},
	)
}

// Maps a Result[T, E] to Result[U, E] by applying a function to a contained Ok value.
// Returns the result produced by the default function if it is Err.
// Default is lazily evaluated.
func MapOrElse[T any, E any, U any](res Result[T, E], defaultFunc func(E) U, f func(T) U) Result[U, E] {
	return Match(res,
		func(ok T) Result[U, E] {
			return Ok[U, E](f(ok))
		},
		func(e E) Result[U, E] {
			return Ok[U, E](defaultFunc(e))
		},
	)
}

// Returns res2 if the result is Err, otherwise returns the Ok value of res1.
// res2 is eagerly evaluated. Consider using OrElse if you are passing the result of a function call.
func Or[T any, E any, E2 any](res1 Result[T, E], res2 Result[T, E2]) Result[T, E2] {
	return Match(res1,
		func(t T) Result[T, E2] {
			return Ok[T, E2](t)
		},
		func(e E) Result[T, E2] {
			return res2
		},
	)
}

// Returns f(T) if the result is Err[T], otherwise returns the Ok value of res1.
// f is lazily evaluated.
func OrElse[T any, E any, E2 any](res1 Result[T, E], f func(E) Result[T, E2]) Result[T, E2] {
	return Match(res1,
		func(t T) Result[T, E2] {
			return Ok[T, E2](t)
		},
		func(e E) Result[T, E2] {
			return f(e)
		},
	)
}

// Match calls okArm if the result is Ok[T] and returns that.
// It calls errArm if the result is Err[E] and returns that instead.
// The two functions must return the same type.
func Match[T any, E any, U any](r Result[T, E], okArm func(t T) U, errArm func(e E) U) U {
	switch r.variant {
	case resultVariantOk:
		return okArm(r.ok)
	case resultVariantErr:
		return errArm(r.err)
	default:
		panic("result type is neither Ok[T, E] nor Err[T, E]") // This should never happen.
	}
}

// Converts an Option[T] to a Result[T, E], mapping Some[T] to Ok[T], and Nothing to Err[err].
// Arguments are eagerly evaluated; consider using OkOrElse if passing the result of a function call.
func OkOr[T any, E any](opt option.Option[T], err E) Result[T, E] {
	return option.Match(opt,
		func(t T) Result[T, E] {
			return Ok[T, E](t)
		},
		func() Result[T, E] {
			return Err[T](err)
		},
	)
}

// Converts an Option[T] to a Result[T, E], mapping Some[T] to Ok[T], and Nothing to Err[f()].
// f is lazily evaluated.
func OkOrElse[T any, E any](opt option.Option[T], f func() E) Result[T, E] {
	return option.Match(opt,
		func(t T) Result[T, E] {
			return Ok[T, E](t)
		},
		func() Result[T, E] {
			return Err[T](f())
		},
	)
}

// Returns Ok[value] if err == nil.
// Returns Err[err] if err != nil.
func From[T any](value T, err error) Result[T, error] {
	if err != nil {
		return Err[T](err)
	}
	return Ok[T, error](value)
}
