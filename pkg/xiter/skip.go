//go:build go1.23

package xiter

import "iter"

// Skip creates an iterator that skips the first n elements.
func Skip[T any](x iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		var i int

		for v := range x {
			if i += 1; i <= n {
				continue
			}

			if !yield(v) {
				break
			}
		}
	}
}

// SkipFunc creates an iterator that skips the first n elements.
func SkipFunc[T any](n int) MappingFunc[T, T] {
	return bind2(Skip[T], n)
}

// Skip2 creates an iterator that skips the first n key-value.
func Skip2[K, V any](x iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var i int

		for k, v := range x {
			if i += 1; i <= n {
				continue
			}

			if !yield(k, v) {
				break
			}
		}
	}
}

// Skip2Func creates an iterator that skips the first n key-value.
func Skip2Func[K, V any](n int) Mapping2Func[K, V, V] {
	return bind2(Skip2[K, V], n)
}

// SkipWhile creates an iterator that skips elements based on a predicate f.
func SkipWhile[T any](x iter.Seq[T], f func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range x {
			if f(v) {
				continue
			}

			if !yield(v) {
				break
			}
		}
	}
}

// SkipWhileFunc creates an iterator that skips elements based on a predicate f.
func SkipWhileFunc[T any](f func(T) bool) MappingFunc[T, T] {
	return bind2(SkipWhile, f)
}

// SkipWhile2 creates an iterator that skips elements based on a predicate f.
func SkipWhile2[K, V any](x iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range x {
			if f(k, v) {
				continue
			}

			if !yield(k, v) {
				break
			}
		}
	}
}

// SkipWhile2Func creates an iterator that skips elements based on a predicate f.
func SkipWhile2Func[K, V any](f func(K, V) bool) Mapping2Func[K, V, V] {
	return bind2(SkipWhile2, f)
}
