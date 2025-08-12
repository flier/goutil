//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/opt"
	"github.com/flier/goutil/pkg/tuple"
)

// Find searches for an element of an iterator that satisfies a predicate f.
func Find[T any](x iter.Seq[T], f func(T) bool) opt.Option[T] {
	for v := range x {
		if f(v) {
			return opt.Some(v)
		}
	}

	return opt.None[T]()
}

// FindFunc searches for an element of an iterator that satisfies a predicate f.
func FindFunc[T any](f func(T) bool) ReductionFunc[T, opt.Option[T]] {
	return func(x iter.Seq[T]) opt.Option[T] {
		for v := range x {
			if f(v) {
				return opt.Some(v)
			}
		}

		return opt.None[T]()
	}
}

// Find2 searches for a key-value of an iterator that satisfies a predicate f.
func Find2[K, V any](x iter.Seq2[K, V], f func(K, V) bool) opt.Option[tuple.Tuple2[K, V]] {
	for k, v := range x {
		if f(k, v) {
			return opt.Some(tuple.New2(k, v))
		}
	}

	return opt.None[tuple.Tuple2[K, V]]()
}

// Find2Func searches for a key-value of an iterator that satisfies a predicate f.
func Find2Func[K, V any](f func(K, V) bool) Reduction2Func[K, V, opt.Option[tuple.Tuple2[K, V]]] {
	return func(x iter.Seq2[K, V]) opt.Option[tuple.Tuple2[K, V]] {
		for k, v := range x {
			if f(k, v) {
				return opt.Some(tuple.New2(k, v))
			}
		}

		return opt.None[tuple.Tuple2[K, V]]()
	}
}

// FindMap applies function f to the elements of iterator and returns the first result that satisfies a predicate f.
func FindMap[T, B any](x iter.Seq[T], f func(T) (B, bool)) opt.Option[B] {
	for v := range x {
		if b, found := f(v); found {
			return opt.Some(b)
		}
	}

	return opt.None[B]()
}

// FindMapFunc applies function f to the elements of iterator and returns the first result that satisfies a predicate f.
func FindMapFunc[T, B any](f func(T) (B, bool)) ReductionFunc[T, opt.Option[B]] {
	return bind2(FindMap, f)
}

// FindMap2 applies function f to the elements of iterator and returns the first result that satisfies a predicate f.
func FindMap2[K, V, B any](x iter.Seq2[K, V], f func(K, V) (B, bool)) opt.Option[tuple.Tuple2[K, B]] {
	for k, v := range x {
		if b, found := f(k, v); found {
			return opt.Some(tuple.New2(k, b))
		}
	}

	return opt.None[tuple.Tuple2[K, B]]()
}

// FindMap2Func applies function f to the elements of iterator and returns the first result that satisfies a predicate f.
func FindMap2Func[K, V, B any](f func(K, V) (B, bool)) Reduction2Func[K, V, opt.Option[tuple.Tuple2[K, B]]] {
	return bind2(FindMap2, f)
}
