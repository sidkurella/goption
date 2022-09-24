package pair

// Pair represents a 2-tuple of values.
type Pair[T any, U any] struct {
	First  T
	Second U
}

// Takes 2 values and returns a pair from them.
func From[T any, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{
		First:  first,
		Second: second,
	}
}

// Returns a tuple of values from the given pair.
func (p Pair[T, U]) Into() (T, U) {
	return p.First, p.Second
}
