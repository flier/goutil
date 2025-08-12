//go:build go1.23

package xiter

import "iter"

// Once creates an iterator that yields an element exactly once.
func Once[T any](v T) iter.Seq[T] {
	return func(yield func(T) bool) {
		yield(v)
	}
}

// Once2 creates an iterator that yields a key-value pair exactly once.
func Once2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(k, v)
	}
}

// OnceWith creates an iterator that lazily generates a value exactly once by invoking the provided function f.
func OnceWith[T any](f func() T) iter.Seq[T] {
	return func(yield func(T) bool) {
		yield(f())
	}
}

// OnceWith2 creates an iterator that lazily generates a key-value exactly once by invoking the provided function f.
func OnceWith2[K, V any](f func() (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(f())
	}
}
