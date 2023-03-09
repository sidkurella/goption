package option

import (
	"fmt"
)

// Implements Monad[Option[T], Option[U], T].
type OptionMonad[T any, U any] struct {
}

func (m OptionMonad[T, U]) Bind(val Option[T], f func(T) Option[U]) Option[U] {
	return Match(val,
		func(t T) Option[U] {
			return f(t)
		},
		func() Option[U] {
			return Nothing[U]()
		},
	)
}

func (m OptionMonad[T, U]) Return(val T) Option[T] {
	return Some(val)
}

//=====================================================

type optionVariant int

const (
	optionVariantNothing optionVariant = iota
	optionVariantSome
)

// Option type. Either contains a value, or does not.
// The default value is Nothing.
type Option[T any] struct {
	variant optionVariant
	value   T
}

//=====================================================

// Creates a Some variant of the option, which holds the value.
func Some[T any](t T) Option[T] {
	return Option[T]{
		variant: optionVariantSome,
		value:   t,
	}
}

// Creates a Nothing variant of the option, which holds no value.
func Nothing[T any]() Option[T] {
	return Option[T]{
		variant: optionVariantNothing,
	}
}

// Returns if the option contains a value (is Some).
func (o Option[T]) IsSome() bool {
	return o.variant == optionVariantSome
}

// Returns if the option contains no value (is Nothing).
func (o Option[T]) IsNothing() bool {
	return o.variant == optionVariantNothing
}

// Returns the value contained by the option, panicking with the given message if it is Nothing.
func (o Option[T]) Expect(msg string) T {
	return Match(o,
		func(t T) T {
			return t
		},
		func() T {
			panic(msg)
		},
	)
}

// Returns the value contained by the option, panicking if it is Nothing.
func (o Option[T]) Unwrap() T {
	return o.Expect("option contains nothing")
}

// Returns the value contained by the option. Returns the provided default if it is Nothing.
// Default value is eagerly evaluated.
func (o Option[T]) UnwrapOr(val T) T {
	return Match(o,
		func(t T) T {
			return t
		},
		func() T {
			return val
		},
	)
}

// Gets the value contained by the option. The second value indicates if the value is valid or not.
func (o Option[T]) Get() (T, bool) {
	return o.value, o.IsSome()
}

// Returns if the option contains a value and the value matches the given predicate.
func (o Option[T]) IsSomeAnd(pred func(*T) bool) bool {
	return o.IsSome() && pred(&o.value)
}

// Returns the option if it contains a value matching the predicate. Otherwise, returns Nothing.
func (o Option[T]) Filter(pred func(*T) bool) Option[T] {
	if o.IsSomeAnd(pred) {
		return o
	}
	return Nothing[T]()
}

// Returns o if the option is Some. Otherwise, returns opt2.
// opt2 is eagerly evaluated. If you need lazy evaluation, use OrElse.
func (o Option[T]) Or(opt2 Option[T]) Option[T] {
	return Match(o,
		func(_ T) Option[T] {
			return o
		},
		func() Option[T] {
			return opt2
		},
	)
}

// Returns self if the option contains a value. Otherwise, calls f and returns the provided option.
// f is not called unless the option is Nothing (lazily-evaluated).
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	return Match(o,
		func(_ T) Option[T] {
			return o
		},
		func() Option[T] {
			return f()
		},
	)
}

// If exactly one of the options contains a value, returns that option. Otherwise, returns Nothing.
func (o Option[T]) Xor(opt2 Option[T]) Option[T] {
	if o.IsSome() {
		if opt2.IsNothing() {
			return o
		}
		return Nothing[T]()
	}
	// o is Nothing.
	if opt2.IsSome() {
		return opt2
	}
	return Nothing[T]()
}

// Returns a string representation of this option.
func (o Option[T]) String() string {
	return Match(o,
		func(t T) string {
			return fmt.Sprintf("Some(%v)", t)
		},
		func() string {
			return "Nothing"
		},
	)
}

//=====================================================

// Returns opt2 if both options contain a value. Otherwise returns Nothing.
func And[T any, U any](opt1 Option[T], opt2 Option[U]) Option[U] {
	return Match(opt1,
		func(_ T) Option[U] {
			return opt2
		},
		func() Option[U] {
			return Nothing[U]()
		},
	)
}

// Calls f with the inner value of the option. Returns Nothing if there is no value.
func AndThen[T any, U any](opt1 Option[T], f func(T) Option[U]) Option[U] {
	return OptionMonad[T, U]{}.Bind(opt1, f)
}

// Flattens an option of type Option[Option[T]] to just Option[T].
func Flatten[T any](opt Option[Option[T]]) Option[T] {
	return Match(opt,
		func(t Option[T]) Option[T] {
			return t
		},
		func() Option[T] {
			return Nothing[T]()
		},
	)
}

// Maps the inner value of an option via f. Returns Nothing if there is no value.
func Map[T any, U any](opt Option[T], f func(T) U) Option[U] {
	return Match(opt,
		func(t T) Option[U] {
			return Some(f(t))
		},
		func() Option[U] {
			return Nothing[U]()
		},
	)
}

// Maps the inner value of an option via f. If the option is Nothing, returns the default.
// Arguments are eagerly evaluated. Consider MapOrElse if you are passing the result of a function call.
func MapOr[T any, U any](opt Option[T], defaultValue U, f func(T) U) Option[U] {
	return Match(opt,
		func(t T) Option[U] {
			return Some(f(t))
		},
		func() Option[U] {
			return Some(defaultValue)
		},
	)
}

// Maps the inner value of an option via f. If the option is Nothing, calls defaultFunc to provide a default.
// defaultFunc() is lazily evaluated.
func MapOrElse[T any, U any](opt Option[T], defaultFunc func() U, f func(T) U) Option[U] {
	return Match(opt,
		func(t T) Option[U] {
			return Some(f(t))
		},
		func() Option[U] {
			return Some(defaultFunc())
		},
	)
}

// Match calls someArm if the option is Some[T] and returns that result.
// It calls nothingArm if the option is Nothing and returns that instead.
// The two functions must return the same type.
func Match[T any, U any](opt Option[T], someArm func(T) U, nothingArm func() U) U {
	switch opt.variant {
	case optionVariantSome:
		return someArm(opt.value)
	case optionVariantNothing:
		return nothingArm()
	default:
		panic("option type is neither Some[T] nor Nothing[T]") // This should never happen.
	}
}

// Returns an option from the provided value and boolean indicating if the value is valid.
// Returns Nothing if ok is false.
func From[T any](val T, ok bool) Option[T] {
	if ok {
		return Some(val)
	}
	return Nothing[T]()
}

// Returns an option from the provided value and error.
// Returns Some[T] if the error is nil.
// Returns Nothing if the error is not nil.
// NOTE: The error value is discarded. Use Either if this is not desired.
func FromError[T any](val T, err error) Option[T] {
	return From(val, err == nil)
}
