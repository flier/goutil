//go:build go1.23

package xiter

import "iter"

// Pipeline applies the given Mapper functors to the input sequence in order.
func Pipeline[T any](s iter.Seq[T], x ...Mapper[T, T]) iter.Seq[T] {
	for _, m := range x {
		s = m.Map(s)
	}

	return s
}

// PipelineFunc applies the given Mapper functors to the input sequence in order.
func PipelineFunc[T any](x ...Mapper[T, T]) MappingFunc[T, T] {
	return bind2rest(Pipeline[T], x)
}

// Pipeline2 applies the given Mapper functors to the input sequence in order.
func Pipeline2[K, V any](s iter.Seq2[K, V], x ...ValueMapper[K, V, V]) iter.Seq2[K, V] {
	for _, m := range x {
		s = m.MapValue(s)
	}

	return s
}

// Pipeline2Func applies the given Mapper functors to the input sequence in order.
func Pipeline2Func[K, V any](x ...ValueMapper[K, V, V]) MappingValueFunc[K, V, V] {
	return bind2rest(Pipeline2[K, V], x)
}
