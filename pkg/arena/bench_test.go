//go:build go1.22

package arena_test

import (
	"testing"

	. "github.com/flier/goutil/pkg/arena"
)

// BenchmarkRecycled_Release benchmarks Recycled release performance
func BenchmarkRecycled_Release(b *testing.B) {
	arena := &Recycled{}

	// Pre-allocate some memory
	pointers := make([]*byte, b.N)
	for i := range pointers {
		pointers[i] = arena.Alloc(64)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		arena.Release(pointers[i], 64)
	}
}

// BenchmarkRecycled_MultipleRecycling benchmarks Recycled with multiple allocation cycles
func BenchmarkRecycled_MultipleRecycling(b *testing.B) {
	arena := &Recycled{}
	sizes := []int{64, 128, 256, 512, 1024}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Allocate a set of blocks
		ptrs := make([]*byte, len(sizes))
		for j, size := range sizes {
			ptrs[j] = arena.Alloc(size)
		}

		// Release them
		for j, ptr := range ptrs {
			arena.Release(ptr, sizes[j])
		}

		// Reallocate them (should be recycled)
		for _, size := range sizes {
			ptr := arena.Alloc(size)
			arena.Release(ptr, size)
		}
	}
}

// BenchmarkComparison_ArenaVsRecycled compares Arena vs Recycled for basic allocation
func BenchmarkComparison_ArenaVsRecycled(b *testing.B) {
	b.Run("Arena", func(b *testing.B) {
		arena := &Arena{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := arena.Alloc(64)
			_ = ptr
		}
	})

	b.Run("Recycled", func(b *testing.B) {
		arena := &Recycled{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := arena.Alloc(64)
			_ = ptr
		}
	})
}

// BenchmarkComparison_Allocation compares Arena vs Recycled for allocation
func BenchmarkComparison_Allocation(b *testing.B) {
	b.Run("Arena_AllocOnly", func(b *testing.B) {
		arena := &Arena{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := arena.Alloc(64)
			_ = ptr
		}
	})

	b.Run("Recycled_AllocAndRelease", func(b *testing.B) {
		arena := &Recycled{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := arena.Alloc(64)
			_ = ptr
		}
	})
}

// BenchmarkComparison_AllocationAndRelease compares Arena vs Recycled for allocation and release
func BenchmarkComparison_AllocationAndRelease(b *testing.B) {
	b.Run("Arena_AllocOnly", func(b *testing.B) {
		arena := &Arena{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := arena.Alloc(64)
			_ = ptr
		}
	})

	b.Run("Recycled_AllocAndRelease", func(b *testing.B) {
		arena := &Recycled{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ptr := arena.Alloc(64)
			arena.Release(ptr, 64)
		}
	})
}

// BenchmarkComparison_MixedSizes compares Arena vs Recycled for mixed size allocations
func BenchmarkComparison_MixedSizes(b *testing.B) {
	sizes := []int{8, 16, 32, 64, 128, 256, 512, 1024}

	b.Run("Arena", func(b *testing.B) {
		arena := &Arena{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			size := sizes[i%len(sizes)]
			ptr := arena.Alloc(size)
			_ = ptr
		}
	})

	b.Run("Recycled", func(b *testing.B) {
		arena := &Recycled{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			size := sizes[i%len(sizes)]
			ptr := arena.Alloc(size)
			arena.Release(ptr, size)
		}
	})
}

// BenchmarkComparison_LargeAllocations compares Arena vs Recycled for large allocations
func BenchmarkComparison_LargeAllocations(b *testing.B) {
	sizes := []int{4096, 8192, 16384, 32768}

	b.Run("Arena", func(b *testing.B) {
		arena := &Arena{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			size := sizes[i%len(sizes)]
			ptr := arena.Alloc(size)
			_ = ptr
		}
	})

	b.Run("Recycled", func(b *testing.B) {
		arena := &Recycled{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			size := sizes[i%len(sizes)]
			ptr := arena.Alloc(size)
			arena.Release(ptr, size)
		}
	})
}
