//go:build go1.23

package xiter

import "iter"

// Take creates an iterator that yields the first n elements, or fewer if the underlying iterator ends sooner.
func Take[T any](x iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}

		var i int

		for v := range x {
			if !yield(v) {
				break
			}

			if i += 1; i >= n {
				break
			}
		}
	}
}

// TakeFunc creates an iterator that yields the first n elements, or fewer if the underlying iterator ends sooner.
func TakeFunc[T any](n int) MappingFunc[T, T] {
	return bind2(Take[T], n)
}

// Take2 creates an iterator that yields the first n key-values, or fewer if the underlying iterator ends sooner.
func Take2[K, V any](x iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			return
		}

		var i int

		for k, v := range x {
			if !yield(k, v) {
				break
			}

			if i += 1; i >= n {
				break
			}
		}
	}
}

// Take2Func creates an iterator that yields the first n key-values, or fewer if the underlying iterator ends sooner.
func Take2Func[K, V any](n int) Mapping2Func[K, V, V] {
	return bind2(Take2[K, V], n)
}

// TakeWhile creates an iterator that yields elements based on a predicate f.
func TakeWhile[T any](x iter.Seq[T], f func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range x {
			if !f(v) {
				break
			}

			if !yield(v) {
				break
			}
		}
	}
}

// TakeWhileFunc creates an iterator that yields elements based on a predicate f.
func TakeWhileFunc[T any](f func(T) bool) MappingFunc[T, T] {
	return bind2(TakeWhile, f)
}

// TakeWhile2 creates an iterator that yields key-values based on a predicate f.
func TakeWhile2[K, V any](x iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range x {
			if !f(k, v) {
				break
			}

			if !yield(k, v) {
				break
			}
		}
	}
}

// TakeWhile creates an iterator that yields key-values based on a predicate f.
func TakeWhile2Func[K, V any](f func(K, V) bool) Mapping2Func[K, V, V] {
	return bind2(TakeWhile2, f)
}
