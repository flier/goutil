//go:build go1.23

package xiter

import (
	"iter"
)

// Map takes a function and creates an iterator which calls that function f on each element.
func Map[T, O any](x iter.Seq[T], f func(T) O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for v := range x {
			if !yield(f(v)) {
				break
			}
		}
	}
}

// MapFunc takes a function and creates an iterator which calls that function f on each element.
func MapFunc[T, O any](f func(T) O) MappingFunc[T, O] {
	return bind2(Map, f)
}

// MapKeyValue takes a function and creates an iterator which calls that function f on each key-value pair.
func MapKeyValue[K, V, O, P any](x iter.Seq2[K, V], f func(K, V) (O, P)) iter.Seq2[O, P] {
	return func(yield func(O, P) bool) {
		for k, v := range x {
			if !yield(f(k, v)) {
				break
			}
		}
	}
}

// Map2Func takes a function and creates an iterator which calls that function f on each key-value pair.
func Map2Func[K, V, O, P any](f func(K, V) (O, P)) MappingKeyValueFunc[K, V, O, P] {
	return bind2(MapKeyValue, f)
}

// MapKey takes a function and creates an iterator which calls that function f on each key-value pair.
func MapKey[K, V, O any](x iter.Seq2[K, V], f func(K, V) O) iter.Seq2[O, V] {
	return func(yield func(O, V) bool) {
		for k, v := range x {
			if !yield(f(k, v), v) {
				break
			}
		}
	}
}

// MapKeyFunc takes a function and creates an iterator which calls that function f on each key-value pair.
func MapKeyFunc[K, V, O any](f func(K, V) O) MappingKeyFunc[K, V, O] {
	return bind2(MapKey, f)
}

// MapValue takes a function and creates an iterator which calls that function f on each key-value pair.
func MapValue[K, V, O any](x iter.Seq2[K, V], f func(K, V) O) iter.Seq2[K, O] {
	return func(yield func(K, O) bool) {
		for k, v := range x {
			if !yield(k, f(k, v)) {
				break
			}
		}
	}
}

// MapValueFunc takes a function and creates an iterator which calls that function f on each key-value pair.
func MapValueFunc[K, V, O any](f func(K, V) O) MappingValueFunc[K, V, O] {
	return bind2(MapValue, f)
}

// FlatMap creates an iterator that works like Map, but flattens nested iterator.
func FlatMap[T, O any](x iter.Seq[T], f func(T) iter.Seq[O]) iter.Seq[O] {
	return func(yield func(O) bool) {
		for v := range x {
			for i := range f(v) {
				if !yield(i) {
					break
				}
			}
		}
	}
}

// FlatMapFunc creates an iterator that works like Map, but flattens nested iterator.
func FlatMapFunc[T, O any](f func(T) iter.Seq[O]) MappingFunc[T, O] {
	return bind2(FlatMap, f)
}

// FlatMap2 creates an iterator that works like Map2, but flattens nested iterator.
func FlatMap2[K, V, O any](x iter.Seq2[K, V], f func(K, V) iter.Seq2[K, O]) iter.Seq2[K, O] {
	return func(yield func(K, O) bool) {
		for k, v := range x {
			for k, o := range f(k, v) {
				if !yield(k, o) {
					break
				}
			}
		}
	}
}

// FlatMap2Func creates an iterator that works like Map2, but flattens nested iterator.
func FlatMap2Func[K, V, O any](f func(K, V) iter.Seq2[K, O]) MappingValueFunc[K, V, O] {
	return bind2(FlatMap2, f)
}

// MapWhile creates an iterator that both yields elements based on a predicate f and maps.
func MapWhile[T, O any](x iter.Seq[T], f func(T) (O, bool)) iter.Seq[O] {
	return func(yield func(O) bool) {
		for v := range x {
			o, ok := f(v)
			if !ok {
				continue
			}

			if !yield(o) {
				break
			}
		}
	}
}

// MapWhileFunc creates an iterator that both yields elements based on a predicate f and maps.
func MapWhileFunc[T, O any](f func(T) (O, bool)) MappingFunc[T, O] {
	return func(x iter.Seq[T]) iter.Seq[O] {
		return func(yield func(O) bool) {
			for v := range x {
				o, ok := f(v)
				if !ok {
					continue
				}

				if !yield(o) {
					break
				}
			}
		}
	}
}

// MapWhile2 creates an iterator that both yields key-values based on a predicate f and maps.
func MapWhile2[K, V, O any](x iter.Seq2[K, V], f func(K, V) (O, bool)) iter.Seq2[K, O] {
	return func(yield func(K, O) bool) {
		for k, v := range x {
			o, ok := f(k, v)
			if !ok {
				continue
			}

			if !yield(k, o) {
				break
			}
		}
	}
}

// MapWhile2Func creates an iterator that both yields key-values based on a predicate f and maps.
func MapWhile2Func[K, V, O any](f func(K, V) (O, bool)) MappingValueFunc[K, V, O] {
	return bind2(MapWhile2, f)
}
