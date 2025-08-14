//go:build go1.23

package xsync

import "iter"

// All returns an iterator over the values in this map, using [sync.Map.Range].
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.impl.Range(func(key, value any) bool {
			return yield(key.(K), value.(V)) //nolint:errcheck
		})
	}
}
