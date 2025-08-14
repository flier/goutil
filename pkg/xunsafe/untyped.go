//go:build go1.23

package xunsafe

import "unsafe"

// ByteAdd adds the given offset to p, without scaling.
//
// It also throws in a cast for free.
func ByteAdd[T any, P ~*E, E any, I Int](p P, n I) *T {
	return (*T)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(n)))
}

// Sub computes the difference between two pointers, without scaling.
func ByteSub[P1 ~*E1, P2 ~*E2, E1, E2 any](p1 P1, p2 P2) int {
	return int(uintptr(unsafe.Pointer(p1)) - uintptr(unsafe.Pointer(p2)))
}

// ByteLoad loads a value of the given type at the given byte offset.
func ByteLoad[T any, P ~*E, E any, I Int](p P, n I) T {
	return *ByteAdd[T](p, n)
}

// ByteLoad stores a value of the given type at the given byte offset.
func ByteStore[T any, P ~*E, E any, I Int](p P, n I, v T) {
	*ByteAdd[T](p, n) = v
}
