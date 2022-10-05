package either

import (
	"fmt"

	"github.com/sidkurella/goption/option"
)

// Implements Monad[Either[L, T1], Either[L, T2], T1].
type EitherMonad[L any, T1 any, T2 any] struct {
}

func (m EitherMonad[L, T1, T2]) Bind(val Either[L, T1], f func(T1) Either[L, T2]) Either[L, T2] {
	return Match(val,
		func(l Left[L, T1]) Either[L, T2] {
			return Left[L, T2]{Value: l.Value}
		},
		func(r Right[L, T1]) Either[L, T2] {
			return f(r.Value)
		},
	)
}

func (m EitherMonad[L, T1, T2]) Return(val T1) Either[L, T1] {
	return Right[L, T1]{Value: val}
}

// Either type. Represents one of two values (Left or Right).
// In the common case of success or failure, success is right (Right is Right) and failure is Left.
type Either[L any, R any] interface {
	// Sentinel method to prevent creation of other either types.
	isEither()

	// Returns true if the either is Right.
	IsRight() bool
	// Returns true if the either is Right and matches the given predicate.
	IsRightAnd(pred func(*R) bool) bool
	// Returns true if the either is Left.
	IsLeft() bool
	// Returns true if the either is Left and matches the given predicate.
	IsLeftAnd(pred func(*L) bool) bool

	// Converts from Either[L, R] to Option[R].
	// Converts self into an Option[R], discarding the left value, if any.
	Right() option.Option[R]
	// Converts from Either[L, R] to Option[L].
	// Converts self into an Option[L], and discarding the right value, if any.
	Left() option.Option[L]

	// Unwrap returns the contained Right value. Panics if it is Left.
	Unwrap() R
	// Returns the contained Right value. If the either is Left, returns the provided default.
	// Default value is eagerly evaluated. Consider using UnwrapOrElse if providing the either of a function call.
	UnwrapOr(defaultValue R) R
	// Returns the contained Right value. If the either is Left, computes the default from the provided closure.
	UnwrapOrElse(defaultFunc func(L) R) R
	// Returns the contained Left value. Panics if it is Right.
	UnwrapLeft() L
	// Returns the contained Left value. If the either is Right, returns the provided default.
	UnwrapLeftOr(defaultValue L) L
	// Returns the contained Left value. If the either is Right, computes the default from the provided closure.
	UnwrapLeftOrElse(defaultFunc func(R) L) L

	// Returns the contained Right value. Panics with the given message if the either is not Right.
	Expect(msg string) R
	// Returns the contained Left value. Panics with the given message if the either is not Left.
	ExpectLeft(msg string) L

	String() string
}

//=====================================================

type Right[L any, R any] struct {
	Value R
}

func (r Right[L, R]) isEither() {
}

func (r Right[L, R]) IsRight() bool {
	return true
}

func (r Right[L, R]) IsLeft() bool {
	return false
}

func (r Right[L, R]) Right() option.Option[R] {
	return option.Some[R]{Value: r.Value}
}

func (r Right[L, R]) IsRightAnd(pred func(*R) bool) bool {
	return pred(&r.Value)
}

func (r Right[L, R]) Left() option.Option[L] {
	return option.Nothing[L]{}
}

func (r Right[L, R]) IsLeftAnd(pred func(*L) bool) bool {
	return false
}

func (r Right[L, R]) Expect(_ string) R {
	return r.Value
}

func (r Right[L, R]) ExpectLeft(msg string) L {
	panic(msg)
}

func (r Right[L, R]) Unwrap() R {
	return r.Value
}

func (r Right[L, R]) UnwrapOr(_ R) R {
	return r.Value
}

func (r Right[L, R]) UnwrapOrElse(_ func(L) R) R {
	return r.Value
}

func (r Right[L, R]) UnwrapLeft() L {
	panic(r.Value)
}

func (r Right[L, R]) UnwrapLeftOr(defaultValue L) L {
	return defaultValue
}

