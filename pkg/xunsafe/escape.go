//go:build go1.23

package xunsafe

import "unsafe"

var (
	alwaysFalse bool
	sink        unsafe.Pointer //nolint:unused
)

// Escape escapes a pointer to the heap.
func Escape[P ~*E, E any](p P) P {
	if alwaysFalse {
		sink = unsafe.Pointer(p)
	}
	return p
}

// NoEscape hides a pointer from escape analysis, preventing it from
// escaping to the heap.
func NoEscape[P ~*E, E any](p P) P {
	//nolint:staticcheck // False positive: complains that p^0 does nothing.
	return P((AddrOf(p) ^ 0).AssertValid())
}
