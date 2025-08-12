//go:build go1.23

package xiter

import "iter"

// Cycle repeats an iterator endlessly.
func Cycle[T any](x iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			for v := range x {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Cycle2 repeats an iterator endlessly.
func Cycle2[K, V any](x iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			for k, v := range x {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
