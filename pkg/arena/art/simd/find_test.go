package simd

import (
	"testing"
)

// Benchmark data setup
func setupBenchmarkData() (*[16]byte, int) {
	keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30}
	return keys, 16
}

// Benchmark scalar implementations
func BenchmarkFindKeyIndexScalar(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 32)
		_ = findKeyIndexScalar(keys, n, key)
	}
}

func BenchmarkFindInsertPositionScalar(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 32)
		_ = findInsertPositionScalar(keys, n, key)
	}
}

// Benchmark SIMD implementations (only on supported architectures)
func BenchmarkFindKeyIndexSIMD(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 32)
		_ = FindKeyIndex(keys, n, key)
	}
}

func BenchmarkFindInsertPositionSIMD(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 32)
		_ = FindInsertPosition(keys, n, key)
	}
}

// Benchmark different key distributions
func BenchmarkFindKeyIndex_FirstElement(b *testing.B) {
	keys, n := setupBenchmarkData()
	key := byte(0) // First element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, n, key)
	}
}

func BenchmarkFindKeyIndex_LastElement(b *testing.B) {
	keys, n := setupBenchmarkData()
	key := byte(30) // Last element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, n, key)
	}
}

func BenchmarkFindKeyIndex_NotFound(b *testing.B) {
	keys, n := setupBenchmarkData()
	key := byte(31) // Not found

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, n, key)
	}
}

func BenchmarkFindInsertPosition_FirstPosition(b *testing.B) {
	keys, n := setupBenchmarkData()
	key := byte(1) // Insert at first position

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindInsertPosition(keys, n, key)
	}
}

func BenchmarkFindInsertPosition_MiddlePosition(b *testing.B) {
	keys, n := setupBenchmarkData()
	key := byte(15) // Insert at middle position

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindInsertPosition(keys, n, key)
	}
}

func BenchmarkFindInsertPosition_LastPosition(b *testing.B) {
	keys, n := setupBenchmarkData()
	key := byte(31) // Insert at last position

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindInsertPosition(keys, n, key)
	}
}

// Benchmark different array sizes
func BenchmarkFindKeyIndex_Size4(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6}
	n := 4

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 8)
		_ = FindKeyIndex(keys, n, key)
	}
}

func BenchmarkFindKeyIndex_Size8(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14}
	n := 8

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 16)
		_ = FindKeyIndex(keys, n, key)
	}
}

func BenchmarkFindKeyIndex_Size16(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 32)
		_ = FindKeyIndex(keys, n, key)
	}
}

// Benchmark scalar vs SIMD for different sizes
func BenchmarkFindKeyIndex_Scalar_Size4(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6}
	n := 4

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 8)
		_ = findKeyIndexScalar(keys, n, key)
	}
}

func BenchmarkFindKeyIndex_SIMD_Size4(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6}
	n := 4

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 8)
		_ = FindKeyIndex(keys, n, key)
	}
}

func BenchmarkFindKeyIndex_Scalar_Size8(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14}
	n := 8

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 16)
		_ = findKeyIndexScalar(keys, n, key)
	}
}

func BenchmarkFindKeyIndex_SIMD_Size8(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14}
	n := 8

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 8)
		_ = FindKeyIndex(keys, n, key)
	}
}

// Benchmark with realistic data patterns
func BenchmarkFindKeyIndex_Sequential(b *testing.B) {
	keys := &[16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	n := 16

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 16)
		_ = FindKeyIndex(keys, n, key)
	}
}

func BenchmarkFindKeyIndex_Sparse(b *testing.B) {
	keys := &[16]byte{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150}
	n := 16

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 160)
		_ = FindKeyIndex(keys, n, key)
	}
}

// Benchmark with different key distributions for insert position
func BenchmarkFindInsertPosition_Sequential(b *testing.B) {
	keys := &[16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	n := 16

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 16)
		_ = FindInsertPosition(keys, n, key)
	}
}

func BenchmarkFindInsertPosition_Sparse(b *testing.B) {
	keys := &[16]byte{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150}
	n := 16

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := byte(i % 160)
		_ = FindInsertPosition(keys, n, key)
	}
}
