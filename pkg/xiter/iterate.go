//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/tuple"
)

// Iterate creates an infinite iterator by repeatedly applying the given function f to the initial value init.
func Iterate[T any](init T, f func(T) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		v := init

		for {
			if !yield(v) {
				return
			}

			v = f(v)
		}
	}
}

// Iterate2 creates an infinite iterator by repeatedly applying the given function f to the initial value init.
func Iterate2[T, K, V any](init tuple.Tuple2[K, V], f func(K, V) (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		k, v := init.Unpack()

		for {
			if !yield(k, v) {
				return
			}

			k, v = f(k, v)
		}
	}
}
