//go:build !amd64
// +build !amd64

package simd

func FindKeyIndex(keys *[16]byte, n int, key byte) int {
	return findKeyIndexScalar(keys, n, key)
}

func FindInsertPosition(keys *[16]byte, n int, key byte) int {
	return findInsertPositionScalar(keys, n, key)
}
