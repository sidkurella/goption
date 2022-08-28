package monad

type Monad[MT any, MU any, T any, U any] interface {
	Bind(MT, func(T) MU) MU
	Return(T) MT
}

func Then[MT any, MU any, T any, U any](monad Monad[MT, MU, T, U], f func(T) MT, t T, v MU) MU {
	f(t)
	return v
}
