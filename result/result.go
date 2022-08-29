package result

import (
	"fmt"

	"github.com/sidkurella/goption/option"
)

// Implements Monad[Result[T1, E], Result[T2, E], T1].
type ResultMonad[T1 any, E any, T2 any] struct {
}

func (m ResultMonad[T1, E, T2]) Bind(val Result[T1, E], f func(T1) Result[T2, E]) Result[T2, E] {
	return Match(val,
		func(o Ok[T1, E]) Result[T2, E] {
			return f(o.Value)
		},
		func(e Err[T1, E]) Result[T2, E] {
			return Err[T2, E]{Value: e.Value}
		},
	)
}

func (m ResultMonad[T1, E, T2]) Return(val T1) Result[T1, E] {
	return Ok[T1, E]{Value: val}
}

// Result type. Represents either success (Ok) or failure (Err).
type Result[T any, E any] interface {
	// Sentinel method to prevent creation of other result types.
	isResult()

	// Returns true if the result is Ok.
	IsOk() bool
	// Returns true if the result is Ok and matches the given predicate.
	IsOkAnd(pred func(*T) bool) bool
	// Returns true if the result is Err.
	IsErr() bool
	// Returns true if the result is Err and matches the given predicate.
	IsErrAnd(pred func(*E) bool) bool

	// Converts from Result[T, E] to Option[T].
	// Converts self into an Option[T], discarding the error, if any.
	Ok() option.Option[T]
	// Converts from Result[T, E] to Option[E].
	// Converts self into an Option[E], and discarding the error, if any.
	Err() option.Option[E]

	// Returns the contained Ok value. Panics if the result is not Ok.
	Expect(msg string) T
	// Returns the contained Err value. Panics if the result is not Err.
	ExpectErr(msg string) E

	// Unwrap returns the contained Ok value. Panics if it is Err.
	Unwrap() T
	// Returns the contained Err value. Panics if it is Ok.
	UnwrapErr() E
	// Returns the contained Ok value. If the result is Err, returns the provided default.
	// Default value is eagerly evaluated. Consider using UnwrapOrElse if providing the result of a function call.
	UnwrapOr(defaultValue T) T
	// Returns the contained Ok value. If the result is Err, computes the default from the provided closure.
	UnwrapOrElse(defaultFunc func(E) T) T

	String() string
}

//=====================================================

type Ok[T any, E any] struct {
	Value T
}

func (o Ok[T, E]) isResult() {
}

func (o Ok[T, E]) IsOk() bool {
	return true
}

func (o Ok[T, E]) IsErr() bool {
	return false
}

func (o Ok[T, E]) Ok() option.Option[T] {
	return option.Some[T]{Value: o.Value}
}

func (o Ok[T, E]) IsOkAnd(pred func(*T) bool) bool {
	return pred(&o.Value)
}

func (o Ok[T, E]) Err() option.Option[E] {
	return option.Nothing[E]{}
}

func (o Ok[T, E]) IsErrAnd(pred func(*E) bool) bool {
	return false
}

func (o Ok[T, E]) Expect(_ string) T {
	return o.Value
}

func (o Ok[T, E]) ExpectErr(msg string) E {
	panic(msg)
}

func (o Ok[T, E]) Unwrap() T {
	return o.Value
}

func (o Ok[T, E]) UnwrapErr() E {
	panic(o.Value)
}

func (o Ok[T, E]) UnwrapOr(_ T) T {
	return o.Value
}

func (o Ok[T, E]) UnwrapOrElse(_ func(E) T) T {
	return o.Value
}

func (o Ok[T, E]) String() string {
	return fmt.Sprintf("Ok(%v)", o.Value)
}

//=====================================================

type Err[T any, E any] struct {
	Value E
}

func (e Err[T, E]) isResult() {
}

func (e Err[T, E]) IsOk() bool {
	return false
}

func (e Err[T, E]) IsOkAnd(pred func(*T) bool) bool {
	return false
}

func (e Err[T, E]) IsErr() bool {
	return true
}

func (e Err[T, E]) IsErrAnd(pred func(*E) bool) bool {
	return pred(&e.Value)
}

func (e Err[T, E]) Ok() option.Option[T] {
	return option.Nothing[T]{}
}

func (e Err[T, E]) Err() option.Option[E] {
	return option.Some[E]{Value: e.Value}
}

func (e Err[T, E]) Expect(msg string) T {
	panic(msg)
}

