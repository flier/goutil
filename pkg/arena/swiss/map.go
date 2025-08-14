package swiss

import (
	"github.com/dolthub/maphash"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

const (
	maxLoadFactor = float32(maxAvgGroupLoad) / float32(groupSize)
)

// Map is an open-addressing hash map
// based on Abseil's flat_hash_map.
type Map[K comparable, V any] struct {
	arena    *arena.Arena
	ctrl     slice.Slice[metadata]
	groups   slice.Slice[group[K, V]]
	hash     maphash.Hasher[K]
	resident uint32
	dead     uint32
	limit    uint32
}

// metadata is the h2 metadata array for a group.
// find operations first probe the controls bytes
// to filter candidates before matching keys
type metadata [groupSize]int8

// group is a group of 16 key-value pairs
type group[K comparable, V any] struct {
	keys   [groupSize]K
	values [groupSize]V
}

const (
	h1Mask    uint64 = 0xffff_ffff_ffff_ff80
	h2Mask    uint64 = 0x0000_0000_0000_007f
	empty     int8   = -128 // 0b1000_0000
	tombstone int8   = -2   // 0b1111_1110
)

// h1 is a 57 bit hash prefix
type h1 uint64

// h2 is a 7 bit hash suffix
type h2 int8

// NewMap constructs a Map.
func NewMap[K comparable, V any](a *arena.Arena, sz uint32) (m *Map[K, V]) {
	groups := numGroups(sz)

	m = arena.New(a, Map[K, V]{
		arena:  a,
		ctrl:   slice.Make[metadata](a, int(groups)),
		groups: slice.Make[group[K, V]](a, int(groups)),
		hash:   maphash.NewHasher[K](),
		limit:  groups * maxAvgGroupLoad,
	})

	for i := 0; i < m.ctrl.Len(); i++ {
		m.ctrl.Store(i, newEmptyMetadata())
	}

	return
}

// Has returns true if |key| is present in |m|.
func (m *Map[K, V]) Has(key K) (ok bool) {
	hi, lo := splitHash(m.hash.Hash(key))
	g := probeStart(hi, m.groups.Len())
	for { // inlined find loop
		matches := metaMatchH2(m.ctrl.Get(int(g)), lo)
		for matches != 0 {
			s := nextMatch(&matches)
			if key == m.groups.Get(int(g)).keys[s] {
				ok = true
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty(m.ctrl.Get(int(g)))
		if matches != 0 {
			ok = false
			return
		}
		g += 1 // linear probing
		if g >= uint32(m.groups.Len()) {
			g = 0
		}
	}
}

// Get returns the |value| mapped by |key| if one exists.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	hi, lo := splitHash(m.hash.Hash(key))
	g := probeStart(hi, m.groups.Len())
	for { // inlined find loop
		matches := metaMatchH2(m.ctrl.Get(int(g)), lo)
		for matches != 0 {
			s := nextMatch(&matches)
			if key == m.groups.Get(int(g)).keys[s] {
				value, ok = m.groups.Get(int(g)).values[s], true
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty(m.ctrl.Get(int(g)))
		if matches != 0 {
			ok = false
			return
		}
		g += 1 // linear probing
		if g >= uint32(m.groups.Len()) {
			g = 0
		}
	}
}

// Put attempts to insert |key| and |value|
func (m *Map[K, V]) Put(key K, value V) {
	if m.resident >= m.limit {
		m.rehash(m.nextSize())
	}
	hi, lo := splitHash(m.hash.Hash(key))
	g := probeStart(hi, m.groups.Len())
	for { // inlined find loop
		matches := metaMatchH2(m.ctrl.Get(int(g)), lo)
		for matches != 0 {
			s := nextMatch(&matches)
			if key == m.groups.Get(int(g)).keys[s] { // update
				m.groups.Get(int(g)).keys[s] = key
				m.groups.Get(int(g)).values[s] = value
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty(m.ctrl.Get(int(g)))
		if matches != 0 { // insert
			s := nextMatch(&matches)
			m.groups.Get(int(g)).keys[s] = key
			m.groups.Get(int(g)).values[s] = value
			m.ctrl.Get(int(g))[s] = int8(lo)
			m.resident++
			return
		}
		g += 1 // linear probing
		if g >= uint32(m.groups.Len()) {
			g = 0
		}
	}
}

// Delete attempts to remove |key|, returns true successful.
func (m *Map[K, V]) Delete(key K) (ok bool) {
	hi, lo := splitHash(m.hash.Hash(key))
	g := probeStart(hi, m.groups.Len())
	for {
		matches := metaMatchH2(m.ctrl.Get(int(g)), lo)
		for matches != 0 {
			s := nextMatch(&matches)
			if key == m.groups.Get(int(g)).keys[s] {
				ok = true
				// optimization: if |m.ctrl.Get(int(g))| contains any empty
				// metadata bytes, we can physically delete |key|
				// rather than placing a tombstone.
				// The observation is that any probes into group |g|
				// would already be terminated by the existing empty
				// slot, and therefore reclaiming slot |s| will not
				// cause premature termination of probes into |g|.
				if metaMatchEmpty(m.ctrl.Get(int(g))) != 0 {
					m.ctrl.Get(int(g))[s] = empty
					m.resident--
				} else {
					m.ctrl.Get(int(g))[s] = tombstone
					m.dead++
				}
				var k K
				var v V
				m.groups.Get(int(g)).keys[s] = k
				m.groups.Get(int(g)).values[s] = v
				return
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty(m.ctrl.Get(int(g)))
		if matches != 0 { // |key| absent
			ok = false
			return
		}
		g += 1 // linear probing
		if g >= uint32(m.groups.Len()) {
			g = 0
		}
	}
}

// Clear removes all elements from the Map.
func (m *Map[K, V]) Clear() {
	for i, c := range m.ctrl.Raw() {
		for j := range c {
			m.ctrl.Get(i)[j] = empty
		}
	}
	var k K
	var v V
	for i := 0; i < m.groups.Len(); i++ {
		g := m.groups.Get(i)
		for i := range g.keys {
			g.keys[i] = k
			g.values[i] = v
		}
	}
	m.resident, m.dead = 0, 0
}

// Count returns the number of elements in the Map.
func (m *Map[K, V]) Count() int {
	return int(m.resident - m.dead)
}

// Capacity returns the number of additional elements
// the can be added to the Map before resizing.
func (m *Map[K, V]) Capacity() int {
	return int(m.limit - m.resident)
}

// find returns the location of |key| if present, or its insertion location if absent.
// for performance, find is manually inlined into public methods.
func (m *Map[K, V]) find(key K, hi h1, lo h2) (g, s uint32, ok bool) {
	g = probeStart(hi, m.groups.Len())
	for {
		matches := metaMatchH2(m.ctrl.Get(int(g)), lo)
		for matches != 0 {
			s = nextMatch(&matches)
			if key == m.groups.Get(int(g)).keys[s] {
				return g, s, true
			}
		}
		// |key| is not in group |g|,
		// stop probing if we see an empty slot
		matches = metaMatchEmpty(m.ctrl.Get(int(g)))
		if matches != 0 {
			s = nextMatch(&matches)
			return g, s, false
		}
		g += 1 // linear probing
		if g >= uint32(m.groups.Len()) {
			g = 0
		}
	}
}

func (m *Map[K, V]) nextSize() (n uint32) {
	n = uint32(m.groups.Len()) * 2
	if m.dead >= (m.resident / 2) {
		n = uint32(m.groups.Len())
	}
	return
}

func (m *Map[K, V]) rehash(n uint32) {
	groups, ctrl := m.groups, m.ctrl
	m.groups = slice.Make[group[K, V]](m.arena, int(n))
	m.ctrl = slice.Make[metadata](m.arena, int(n))
	for i := 0; i < m.ctrl.Len(); i++ {
		m.ctrl.Store(i, newEmptyMetadata())
	}
	m.hash = maphash.NewSeed(m.hash)
	m.limit = n * maxAvgGroupLoad
	m.resident, m.dead = 0, 0
	for g := 0; g < ctrl.Len(); g++ {
		for s := range ctrl.Get(int(g)) {
			c := ctrl.Get(g)[s]
			if c == empty || c == tombstone {
				continue
			}
			m.Put(groups.Get(g).keys[s], groups.Get(g).values[s])
		}
	}
}

func (m *Map[K, V]) loadFactor() float32 {
	slots := float32(m.groups.Len() * groupSize)
	return float32(m.resident-m.dead) / slots
}

// numGroups returns the minimum number of groups needed to store |n| elems.
func numGroups(n uint32) (groups uint32) {
	groups = (n + maxAvgGroupLoad - 1) / maxAvgGroupLoad
	if groups == 0 {
		groups = 1
	}
	return
}

func newEmptyMetadata() (meta metadata) {
	for i := range meta {
		meta[i] = empty
	}
	return
}

func splitHash(h uint64) (h1, h2) {
	return h1((h & h1Mask) >> 7), h2(h & h2Mask)
}

func probeStart(hi h1, groups int) uint32 {
	return fastModN(uint32(hi), uint32(groups))
}

// lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
func fastModN(x, n uint32) uint32 {
	return uint32((uint64(x) * uint64(n)) >> 32)
}
