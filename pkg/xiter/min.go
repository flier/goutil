//go:build go1.23

package xiter

import (
	"cmp"
	"iter"
)

// Min returns the minimum element of an iterator.
//
// If several elements are equally minimum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func Min[T cmp.Ordered](x iter.Seq[T]) (r T) {
	var p *T

	for v := range x {
		if p == nil {
			p = &v
		} else {
			*p = min(*p, v)
		}
	}

	if p != nil {
		r = *p
	}

	return
}

// MinBy returns the element that gives the minimum value with respect to the specified comparison function.
//
// If several elements are equally minimum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MinBy[T any](x iter.Seq[T], f func(T, T) int) (r T) {
	var p *T

	for v := range x {
		if p == nil || f(*p, v) > 0 {
			p = &v
		}
	}

	if p != nil {
		r = *p
	}

	return
}

// MinByFunc returns the element that gives the minimum value with respect to the specified comparison function.
//
// If several elements are equally minimum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MinByFunc[T any](f func(T, T) int) ReductionFunc[T, T] {
	return bind2(MinBy, f)
}

// MinByKey returns the element that gives the minimum value from the specified function.
//
// If several elements are equally minimum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MinByKey[T any, B cmp.Ordered](x iter.Seq[T], f func(T) B) (r T) {
	var p *T

	for v := range x {
		if p == nil {
			p = &v
		} else if f(*p) >= f(v) {
			p = &v
		}
	}

	if p != nil {
		r = *p
	}

	return
}

// MinByKeyFunc returns the element that gives the minimum value from the specified function.
//
// If several elements are equally minimum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MinByKeyFunc[T any, B cmp.Ordered](f func(T) B) ReductionFunc[T, T] {
	return bind2(MinByKey, f)
}
