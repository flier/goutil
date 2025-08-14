//go:build go1.23

package xsync

import (
	"iter"
	"sync"
)

// Set is a strongly-typed wrapper over sync.Map, used as a set.
type Set[K comparable] struct {
	impl sync.Map
}

// Load forwards to [sync.Map.Load].
func (s *Set[K]) Load(k K) bool {
	_, ok := s.impl.Load(k)
	return ok
}

// Store forwards to [sync.Map.Store].
func (s *Set[K]) Store(k K) {
	s.impl.Store(k, nil)
}

// All returns an iterator over the values in this set, using [sync.Map.Range].
func (s *Set[K]) All() iter.Seq[K] {
	return func(yield func(K) bool) {
		s.impl.Range(func(key, _ any) bool {
			return yield(key.(K)) //nolint:errcheck
		})
	}
}
