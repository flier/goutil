//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/opt"
	"github.com/flier/goutil/pkg/tuple"
)

// First returns the first element.
func First[T any](x iter.Seq[T]) opt.Option[T] {
	next, stop := iter.Pull(x)
	defer stop()

	v, ok := next()
	if ok {
		return opt.Some(v)
	}

	return opt.None[T]()
}

// First2 returns the last key-value.
func First2[K, V any](x iter.Seq2[K, V]) opt.Option[tuple.Tuple2[K, V]] {
	next, stop := iter.Pull2(x)
	defer stop()

	k, v, ok := next()
	if ok {
		return opt.Some(tuple.New2(k, v))
	}

	return opt.None[tuple.Tuple2[K, V]]()
}
