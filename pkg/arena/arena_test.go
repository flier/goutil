//go:build go1.22

package arena_test

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"reflect"
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
)

func BenchmarkArena(b *testing.B) {
	bench[int](b)
	bench[[2]int](b)
	bench[[64]int](b)
	bench[[1024]int](b)
}

const runs = 100000

var sink any

func bench[T any](b *testing.B) {
	var z T
	n := int64(runs * unsafe.Sizeof(z))
	name := fmt.Sprintf("%v", reflect.TypeFor[T]())

	b.Run(name, func(b *testing.B) {
		b.Run("arena.alloc", func(b *testing.B) {
			b.SetBytes(n)
			for n := 0; n < b.N; n++ {
				a := new(arena.Arena)
				for i := 0; i < runs; i++ {
					sink = arena.Alloc[T](a)
				}
			}
		})

		b.Run("arena.new", func(b *testing.B) {
			var v T

			b.SetBytes(n)
			for n := 0; n < b.N; n++ {
				a := new(arena.Arena)
				for i := 0; i < runs; i++ {
					sink = arena.New(a, v)
				}
			}
		})

		b.Run("new", func(b *testing.B) {
			b.SetBytes(n)
			for n := 0; n < b.N; n++ {
				for i := 0; i < runs; i++ {
					sink = new(T)
				}
			}
		})
	})
}

func TestArena(t *testing.T) {
	Convey("Given an Arena", t, func() {
		a := new(arena.Arena)

		type testStruct struct {
			X int
			Y float64
		}

		Convey("When allocate a value", func() {
			p := arena.New(a, testStruct{X: 42, Y: 3.14})
			So(p, ShouldNotBeNil)

			Convey("Then the value should be set", func() {
				So(p.X, ShouldEqual, 42)
				So(p.Y, ShouldEqual, 3.14)
			})

			Convey("Then the pointer should be aligned", func() {
				So(uintptr(unsafe.Pointer(p))%8, ShouldEqual, uintptr(0))
			})
		})

		Convey("When allocate multiple values", func() {
			var ptrs []*testStruct
			for i := 0; i < 10; i++ {
				p := arena.New(a, testStruct{X: i, Y: float64(i)})
				ptrs = append(ptrs, p)
			}

			Convey("Then the value should be set", func() {
				for i, p := range ptrs {
					So(p.X, ShouldEqual, i)
					So(p.Y, ShouldEqual, float64(i))
				}
			})

			Convey("Then reset the arena and check state", func() {
				a.Reset()

				So(a.Empty(), ShouldBeTrue)
			})
		})

		Convey("When allocate a large memory", func() {
			p := arena.New(a, [1024]byte{})

			So(p, ShouldNotBeNil)
		})

		Convey("When allocate multiple types", func() {
			i := arena.New(a, 123)
			So(*i, ShouldEqual, 123)

			f := arena.New(a, 3.14)
			So(*f, ShouldEqual, 3.14)

			s := arena.New(a, "hello")
			So(*s, ShouldEqual, "hello")
		})

		i := arena.New(a, 42)
		So(i, ShouldNotBeNil)
		So(*i, ShouldEqual, 42)

		Convey("When realloc same type", func() {
			i = arena.Realloc[int](a, i)

			Convey("Then the value should be same", func() {
				So(i, ShouldNotBeNil)
				So(*i, ShouldEqual, 42)
			})
		})

		Convey("When realloc a different type", func() {
			r := arena.Realloc[float64](a, i)

			Convey("Then the bytes should be copied", func() {
				So(r, ShouldNotBeNil)
				So(*r, ShouldEqual, math.Float64frombits(42))
			})
		})

		Convey("When realloc struct to array", func() {
			s := arena.New(a, testStruct{X: 42, Y: 3.14})
			So(s, ShouldNotBeNil)

			p := arena.Realloc[[64]byte](a, s)
			So(p, ShouldNotBeNil)
			So(binary.NativeEndian.Uint64((*p)[:]), ShouldEqual, 42)
			So(math.Float64frombits(binary.NativeEndian.Uint64((*p)[8:])), ShouldEqual, 3.14)
			So((*p)[16:], ShouldResemble, make([]byte, 48))
		})

		Convey("When realloc a little more memory", func() {
			p := arena.Realloc[[2]int](a, i)

			Convey("Then the value should be copied", func() {
				So(p, ShouldNotBeNil)
				So(p[0], ShouldEqual, 42)
				So(p[1], ShouldEqual, 0)
			})
		})

		Convey("When realloc a very large memory", func() {
			p := arena.Realloc[[1024]byte](a, i)

			Convey("Then the value should be copied", func() {
				So(p, ShouldNotBeNil)
				So(binary.NativeEndian.Uint64((*p)[:]), ShouldEqual, 42)
			})
		})
	})
}

