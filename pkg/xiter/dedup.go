//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/tuple"
)

// Dedup creates an iterator that only emits elements if they are different from the last emitted element.
func Dedup[T comparable](x iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var prev *T

		for v := range x {
			if prev == nil {
				prev = new(T)
			} else if *prev == v {
				continue
			}

			if !yield(v) {
				break
			}

			*prev = v
		}
	}
}

// DedupBy creates an iterator that only emits elements if they are different from the last emitted element,
// as determined by the provided comparison function f.
//
// The comparison function f should return true if the two elements are considered equal, and false otherwise.
func DedupBy[T any](x iter.Seq[T], f func(T, T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		var prev *T

		for v := range x {
			if prev == nil {
				prev = new(T)
			} else if f(*prev, v) {
				continue
			}

			if !yield(v) {
				break
			}

			*prev = v
		}
	}
}

// DedupByFunc creates an iterator that only emits elements if they are different from the last emitted element,
// as determined by the provided comparison function f.
//
// The comparison function f should return true if the two elements are considered equal, and false otherwise.
func DedupByFunc[T any](f func(T, T) bool) MappingFunc[T, T] {
	return bind2(DedupBy, f)
}

// DedupByKey creates an iterator that only emits elements if they are different from the last emitted element.
func DedupByKey[T any, B comparable](x iter.Seq[T], f func(T) B) iter.Seq[T] {
	return func(yield func(T) bool) {
		var prev *T

		for v := range x {
			if prev == nil {
				prev = new(T)
			} else if f(*prev) == f(v) {
				continue
			}

			if !yield(v) {
				break
			}

			*prev = v
		}
	}
}

// DedupByKeyFunc creates an iterator that only emits elements if they are different from the last emitted element.
func DedupByKeyFunc[T any, B comparable](f func(T) B) MappingFunc[T, T] {
	return bind2(DedupByKey, f)
}

// DedupByKey2Func creates an iterator that only emits elements if they are different from the last emitted element.
func DedupByKey2[K, V any, B comparable](x iter.Seq2[K, V], f func(K, V) B) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var prev *tuple.Tuple2[K, V]

		for k, v := range x {
			if prev == nil {
				prev = new(tuple.Tuple2[K, V])
			} else if f(prev.Unpack()) == f(k, v) {
				continue
			}

			if !yield(k, v) {
				break
			}

			*prev = tuple.New2(k, v)
		}
	}
}

// DedupByKey2Func creates an iterator that only emits elements if they are different from the last emitted element.
func DedupByKey2Func[K, V any, B comparable](f func(K, V) B) MappingValueFunc[K, V, V] {
	return bind2(DedupByKey2, f)
}
