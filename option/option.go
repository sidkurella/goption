package option

// Implements Monad[Option[T], Option[U], T, U].
type OptionMonad[T any, U any] struct {
}

func (m OptionMonad[T, U]) Bind(val Option[T], f func(T) Option[U]) Option[U] {
	if val.IsNothing() {
		return Nothing[U]{}
	}

	inner := val.(Some[T]).Value
	return f(inner)
}

func (m OptionMonad[T, U]) Return(val T) Option[T] {
	return Some[T]{Value: val}
}

//=====================================================

// Option type. Either contains a value, or does not.
type Option[T any] interface {
	// Returns true if the option has a value.
	IsSome() bool
	// Returns true if the option has no value.
	IsNothing() bool
	// Returns the inner value of the option. Panics if there is no value.
	Unwrap() T
	// Returns the inner value of the option. If there is no value, returns the provided default.
	UnwrapOr(defaultValue T) T
	// Gets the inner value of the option. The second value indicates success or failure.
	Get() (T, bool)
	// Returns the inner value of the option. Panics with the provided message if there is no value.
	Expect(msg string) T
	// Returns true if there if the option contains a value and the value passes the provided predicate.
	IsSomeAnd(predicate func(*T) bool) bool
	// Returns Some[T] if there if the option contains a value and the value passes the provided predicate.
	// Returns Nothing otherwise.
	Filter(predicate func(*T) bool) Option[T]
	// Returns opt2 if both options contain a value.
	And(opt2 Option[T]) Option[T]
	// Calls f with the inner value of the option. Returns Nothing if there is no value.
	AndThen(f func(T) Option[T]) Option[T]
	// Returns self if the option contains a value. Otherwise, returns opt2.
	Or(opt2 Option[T]) Option[T]
	// Returns self if the option contains a value. Otherwise, returns f.
	OrElse(f func() Option[T]) Option[T]
	// If exactly one of the options contains a value, returns that option. Otherwise, returns Nothing.
	Xor(opt2 Option[T]) Option[T]
}

//=====================================================

type Some[T any] struct {
	Value T
}

func (s Some[T]) IsSome() bool {
	return true
}

func (s Some[T]) IsNothing() bool {
	return false
}

func (s Some[T]) Unwrap() T {
	return s.Value
}

func (s Some[T]) UnwrapOr(val T) T {
	return s.Value
}

func (s Some[T]) Get() (T, bool) {
	return s.Value, true
}

func (s Some[T]) Expect(msg string) T {
	return s.Value
}

func (s Some[T]) IsSomeAnd(pred func(*T) bool) bool {
	return pred(&s.Value)
}

func (s Some[T]) Filter(pred func(*T) bool) Option[T] {
	if pred(&s.Value) {
		return s
	}
	return Nothing[T]{}
}

func (s Some[T]) And(opt2 Option[T]) Option[T] {
	return And[T](s, opt2)
}

func (s Some[T]) AndThen(f func(T) Option[T]) Option[T] {
	return AndThen[T](s, f)
}

func (s Some[T]) Or(opt2 Option[T]) Option[T] {
	return s
}

func (s Some[T]) OrElse(f func() Option[T]) Option[T] {
	return s
}

func (s Some[T]) Xor(opt2 Option[T]) Option[T] {
	if !opt2.IsSome() {
		return s
	}
	return Nothing[T]{}
}

//=====================================================

type Nothing[T any] struct {
}

func (n Nothing[T]) IsSome() bool {
	return false
}

func (n Nothing[T]) IsNothing() bool {
	return true
}

func (n Nothing[T]) UnwrapOr(val T) T {
	return val
}

func (n Nothing[T]) Get() (T, bool) {
	return *new(T), false
}

func (n Nothing[T]) Expect(msg string) T {
	panic(msg)
}

func (n Nothing[T]) Unwrap() T {
	return n.Expect("option type contains nothing")
}

func (n Nothing[T]) IsSomeAnd(pred func(*T) bool) bool {
	return false
}

func (n Nothing[T]) Filter(pred func(*T) bool) Option[T] {
	return n
}

func (n Nothing[T]) And(opt2 Option[T]) Option[T] {
	return n
}

func (n Nothing[T]) AndThen(f func(T) Option[T]) Option[T] {
	return n
}

func (n Nothing[T]) Or(opt2 Option[T]) Option[T] {
	return opt2
}

func (n Nothing[T]) OrElse(f func() Option[T]) Option[T] {
	return f()
}

func (n Nothing[T]) Xor(opt2 Option[T]) Option[T] {
	return opt2
}

//=====================================================

// Returns opt2 if both options contain a value. Otherwise returns Nothing.
func And[T any, U any](opt1 Option[T], opt2 Option[U]) Option[U] {
	if opt1.IsNothing() {
		return Nothing[U]{}
	}
	return opt2
}

// Calls f with the inner value of the option. Returns Nothing if there is no value.
func AndThen[T any, U any](opt1 Option[T], f func(T) Option[U]) Option[U] {
	return OptionMonad[T, U]{}.Bind(opt1, f)
}

// Flattens an option of type Option[Option[T]] to just Option[T].
func Flatten[T any](opt Option[Option[T]]) Option[T] {
	val, ok := opt.Get()
	if !ok {
		return Nothing[T]{}
	}
	return val
}

// Maps the inner value of an option via f. Returns Nothing if there is no value.
func Map[T any, U any](opt Option[T], f func(T) U) Option[U] {
	val, ok := opt.Get()
	if !ok {
		return Nothing[U]{}
	}
	return Some[U]{Value: f(val)}
}

// Maps the inner value of an option via f. If the option is Nothing, returns the default.
// Arguments are eagerly evaluated. Consider MapOrElse if you are passing the result of a function call.
func MapOr[T any, U any](opt Option[T], defaultValue U, f func(T) U) Option[U] {
	val, ok := opt.Get()
	if !ok {
		return Some[U]{Value: defaultValue}
	}
	return Some[U]{Value: f(val)}
}

// Maps the inner value of an option via f. If the option is Nothing, calls defaultFunc to provide a default.
// defaultFunc() is lazily evaluated.
func MapOrElse[T any, U any](opt Option[T], defaultFunc func() U, f func(T) U) Option[U] {
	val, ok := opt.Get()
	if !ok {
		return Some[U]{Value: defaultFunc()}
	}
	return Some[U]{Value: f(val)}
}
