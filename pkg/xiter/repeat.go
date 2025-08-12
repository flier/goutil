//go:build go1.23

package xiter

import "iter"

// Repeat creates a new iterator that endlessly repeats a single element.
func Repeat[T any](e T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for yield(e) {
		}
	}
}

// Repeat2 creates a new iterator that endlessly repeats a single key-value.
func Repeat2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for yield(k, v) {
		}
	}
}

// RepeatN creates a new iterator that repeats a single element a given number of times.
func RepeatN[T any](e T, n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < n; i++ {
			if !yield(e) {
				break
			}
		}
	}
}

// RepeatN creates a new iterator that repeats a single key-value a given number of times.
func RepeatN2[K, V any](k K, v V, n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := 0; i < n; i++ {
			if !yield(k, v) {
				break
			}
		}
	}
}

// RepeatWith creates a new iterator that repeats elements of type T endlessly by applying the provided function f
func RepeatWith[T any](f func() T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for yield(f()) {
		}
	}
}

// RepeatWith creates a new iterator that repeats key-value endlessly by applying the provided function f
func RepeatWith2[K, V any](f func() (K, V)) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for yield(f()) {
		}
	}
}
