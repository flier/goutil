//go:build go1.23

package xiter

import "iter"

// Range returns a sequence of numbers from start (inclusive) to stop (exclusive).
func Range[T Integer](start, stop T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := start; i < stop; i++ {
			if !yield(i) {
				break
			}
		}
	}
}

// RangeFrom returns an infinite iterator of numbers starting from the given index start (inclusive).
func RangeFrom[T Integer](start T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for yield(start) {
			start += 1
		}
	}
}

// RangeTo returns a sequence of numbers from 0 (inclusive) to the given stop value n (exclusive).
func RangeTo[T Integer](stop T) iter.Seq[T] {
	return func(yield func(T) bool) {
		var i T
		for ; i < stop; i++ {
			if !yield(i) {
				break
			}
		}
	}
}
