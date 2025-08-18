package simd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Unit tests for FindKeyIndex using GoConvey
func TestFindKeyIndex(t *testing.T) {
	Convey("Given FindKeyIndex function", t, func() {
		Convey("When searching in an empty array", func() {
			keys := &[16]byte{}

			result := FindKeyIndex(keys, 0, 42)

			So(result, ShouldEqual, -1)
		})

		Convey("When searching in a single element array", func() {
			keys := &[16]byte{42}

			Convey("And the key is found", func() {
				result := FindKeyIndex(keys, 1, 42)

				So(result, ShouldEqual, 0)
			})

			Convey("And the key is not found", func() {
				result := FindKeyIndex(keys, 1, 24)

				So(result, ShouldEqual, -1)
			})
		})

		Convey("When searching in a multiple element array", func() {
			keys := &[16]byte{1, 2, 3, 4, 5}

			Convey("And searching for the first element", func() {
				result := FindKeyIndex(keys, 5, 1)

				So(result, ShouldEqual, 0)
			})

			Convey("And searching for a middle element", func() {
				result := FindKeyIndex(keys, 5, 3)

				So(result, ShouldEqual, 2)
			})

			Convey("And searching for the last element", func() {
				result := FindKeyIndex(keys, 5, 5)

				So(result, ShouldEqual, 4)
			})

			Convey("And searching for a non-existent element", func() {
				result := FindKeyIndex(keys, 5, 6)

				So(result, ShouldEqual, -1)
			})
		})

		Convey("When searching in a full 16-byte array", func() {
			keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30}

			Convey("And searching for the first element", func() {
				result := FindKeyIndex(keys, 16, 0)

				So(result, ShouldEqual, 0)
			})

			Convey("And searching for a middle element", func() {
				result := FindKeyIndex(keys, 16, 16)

				So(result, ShouldEqual, 8)
			})

			Convey("And searching for the last element", func() {
				result := FindKeyIndex(keys, 16, 30)

				So(result, ShouldEqual, 15)
			})

			Convey("And searching for a non-existent element", func() {
				result := FindKeyIndex(keys, 16, 31)

				So(result, ShouldEqual, -1)
			})
		})

		Convey("When searching in an array with duplicate elements", func() {
			keys := &[16]byte{1, 1, 2, 3, 4}

			result := FindKeyIndex(keys, 5, 1)

			So(result, ShouldEqual, 0) // Should return first occurrence
		})

		Convey("When searching in an array with all same elements", func() {
			keys := &[16]byte{42, 42, 42, 42, 42}

			Convey("And searching for the existing element", func() {
				result := FindKeyIndex(keys, 5, 42)

				So(result, ShouldEqual, 0)
			})

			Convey("And searching for a non-existent element", func() {
				result := FindKeyIndex(keys, 5, 24)

				So(result, ShouldEqual, -1)
			})
		})
	})
}

// Unit tests for FindInsertPosition using GoConvey
func TestFindInsertPosition(t *testing.T) {
	Convey("Given FindInsertPosition function", t, func() {
		Convey("When inserting into an empty array", func() {
			keys := &[16]byte{}

			result := FindInsertPosition(keys, 0, 42)

			So(result, ShouldEqual, 0)
		})

		Convey("When inserting into a single element array", func() {
			keys := &[16]byte{5}

			Convey("And inserting before the existing element", func() {
				result := FindInsertPosition(keys, 1, 3)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting after the existing element", func() {
				result := FindInsertPosition(keys, 1, 7)

				So(result, ShouldEqual, 1)
			})
		})

		Convey("When inserting into a multiple element array", func() {
			keys := &[16]byte{2, 4, 6, 8}

			Convey("And inserting at the beginning", func() {
				result := FindInsertPosition(keys, 4, 1)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting in the middle", func() {
				result := FindInsertPosition(keys, 4, 5)

				So(result, ShouldEqual, 2)
			})

			Convey("And inserting at the end", func() {
				result := FindInsertPosition(keys, 4, 9)

				So(result, ShouldEqual, 4)
			})
		})

		Convey("When inserting into a full 16-byte array", func() {
			keys := &[16]byte{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32}

			Convey("And inserting at the beginning", func() {
				result := FindInsertPosition(keys, 16, 1)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting in the middle", func() {
				result := FindInsertPosition(keys, 16, 15)

				So(result, ShouldEqual, 7)
			})

			Convey("And inserting at the end", func() {
				result := FindInsertPosition(keys, 16, 33)

				So(result, ShouldEqual, 16)
			})
		})

		Convey("When inserting into an array with duplicate elements", func() {
			keys := &[16]byte{2, 2, 4, 6}

			Convey("And inserting before the first duplicate", func() {
				result := FindInsertPosition(keys, 4, 1)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting between duplicates", func() {
				result := FindInsertPosition(keys, 4, 3)

				So(result, ShouldEqual, 2)
			})
		})

		Convey("When inserting into an array with all same elements", func() {
			keys := &[16]byte{5, 5, 5, 5}

			Convey("And inserting before all elements", func() {
				result := FindInsertPosition(keys, 4, 3)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting after all elements", func() {
				result := FindInsertPosition(keys, 4, 7)

				So(result, ShouldEqual, 4)
			})
		})
	})
}

