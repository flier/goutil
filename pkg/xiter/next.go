//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/opt"
	"github.com/flier/goutil/pkg/tuple"
)

// Next returns the next value.
func Next[T any](s iter.Seq[T]) opt.Option[T] {
	next, stop := iter.Pull(s)
	defer stop()

	v, ok := next()
	if ok {
		return opt.Some(v)
	}

	return opt.None[T]()
}

// Next2 returns the next value.
func Next2[K, V any](s iter.Seq2[K, V]) opt.Option[tuple.Tuple2[K, V]] {
	next, stop := iter.Pull2(s)
	defer stop()

	k, v, ok := next()
	if ok {
		return opt.Some(tuple.New2(k, v))
	}

	return opt.None[tuple.Tuple2[K, V]]()
}
