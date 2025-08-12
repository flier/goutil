//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/opt"
	"github.com/flier/goutil/pkg/tuple"
)

// Last returns the last element.
func Last[T any](x iter.Seq[T]) opt.Option[T] {
	var last *T
	for v := range x {
		last = &v
	}

	return opt.Wrap(last)
}

// Last2 returns the last key-value.
func Last2[K, V any](x iter.Seq2[K, V]) opt.Option[tuple.Tuple2[K, V]] {
	var last *tuple.Tuple2[K, V]
	for k, v := range x {
		p := tuple.New2(k, v)
		last = &p
	}

	return opt.Wrap(last)
}
