//go:build go1.23

package xunsafe

import (
	"unsafe"

	"github.com/flier/goutil/pkg/xunsafe/layout"
)

// Cast casts one pointer type to another.
func Cast[To, From any](p *From) *To {
	return (*To)(unsafe.Pointer(p))
}

// Add adds the given offset to p, scaled by the size of T.
func Add[P ~*E, E any, I Int](p P, n I) P {
	size := layout.Size[E]()
	return P(unsafe.Add(unsafe.Pointer(p), uintptr(size)*uintptr(n)))
}

// Sub computes the difference between two pointers, scaled by the size of T.
func Sub[P ~*E, E any](p1, p2 P) int {
	size := layout.Size[E]()
	return int(uintptr(unsafe.Pointer(p1))-uintptr(unsafe.Pointer(p2))) / size
}

// Load loads a value of the given type at the given index.
func Load[P ~*E, E any, I Int](p P, n I) E {
	return *Add(p, n)
}

// Store stores a value at the given index.
func Store[P ~*E, E any, I Int](p P, n I, v E) {
	*Add(p, n) = v
}

// StoreNoWB performs a store without generating any write barriers.
func StoreNoWB[P ~*E, E any](p *P, q P) {
	*Cast[uintptr](p) = uintptr(unsafe.Pointer(q))
}

// StoreNoWBUntyped performs a store without generating any write barriers.
func StoreNoWBUntyped[P ~unsafe.Pointer](p *P, q P) {
	*Cast[uintptr](p) = uintptr(q)
}

// Copy copies n elements from one pointer to the other.
func Copy[P ~*E, E any, I Int](dst, src P, n I) {
	copy(unsafe.Slice(dst, n), unsafe.Slice(src, n))
}

// Clear zeros n elements at p.
func Clear[P ~*E, E any, I Int](p P, n I) {
	clear(unsafe.Slice(p, n))
}
