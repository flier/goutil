//go:build go1.23

package swiss

import "iter"

// Iter iterates the elements of the Map, passing them to the callback.
// It guarantees that any key in the Map will be visited only once, and
// for un-mutated Maps, every key will be visited once. If the Map is
// Mutated during iteration, mutations will be reflected on return from
// Iter, but the set of keys visited by Iter is non-deterministic.
func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// take a consistent view of the table in case
		// we rehash during iteration
		ctrl, groups := m.ctrl, m.groups
		// pick a random starting group
		g := randIntN(groups.Len())
		for n := 0; n < groups.Len(); n++ {
			for s, c := range ctrl.Get(int(g)) {
				if c == empty || c == tombstone {
					continue
				}
				group := groups.Get(int(g))
				k, v := group.keys[s], group.values[s]
				if !yield(k, v) {
					return
				}
			}
			g++
			if g >= uint32(groups.Len()) {
				g = 0
			}
		}
	}
}
