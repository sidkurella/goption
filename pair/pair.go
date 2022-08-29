package pair

// Pair represents a 2-tuple of values.
type Pair[T any, U any] struct {
	First  T
	Second U
}
