//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/tuple"
)

// Zip converts the arguments to iterators and zips them.
func Zip[K, V any](k iter.Seq[K], v iter.Seq[V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		kn, ks := iter.Pull(k)
		vn, vs := iter.Pull(v)

		defer ks()
		defer vs()

		for {
			k, ok := kn()
			if !ok {
				break
			}

			v, ok := vn()
			if !ok {
				break
			}

			if !yield(k, v) {
				break
			}
		}
	}
}

// ZipWith takes two iterators and a function, and returns a new iterator that
// applies the function to the corresponding elements of the input iterators and yields the results.
func ZipWith[K, V, B any](k iter.Seq[K], v iter.Seq[V], f func(K, V) B) iter.Seq[B] {
	return func(yield func(B) bool) {
		kn, ks := iter.Pull(k)
		vn, vs := iter.Pull(v)

		defer ks()
		defer vs()

		for {
			k, ok := kn()
			if !ok {
				break
			}

			v, ok := vn()
			if !ok {
				break
			}

			if !yield(f(k, v)) {
				break
			}
		}
	}
}

// ZipWithFunc takes two iterators and a function, and returns a new iterator that
// applies the function to the corresponding elements of the input iterators
// and yields the results.
func ZipWithFunc[K, V, B any](f func(K, V) B) func(iter.Seq[K], iter.Seq[V]) iter.Seq[B] {
	return bind3(ZipWith, f)
}

// Unzip converts an iterator of key-values into a pair of containers.
func Unzip[K, V any](x iter.Seq2[K, V]) tuple.Tuple2[iter.Seq[K], iter.Seq[V]] {
	next, stop := iter.Pull2(x)

	var keys []K
	var values []V
	var done int

	l := func(yield func(K) bool) {
		for {
			if len(keys) > 0 {
				if !yield(keys[0]) {
					break
				}

				keys = keys[1:]
			} else {
				k, v, ok := next()
				if !ok {
					break
				}

				if !yield(k) {
					break
				}

				values = append(values, v)
			}
		}

		if done += 1; done == 2 {
			stop()
		}
	}

	r := func(yield func(V) bool) {
		for {
			if len(values) > 0 {
				if !yield(values[0]) {
					break
				}

				values = values[1:]
			} else {
				k, v, ok := next()
				if !ok {
					break
				}

				if !yield(v) {
					break
				}

				keys = append(keys, k)
			}
		}

		if done += 1; done == 2 {
			stop()
		}
	}

	return tuple.New2[iter.Seq[K], iter.Seq[V]](l, r)
}
