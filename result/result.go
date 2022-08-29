package result

// Implements Monad[Result[T1, E], Result[T2, E], T1].
type ResultMonad[T1 any, E any, T2 any] struct {
}

func (m ResultMonad[T1, E, T2]) Bind(val Result[T1, E], f func(T1) Result[T2, E]) Result[T2, E] {
	return Match(val,
		func(o Ok[T1, E]) Result[T2, E] {
			return f(o.Value)
		},
		func(e Err[T1, E]) Result[T2, E] {
			return e
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
	// Returns true if the result is Err.
	IsErr() bool
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

//=====================================================

type Err[T any, E any] struct {
	Value E
}

func (e Err[T, E]) isResult() {
}

func (e Err[T, E]) IsOk() bool {
	return false
}

func (e Err[T, E]) IsErr() bool {
	return true
}

//=====================================================

// Match calls okArm if the result is Ok[T] and returns that result.
// It calls errArm if the result is Err[E] and returns that instead.
// The two functions must return the same type.
func Match[T any, E any, U any](opt Result[T, E], okArm func(Ok[T, E]) U, errArm func(Err[T, E]) U) U {
	switch inner := opt.(type) {
	case Ok[T, E]:
		return okArm(inner)
	case Err[T, E]:
		return errArm(inner)
	default:
		panic("result type is neither Ok[T, E] nor Err[T, E]") // This should never happen.
	}
}
