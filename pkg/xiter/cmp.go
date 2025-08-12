//go:build go1.23

package xiter

import (
	"cmp"
	"iter"
)

// Compare compares the elements of tow iterators.
func Compare[T cmp.Ordered](l, r iter.Seq[T]) int {
	ln, ls := iter.Pull(l)
	rn, rs := iter.Pull(r)

	defer ls()
	defer rs()

	for {
		lv, lok := ln()
		rv, rok := rn()

		if !lok && !rok {
			return 0
		} else if lok && !rok {
			return 1
		} else if !lok && rok {
			return -1
		}

		r := cmp.Compare(lv, rv)
		if r != 0 {
			return r
		}
	}
}

// CompareBy compares the elements of tow iterators with respect to the specified comparison function f.
func CompareBy[T any](l, r iter.Seq[T], f func(T, T) int) int {
	ln, ls := iter.Pull(l)
	rn, rs := iter.Pull(r)

	defer ls()
	defer rs()

	for {
		lv, lok := ln()
		rv, rok := rn()

		if !lok && !rok {
			return 0
		} else if lok && !rok {
			return 1
		} else if !lok && rok {
			return -1
		}

		r := f(lv, rv)
		if r != 0 {
			return r
		}
	}
}

// CompareByFunc compares the elements of tow iterators with respect to the specified comparison function f.
func CompareByFunc[T any](f func(T, T) int) CompareFunc[T] {
	return bind3(CompareBy, f)
}

// CompareByKey compares two sequences l and r element-wise using a key extraction function f.
//
// The comparison is performed by applying f to each element and comparing the resulting keys
// using cmp.Compare. The function returns 0 if both sequences are equal, 1 if l is greater,
// and -1 if r is greater. The comparison stops at the first difference or when one sequence ends.
//
// T: the element type of the sequences.
// B: the key type, which must be cmp.Ordered.
// l, r: input sequences to compare.
// f: function to extract the key from each element.
//
// Returns:
//
//	 0 if sequences are equal,
//	 1 if l is greater,
//	-1 if r is greater.
func CompareByKey[T any, B cmp.Ordered](l, r iter.Seq[T], f func(T) B) int {
	ln, ls := iter.Pull(l)
	rn, rs := iter.Pull(r)

	defer ls()
	defer rs()

	for {
		lv, lok := ln()
		rv, rok := rn()

		if !lok && !rok {
			return 0
		} else if lok && !rok {
			return 1
		} else if !lok && rok {
			return -1
		}

		r := cmp.Compare(f(lv), f(rv))
		if r != 0 {
			return r
		}
	}
}

// CompareByKeyFunc compares the elements of tow iterators that gives the value from the specified function f.
func CompareByKeyFunc[T any, B cmp.Ordered](f func(T) B) CompareFunc[T] {
	return bind3(CompareByKey, f)
}

// Equal compares two iterators and returns true if they contain the same elements in the same order.
func Equal[T cmp.Ordered](l, r iter.Seq[T]) bool {
	return Compare(l, r) == 0
}

// EqualByFunc returns a function that compares two iterators using the provided comparison function f.
func EqualBy[T any](x, y iter.Seq[T], f func(T, T) bool) bool {
	cmp := CompareByFunc(func(x, y T) int {
		if f(x, y) {
			return 0
		}

		return 1
	})

	return cmp(x, y) == 0
}

// EqualByFunc returns a function that compares two iterators using the provided comparison function f.
func EqualByFunc[T any](f func(T, T) bool) func(iter.Seq[T], iter.Seq[T]) bool {
	return bind3(EqualBy, f)
}

func EqualByKey[T any, B cmp.Ordered](x, y iter.Seq[T], f func(T) B) bool {
	return CompareByKey(x, y, f) == 0
}

// EqualByKeyFunc returns a function that compares two iterators using the provided comparison function f.
func EqualByKeyFunc[T any, B cmp.Ordered](f func(T) B) func(iter.Seq[T], iter.Seq[T]) bool {
	return bind3(EqualByKey, f)
}

// NotEqual compares two iterators and returns true if they do not contain the same elements in the same order.
func NotEqual[T cmp.Ordered](l, r iter.Seq[T]) bool {
	return Compare(l, r) != 0
}

// LessThan compares two iterators and returns true if the elements in the first sequence are lexicographically less than the elements in the second sequence.
func LessThan[T cmp.Ordered](l, r iter.Seq[T]) bool {
	return Compare(l, r) < 0
}

// LessOrEqual compares two iterators and returns true if the elements in the first sequence are lexicographically less than or equal to the elements in the second sequence.
func LessOrEqual[T cmp.Ordered](l, r iter.Seq[T]) bool {
	return Compare(l, r) <= 0
}

// GreaterThan compares two iterators and returns true if the elements in the first sequence are lexicographically greater than the elements in the second sequence.
func GreaterThan[T cmp.Ordered](l, r iter.Seq[T]) bool {
	return Compare(l, r) > 0
}

// GreaterOrEqual compares two iterators and returns true if the elements in the first sequence are lexicographically greater than or equal to the elements in the second sequence.
func GreaterOrEqual[T cmp.Ordered](l, r iter.Seq[T]) bool {
	return Compare(l, r) >= 0
}