// Unit tests for scalar implementations using GoConvey
func TestFindKeyIndexScalar(t *testing.T) {
	Convey("Given findKeyIndexScalar function", t, func() {
		Convey("When searching in an empty array", func() {
			keys := &[16]byte{}

			result := findKeyIndexScalar(keys, 0, 42)

			So(result, ShouldEqual, -1)
		})

		Convey("When searching in a single element array", func() {
			keys := &[16]byte{42}

			Convey("And the key is found", func() {
				result := findKeyIndexScalar(keys, 1, 42)

				So(result, ShouldEqual, 0)
			})

			Convey("And the key is not found", func() {
				result := findKeyIndexScalar(keys, 1, 24)

				So(result, ShouldEqual, -1)
			})
		})

		Convey("When searching in a multiple element array", func() {
			keys := &[16]byte{1, 2, 3, 4, 5}

			Convey("And searching for an existing element", func() {
				result := findKeyIndexScalar(keys, 5, 3)

				So(result, ShouldEqual, 2)
			})

			Convey("And searching for a non-existent element", func() {
				result := findKeyIndexScalar(keys, 5, 6)

				So(result, ShouldEqual, -1)
			})
		})
	})
}

func TestFindInsertPositionScalar(t *testing.T) {
	Convey("Given findInsertPositionScalar function", t, func() {
		Convey("When inserting into an empty array", func() {
			keys := &[16]byte{}

			result := findInsertPositionScalar(keys, 0, 42)

			So(result, ShouldEqual, 0)
		})

		Convey("When inserting into a single element array", func() {
			keys := &[16]byte{5}

			Convey("And inserting before the existing element", func() {
				result := findInsertPositionScalar(keys, 1, 3)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting after the existing element", func() {
				result := findInsertPositionScalar(keys, 1, 7)

				So(result, ShouldEqual, 1)
			})
		})

		Convey("When inserting into a multiple element array", func() {
			keys := &[16]byte{2, 4, 6, 8}

			Convey("And inserting in the middle", func() {
				result := findInsertPositionScalar(keys, 4, 5)

				So(result, ShouldEqual, 2)
			})

			Convey("And inserting at the end", func() {
				result := findInsertPositionScalar(keys, 4, 9)

				So(result, ShouldEqual, 4)
			})
		})
	})
}

// Edge case tests using GoConvey
func TestEdgeCases(t *testing.T) {
	Convey("Given edge cases", t, func() {
		Convey("When testing boundary conditions", func() {
			keys := &[16]byte{255, 254, 253, 252}

			Convey("And testing maximum byte value", func() {
				result := FindKeyIndex(keys, 4, 255)
				So(result, ShouldEqual, 0)
			})

			Convey("And testing minimum byte value", func() {
				result := FindKeyIndex(keys, 4, 0)
				So(result, ShouldEqual, -1)
			})

			Convey("And testing insert position with maximum values", func() {
				result := FindInsertPosition(keys, 4, 255)
				So(result, ShouldEqual, 4)
			})

			Convey("And testing insert position with minimum values", func() {
				result := FindInsertPosition(keys, 4, 0)
				So(result, ShouldEqual, 0)
			})
		})

		Convey("When testing array bounds", func() {
			keys := &[16]byte{1, 2, 3, 4, 5}

			Convey("And testing with n = 0", func() {
				result := FindKeyIndex(keys, 0, 1)
				So(result, ShouldEqual, -1)

				result = FindInsertPosition(keys, 0, 1)
				So(result, ShouldEqual, 0)
			})

			Convey("And testing with n = 1", func() {
				result := FindKeyIndex(keys, 1, 1)
				So(result, ShouldEqual, 0)

				result = FindInsertPosition(keys, 1, 0)
				So(result, ShouldEqual, 0)
			})
		})
	})
}

// Test data setup functions
func setupBenchmarkData() (*[16]byte, int) {
	keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30}
	return keys, 16
}

