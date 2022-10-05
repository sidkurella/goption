package either

import (
	"fmt"

	"github.com/sidkurella/goption/option"
)

// Implements Monad[Either[T1, E], Either[T2, E], T1].
type EitherMonad[T1 any, E any, T2 any] struct {
}

func (m EitherMonad[T1, E, T2]) Bind(val Either[T1, E], f func(T1) Either[T2, E]) Either[T2, E] {
	return Match(val,
		func(o First[T1, E]) Either[T2, E] {
			return f(o.Value)
		},
		func(e Second[T1, E]) Either[T2, E] {
			return Second[T2, E]{Value: e.Value}
		},
	)
}

func (m EitherMonad[T1, E, T2]) Return(val T1) Either[T1, E] {
	return First[T1, E]{Value: val}
}

// Either type. Represents one of two possible values.
// In the common success/failure case, represents either success (First) or failure (Second).
type Either[T any, E any] interface {
	// Sentinel method to prevent creation of other either types.
	isEither()

	// Returns true if the either is First.
	IsFirst() bool
	// Returns true if the either is First and matches the given predicate.
	IsFirstAnd(pred func(*T) bool) bool
	// Returns true if the either is Second.
	IsSecond() bool
	// Returns true if the either is Second and matches the given predicate.
	IsSecondAnd(pred func(*E) bool) bool

	// Converts from Either[T, E] to Option[T].
	// Converts self into an Option[T], discarding the second value, if any.
	First() option.Option[T]
	// Converts from Either[T, E] to Option[E].
	// Converts self into an Option[E], and discarding the first value, if any.
	Second() option.Option[E]

	// Unwrap returns the contained First value. Panics if it is Second.
	Unwrap() T
	// Returns the contained First value. If the either is Second, returns the provided default.
	// Default value is eagerly evaluated. Consider using UnwrapOrElse if providing the either of a function call.
	UnwrapOr(defaultValue T) T
	// Returns the contained First value. If the either is Second, computes the default from the provided closure.
	UnwrapOrElse(defaultFunc func(E) T) T
	// Returns the contained Second value. Panics if it is First.
	UnwrapSecond() E
	// Returns the contained Second value. If the either is First, returns the provided default.
	UnwrapSecondOr(defaultValue E) E
	// Returns the contained Second value. If the either is First, computes the default from the provided closure.
	UnwrapSecondOrElse(defaultFunc func(T) E) E

	// Returns the contained First value. Panics with the given message if the either is not First.
	Expect(msg string) T
	// Returns the contained Second value. Panics with the given message if the either is not Second.
	ExpectSecond(msg string) E

	String() string
}

//=====================================================

type First[T any, E any] struct {
	Value T
}

func (o First[T, E]) isEither() {
}

func (o First[T, E]) IsFirst() bool {
	return true
}

func (o First[T, E]) IsSecond() bool {
	return false
}

func (o First[T, E]) First() option.Option[T] {
	return option.Some[T]{Value: o.Value}
}

func (o First[T, E]) IsFirstAnd(pred func(*T) bool) bool {
	return pred(&o.Value)
}

func (o First[T, E]) Second() option.Option[E] {
	return option.Nothing[E]{}
}

func (o First[T, E]) IsSecondAnd(pred func(*E) bool) bool {
	return false
}

func (o First[T, E]) Expect(_ string) T {
	return o.Value
}

func (o First[T, E]) ExpectSecond(msg string) E {
	panic(msg)
}

func (o First[T, E]) Unwrap() T {
	return o.Value
}

func (o First[T, E]) UnwrapOr(_ T) T {
	return o.Value
}

func (o First[T, E]) UnwrapOrElse(_ func(E) T) T {
	return o.Value
}

func (o First[T, E]) UnwrapSecond() E {
	panic(o.Value)
}

func (o First[T, E]) UnwrapSecondOr(defaultValue E) E {
	return defaultValue
}

func (o First[T, E]) UnwrapSecondOrElse(defaultFunc func(T) E) E {
	return defaultFunc(o.Value)
}

func (o First[T, E]) String() string {
	return fmt.Sprintf("First(%v)", o.Value)
}

//=====================================================

type Second[T any, E any] struct {
	Value E
}

func (e Second[T, E]) isEither() {
}

func (e Second[T, E]) IsFirst() bool {
	return false
}

func (e Second[T, E]) IsFirstAnd(pred func(*T) bool) bool {
	return false
}

func (e Second[T, E]) IsSecond() bool {
	return true
}

func (e Second[T, E]) IsSecondAnd(pred func(*E) bool) bool {
	return pred(&e.Value)
}

func (e Second[T, E]) First() option.Option[T] {
	return option.Nothing[T]{}
}

func (e Second[T, E]) Second() option.Option[E] {
	return option.Some[E]{Value: e.Value}
}

func (e Second[T, E]) Expect(msg string) T {
	panic(msg)
}

func (e Second[T, E]) ExpectSecond(_ string) E {
	return e.Value
}

func (e Second[T, E]) Unwrap() T {
	panic(e.Value)
}

func (e Second[T, E]) UnwrapOr(defaultValue T) T {
	return defaultValue
}

func (e Second[T, E]) UnwrapOrElse(f func(E) T) T {
	return f(e.Value)
}

func (e Second[T, E]) UnwrapSecond() E {
	return e.Value
}

func (e Second[T, E]) UnwrapSecondOr(_ E) E {
	return e.Value
}

func (e Second[T, E]) UnwrapSecondOrElse(_ func(T) E) E {
	return e.Value
}

func (e Second[T, E]) String() string {
	return fmt.Sprintf("Second(%v)", e.Value)
}

//=====================================================

