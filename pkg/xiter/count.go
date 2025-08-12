//go:build go1.23

package xiter

import "iter"

// Count returns the number of iterations.
func Count[T any](x ...iter.Seq[T]) (n int) {
	for _, i := range x {
		for range i {
			n += 1
		}
	}

	return
}

// Count2 returns the number of iterations.
func Count2[K, V any](x ...iter.Seq2[K, V]) (n int) {
	for _, i := range x {
		for range i {
			n += 1
		}
	}

	return
}
