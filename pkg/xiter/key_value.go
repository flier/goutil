//go:build go1.23

package xiter

import (
	"iter"
)

// Keys returns an iterator of keys from the given iterator of key-value pairs.
func Keys[K, V any](x iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range x {
			if !yield(k) {
				break
			}
		}
	}
}

// Values returns an iterator of values from the given iterator of key-value pairs.
func Values[K, V any](x iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range x {
			if !yield(v) {
				break
			}
		}
	}
}

// Swap returns an iterator of value-key pairs from the given iterator of key-value pairs.
func Swap[K, V any](x iter.Seq2[K, V]) iter.Seq2[V, K] {
	return func(yield func(V, K) bool) {
		for k, v := range x {
			if !yield(v, k) {
				break
			}
		}
	}
}
