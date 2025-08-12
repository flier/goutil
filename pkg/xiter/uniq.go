//go:build go1.23

package xiter

import (
	"iter"
)

// Uniq creates a stream that only emits elements if they are unique.
//
// Keep in mind that, in order to know if an element is unique or not,
// this function needs to store all unique values emitted by the stream.
// Therefore, if the stream is infinite, the number of elements stored will grow infinitely,
// never being garbage-collected.
func Uniq[T comparable](x iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		m := make(map[T]struct{})

		for v := range x {
			if _, exists := m[v]; exists {
				continue
			}

			if !yield(v) {
				break
			}

			m[v] = struct{}{}
		}
	}
}

// UniqByKey creates a stream that only emits elements if they are unique,
// by removing the elements for which function fun returned duplicate elements.
//
// The function fun maps every element to a value which is used to determine if two elements are duplicates.
//
// Keep in mind that, in order to know if an element is unique or not,
// this function needs to store all unique values emitted by the stream.
// Therefore, if the stream is infinite, the number of elements stored will grow infinitely,
// never being garbage-collected.
func UniqByKey[T any, B comparable](x iter.Seq[T], f func(T) B) iter.Seq[T] {
	return func(yield func(T) bool) {
		m := make(map[B]struct{})

		for v := range x {
			b := f(v)
			if _, exists := m[b]; exists {
				continue
			}

			if !yield(v) {
				break
			}

			m[b] = struct{}{}
		}
	}
}

// UniqByKeyFunc creates a stream that only emits elements if they are unique,
// by removing the elements for which function fun returned duplicate elements.
//
// The function fun maps every element to a value which is used to determine if two elements are duplicates.
//
// Keep in mind that, in order to know if an element is unique or not,
// this function needs to store all unique values emitted by the stream.
// Therefore, if the stream is infinite, the number of elements stored will grow infinitely,
// never being garbage-collected.
func UniqByKeyFunc[T any, B comparable](f func(T) B) MappingFunc[T, T] {
	return bind2(UniqByKey, f)
}

// UniqByKey2 creates a stream that only emits elements if they are unique,
// by removing the elements for which function fun returned duplicate elements.
//
// The function fun maps every element to a value which is used to determine if two elements are duplicates.
//
// Keep in mind that, in order to know if an element is unique or not,
// this function needs to store all unique values emitted by the stream.
// Therefore, if the stream is infinite, the number of elements stored will grow infinitely,
// never being garbage-collected.
func UniqByKey2[K, V any, B comparable](x iter.Seq2[K, V], f func(K, V) B) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m := make(map[B]struct{})

		for k, v := range x {
			b := f(k, v)
			if _, exists := m[b]; exists {
				continue
			}

			if !yield(k, v) {
				break
			}

			m[b] = struct{}{}
		}
	}
}

// UniqByKey2Func creates a stream that only emits elements if they are unique,
// by removing the elements for which function fun returned duplicate elements.
//
// The function fun maps every element to a value which is used to determine if two elements are duplicates.
//
// Keep in mind that, in order to know if an element is unique or not,
// this function needs to store all unique values emitted by the stream.
// Therefore, if the stream is infinite, the number of elements stored will grow infinitely,
// never being garbage-collected.
func UniqByKey2Func[K, V any, B comparable](f func(K, V) B) Reduction2Func[K, V, iter.Seq2[K, V]] {
	return bind2(UniqByKey2, f)
}
