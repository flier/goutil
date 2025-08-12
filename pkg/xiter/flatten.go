//go:build go1.23

package xiter

import "iter"

// Flatten creates an iterator that flattens nested iterators.
func Flatten[T iter.Seq[V], V any](x iter.Seq[T]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range x {
			for i := range v {
				if !yield(i) {
					break
				}
			}
		}
	}
}

// Flatten2 creates an iterator that flattens nested iterators.
func Flatten2[T iter.Seq2[K, V], K, V any](x iter.Seq[T]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := range x {
			for k, v := range i {
				if !yield(k, v) {
					break
				}
			}
		}
	}
}
