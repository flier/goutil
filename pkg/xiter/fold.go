//go:build go1.23

package xiter

import (
	"iter"
)

// Fold folds every element into an accumulator by applying an operation f, returning the final result.
func Fold[T, B any](x iter.Seq[T], init B, f func(B, T) B) B {
	acc := init

	for v := range x {
		acc = f(acc, v)
	}

	return acc
}

// FoldFunc folds every element into an accumulator by applying an operation f, returning the final result.
func FoldFunc[T, B any](init B, f func(B, T) B) ReductionFunc[T, B] {
	return bind23(Fold, init, f)
}

// Fold2 folds every key-value into an accumulator by applying an operation f, returning the final result.
func Fold2[K, V, B any](x iter.Seq2[K, V], init B, f func(B, K, V) B) B {
	acc := init

	for k, v := range x {
		acc = f(acc, k, v)
	}

	return acc
}

// Fold2Func folds every key-value into an accumulator by applying an operation f, returning the final result.
func Fold2Func[K, V, B any](init B, f func(B, K, V) B) Reduction2Func[K, V, B] {
	return bind23(Fold2, init, f)
}
