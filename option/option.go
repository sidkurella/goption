package option

// Implements Monad[Option[T], Option[U], T, U].
type OptionMonad[T any, U any] struct {
}

// Option type. Either contains a value, or does not.
type Option[T any] interface {
	IsSome() bool
	IsNothing() bool
	OrElse(T) T
	Get() (T, bool)
}

type Some[T any] struct {
	Value T
}

func (s Some[T]) IsSome() bool {
	return true
}

func (s Some[T]) IsNothing() bool {
	return false
}

func (s Some[T]) OrElse(val T) T {
	return s.Value
}

func (s Some[T]) Get() (T, bool) {
	return s.Value, true
}

type Nothing[T any] struct {
}

func (n Nothing[T]) IsSome() bool {
	return false
}

func (n Nothing[T]) IsNothing() bool {
	return true
}

func (n Nothing[T]) OrElse(val T) T {
	return val
}

func (n Nothing[T]) Get() (T, bool) {
	return *new(T), false
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
