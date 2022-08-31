package functools

// Curry takes a function of 2 arguments and curries it.
// It returns a function of 1 argument that returns a function of the second argument.
func Curry[A any, B any, C any](f func(A, B) C) func(A) func(B) C {
	return func(a A) func(b B) C {
		return func(b B) C {
			return f(a, b)
		}
	}
}

// Uncurry is the inverse of Curry.
// It takes a curried function and returns it to a function on 2 arguments.
func Uncurry[A any, B any, C any](f func(A) func(B) C) func(A, B) C {
	return func(a A, b B) C {
		return f(a)(b)
	}
}

// Memoize caches all call-result pairs from the given function.
// It will then use the cache to return the information from subsequent calls.
// NOTE: This cache grows without bound.
func Memoize[A comparable, B any](f func(A) B) func(A) B {
	m := map[A]B{}
	return func(a A) B {
		saved, ok := m[a]
		if ok {
			return saved
		}

		b := f(a)
		m[a] = b
		return b
	}
}
