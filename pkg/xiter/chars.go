//go:build go1.23

package xiter

import (
	"iter"
	"unicode/utf8"
)

// Chars returns an iterator sequence over the runes in the given byte slice.
//
// It decodes UTF-8 encoded runes from the slice and yields each rune to the provided function.
// Iteration stops if the yield function returns false or the end of the slice is reached.
func Chars(b []byte) iter.Seq[rune] {
	return func(yield func(rune) bool) {
		for len(b) > 0 {
			r, size := utf8.DecodeRune(b)
			if !yield(r) {
				return
			}

			b = b[size:]
		}
	}
}