// Returns res2 if res1 is First, otherwise returns the Second value of res1.
func And[T any, E any, U any](res1 Either[T, E], res2 Either[U, E]) Either[U, E] {
	return EitherMonad[T, E, U]{}.Bind(
		res1,
		func(_ T) Either[U, E] {
			return res2
		},
	)
}

// Returns f(T) if res1 is First[T], otherwise returns the Second value of res1.
func AndThen[T any, E any, U any](res1 Either[T, E], f func(T) Either[U, E]) Either[U, E] {
	return EitherMonad[T, E, U]{}.Bind(res1, f)
}

// Flattens a either of type Either[Either[T, E], E] to just Either[T, E].
func Flatten[T any, E any](res Either[Either[T, E], E]) Either[T, E] {
	return Match(res,
		func(o First[Either[T, E], E]) Either[T, E] {
			return o.Value
		},
		func(e Second[Either[T, E], E]) Either[T, E] {
			return Second[T, E]{Value: e.Value}
		},
	)
}

// Maps a Either[T, E] to Either[U, E] by applying a function to a contained First value.
// Leaves an Second value untouched.
func Map[T any, E any, U any](res Either[T, E], f func(T) U) Either[U, E] {
	return Match(res,
		func(o First[T, E]) Either[U, E] {
			return First[U, E]{Value: f(o.Value)}
		},
		func(e Second[T, E]) Either[U, E] {
			return Second[U, E]{Value: e.Value}
		},
	)
}

// Maps a Either[T, E] to Either[T, F] by applying a function to a contained Second value.
// Leaves an First value untouched.
func MapSecond[T any, E any, F any](res Either[T, E], f func(E) F) Either[T, F] {
	return Match(res,
		func(o First[T, E]) Either[T, F] {
			return First[T, F]{Value: o.Value}
		},
		func(e Second[T, E]) Either[T, F] {
			return Second[T, F]{Value: f(e.Value)}
		},
	)
}

// Maps a Either[T, E] to Either[U, E] by applying a function to a contained First value.
// Returns the provided default if it is Second.
// Default value is eagerly evaluated. Consider using MapOrElse if you are passing the either of a function call.
func MapOr[T any, E any, U any](res Either[T, E], defaultValue U, f func(T) U) Either[U, E] {
	return Match(res,
		func(o First[T, E]) Either[U, E] {
			return First[U, E]{Value: f(o.Value)}
		},
		func(e Second[T, E]) Either[U, E] {
			return First[U, E]{Value: defaultValue}
		},
	)
}

// Maps a Either[T, E] to Either[U, E] by applying a function to a contained First value.
// Returns the either produced by the default function if it is Second.
// Default is lazily evaluated.
func MapOrElse[T any, E any, U any](res Either[T, E], defaultFunc func(E) U, f func(T) U) Either[U, E] {
	return Match(res,
		func(o First[T, E]) Either[U, E] {
			return First[U, E]{Value: f(o.Value)}
		},
		func(e Second[T, E]) Either[U, E] {
			return First[U, E]{Value: defaultFunc(e.Value)}
		},
	)
}

// Returns res2 if the either is Second, otherwise returns the First value of res1.
// res2 is eagerly evaluated. Consider using OrElse if you are passing the either of a function call.
func Or[T any, E any, F any](res1 Either[T, E], res2 Either[T, F]) Either[T, F] {
	return Match(res1,
		func(o First[T, E]) Either[T, F] {
			return First[T, F]{Value: o.Value}
		},
		func(_ Second[T, E]) Either[T, F] {
			return res2
		},
	)
}

// Returns f(E) if the either is Second[E], otherwise returns the First value of res1.
// f is lazily evaluated.
func OrElse[T any, E any, F any](res1 Either[T, E], f func(E) Either[T, F]) Either[T, F] {
	return Match(res1,
		func(o First[T, E]) Either[T, F] {
			return First[T, F]{Value: o.Value}
		},
		func(e Second[T, E]) Either[T, F] {
			return f(e.Value)
		},
	)
}

// Match calls firstArm if the either is First[T] and returns that either.
// It calls secondArm if the either is Second[E] and returns that instead.
// The two functions must return the same type.
func Match[T any, E any, U any](res Either[T, E], firstArm func(First[T, E]) U, secondArm func(Second[T, E]) U) U {
	switch inner := res.(type) {
	case First[T, E]:
		return firstArm(inner)
	case Second[T, E]:
		return secondArm(inner)
	default:
		panic("either type is neither First[T, E] nor Second[T, E]") // This should never happen.
	}
}

// Converts an Option[T] to a Either[T, E], mapping Some[T] to First[T], and Nothing to Second[second].
// Arguments are eagerly evaluated; consider using FirstOrElse if passing the either of a function call.
func FirstOr[T any, E any](opt option.Option[T], second E) Either[T, E] {
	return option.Match(opt,
		func(s option.Some[T]) Either[T, E] {
			return First[T, E]{Value: s.Value}
		},
		func(n option.Nothing[T]) Either[T, E] {
			return Second[T, E]{Value: second}
		},
	)
}

// Converts an Option[T] to a Either[T, E], mapping Some[T] to First[T], and Nothing to Second[f()].
// f is lazily evaluated.
func FirstOrElse[T any, E any](opt option.Option[T], f func() E) Either[T, E] {
	return option.Match(opt,
		func(s option.Some[T]) Either[T, E] {
			return First[T, E]{Value: s.Value}
		},
		func(n option.Nothing[T]) Either[T, E] {
			return Second[T, E]{Value: f()}
		},
	)
}

// Returns First[r] if err == nil.
// Returns Second[err] if err != nil.
func From[T any](value T, err error) Either[T, error] {
	if err != nil {
		return Second[T, error]{Value: err}
	}
	return First[T, error]{Value: value}
}