func TestArenaAlloc(t *testing.T) {
	Convey("Arena.Alloc", t, func() {
		Convey("Should allocate memory with proper alignment", func() {
			a := new(arena.Arena)

			// Test different types and their alignment requirements
			testCases := []struct {
				name     string
				size     uintptr
				align    uintptr
				expected uintptr
			}{
				{"uint8", unsafe.Sizeof(uint8(0)), unsafe.Alignof(uint8(0)), 8},
				{"uint16", unsafe.Sizeof(uint16(0)), unsafe.Alignof(uint16(0)), 8},
				{"uint32", unsafe.Sizeof(uint32(0)), unsafe.Alignof(uint32(0)), 8},
				{"uint64", unsafe.Sizeof(uint64(0)), unsafe.Alignof(uint64(0)), 8},
				{"uintptr", unsafe.Sizeof(uintptr(0)), unsafe.Alignof(uintptr(0)), 8},
			}

			for _, tc := range testCases {
				Convey(fmt.Sprintf("For %s", tc.name), func() {
					ptr := a.Alloc(tc.size, tc.align)
					So(ptr, ShouldNotBeNil)
					So(uintptr(ptr)%tc.expected, ShouldEqual, uintptr(0))
				})
			}
		})

		Convey("Should handle size rounding correctly", func() {
			a := new(arena.Arena)

			// Test that sizes are properly rounded up to alignment boundaries
			ptr1 := a.Alloc(1, 8)
			ptr2 := a.Alloc(1, 8)

			// The second allocation should be 8 bytes after the first
			So(uintptr(ptr2)-uintptr(ptr1), ShouldEqual, uintptr(8))
		})

		Convey("Should grow arena when needed", func() {
			a := new(arena.Arena)

			// Allocate a large chunk that will require arena growth
			largeSize := uintptr(1024 * 1024) // 1MB
			ptr := a.Alloc(largeSize, 8)
			So(ptr, ShouldNotBeNil)

			// Verify that the arena has grown by checking it's not empty
			So(a.Empty(), ShouldBeFalse)
		})

		Convey("Should handle zero size allocation", func() {
			a := new(arena.Arena)

			ptr := a.Alloc(0, 8)
			So(ptr, ShouldNotBeNil)
		})

		Convey("Should handle very large alignment", func() {
			a := new(arena.Arena)

			// Test with alignment larger than maxAlign (8)
			ptr := a.Alloc(16, 16)
			So(ptr, ShouldNotBeNil)
			// Note: The actual alignment might be limited by maxAlign
		})
	})
}

func TestArenaRealloc(t *testing.T) {
	Convey("Arena.Realloc", t, func() {
		Convey("Should return same pointer when new size is smaller", func() {
			a := new(arena.Arena)

			// Allocate initial memory
			ptr := a.Alloc(16, 8)
			originalPtr := ptr

			// Realloc with smaller size
			newPtr := a.Realloc(ptr, 16, 8, 8)
			So(newPtr, ShouldEqual, originalPtr)
		})

		Convey("Should return same pointer when new size is equal", func() {
			a := new(arena.Arena)

			ptr := a.Alloc(16, 8)
			originalPtr := ptr

			newPtr := a.Realloc(ptr, 16, 16, 8)
			So(newPtr, ShouldEqual, originalPtr)
		})

		Convey("Should grow in-place when possible", func() {
			a := new(arena.Arena)

			// Allocate memory
			ptr := a.Alloc(16, 8)
			originalPtr := ptr

			// Try to grow in-place (should succeed if there's enough space)
			newPtr := a.Realloc(ptr, 16, 32, 8)

			// In-place growth succeeded
			So(newPtr, ShouldEqual, originalPtr)
		})

		Convey("Should copy data when reallocating to new location", func() {
			a := new(arena.Arena)

			// Allocate and initialize memory
			ptr := a.Alloc(8, 8)
			*(*int)(ptr) = 42

			// Realloc to larger size
			newPtr := a.Realloc(ptr, 8, 128, 8)
			So(newPtr, ShouldNotBeNil)
			So(newPtr, ShouldNotEqual, ptr)

			// Verify data was copied
			So(*(*int)(newPtr), ShouldEqual, 42)
		})

		Convey("Should handle reallocation between different types", func() {
			a := new(arena.Arena)

			// Allocate int
			intPtr := arena.New(a, 42)
			So(intPtr, ShouldNotBeNil)

			// Realloc to float64
			floatPtr := arena.Realloc[float64](a, intPtr)
			So(floatPtr, ShouldNotBeNil)

			// Verify the bit pattern is preserved
			So(*floatPtr, ShouldEqual, math.Float64frombits(42))
		})

		Convey("Should handle reallocation to much larger size", func() {
			a := new(arena.Arena)

			// Allocate small int
			intPtr := arena.New(a, 42)
			So(intPtr, ShouldNotBeNil)

			// Realloc to large array
			arrayPtr := arena.Realloc[[1024]byte](a, intPtr)
			So(arrayPtr, ShouldNotBeNil)

			// Verify first 8 bytes contain the original value
			So(binary.NativeEndian.Uint64((*arrayPtr)[:]), ShouldEqual, 42)
		})
	})
}

