//go:build go1.23

package xiter

import (
	"iter"

	"github.com/flier/goutil/pkg/xiter/inspect"
)

// Inspect and writes the given item to the standard output.
func Inspect[T any](x iter.Seq[T], opts ...inspect.Option) iter.Seq[T] {
	i := inspect.New(opts)

	return func(yield func(T) bool) {
		i.Start()
		defer i.Stop()

		for v := range x {
			if !yield(v) {
				break
			}

			i.Inspect(v)
		}
	}
}

// InspectFunc and writes the given item to the standard output.
func InspectFunc[T any](opts ...inspect.Option) MappingFunc[T, T] {
	return bind2rest(Inspect[T], opts)
}

// Inspect2 inspects and writes the given item to the standard output.
func Inspect2[K, V any](x iter.Seq2[K, V], opts ...inspect.Option) iter.Seq2[K, V] {
	i := inspect.New(opts)

	return func(yield func(K, V) bool) {
		i.Start()
		defer i.Stop()

		for k, v := range x {
			if !yield(k, v) {
				break
			}

			i.Inspect2(k, v)
		}
	}
}

// Inspect2Func inspects and writes the given item to the standard output.
func Inspect2Func[K, V any](opts ...inspect.Option) MappingValueFunc[K, V, V] {
	return bind2rest(Inspect2[K, V], opts)
}
