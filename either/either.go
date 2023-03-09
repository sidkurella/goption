package either

import (
	"fmt"

	"github.com/sidkurella/goption/option"
)

// Implements Monad[Either[F1, E], Either[F2, E], F1].
type EitherMonad[F1 any, S any, F2 any] struct {
}

func (m EitherMonad[F1, S, F2]) Bind(val Either[F1, S], f func(F1) Either[F2, S]) Either[F2, S] {
	return Match(val,
		func(first F1) Either[F2, S] {
			return f(first)
		},
		func(s S) Either[F2, S] {
			return Second[F2](s)
		},
	)
}

func (m EitherMonad[T1, E, T2]) Return(val T1) Either[T1, E] {
	return First[T1, E](val)
}

type eitherVariant int

const (
	eitherVariantFirst eitherVariant = iota
	eitherVariantSecond
)

// Either type. Represents one of two possible values.
// The default value is First(*new(F)) (i.e. First variant containing the zero value of F).
type Either[F any, S any] struct {
	variant eitherVariant
	first   F
	second  S
}

//=====================================================

func First[F any, S any](f F) Either[F, S] {
	return Either[F, S]{
		variant: eitherVariantFirst,
		first:   f,
	}
}

func Second[F any, S any](s S) Either[F, S] {
	return Either[F, S]{
		variant: eitherVariantSecond,
		second:  s,
	}
}

// Returns true if the either is First.
func (e Either[F, S]) IsFirst() bool {
	return e.variant == eitherVariantFirst
}

// Returns true if the either is First and matches the given predicate.
func (e Either[F, S]) IsFirstAnd(pred func(*F) bool) bool {
	return e.IsFirst() && pred(&e.first)
}

// Returns true if the either is Second.
func (e Either[F, S]) IsSecond() bool {
	return e.variant == eitherVariantSecond
}

// Returns true if the either is Second and matches the given predicate.
func (e Either[F, S]) IsSecondAnd(pred func(*S) bool) bool {
	return e.IsSecond() && pred(&e.second)
}

// Converts from Either[F, S] to Option[F].
// Converts self into an Option[F], discarding the second value, if any.
func (e Either[F, S]) First() option.Option[F] {
	if e.IsFirst() {
		return option.Some(e.first)
	}
	return option.Nothing[F]()
}

// Converts from Either[F, S] to Option[S].
// Converts self into an Option[S], and discarding the first value, if any.
func (e Either[F, S]) Second() option.Option[S] {
	if e.IsSecond() {
		return option.Some(e.second)
	}
	return option.Nothing[S]()
}

// Returns the contained First value. Panics with the given message if the either is not First.
func (e Either[F, S]) Expect(msg string) F {
	return Match(e,
		func(f F) F {
			return f
		},
		func(s S) F {
			panic(msg)
		},
	)
}

// Returns the contained Second value. Panics with the given message if the either is not Second.
func (e Either[F, S]) ExpectSecond(msg string) S {
	return Match(e,
		func(f F) S {
			panic(msg)
		},
		func(s S) S {
			return s
		},
	)
}

// Unwrap returns the contained First value. Panics if it is Second.
func (e Either[F, S]) Unwrap() F {
	return e.Expect("either was Second")
}

// Returns the contained First value. If the either is Second, returns the provided default.
// Default value is eagerly evaluated. Consider using UnwrapOrElse if providing the return value of a function call.
func (e Either[F, S]) UnwrapOr(defaultValue F) F {
	return Match(e,
		func(f F) F {
			return f
		},
		func(s S) F {
			return defaultValue
		},
	)
}

// Returns the contained First value. If the either is Second, computes the default from the provided closure.
func (e Either[F, S]) UnwrapOrElse(defaultFunc func(S) F) F {
	return Match(e,
		func(f F) F {
			return f
		},
		func(s S) F {
			return defaultFunc(s)
		},
	)
}

// Returns the contained Second value. Panics if it is First.
func (e Either[F, S]) UnwrapSecond() S {
	return e.ExpectSecond("either was First")
}

// Returns the contained Second value. If the either is First, returns the provided default.
func (e Either[F, S]) UnwrapSecondOr(defaultValue S) S {
	return Match(e,
		func(f F) S {
			return defaultValue
		},
		func(s S) S {
			return s
		},
	)
}

// Returns the contained Second value. If the either is First, computes the default from the provided closure.
func (e Either[F, S]) UnwrapSecondOrElse(defaultFunc func(F) S) S {
	return Match(e,
		func(f F) S {
			return defaultFunc(f)
		},
		func(s S) S {
			return s
		},
	)
}

// Returns a string representation of this Either.
func (e Either[F, S]) String() string {
	return Match(e,
		func(f F) string {
			return fmt.Sprintf("First(%v)", f)
		},
		func(s S) string {
			return fmt.Sprintf("Second(%v)", s)
		},
	)
}

//=====================================================

// Returns res2 if res1 is First, otherwise returns the Second value of res1.
func And[F any, S any, U any](res1 Either[F, S], res2 Either[U, S]) Either[U, S] {
	return EitherMonad[F, S, U]{}.Bind(
		res1,
		func(_ F) Either[U, S] {
			return res2
		},
	)
}