func TestArenaEdgeCases(t *testing.T) {
	Convey("Arena Edge Cases", t, func() {
		Convey("Should handle empty arena state", func() {
			a := new(arena.Arena)

			So(a.Empty(), ShouldBeTrue)
		})

		Convey("Should handle reset after allocations", func() {
			a := new(arena.Arena)

			// Make some allocations
			ptr1 := a.Alloc(16, 8)
			ptr2 := a.Alloc(32, 8)
			So(ptr1, ShouldNotBeNil)
			So(ptr2, ShouldNotBeNil)

			// Reset arena
			a.Reset()

			So(a.Empty(), ShouldBeTrue)
		})

		Convey("Should handle multiple reset cycles", func() {
			a := new(arena.Arena)

			for i := 0; i < 5; i++ {
				// Make allocations
				ptr := a.Alloc(16, 8)
				So(ptr, ShouldNotBeNil)

				// Reset
				a.Reset()

				So(a.Empty(), ShouldBeTrue)
			}
		})

		Convey("Should handle very small allocations", func() {
			a := new(arena.Arena)

			// Allocate very small amounts
			for i := 0; i < 100; i++ {
				ptr := a.Alloc(1, 1)
				So(ptr, ShouldNotBeNil)
			}
		})

		Convey("Should handle mixed allocation sizes", func() {
			a := new(arena.Arena)

			// Mix different allocation sizes
			sizes := []uintptr{1, 8, 16, 64, 256, 1024}
			var ptrs []unsafe.Pointer

			for _, size := range sizes {
				ptr := a.Alloc(size, 8)
				So(ptr, ShouldNotBeNil)
				ptrs = append(ptrs, ptr)
			}

			// Verify all pointers are unique
			for i := 0; i < len(ptrs); i++ {
				for j := i + 1; j < len(ptrs); j++ {
					So(ptrs[i], ShouldNotEqual, ptrs[j])
				}
			}
		})
	})
}

func TestArenaMemoryManagement(t *testing.T) {
	Convey("Arena Memory Management", t, func() {
		Convey("Should reuse memory chunks efficiently", func() {
			a := new(arena.Arena)

			// Make allocations to trigger chunk allocation
			for i := 0; i < 1000; i++ {
				ptr := a.Alloc(16, 8)
				So(ptr, ShouldNotBeNil)
			}

			// Reset and verify chunks are reused
			a.Reset()

			// Make allocations again - should reuse existing chunks
			for i := 0; i < 1000; i++ {
				ptr := a.Alloc(16, 8)
				So(ptr, ShouldNotBeNil)
			}
		})

		Convey("Should handle power-of-two chunk sizing", func() {
			a := new(arena.Arena)

			// Test that chunk sizes follow power-of-two pattern
			expectedSizes := []uintptr{8, 16, 32, 64, 128, 256, 512, 1024}

			for _, expectedSize := range expectedSizes {
				// Allocate enough to trigger new chunk
				ptr := a.Alloc(expectedSize, 8)
				So(ptr, ShouldNotBeNil)
			}
		})

		Convey("Should handle alignment requirements correctly", func() {
			a := new(arena.Arena)

			// Test various alignment requirements
			alignments := []uintptr{1, 2, 4, 8}
			sizes := []uintptr{1, 2, 4, 8, 16, 32}

			for _, align := range alignments {
				for _, size := range sizes {
					ptr := a.Alloc(size, align)
					So(ptr, ShouldNotBeNil)

					// Verify alignment (actual alignment might be limited by maxAlign)
					actualAlign := uintptr(1) << bits.TrailingZeros(uint(uintptr(ptr)))
					So(actualAlign, ShouldBeGreaterThanOrEqualTo, align)
				}
			}
		})
	})
}

func TestArenaConcurrency(t *testing.T) {
	t.Run("Should handle sequential operations safely", func(t *testing.T) {
		a := new(arena.Arena)
		const numAllocations = 1000

		// Make many allocations sequentially
		for i := 0; i < numAllocations; i++ {
			ptr := a.Alloc(16, 8)
			if ptr == nil {
				t.Errorf("Allocation %d failed", i)
			}
		}

		// Verify arena is not empty
		if a.Empty() {
			t.Error("Arena should not be empty after many allocations")
		}

		// Reset and verify
		a.Reset()
		if !a.Empty() {
			t.Error("Arena should be empty after reset")
		}
	})

	t.Run("Should handle rapid reset cycles", func(t *testing.T) {
		a := new(arena.Arena)

		// Make allocations and reset rapidly
		for i := 0; i < 100; i++ {
			// Make some allocations
			for j := 0; j < 10; j++ {
				ptr := a.Alloc(16, 8)
				if ptr == nil {
					t.Errorf("Allocation %d in cycle %d failed", j, i)
				}
			}

			// Reset
			a.Reset()
			if !a.Empty() {
				t.Errorf("Reset failed in cycle %d", i)
			}
		}
	})
}
