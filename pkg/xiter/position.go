//go:build go1.23

package xiter

import "iter"

// Position searches for an element in an iterator, returning its index, or -1 if not present.
func Position[T any](x iter.Seq[T], f func(T) bool) int {
	var i int
	for v := range x {
		if f(v) {
			return i
		}

		i += 1
	}

	return -1
}

// PositionFunc searches for an element in an iterator, returning its index, or -1 if not present.
func PositionFunc[T any](f func(T) bool) ReductionFunc[T, int] {
	return bind2(Position, f)
}

// Position2 searches for a key-value in an iterator, returning its index, or -1 if not present.
func Position2[K, V any](x iter.Seq2[K, V], f func(K, V) bool) int {
	var i int
	for k, v := range x {
		if f(k, v) {
			return i
		}

		i += 1
	}

	return -1
}

// Position2Func searches for a key-value in an iterator, returning its index, or -1 if not present.
func Position2Func[K, V any](f func(K, V) bool) Reduction2Func[K, V, int] {
	return bind2(Position2[K, V], f)
}