// Performance comparison test using GoConvey
func TestPerformanceComparison(t *testing.T) {
	Convey("Given performance comparison", t, func() {
		keys, n := setupBenchmarkData()
		iterations := 1000000

		Convey("When comparing FindKeyIndex performance", func() {
			scalarTime := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < iterations; i++ {
					_ = findKeyIndexScalar(keys, n, byte(i%32))
				}
			})

			simdTime := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < iterations; i++ {
					_ = FindKeyIndex(keys, n, byte(i%32))
				}
			})

			t.Logf("=== FindKeyIndex Performance Test ===")
			t.Logf("Scalar: %d iterations in %v (%.2f ns/op)", iterations, scalarTime.T, float64(scalarTime.T.Nanoseconds())/float64(iterations))
			t.Logf("SIMD:   %d iterations in %v (%.2f ns/op)", iterations, simdTime.T, float64(simdTime.T.Nanoseconds())/float64(iterations))

			ratio := float64(simdTime.T.Nanoseconds()) / float64(scalarTime.T.Nanoseconds())
			if ratio > 1 {
				t.Logf("Ratio:  SIMD is %.2fx slower than Scalar", ratio)
			} else {
				t.Logf("Ratio:  SIMD is %.2fx faster than Scalar", 1/ratio)
			}

			So(ratio, ShouldBeGreaterThan, 0) // Basic validation
		})

		Convey("When comparing FindInsertPosition performance", func() {
			scalarTime := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < iterations; i++ {
					_ = findInsertPositionScalar(keys, n, byte(i%32))
				}
			})

			simdTime := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < iterations; i++ {
					_ = FindInsertPosition(keys, n, byte(i%32))
				}
			})

			t.Logf("\n=== FindInsertPosition Performance Test ===")
			t.Logf("Scalar: %d iterations in %v (%.2f ns/op)", iterations, scalarTime.T, float64(scalarTime.T.Nanoseconds())/float64(iterations))
			t.Logf("SIMD:   %d iterations in %v (%.2f ns/op)", iterations, simdTime.T, float64(simdTime.T.Nanoseconds())/float64(iterations))

			ratio := float64(simdTime.T.Nanoseconds()) / float64(scalarTime.T.Nanoseconds())
			if ratio > 1 {
				t.Logf("Ratio:  SIMD is %.2fx slower than Scalar", ratio)
			} else {
				t.Logf("Ratio:  SIMD is %.2fx faster than Scalar", 1/ratio)
			}

			So(ratio, ShouldBeGreaterThan, 0) // Basic validation
		})
	})
}

// Benchmark scalar implementations
func BenchmarkFindKeyIndexScalar(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = findKeyIndexScalar(keys, n, byte(i%32))
	}
}

func BenchmarkFindInsertPositionScalar(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = findInsertPositionScalar(keys, n, byte(i%32))
	}
}

// Benchmark SIMD implementations (only on supported architectures)
func BenchmarkFindKeyIndexSIMD(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, n, byte(i%32))
	}
}

func BenchmarkFindInsertPositionSIMD(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindInsertPosition(keys, n, byte(i%32))
	}
}

// Benchmark different key distributions
func BenchmarkFindKeyIndex_FirstElement(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, n, 0) // First element
	}
}

func BenchmarkFindKeyIndex_LastElement(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, n, 30) // Last element
	}
}

func BenchmarkFindKeyIndex_NotFound(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, n, 31) // Not found
	}
}

func BenchmarkFindInsertPosition_FirstPosition(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindInsertPosition(keys, n, 1) // Insert at first position
	}
}

func BenchmarkFindInsertPosition_MiddlePosition(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindInsertPosition(keys, n, 15) // Insert at middle position
	}
}

func BenchmarkFindInsertPosition_LastPosition(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindInsertPosition(keys, n, 31) // Insert at last position
	}
}

// Benchmark different array sizes
func BenchmarkFindKeyIndex_Size4(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, 4, byte(i%8))
	}
}

func BenchmarkFindKeyIndex_Size8(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, 8, byte(i%16))
	}
}

func BenchmarkFindKeyIndex_Size16(b *testing.B) {
	keys, n := setupBenchmarkData()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, n, byte(i%32))
	}
}

// Benchmark scalar vs SIMD for different sizes
func BenchmarkFindKeyIndex_Scalar_Size4(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = findKeyIndexScalar(keys, 4, byte(i%8))
	}
}

func BenchmarkFindKeyIndex_SIMD_Size4(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, 4, byte(i%8))
	}
}

func BenchmarkFindKeyIndex_Scalar_Size8(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = findKeyIndexScalar(keys, 8, byte(i%16))
	}
}

func BenchmarkFindKeyIndex_SIMD_Size8(b *testing.B) {
	keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, 8, byte(i%16))
	}
}

// Benchmark with realistic data patterns
func BenchmarkFindKeyIndex_Sequential(b *testing.B) {
	keys := &[16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, 16, byte(i%16))
	}
}

func BenchmarkFindKeyIndex_Sparse(b *testing.B) {
	keys := &[16]byte{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindKeyIndex(keys, 16, byte(i%160))
	}
}

// Benchmark with different key distributions for insert position
func BenchmarkFindInsertPosition_Sequential(b *testing.B) {
	keys := &[16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindInsertPosition(keys, 16, byte(i%16))
	}
}

func BenchmarkFindInsertPosition_Sparse(b *testing.B) {
	keys := &[16]byte{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindInsertPosition(keys, 16, byte(i%160))
	}
}
