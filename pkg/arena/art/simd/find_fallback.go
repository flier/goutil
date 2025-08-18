//go:build !amd64
// +build !amd64

package simd

func FindKeyIndex(keys *[16]byte, n int, key byte) int {
	return findKeyIndexScalar(keys, n, key)
}

func FindInsertPosition(keys *[16]byte, n int, key byte) int {
	return findInsertPositionScalar(keys, n, key)
}

func FindNonZeroKeyIndex(keys *[256]byte) int {
	return findNonZeroKeyIndexScalar(keys)
}

func FindLastNonZeroKeyIndex(keys *[256]byte) int {
	return findLastNonZeroKeyIndexScalar(keys)
}
