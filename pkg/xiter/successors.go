//go:build go1.23

package xiter

import "iter"

// Successors creates a new iterator where each successive item is computed based on the preceding one.
func Successors[T any](v T, f func(T) (T, bool)) iter.Seq[T] {
	return func(yield func(T) bool) {
		if !yield(v) {
			return
		}

		for {
			var ok bool
			if v, ok = f(v); !ok {
				break
			}

			if !yield(v) {
				break
			}
		}
	}
}
