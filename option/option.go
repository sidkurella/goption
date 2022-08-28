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
	IsSome() bool
	IsNothing() bool

	Unwrap() T
	UnwrapOr(T) T

	Get() (T, bool)
	Expect(string) T

	IsSomeAnd(func(*T) bool) bool
	Filter(func(*T) bool) Option[T]

	And(Option[T]) Option[T]
	AndThen(func(T) Option[T]) Option[T]
	Or(Option[T]) Option[T]
	OrElse(func() Option[T]) Option[T]
	Xor(Option[T]) Option[T]
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

func And[T any, U any](opt1 Option[T], opt2 Option[U]) Option[U] {
	if opt1.IsNothing() {
		return Nothing[U]{}
	}
	return opt2
}

func AndThen[T any, U any](opt1 Option[T], f func(T) Option[U]) Option[U] {
	return OptionMonad[T, U]{}.Bind(opt1, f)
}

func Flatten[T any](opt Option[Option[T]]) Option[T] {
	val, ok := opt.Get()
	if !ok {
		return Nothing[T]{}
	}
	return val
}
