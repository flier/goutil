//go:build go1.23

package xsync

import (
	"iter"
	"sync"
)

// Map is a strongly-typed wrapper over sync.Map.
type Map[K comparable, V any] struct {
	impl sync.Map
}

// Load forwards to [sync.Map.Load].
func (m *Map[K, V]) Load(k K) (V, bool) {
	v, ok := m.impl.Load(k)
	if !ok {
		var z V
		return z, ok
	}

	return v.(V), ok //nolint:errcheck
}

// Store forwards to [sync.Map.Store].
func (m *Map[K, V]) Store(k K, v V) {
	m.impl.Store(k, v)
}

// LoadOrStore loads a value if its present, or constructs it with make and
// inserts it.
//
// There is a possibility that make is called, but the return value is not
// inserted.
func (m *Map[K, V]) LoadOrStore(k K, make func() V) (actual V, loaded bool) {
	v, ok := m.Load(k)
	if ok {
		return v, true
	}
	w, ok := m.impl.LoadOrStore(k, make())
	return w.(V), ok //nolint:errcheck
}

// All returns an iterator over the values in this map, using [sync.Map.Range].
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.impl.Range(func(key, value any) bool {
			return yield(key.(K), value.(V)) //nolint:errcheck
		})
	}
}
