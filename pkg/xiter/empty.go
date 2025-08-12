//go:build go1.23

package xiter

import "iter"

// Empty creates an iterator that yields nothing.
func Empty[T any]() iter.Seq[T] { return func(func(T) bool) {} }

// Empty2 creates an iterator that yields nothing.
func Empty2[K, V any]() iter.Seq2[K, V] { return func(func(K, V) bool) {} }
