//go:build go1.23

package xiter

import (
	"iter"
)

// Chunks splits the input sequence x into consecutive chunks of size n.
//
// Each chunk is yielded as a slice of T. The last chunk may contain fewer than n elements
// if the total number of elements in x is not a multiple of n.
// The function returns an iter.Seq of slices, where each slice contains up to n elements.
// If n is less than or equal to zero, no chunks will be yielded.
//
// Example usage:
//
//	seq := slices.Value([]int{1, 2, 3, 4, 5})
//
//	for chunk := range Chunks(seq, 2) {
//	    fmt.Println(chunk) // Output: [1 2], [3 4], [5]
//	}
func Chunks[T any](x iter.Seq[T], n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if n <= 0 {
			return
		}

		chunk := make([]T, 0, n)

		for v := range x {
			if len(chunk) < n {
				chunk = append(chunk, v)
			}

			if len(chunk) == n {
				if !yield(chunk) {
					return
				}

				chunk = make([]T, 0, n)
			}
		}

		if len(chunk) > 0 {
			yield(chunk)
		}
	}
}

// ChunkBy groups consecutive elements from the input sequence x into slices,
// starting a new chunk each time the predicate function f returns true for an element.
//
// It returns a sequence of slices, where each slice contains a chunk of elements.
// The chunk boundary is determined by the predicate function `f`.
func ChunkBy[T any](x iter.Seq[T], f func(T) bool) iter.Seq[[]T] {
	return ChunkByKey(x, f)
}

// ChunkByKey groups elements from the input sequence x into contiguous chunks,
// where each chunk contains consecutive elements that share the same key as determined
// by the function f. The function returns a sequence of slices, each representing a chunk.
//
// Type Parameters:
//
//	T - the type of elements in the input sequence.
//	B - the type of the key, which must be comparable.
//
// Parameters:
//
//	x - an input sequence of elements of type T.
//	f - a function that maps each element to a key of type B.
//
// Returns:
//
//	An iter.Seq[T], where each slice contains consecutive elements
//	from the input sequence that have the same key.
func ChunkByKey[T any, B comparable](x iter.Seq[T], f func(T) B) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		var chunk []T
		var cur *B

		for v := range x {
			b := f(v)

			if cur == nil {
				cur = &b
			} else if *cur != b {
				if !yield(chunk) {
					return
				}

				cur = &b
				chunk = nil
			}

			chunk = append(chunk, v)
		}

		if len(chunk) > 0 {
			yield(chunk)
		}
	}
}

// ChunksFunc returns a MapFunc that splits an input sequence into consecutive chunks of size n.
// Each chunk is represented as a slice of type T. If the input sequence length is not a multiple of n,
// the final chunk will contain the remaining elements.
//
// Example usage:
//
//	chunks := ChunksFunc[int](3)
//	seq := slices.Values([]int{1, 2, 3, 4, 5, 6, 7})
//	for chunk := range chunks(seq) {
//	    fmt.Println(chunk) // Output: [1 2 3], [4 5 6], [7]
//	}
//
// Parameters:
//
//	n int: The size of each chunk.
//
// Returns:
//
//	MapFunc[T, []T]: A function that maps a sequence of T to a sequence of []T (chunks).
func ChunksFunc[T any](n int) MappingFunc[T, []T] {
	return bind2(Chunks[T], n)
}

// ChunkByFunc returns an iterator that groups the elements of the input iterator x into
// chunks based on the result of the provided function f.
func ChunkByFunc[T any](f func(T) bool) MappingFunc[T, []T] {
	return ChunkByKeyFunc(f)
}

// ChunkByKeyFunc returns an iterator that groups the elements of the input iterator x into
// chunks based on the result of the provided function f.
func ChunkByKeyFunc[T any, B comparable](f func(T) B) MappingFunc[T, []T] {
	return bind2(ChunkByKey, f)
}
