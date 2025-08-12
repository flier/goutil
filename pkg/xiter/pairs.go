//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/tuple"
)

// Pairs returns an iterator of pairs from the given iterator of key-values.
func Pairs[K, V any](x iter.Seq2[K, V]) iter.Seq[tuple.Tuple2[K, V]] {
	return func(yield func(tuple.Tuple2[K, V]) bool) {
		for k, v := range x {
			if !yield(tuple.New2(k, v)) {
				break
			}
		}
	}
}

// Unpairs returns an iterator of key-values from the given iterator of pairs.
func Unpairs[K, V any](x iter.Seq[tuple.Tuple2[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for t := range x {
			if !yield(t.Unpack()) {
				break
			}
		}
	}
}