func (r Right[L, R]) UnwrapLeftOrElse(defaultFunc func(R) L) L {
	return defaultFunc(r.Value)
}

func (r Right[L, R]) String() string {
	return fmt.Sprintf("Right(%v)", r.Value)
}

//=====================================================

type Left[L any, R any] struct {
	Value L
}

func (l Left[L, R]) isEither() {
}

func (l Left[L, R]) IsRight() bool {
	return false
}

func (l Left[L, R]) IsRightAnd(pred func(*R) bool) bool {
	return false
}

func (l Left[L, R]) IsLeft() bool {
	return true
}

func (l Left[L, R]) IsLeftAnd(pred func(*L) bool) bool {
	return pred(&l.Value)
}

func (l Left[L, R]) Right() option.Option[R] {
	return option.Nothing[R]{}
}

func (l Left[L, R]) Left() option.Option[L] {
	return option.Some[L]{Value: l.Value}
}

func (l Left[L, R]) Expect(msg string) R {
	panic(msg)
}

func (l Left[L, R]) ExpectLeft(_ string) L {
	return l.Value
}

func (l Left[L, R]) Unwrap() R {
	panic(l.Value)
}

func (l Left[L, R]) UnwrapOr(defaultValue R) R {
	return defaultValue
}

func (l Left[L, R]) UnwrapOrElse(f func(L) R) R {
	return f(l.Value)
}

func (l Left[L, R]) UnwrapLeft() L {
	return l.Value
}

func (l Left[L, R]) UnwrapLeftOr(_ L) L {
	return l.Value
}

func (l Left[L, R]) UnwrapLeftOrElse(_ func(R) L) L {
	return l.Value
}

func (l Left[L, R]) String() string {
	return fmt.Sprintf("Left(%v)", l.Value)
}

//=====================================================

// Returns e2 if e1 is Right, otherwise returns the Left value of e1.
func And[L any, R any, R2 any](e1 Either[L, R], e2 Either[L, R2]) Either[L, R2] {
	return EitherMonad[L, R, R2]{}.Bind(
		e1,
		func(_ R) Either[L, R2] {
			return e2
		},
	)
}

// Returns f(R) if e1 is Right[R], otherwise returns the Left value of e1.
func AndThen[L any, R any, R2 any](e1 Either[L, R], f func(R) Either[L, R2]) Either[L, R2] {
	return EitherMonad[L, R, R2]{}.Bind(e1, f)
}

// Flattens a either of type Either[L, Either[L, R]] to just Either[L, R].
func Flatten[L any, R any](e Either[L, Either[L, R]]) Either[L, R] {
	return Match(e,
		func(l Left[L, Either[L, R]]) Either[L, R] {
			return Left[L, R]{Value: l.Value}
		},
		func(r Right[L, Either[L, R]]) Either[L, R] {
			return r.Value
		},
	)
}

// Maps a Either[L, R] to Either[L, R2] by applying a function to a contained Right value.
// Leaves an Left value untouched.
func Map[L any, R any, R2 any](e Either[L, R], f func(R) R2) Either[L, R2] {
	return Match(e,
		func(l Left[L, R]) Either[L, R2] {
			return Left[L, R2]{Value: l.Value}
		},
		func(r Right[L, R]) Either[L, R2] {
			return Right[L, R2]{Value: f(r.Value)}
		},
	)
}

// Maps a Either[L, R] to Either[L2, R] by applying a function to a contained Left value.
// Leaves an Right value untouched.
func MapLeft[L any, L2 any, R any](e Either[L, R], f func(L) L2) Either[L2, R] {
	return Match(e,
		func(l Left[L, R]) Either[L2, R] {
			return Left[L2, R]{Value: f(l.Value)}
		},
		func(r Right[L, R]) Either[L2, R] {
			return Right[L2, R]{Value: r.Value}
		},
	)
}

