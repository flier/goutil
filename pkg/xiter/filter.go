//go:build go1.23

package xiter

import "iter"

// Filter creates an iterator which uses a function f to determine if an element should be yielded.
func Filter[T any](x iter.Seq[T], f func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range x {
			if !f(v) {
				continue
			}

			if !yield(v) {
				break
			}
		}
	}
}

// FilterFunc creates an iterator which uses a function f to determine if an element should be yielded.
func FilterFunc[T any](f func(T) bool) MappingFunc[T, T] {
	return bind2(Filter, f)
}

// Filter2 creates an iterator which uses a function f to determine if a key-value should be yielded.
func Filter2[K, V any](x iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range x {
			if !f(k, v) {
				continue
			}

			if !yield(k, v) {
				break
			}
		}
	}
}

// Filter2Func creates an iterator which uses a function f to determine if a key-value should be yielded.
func Filter2Func[K, V any](f func(K, V) bool) MappingValueFunc[K, V, V] {
	return bind2(Filter2, f)
}

// FilterMap creates an iterator that both filters and maps.
func FilterMap[T, B any](x iter.Seq[T], f func(T) (B, bool)) iter.Seq[B] {
	return func(yield func(B) bool) {
		for v := range x {
			b, ok := f(v)
			if !ok {
				continue
			}

			if !yield(b) {
				break
			}
		}
	}
}

// FilterMapFunc creates an iterator that both filters and maps.
func FilterMapFunc[T, B any](f func(T) (B, bool)) MappingFunc[T, B] {
	return bind2(FilterMap, f)
}

// FilterMap2 creates an iterator that both filters and maps.
func FilterMap2[K, V, B any](x iter.Seq2[K, V], f func(K, V) (B, bool)) iter.Seq2[K, B] {
	return func(yield func(K, B) bool) {
		for k, v := range x {
			b, ok := f(k, v)

			if !ok {
				continue
			}

			if !yield(k, b) {
				break
			}
		}
	}
}

// FilterMap2Func creates an iterator that both filters and maps.
func FilterMap2Func[K, V, B any](f func(K, V) (B, bool)) MappingValueFunc[K, V, B] {
	return bind2(FilterMap2, f)
}
