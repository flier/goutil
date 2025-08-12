//go:build go1.23

package xiter

import "iter"

// Enumerate creates an iterator which gives the current iteration count as well as the next value.
func Enumerate[T any](x iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var i int
		for v := range x {
			if !yield(i, v) {
				break
			}

			i += 1
		}
	}
}