// Maps a Either[L, R] to Either[L, R2] by applying a function to a contained Right value.
// Returns the provided default if it is Left.
// Default value is eagerly evaluated. Consider using MapOrElse if you are passing the either of a function call.
func MapOr[L any, R any, R2 any](e Either[L, R], defaultValue R2, f func(R) R2) Either[L, R2] {
	return Match(e,
		func(l Left[L, R]) Either[L, R2] {
			return Right[L, R2]{Value: defaultValue}
		},
		func(r Right[L, R]) Either[L, R2] {
			return Right[L, R2]{Value: f(r.Value)}
		},
	)
}

// Maps a Either[L, R] to Either[L, R2] by applying a function to a contained Right value.
// Returns the either produced by the default function if it is Left.
// Default is lazily evaluated.
func MapOrElse[L any, R any, R2 any](e Either[L, R], defaultFunc func(L) R2, f func(R) R2) Either[L, R2] {
	return Match(e,
		func(l Left[L, R]) Either[L, R2] {
			return Right[L, R2]{Value: defaultFunc(l.Value)}
		},
		func(r Right[L, R]) Either[L, R2] {
			return Right[L, R2]{Value: f(r.Value)}
		},
	)
}

// Returns e2 if the either is Left, otherwise returns the Right value of e1.
// e2 is eagerly evaluated. Consider using OrElse if you are passing the either of a function call.
func Or[L any, L2 any, R any](e1 Either[L, R], e2 Either[L2, R]) Either[L2, R] {
	return Match(e1,
		func(l Left[L, R]) Either[L2, R] {
			return e2
		},
		func(r Right[L, R]) Either[L2, R] {
			return Right[L2, R]{Value: r.Value}
		},
	)
}

// Returns f(L) if the either is Left[L], otherwise returns the Right value of e.
// f is lazily evaluated.
func OrElse[L any, L2 any, R any](e Either[L, R], f func(L) Either[L2, R]) Either[L2, R] {
	return Match(e,
		func(l Left[L, R]) Either[L2, R] {
			return f(l.Value)
		},
		func(r Right[L, R]) Either[L2, R] {
			return Right[L2, R]{Value: r.Value}
		},
	)
}

// Match calls leftArm if the either is Left[L] and returns that.
// It calls rightArm if the either is Right[R] and returns that instead.
// The two functions must return the same type.
func Match[L any, R any, U any](e Either[L, R], leftArm func(Left[L, R]) U, rightArm func(Right[L, R]) U) U {
	switch inner := e.(type) {
	case Left[L, R]:
		return leftArm(inner)
	case Right[L, R]:
		return rightArm(inner)
	default:
		panic("either type is neither Left[L, R] nor Right[L, R]") // This should never happen.
	}
}

// Converts an Option[R] to a Either[L, R], mapping Some[R] to Right[R], and Nothing to Left[err].
// Arguments are eagerly evaluated; consider using RightOrElse if passing the either of a function call.
func RightOr[L any, R any](opt option.Option[R], err L) Either[L, R] {
	return option.Match(opt,
		func(s option.Some[R]) Either[L, R] {
			return Right[L, R]{Value: s.Value}
		},
		func(_ option.Nothing[R]) Either[L, R] {
			return Left[L, R]{Value: err}
		},
	)
}

// Converts an Option[R] to a Either[L, R], mapping Some[R] to Right[R], and Nothing to Left[f()].
// f is lazily evaluated.
func RightOrElse[L any, R any](opt option.Option[R], f func() L) Either[L, R] {
	return option.Match(opt,
		func(s option.Some[R]) Either[L, R] {
			return Right[L, R]{Value: s.Value}
		},
		func(n option.Nothing[R]) Either[L, R] {
			return Left[L, R]{Value: f()}
		},
	)
}

// Returns Left[err] if err != nil. Returns Right[r] if err == nil.
func From[R any](r R, err error) Either[error, R] {
	if err != nil {
		return Left[error, R]{Value: err}
	}
	return Right[error, R]{Value: r}
}