func (e Err[T, E]) ExpectErr(_ string) E {
	return e.Value
}

func (e Err[T, E]) Unwrap() T {
	panic(e.Value)
}

func (e Err[T, E]) UnwrapErr() E {
	return e.Value
}

func (e Err[T, E]) UnwrapOr(defaultValue T) T {
	return defaultValue
}

func (e Err[T, E]) UnwrapOrElse(f func(E) T) T {
	return f(e.Value)
}

func (e Err[T, E]) String() string {
	return fmt.Sprintf("Err(%v)", e.Value)
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
		func(o Ok[Result[T, E], E]) Result[T, E] {
			return o.Value
		},
		func(e Err[Result[T, E], E]) Result[T, E] {
			return Err[T, E]{Value: e.Value}
		},
	)
}

// Maps a Result[T, E] to Result[U, E] by applying a function to a contained Ok value.
// Leaves an Err value untouched.
func Map[T any, E any, U any](res Result[T, E], f func(T) U) Result[U, E] {
	return Match(res,
		func(o Ok[T, E]) Result[U, E] {
			return Ok[U, E]{Value: f(o.Value)}
		},
		func(e Err[T, E]) Result[U, E] {
			return Err[U, E]{Value: e.Value}
		},
	)
}

// Maps a Result[T, E] to Result[T, F] by applying a function to a contained Err value.
// Leaves an Ok value untouched.
func MapErr[T any, E any, F any](res Result[T, E], f func(E) F) Result[T, F] {
	return Match(res,
		func(o Ok[T, E]) Result[T, F] {
			return Ok[T, F]{Value: o.Value}
		},
		func(e Err[T, E]) Result[T, F] {
			return Err[T, F]{Value: f(e.Value)}
		},
	)
}

// Maps a Result[T, E] to Result[U, E] by applying a function to a contained Ok value.
// Returns the provided default if it is Err.
// Default value is eagerly evaluated. Consider using MapOrElse if you are passing the result of a function call.
func MapOr[T any, E any, U any](res Result[T, E], defaultValue U, f func(T) U) Result[U, E] {
	return Match(res,
		func(o Ok[T, E]) Result[U, E] {
			return Ok[U, E]{Value: f(o.Value)}
		},
		func(e Err[T, E]) Result[U, E] {
			return Ok[U, E]{Value: defaultValue}
		},
	)
}

// Maps a Result[T, E] to Result[U, E] by applying a function to a contained Ok value.
// Returns the result produced by the default function if it is Err.
// Default is lazily evaluated.
func MapOrElse[T any, E any, U any](res Result[T, E], defaultFunc func() U, f func(T) U) Result[U, E] {
	return Match(res,
		func(o Ok[T, E]) Result[U, E] {
			return Ok[U, E]{Value: f(o.Value)}
		},
		func(e Err[T, E]) Result[U, E] {
			return Ok[U, E]{Value: defaultFunc()}
		},
	)
}

// Returns res2 if the result is Err, otherwise returns the Ok value of res1.
// res2 is eagerly evaluated. Consider using OrElse if you are passing the result of a function call.
func Or[T any, E any, F any](res1 Result[T, E], res2 Result[T, F]) Result[T, F] {
	return Match(res1,
		func(o Ok[T, E]) Result[T, F] {
			return Ok[T, F]{Value: o.Value}
		},
		func(_ Err[T, E]) Result[T, F] {
			return res2
		},
	)
}

// Returns f(E) if the result is Err[E], otherwise returns the Ok value of res1.
// f is lazily evaluated.
func OrElse[T any, E any, F any](res1 Result[T, E], f func(E) Result[T, F]) Result[T, F] {
	return Match(res1,
		func(o Ok[T, E]) Result[T, F] {
			return Ok[T, F]{Value: o.Value}
		},
		func(e Err[T, E]) Result[T, F] {
			return f(e.Value)
		},
	)
}

// Match calls okArm if the result is Ok[T] and returns that result.
// It calls errArm if the result is Err[E] and returns that instead.
// The two functions must return the same type.
func Match[T any, E any, U any](res Result[T, E], okArm func(Ok[T, E]) U, errArm func(Err[T, E]) U) U {
	switch inner := res.(type) {
	case Ok[T, E]:
		return okArm(inner)
	case Err[T, E]:
		return errArm(inner)
	default:
		panic("result type is neither Ok[T, E] nor Err[T, E]") // This should never happen.
	}
}
