//go:build go1.23

package xiter

import "iter"

// Accumulate makes an iterator that returns accumulated sums.
func Accumulate[T Number](x iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var acc *T

		for v := range x {
			if acc == nil {
				acc = &v
			} else {
				*acc += v
			}

			if !yield(*acc) {
				return
			}
		}
	}
}

// AccumulateBy makes an iterator that returns accumulated results from other binary functions.
func AccumulateBy[T Number](x iter.Seq[T], f func(T, T) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		var acc *T

		for v := range x {
			if acc == nil {
				acc = &v
			} else {
				*acc = f(*acc, v)
			}

			if !yield(*acc) {
				return
			}
		}
	}
}

// AccumulateByFunc makes an iterator that returns accumulated results from other binary functions.
func AccumulateByFunc[T Number](f func(T, T) T) MappingFunc[T, T] {
	return bind2(AccumulateBy, f)
}
