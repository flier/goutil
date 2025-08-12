//go:build go1.23

package xiter

import "iter"

// ForEachFunc calls a function f on each element of an iterator.
func ForEach[T any](x iter.Seq[T], f func(T)) {
	for v := range x {
		f(v)
	}
}

// ForEachFunc calls a function f on each element of an iterator.
func ForEachFunc[T any](f func(T)) func(iter.Seq[T]) {
	return func(x iter.Seq[T]) {
		for v := range x {
			f(v)
		}
	}
}

// ForEach2 calls a function f on each key-value of an iterator.
func ForEach2[K, V any](x iter.Seq2[K, V], f func(K, V)) {
	for k, v := range x {
		f(k, v)
	}
}

// ForEach2Func calls a function f on each key-value of an iterator.
func ForEach2Func[K, V any](f func(K, V)) func(iter.Seq2[K, V]) {
	return func(x iter.Seq2[K, V]) {
		for k, v := range x {
			f(k, v)
		}
	}
}
