//go:build amd64

package simd

import (
	"strconv"
	"testing"
)

// Test data setup functions
func setupBenchmarkData() (*[16]byte, int) {
	keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30}
	return keys, 16
}

func BenchmarkFindKeyIndex(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()

	b.Run("findKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexScalar(keys, n, byte(i%32))
		}
	})

	b.Run("findKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexAVX2(keys, byte(i%32))
		}
	})
}

func BenchmarkFindInsertPosition(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()

	b.Run("findInsertPositionScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionScalar(keys, n, byte(i%32))
		}
	})

	b.Run("findInsertPositionAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionAVX2(keys, byte(i%32))
		}
	})
}

// Benchmark different key distributions
func BenchmarkFindKeyIndex_FirstElement(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()

	b.Run("findKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexScalar(keys, n, 0) // First element
		}
	})

	b.Run("findKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexAVX2(keys, byte(i%32))
		}
	})
}

func BenchmarkFindKeyIndex_LastElement(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()

	b.Run("findKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexScalar(keys, n, 30) // Last element
		}
	})

	b.Run("findKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexAVX2(keys, byte(i%32))
		}
	})
}

func BenchmarkFindKeyIndex_NotFound(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()

	b.Run("findKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexScalar(keys, n, 31) // Not found
		}
	})

	b.Run("findKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexAVX2(keys, byte(i%32))
		}
	})
}

func BenchmarkFindInsertPosition_FirstPosition(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()

	b.Run("findInsertPositionScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionScalar(keys, n, 1) // Insert at first position
		}
	})

	b.Run("findInsertPositionAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionAVX2(keys, byte(i%32))
		}
	})
}

func BenchmarkFindInsertPosition_MiddlePosition(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()

	b.Run("findInsertPositionScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionScalar(keys, n, 15) // Insert at middle position
		}
	})

	b.Run("findInsertPositionAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionAVX2(keys, byte(i%32))
		}
	})
}

func BenchmarkFindInsertPosition_LastPosition(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()

	b.Run("findInsertPositionScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionScalar(keys, n, 31) // Insert at last position
		}
	})

	b.Run("findInsertPositionAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionAVX2(keys, byte(i%32))
		}
	})
}

// Benchmark different array sizes
func BenchmarkFindKeyIndex_Size4(b *testing.B) {
	cases := map[int][16]byte{
		4:  {0, 2, 4, 6},
		8:  {0, 2, 4, 6, 8, 10, 12, 14},
		16: {0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30},
	}

	b.ResetTimer()

	for size, keys := range cases {
		b.Run("findKeyIndexScalar_"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = findKeyIndexScalar(&keys, size, byte(i%8))
			}
		})

		b.Run("findKeyIndexAVX2_"+strconv.Itoa(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = findKeyIndexAVX2(&keys, byte(i%8))
			}
		})
	}
}

// Benchmark with realistic data patterns
func BenchmarkFindKeyIndex_Sequential(b *testing.B) {
	keys := &[16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	b.ResetTimer()

	b.Run("findKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexScalar(keys, 16, byte(i%16))
		}
	})

	b.Run("findKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexAVX2(keys, byte(i%16))
		}
	})
}

func BenchmarkFindKeyIndex_Sparse(b *testing.B) {
	keys := &[16]byte{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150}

	b.ResetTimer()

	b.Run("findKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexScalar(keys, 16, byte(i%160))
		}
	})

	b.Run("findKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findKeyIndexAVX2(keys, byte(i%160))
		}
	})
}

// Benchmark with different key distributions for insert position
func BenchmarkFindInsertPosition_Sequential(b *testing.B) {
	keys := &[16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	b.ResetTimer()

	b.Run("findInsertPositionScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionScalar(keys, 16, byte(i%16))
		}
	})

	b.Run("findInsertPositionAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionAVX2(keys, byte(i%16))
		}
	})
}

