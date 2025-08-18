//go:build amd64
// +build amd64

package simd

import "math/bits"

//go:noescape
func findKeyIndexAVX2(keys *[16]byte, key byte) int

//go:noescape
func findInsertPositionAVX2(keys *[16]byte, key byte) int

func FindKeyIndex(keys *[16]byte, n int, key byte) int {
	res := findKeyIndexAVX2(keys, key)
	res &= (1 << n) - 1

	return bits.TrailingZeros(uint(res))
}

func FindInsertPosition(keys *[16]byte, n int, key byte) int {
	res := findInsertPositionAVX2(keys, key)
	res &= (1 << n) - 1

	return bits.TrailingZeros(uint(res))
}
