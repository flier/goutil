//go:build go1.23

package xiter

import (
	"cmp"
	"iter"
)

// Max returns the maximum element of an iterator.
//
// If several elements are equally maximum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func Max[T cmp.Ordered](x iter.Seq[T]) (r T) {
	var p *T

	for v := range x {
		if p == nil {
			p = &v
		} else {
			*p = max(*p, v)
		}
	}

	if p != nil {
		r = *p
	}

	return
}

// MaxBy returns the element that gives the maximum value with respect to the specified comparison function.
//
// If several elements are equally maximum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MaxBy[T any](x iter.Seq[T], f func(T, T) int) (r T) {
	var p *T

	for v := range x {
		if p == nil || f(*p, v) <= 0 {
			p = &v
		}
	}

	if p != nil {
		r = *p
	}

	return
}

// MaxByFunc returns the element that gives the maximum value with respect to the specified comparison function.
//
// If several elements are equally maximum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MaxByFunc[T any](f func(T, T) int) ReductionFunc[T, T] {
	return bind2(MaxBy, f)
}

// MaxByKey returns the element that gives the maximum value from the specified function.
//
// If several elements are equally maximum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MaxByKey[T any, B cmp.Ordered](x iter.Seq[T], f func(T) B) (r T) {
	var p *T

	for v := range x {
		if p == nil {
			p = &v
		} else if f(*p) <= f(v) {
			p = &v
		}
	}

	if p != nil {
		r = *p
	}

	return
}

// MaxByKeyFunc returns the element that gives the maximum value from the specified function.
//
// If several elements are equally maximum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MaxByKeyFunc[T any, B cmp.Ordered](f func(T) B) ReductionFunc[T, T] {
	return bind2(MaxByKey, f)
}