func BenchmarkFindInsertPosition_Sparse(b *testing.B) {
	keys := &[16]byte{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150}

	b.ResetTimer()

	b.Run("findInsertPositionScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionScalar(keys, 16, byte(i%160))
		}
	})

	b.Run("findInsertPositionAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findInsertPositionAVX2(keys, byte(i%160))
		}
	})
}

// Benchmark for FindNonZeroKeyIndex
func BenchmarkFindNonZeroKeyIndex(b *testing.B) {
	keys := &[256]byte{}
	// Set some non-zero elements for testing
	keys[0] = 1
	keys[128] = 42
	keys[255] = 100

	b.ResetTimer()

	b.Run("findNonZeroKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findNonZeroKeyIndexScalar(keys)
		}
	})

	b.Run("findNonZeroKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findNonZeroKeyIndexAVX2(keys)
		}
	})
}

// Benchmark for FindLastNonZeroKeyIndex
func BenchmarkFindLastNonZeroKeyIndex_SIMD(b *testing.B) {
	keys := &[256]byte{}
	// Set some non-zero elements for testing
	keys[0] = 1
	keys[128] = 42
	keys[255] = 100

	b.ResetTimer()

	b.Run("findLastNonZeroKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findLastNonZeroKeyIndexScalar(keys)
		}
	})

	b.Run("findLastNonZeroKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findLastNonZeroKeyIndexAVX2(keys)
		}
	})
}

// Benchmark different patterns for non-zero search
func BenchmarkFindNonZeroKeyIndex_FirstElement(b *testing.B) {
	keys := &[256]byte{}
	keys[0] = 1 // Only first element is non-zero

	b.ResetTimer()

	b.Run("findNonZeroKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findNonZeroKeyIndexScalar(keys)
		}
	})

	b.Run("findNonZeroKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findNonZeroKeyIndexAVX2(keys)
		}
	})
}

func BenchmarkFindNonZeroKeyIndex_LastElement(b *testing.B) {
	keys := &[256]byte{}
	keys[255] = 1 // Only last element is non-zero

	b.ResetTimer()

	b.Run("findNonZeroKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findNonZeroKeyIndexScalar(keys)
		}
	})

	b.Run("findNonZeroKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findNonZeroKeyIndexAVX2(keys)
		}
	})
}

func BenchmarkFindNonZeroKeyIndex_MiddleElement(b *testing.B) {
	keys := &[256]byte{}
	keys[128] = 1 // Only middle element is non-zero

	b.ResetTimer()

	b.Run("findNonZeroKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findNonZeroKeyIndexScalar(keys)
		}
	})

	b.Run("findNonZeroKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findNonZeroKeyIndexAVX2(keys)
		}
	})
}

func BenchmarkFindLastNonZeroKeyIndex_FirstElement(b *testing.B) {
	keys := &[256]byte{}
	keys[0] = 1 // Only first element is non-zero

	b.ResetTimer()

	b.Run("findLastNonZeroKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findLastNonZeroKeyIndexScalar(keys)
		}
	})

	b.Run("findLastNonZeroKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findLastNonZeroKeyIndexAVX2(keys)
		}
	})
}

func BenchmarkFindLastNonZeroKeyIndex_LastElement(b *testing.B) {
	keys := &[256]byte{}
	keys[255] = 1 // Only last element is non-zero

	b.ResetTimer()

	b.Run("findLastNonZeroKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findLastNonZeroKeyIndexScalar(keys)
		}
	})

	b.Run("findLastNonZeroKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findLastNonZeroKeyIndexAVX2(keys)
		}
	})
}

func BenchmarkFindLastNonZeroKeyIndex_MiddleElement(b *testing.B) {
	keys := &[256]byte{}
	keys[128] = 1 // Only middle element is non-zero

	b.ResetTimer()

	b.Run("findLastNonZeroKeyIndexScalar", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findLastNonZeroKeyIndexScalar(keys)
		}
	})

	b.Run("findLastNonZeroKeyIndexAVX2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findLastNonZeroKeyIndexAVX2(keys)
		}
	})
}
