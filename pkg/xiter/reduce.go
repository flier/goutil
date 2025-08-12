//go:build go1.23

package xiter

import (
	"iter"
)

// Reduce reduces the elements to a single one, by repeatedly applying a reducing operation f.
func Reduce[T any](x iter.Seq[T], f func(T, T) T) (r T) {
	var last *T

	for i := range x {
		if last == nil {
			last = &i
		} else {
			*last = f(*last, i)
		}
	}

	if last != nil {
		r = *last
	}

	return r
}

// ReduceFunc reduces the elements to a single one, by repeatedly applying a reducing operation f.
func ReduceFunc[T any](f func(T, T) T) ReductionFunc[T, T] {
	return bind2(Reduce, f)
}
