//go:build go1.23

package xiter

import "iter"

// Scan applies the provided function f to each element in the input iterator x,
// yielding a new iterator of the results of applying f.
func Scan[C, T, B any](x iter.Seq[T], ctx C, f func(C, T) (B, bool)) iter.Seq[B] {
	return func(yield func(B) bool) {
		for v := range x {
			b, ok := f(ctx, v)
			if !ok {
				continue
			}

			if !yield(b) {
				break
			}
		}
	}
}

// ScanFunc applies the provided function f to each element in the input iterator x,
// yielding a new iterator of the results of applying f.
func ScanFunc[C, T, B any](ctx C, f func(C, T) (B, bool)) MappingFunc[T, B] {
	return bind23(Scan, ctx, f)
}

// Scan2 applies the provided function f to each key-value in the input iterator x,
// yielding a new iterator of the results of applying f.
func Scan2[C, K, V, B any](x iter.Seq2[K, V], ctx C, f func(C, K, V) (B, bool)) iter.Seq2[K, B] {
	return func(yield func(K, B) bool) {
		for k, v := range x {
			b, ok := f(ctx, k, v)
			if !ok {
				continue
			}

			if !yield(k, b) {
				break
			}
		}
	}
}

// Scan2Func applies the provided function f to each key-value in the input iterator x,
// yielding a new iterator of the results of applying f.
func Scan2Func[C, K, V, B any](ctx C, f func(C, K, V) (B, bool)) MappingValueFunc[K, V, B] {
	return bind23(Scan2, ctx, f)
}
