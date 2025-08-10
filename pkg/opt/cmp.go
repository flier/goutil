//go:build go1.21

package opt

import (
	"cmp"
	"sort"
)

// OptionSlice attaches the methods of [sort.Interface] to []Option[T], sorting in increasing order.
type OptionSlice[T cmp.Ordered] []Option[T]

// Len is the number of elements in the collection.
func (s OptionSlice[T]) Len() int { return len(s) }

// Less reports whether the element with index i must sort before the element with index j.
func (s OptionSlice[T]) Less(i, j int) bool { return Less(s[i], s[j]) }

// Swap swaps the elements with indexes i and j.
func (s OptionSlice[T]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (s OptionSlice[T]) Sort() { sort.Sort(s) }

// IsSorted reports whether data is sorted.
func (s OptionSlice[T]) IsSorted() bool { return sort.IsSorted(s) }

// Sort a slice of options in increasing order.
func Sort[T cmp.Ordered](s []Option[T]) { OptionSlice[T](s).Sort() }

// IsSorted reports whether the slice s is sorted in increasing order.
func IsSorted[T cmp.Ordered](s []Option[T]) bool { return OptionSlice[T](s).IsSorted() }

// Compare compares x and y and returns:
//
//	-1 if x < y;
//	0 if x == y;
//	+1 if x > y.
func Compare[T cmp.Ordered](x, y Option[T]) int {
	if x.val == nil && y.val == nil {
		return 0
	}

	if x.val != nil && y.val != nil {
		return cmp.Compare(*x.val, *y.val)
	}

	if x.val == nil {
		return -1
	}

	return 1
}

// Less reports whether x is less than y.
func Less[T cmp.Ordered](x, y Option[T]) bool {
	return Compare(x, y) < 0
}

// Equal compares two Option values for equality.
//
// It returns true if both options contain the same value, or if both options are None.
func Equal[T cmp.Ordered](x, y Option[T]) bool {
	return Compare(x, y) == 0
}
