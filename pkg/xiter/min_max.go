//go:build go1.23

package xiter

import (
	"cmp"
	"iter"

	"github.com/flier/goutil/pkg/tuple"
)

// MinMax returns the minimum and maximum elements of an iterator.
//
// If several elements are equally minimum or maximum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MinMax[T cmp.Ordered](x iter.Seq[T]) tuple.Tuple2[T, T] {
	var lo, hi *T

	for v := range x {
		if lo == nil {
			lo = new(T)
			*lo = v
		} else {
			*lo = min(*lo, v)
		}

		if hi == nil {
			hi = new(T)
			*hi = v
		} else {
			*hi = max(*hi, v)
		}
	}

	if lo == nil || hi == nil {
		var v T

		return tuple.New2(v, v)
	}

	return tuple.New2(*lo, *hi)
}

// MinMaxBy returns the minimum and maximum elements of an iterator with respect to the specified comparison function.
//
// If several elements are equally minimum or maximum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MinMaxBy[T cmp.Ordered](x iter.Seq[T], f func(T, T) int) tuple.Tuple2[T, T] {
	var lo, hi *T

	for v := range x {
		if lo == nil {
			lo = new(T)
			*lo = v
		} else if f(*lo, v) >= 0 {
			*lo = v
		}

		if hi == nil {
			hi = new(T)
			*hi = v
		} else if f(*hi, v) <= 0 {
			*hi = v
		}
	}

	if lo == nil || hi == nil {
		var v T

		return tuple.New2(v, v)
	}

	return tuple.New2(*lo, *hi)
}

// MinMaxByFunc returns the minimum and maximum elements of an iterator with respect to the specified comparison function.
//
// If several elements are equally minimum or maximum, the last element is returned.
// If the iterator is empty, an empty value is returned.
func MinMaxByFunc[T cmp.Ordered](f func(T, T) int) ReductionFunc[T, tuple.Tuple2[T, T]] {
	return bind2(MinMaxBy[T], f)
}

// MinMaxByKey returns the element that gives the minimum and maximum value from the specified function.
func MinMaxByKey[T any, B cmp.Ordered](x iter.Seq[T], f func(T) B) tuple.Tuple2[T, T] {
	var lo, hi *T

	for v := range x {
		b := f(v)

		if lo == nil {
			lo = new(T)
			*lo = v
		} else if f(*lo) >= b {
			*lo = v
		}

		if hi == nil {
			hi = new(T)
			*hi = v
		} else if f(*hi) <= b {
			*hi = v
		}
	}

	if lo == nil || hi == nil {
		var v T

		return tuple.New2(v, v)
	}

	return tuple.New2(*lo, *hi)
}

// MinMaxByKeyFunc returns the element that gives the minimum and maximum value from the specified function.
func MinMaxByKeyFunc[T any, B cmp.Ordered](f func(T) B) ReductionFunc[T, tuple.Tuple2[T, T]] {
	return bind2(MinMaxByKey[T, B], f)
}
