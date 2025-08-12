//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/tuple"
)

// Partition takes an iterator of values and a predicate function f,
// and returns two slices containing the values for which the predicate returned true of false.
func Partition[T any](x iter.Seq[T], f func(T) bool) tuple.Tuple2[iter.Seq[T], iter.Seq[T]] {
	next, stop := iter.Pull(x)

	var lvals, rvals []T
	var done int

	l := func(yield func(T) bool) {
		for {
			if len(lvals) > 0 {
				if !yield(lvals[0]) {
					break
				}

				lvals = lvals[1:]
			} else {
				v, ok := next()
				if !ok {
					break
				}

				if f(v) {
					if !yield(v) {
						break
					}
				} else {
					rvals = append(rvals, v)
				}
			}
		}

		if done += 1; done == 2 {
			stop()
		}
	}

	r := func(yield func(T) bool) {
		for {
			if len(rvals) > 0 {
				if !yield(rvals[0]) {
					break
				}

				rvals = rvals[1:]
			} else {
				v, ok := next()
				if !ok {
					break
				}

				if f(v) {
					lvals = append(lvals, v)
				} else {
					if !yield(v) {
						break
					}
				}
			}
		}

		if done += 1; done == 2 {
			stop()
		}
	}

	return tuple.New2[iter.Seq[T], iter.Seq[T]](l, r)
}

// PartitionFunc takes an iterator of values and a predicate function f,
// and returns two slices containing the values for which the predicate returned true of false.
func PartitionFunc[T any](f func(T) bool) ReductionFunc[T, tuple.Tuple2[iter.Seq[T], iter.Seq[T]]] {
	return bind2(Partition, f)
}

// Partition2 takes an iterator of values and a predicate function f,
// and returns two slices containing the values for which the predicate returned true of false.
func Partition2[K, V any](x iter.Seq2[K, V], f func(K, V) bool) tuple.Tuple2[iter.Seq2[K, V], iter.Seq2[K, V]] {
	next, stop := iter.Pull2(x)

	var lvals, rvals []tuple.Tuple2[K, V]
	var done int

	l := func(yield func(K, V) bool) {
		for {
			if len(lvals) > 0 {
				if !yield(lvals[0].Unpack()) {
					break
				}

				lvals = lvals[1:]
			} else {
				k, v, ok := next()
				if !ok {
					break
				}

				if f(k, v) {
					if !yield(k, v) {
						break
					}
				} else {
					rvals = append(rvals, tuple.New2(k, v))
				}
			}
		}

		if done += 1; done == 2 {
			stop()
		}
	}

	r := func(yield func(K, V) bool) {
		for {
			if len(rvals) > 0 {
				if !yield(rvals[0].Unpack()) {
					break
				}

				rvals = rvals[1:]
			} else {
				k, v, ok := next()
				if !ok {
					break
				}

				if f(k, v) {
					lvals = append(lvals, tuple.New2(k, v))
				} else {
					if !yield(k, v) {
						break
					}
				}
			}
		}

		if done += 1; done == 2 {
			stop()
		}
	}

	return tuple.New2[iter.Seq2[K, V], iter.Seq2[K, V]](l, r)
}

// Partition2Func takes an iterator of values and a predicate function f,
// and returns two slices containing the values for which the predicate returned true of false.
func Partition2Func[K, V any](f func(K, V) bool) Reduction2Func[K, V, tuple.Tuple2[iter.Seq2[K, V], iter.Seq2[K, V]]] {
	return bind2(Partition2, f)
}
