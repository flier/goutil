//go:build go1.23

package xiter

import "iter"

// All returns true if all elements in the provided sequence x satisfy the predicate function f.
//
// It iterates over each element in x, applying f to each element. If f returns false for any element,
// All returns false immediately. If f returns true for all elements, All returns true.
//
// T is the type of elements in the sequence.
//
// Parameters:
//
//	x - an iter.Seq[T], the sequence to iterate over.
//	f - a predicate function that takes an element of type T and returns a bool.
//
// Returns:
//
//	bool - true if all elements satisfy f, false otherwise.
func All[T any](x iter.Seq[T], f func(T) bool) bool {
	for v := range x {
		if !f(v) {
			return false
		}
	}

	return true
}

// All2 returns true if the predicate function f returns true for all key-value pairs in the given iterator x.
//
// It iterates over each (K, V) pair in x, applying f to each pair.
// If f returns false for any pair, All2 returns false immediately.
// Otherwise, it returns true after checking all pairs.
//
// T is the type of elements in the sequence.
//
// Parameters:
//
//	x - an iter.Seq2[K, V], the sequence to iterate over.
//	f - a predicate function that takes an element of type K, V and returns a bool.
//
// Returns:
//
//	bool - true if all elements satisfy f, false otherwise.
func All2[K, V any](x iter.Seq2[K, V], f func(K, V) bool) bool {
	for k, v := range x {
		if !f(k, v) {
			return false
		}
	}

	return true
}

// AllFunc tests if every element of the iterator matches a predicate f.
func AllFunc[T any](f func(T) bool) ReductionFunc[T, bool] {
	return bind2(All, f)
}

// All tests if every key-value of the iterator matches a predicate f.
func All2Func[K, V any](f func(K, V) bool) Reduction2Func[K, V, bool] {
	return bind2(All2, f)
}
