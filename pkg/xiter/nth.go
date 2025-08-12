//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/opt"
	"github.com/flier/goutil/pkg/tuple"
)

// Nth returns the nth element of the iterator.
func Nth[T any](x iter.Seq[T], n int) opt.Option[T] {
	var i int
	for v := range x {
		if i += 1; i > n {
			return opt.Some(v)
		}
	}

	return opt.None[T]()
}

// NthFunc returns the nth element of the iterator.
func NthFunc[T any](n int) ReductionFunc[T, opt.Option[T]] {
	return bind2(Nth[T], n)
}

// Nth2 returns the nth key-value of the iterator.
func Nth2[K, V any](x iter.Seq2[K, V], n int) opt.Option[tuple.Tuple2[K, V]] {
	var i int
	for k, v := range x {
		if i += 1; i > n {
			return opt.Some(tuple.New2(k, v))
		}
	}

	return opt.None[tuple.Tuple2[K, V]]()
}

// Nth2Func returns the nth key-value of the iterator.
func Nth2Func[K, V any](n int) Reduction2Func[K, V, opt.Option[tuple.Tuple2[K, V]]] {
	return bind2(Nth2[K, V], n)
}