// Returns f(F) if res1 is First[F], otherwise returns the Second value of res1.
func AndThen[F any, S any, U any](res1 Either[F, S], f func(F) Either[U, S]) Either[U, S] {
	return EitherMonad[F, S, U]{}.Bind(res1, f)
}

// Flattens a either of type Either[Either[F, S], S] to just Either[F, S].
func Flatten[F any, S any](res Either[Either[F, S], S]) Either[F, S] {
	return Match(res,
		func(f Either[F, S]) Either[F, S] {
			return f
		},
		func(s S) Either[F, S] {
			return Second[F](s)
		},
	)
}

// Maps a Either[F, S] to Either[U, S] by applying a function to a contained First value.
// Leaves an Second value untouched.
func Map[F any, S any, U any](res Either[F, S], f func(F) U) Either[U, S] {
	return Match(res,
		func(first F) Either[U, S] {
			return First[U, S](f(first))
		},
		func(s S) Either[U, S] {
			return Second[U](s)
		},
	)
}

// Maps a Either[F, S] to Either[F, S2] by applying a function to a contained Second value.
// Leaves an First value untouched.
func MapSecond[F any, S any, S2 any](res Either[F, S], f func(S) S2) Either[F, S2] {
	return Match(res,
		func(first F) Either[F, S2] {
			return First[F, S2](first)
		},
		func(s S) Either[F, S2] {
			return Second[F](f(s))
		},
	)
}

// Maps a Either[F, S] to Either[U, S] by applying a function to a contained First value.
// Returns the provided default if it is Second.
// Default value is eagerly evaluated. Consider using MapOrElse if you are passing the return value of a function call.
func MapOr[F any, S any, U any](res Either[F, S], defaultValue U, f func(F) U) Either[U, S] {
	return Match(res,
		func(first F) Either[U, S] {
			return First[U, S](f(first))
		},
		func(s S) Either[U, S] {
			return First[U, S](defaultValue)
		},
	)
}

// Maps a Either[F, S] to Either[U, S] by applying a function to a contained First value.
// Returns the either produced by the default function if it is Second.
// Default is lazily evaluated.
func MapOrElse[F any, S any, U any](res Either[F, S], defaultFunc func(S) U, f func(F) U) Either[U, S] {
	return Match(res,
		func(first F) Either[U, S] {
			return First[U, S](f(first))
		},
		func(s S) Either[U, S] {
			return First[U, S](defaultFunc(s))
		},
	)
}

// Returns res2 if the either is Second, otherwise returns the First value of res1.
// res2 is eagerly evaluated. Consider using OrElse if you are passing the either of a function call.
func Or[F any, S any, S2 any](res1 Either[F, S], res2 Either[F, S2]) Either[F, S2] {
	return Match(res1,
		func(f F) Either[F, S2] {
			return First[F, S2](f)
		},
		func(s S) Either[F, S2] {
			return res2
		},
	)
}

// Returns f(S) if the either is Second[S], otherwise returns the First value of res1.
// f is lazily evaluated.
func OrElse[F any, S any, S2 any](res1 Either[F, S], f func(S) Either[F, S2]) Either[F, S2] {
	return Match(res1,
		func(f F) Either[F, S2] {
			return First[F, S2](f)
		},
		func(s S) Either[F, S2] {
			return f(s)
		},
	)
}

// Match calls firstArm if the either is First[F] and returns that.
// It calls secondArm if the either is Second[S] and returns that instead.
// The two functions must return the same type.
func Match[F any, S any, T any](e Either[F, S], firstArm func(f F) T, secondArm func(s S) T) T {
	switch e.variant {
	case eitherVariantFirst:
		return firstArm(e.first)
	case eitherVariantSecond:
		return secondArm(e.second)
	default:
		panic("either type is neither First[F, S] nor Second[F, S]") // This should never happen.
	}
}

// Converts an Option[T] to a Either[T, E], mapping Some[T] to First[T], and Nothing to Second[second].
// Arguments are eagerly evaluated; consider using FirstOrElse if passing the either of a function call.
func FirstOr[T any, E any](opt option.Option[T], second E) Either[T, E] {
	return option.Match(opt,
		func(t T) Either[T, E] {
			return First[T, E](t)
		},
		func() Either[T, E] {
			return Second[T](second)
		},
	)
}

// Converts an Option[T] to a Either[T, E], mapping Some[T] to First[T], and Nothing to Second[f()].
// f is lazily evaluated.
func FirstOrElse[T any, E any](opt option.Option[T], f func() E) Either[T, E] {
	return option.Match(opt,
		func(t T) Either[T, E] {
			return First[T, E](t)
		},
		func() Either[T, E] {
			return Second[T](f())
		},
	)
}

// Returns First[r] if err == nil.
// Returns Second[err] if err != nil.
func From[T any](value T, err error) Either[T, error] {
	if err != nil {
		return Second[T](err)
	}
	return First[T, error](value)
}
