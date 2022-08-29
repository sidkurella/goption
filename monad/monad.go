package monad

type Monad[MT any, MU any, T any] interface {
	Bind(MT, func(T) MU) MU
	Return(T) MT
}

func Then[MT any, MU any, T any](monad Monad[MT, MU, T], f func(T) MT, t T, v MU) MU {
	f(t)
	return v
}
