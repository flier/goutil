//go:build go1.23

package xiter

import "iter"

// Sum sums the elements of an iterator.
func Sum[T Number](x ...iter.Seq[T]) (r T) {
	for _, i := range x {
		for v := range i {
			r += v
		}
	}

	return
}

// SumBy sums the element that gives the value from the specified function.
func SumBy[T any, B Number](x iter.Seq[T], f func(T) B) (r B) {
	for v := range x {
		r += f(v)
	}

	return
}

// SumByFunc sums the element that gives the value from the specified function.
func SumByFunc[T any, B Number](f func(T) B) ReductionFunc[T, B] {
	return bind2(SumBy[T, B], f)
}
