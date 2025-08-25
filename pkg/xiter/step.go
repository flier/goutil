//go:build go1.23

package xiter

import "iter"

// StepBy creates an iterator starting at the same point, but stepping by the given amount at each iteration.
func StepBy[T any](x iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		var i int

		for v := range x {
			if i%n == 0 {
				if !yield(v) {
					break
				}
			}

			i += 1
		}
	}
}

// StepByFunc creates an iterator starting at the same point, but stepping by the given amount at each iteration.
func StepByFunc[T any](n int) MappingFunc[T, T] {
	return bind2(StepBy[T], n)
}

// StepBy2 creates an iterator starting at the same point, but stepping by the given amount at each iteration.
func StepBy2[K, V any](x iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var i int

		for k, v := range x {
			if i%n == 0 {
				if !yield(k, v) {
					break
				}
			}

			i += 1
		}
	}
}

// StepBy2Func creates an iterator starting at the same point, but stepping by the given amount at each iteration.
func StepBy2Func[K, V any](n int) MappingValueFunc[K, V, V] {
	return bind2(StepBy2[K, V], n)
}
