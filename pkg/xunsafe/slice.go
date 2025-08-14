package xunsafe

import (
	"math"
	"unsafe"

	"github.com/flier/goutil/pkg/xunsafe/layout"
)

// BoundsCheck emulates a bounds check on a slice with the given index and
// length.
func BoundsCheck(n, len int) {
	dummy := unsafe.Slice(&struct{}{}, len&^math.MinInt)
	_ = dummy[n]
}

// Bytes converts a pointer into a slice of its contents.
func Bytes[P ~*E, E any](p P) []byte {
	size := layout.Size[E]()
	return unsafe.Slice(Cast[byte](p), size)
}

// LoadSlice loads a slice without performing a bounds check.
func LoadSlice[S ~[]E, E any, I Int](s S, n I) E {
	return Load(unsafe.SliceData(s), n)
}

// SliceToString converts a slice into a string, multiplying the slice length
// as appropriate.
func SliceToString[S ~[]E, E any](s S) string {
	size := layout.Size[E]()
	str := struct {
		ptr *E
		len int
	}{unsafe.SliceData(s), len(s) * size}
	return BitCast[string](str)
}

// StringToSlice converts a string into a slice, multiplying the slice length
// as appropriate.
func StringToSlice[S ~[]E, E any](s string) S {
	size := layout.Size[E]()
	return unsafe.Slice(Cast[E](unsafe.StringData(s)), len(s)/size)
}
