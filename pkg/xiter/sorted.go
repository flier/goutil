//go:build go1.23

package xiter

import (
	"cmp"
	"iter"
)

// IsSorted reports whether x is sorted in ascending order.
func IsSorted[T cmp.Ordered](x iter.Seq[T]) bool {
	var last *T

	for v := range x {
		if last == nil || cmp.Compare(*last, v) <= 0 {
			last = &v
		} else {
			return false
		}
	}

	return true
}

// IsSortedBy reports whether x is sorted in ascending order using the given comparator function.
func IsSortedBy[T any](x iter.Seq[T], f func(T, T) bool) bool {
	var last *T

	for v := range x {
		if last == nil || f(*last, v) {
			last = &v
		} else {
			return false
		}
	}

	return true
}

// IsSortedByFunc reports whether x is sorted in ascending order using the given comparator function.
func IsSortedByFunc[T any](f func(T, T) bool) ReductionFunc[T, bool] {
	return bind2(IsSortedBy, f)
}

// IsSortedByKey reports whether x is sorted in ascending order using the given key extraction function.
func IsSortedByKey[T any, B cmp.Ordered](x iter.Seq[T], f func(T) B) bool {
	var last *T

	for v := range x {
		if last == nil || cmp.Compare(f(*last), f(v)) <= 0 {
			last = &v
		} else {
			return false
		}
	}

	return true
}

// IsSortedByKeyFunc reports whether x is sorted in ascending order using the given key extraction function.
func IsSortedByKeyFunc[T any, B cmp.Ordered](f func(T) B) ReductionFunc[T, bool] {
	return bind2(IsSortedByKey, f)
}
