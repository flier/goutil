package simd

// FindKeyIndex is a scalar fallback implementation for finding key index.
//
// This function is used by all architectures when SIMD is not available.
func findKeyIndexScalar(keys *[16]byte, n int, key byte) int {
	for i := 0; i < n; i++ {
		if keys[i] == key {
			return i
		}
	}

	return -1
}

// FindInsertPosition is a scalar fallback implementation for finding insert position.
//
// This function is used by all architectures when SIMD is not available.
func findInsertPositionScalar(keys *[16]byte, n int, key byte) int {
	for i := 0; i < n; i++ {
		if key < keys[i] {
			return i
		}
	}

	return n
}
