//go:build go1.23

package xiter

import "iter"

// Chain converts the arguments to iterators and links them together, in a chain.
func Chain[T any](x ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, i := range x {
			for v := range i {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Chain2 converts the arguments to iterators and links them together, in a chain.
func Chain2[K, V any](x ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, i := range x {
			for k, v := range i {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
