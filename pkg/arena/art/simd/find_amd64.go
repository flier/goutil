//go:build amd64
// +build amd64

package simd

//go:noescape
func findKeyIndexAVX2(keys *[16]byte, key byte) int

//go:noescape
func findInsertPositionAVX2(keys *[16]byte, key byte) int

//go:noescape
func findNonZeroKeyIndexAVX2(keys *[256]byte) int

//go:noescape
func findLastNonZeroKeyIndexAVX2(keys *[256]byte) int

func FindKeyIndex(keys *[16]byte, n int, key byte) int {
	res := findKeyIndexAVX2(keys, key)

	// Check if the result is within the valid range
	if res >= n {
		return -1
	}

	return res
}

func FindInsertPosition(keys *[16]byte, n int, key byte) int {
	// Temporary: fallback to scalar for correctness until SIMD version is fixed
	return findInsertPositionScalar(keys, n, key)
}

func FindNonZeroKeyIndex(keys *[256]byte) int {
	return findNonZeroKeyIndexAVX2(keys)
}

func FindLastNonZeroKeyIndex(keys *[256]byte) int {
	return findLastNonZeroKeyIndexAVX2(keys)
}
