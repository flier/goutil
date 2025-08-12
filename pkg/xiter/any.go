//go:build go1.23

package xiter

import "iter"

// Any returns true if any element in the provided sequence x satisfies the predicate function f.
//
// It iterates over the sequence x, applying f to each element, and returns true upon the first match.
// If no elements satisfy the predicate, it returns false.
//
// Example usage:
//
//	found := Any(seq, func(v T) bool { return v == target })
//
// Parameters:
//
//	x - an iterable sequence of type T
//	f - a predicate function that takes an element of type T and returns a bool
//
// Returns:
//
//	bool - true if any element satisfies f, false otherwise
func Any[T any](x iter.Seq[T], f func(T) bool) bool {
	for v := range x {
		if f(v) {
			return true
		}
	}

	return false
}

// Any2 iterates over a two-value sequence x and applies the predicate function f to each key-value pair.
//
// It returns true if f returns true for any pair; otherwise, it returns false.
// This function is generic over key (K) and value (V) types.
func Any2[K, V any](x iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range x {
		if f(k, v) {
			return true
		}
	}

	return false
}

// AnyFunc tests if any element of the iterator matches a predicate f.
func AnyFunc[T any](f func(T) bool) ReductionFunc[T, bool] {
	return bind2(Any, f)
}

// Any2Func tests if any key-value of the iterator matches a predicate f.
func Any2Func[K, V any](f func(K, V) bool) Reduction2Func[K, V, bool] {
	return bind2(Any2, f)
}
