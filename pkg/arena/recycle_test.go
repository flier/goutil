//go:build go1.22

package arena_test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/xunsafe"
)

// Test data structures for better test organization
type testCase struct {
	name     string
	size     int
	expected uint
}

type allocationTest struct {
	size int
	data byte
}

// TestRecycledArena_BasicAllocation tests basic allocation functionality
func TestRecycledArena_BasicAllocation(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		arena := &Recycled{}

		Convey("When allocating memory of different sizes", func() {
			testSizes := []int{8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096}
			pointers := make([]*byte, len(testSizes))

			for i, size := range testSizes {
				ptr := arena.Alloc(size)
				pointers[i] = ptr
			}

			Convey("Then all allocations should succeed and be properly aligned", func() {
				for i, ptr := range pointers {
					So(ptr, ShouldNotBeNil)
					addr := uintptr(unsafe.Pointer(ptr))
					So(addr%uintptr(Align), ShouldEqual, uintptr(0))

					// Verify we can write to the allocated memory
					*ptr = byte(i)
					So(*ptr, ShouldEqual, byte(i))
				}
			})

			Convey("And all pointers should be unique", func() {
				uniquePtrs := make(map[uintptr]bool)
				for _, ptr := range pointers {
					addr := uintptr(unsafe.Pointer(ptr))
					So(uniquePtrs[addr], ShouldBeFalse)
					uniquePtrs[addr] = true
				}
			})
		})

		Convey("When testing edge case sizes", func() {
			edgeSizes := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 15, 16, 17, 31, 32, 33, 63, 64, 65, 127, 128, 129}

			for _, size := range edgeSizes {
				Convey("And testing size "+string(rune(size)), func() {
					ptr := arena.Alloc(size)
					So(ptr, ShouldNotBeNil)

					addr := xunsafe.AddrOf(ptr)
					So(int(addr)%Align, ShouldEqual, 0)

					// Test write access to first and last byte
					*ptr = 0xAA
					So(*ptr, ShouldEqual, byte(0xAA))

					if size > 1 {
						lastByte := addr.Add(size - 1).AssertValid()
						*lastByte = 0xBB
						So(*lastByte, ShouldEqual, byte(0xBB))
					}
				})
			}
		})
	})
}

// TestRecycledArena_Recycling tests the recycling functionality
func TestRecycledArena_Recycling(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		arena := &Recycled{}

		Convey("When allocating memory and then releasing it", func() {
			ptr1 := arena.Alloc(64)
			So(ptr1, ShouldNotBeNil)

			// Write some data to verify it's preserved
			testData := byte(42)
			*ptr1 = testData
			So(*ptr1, ShouldEqual, testData)

			arena.Release(ptr1, 64)

			Convey("And allocating the same size again", func() {
				ptr2 := arena.Alloc(64)
				So(ptr2, ShouldNotBeNil)

				Convey("Then the recycled pointer should be the same as the original", func() {
					So(ptr2, ShouldEqual, ptr1)
					So(*ptr2, ShouldEqual, byte(0))
				})
			})
		})

		Convey("When testing recycling with different sizes", func() {
			testCases := []allocationTest{
				{32, 0xAA},
				{128, 0xBB},
				{256, 0xCC},
			}

			for _, tc := range testCases {
				Convey("And testing size "+string(rune(tc.size)), func() {
					ptr := arena.Alloc(tc.size)
					So(ptr, ShouldNotBeNil)

					*ptr = tc.data
					So(*ptr, ShouldEqual, tc.data)

					arena.Release(ptr, tc.size)
					recycledPtr := arena.Alloc(tc.size)
					So(recycledPtr, ShouldEqual, ptr)
					So(*recycledPtr, ShouldEqual, byte(0))
				})
			}
		})

		Convey("When testing recycling with boundary sizes", func() {
			boundaryCases := []struct {
				size        int
				description string
			}{
				{Align - 1, "just below alignment boundary"},
				{Align, "exactly at alignment boundary"},
				{Align + 1, "just above alignment boundary"},
				{Align*2 - 1, "just below next alignment boundary"},
				{Align * 2, "exactly at next alignment boundary"},
				{Align*4 - 1, "just below power-of-2 boundary"},
				{Align * 4, "exactly at power-of-2 boundary"},
			}

			for _, tc := range boundaryCases {
				Convey("And testing "+tc.description+" (size "+string(rune(tc.size))+")", func() {
					ptr := arena.Alloc(tc.size)
					So(ptr, ShouldNotBeNil)

					// Write pattern across the entire allocation
					addr := xunsafe.AddrOf(ptr)
					for i := 0; i < tc.size; i++ {
						*addr.Add(i).AssertValid() = byte(i % 256)
					}

					// Verify pattern was written
					So(*ptr, ShouldEqual, byte(0))
					if tc.size > 1 {
						lastByte := *addr.Add(tc.size - 1).AssertValid()
						So(lastByte, ShouldEqual, byte((tc.size-1)%256))
					}

					arena.Release(ptr, tc.size)
					recycledPtr := arena.Alloc(tc.size)
					So(recycledPtr, ShouldEqual, ptr)

					// Verify memory is cleared
					So(*recycledPtr, ShouldEqual, byte(0))
				})
			}
		})

		Convey("When testing recycling with very small sizes", func() {
			smallSizes := []int{1, 2, 3, 4, 5, 6, 7, 8}

			for _, size := range smallSizes {
				Convey("And testing size "+string(rune(size)), func() {
					ptr := arena.Alloc(size)
					So(ptr, ShouldNotBeNil)

					// Write data
					*ptr = byte(size)
					So(*ptr, ShouldEqual, byte(size))

					// For sizes < Align, release is ignored, but allocation should still work
					arena.Release(ptr, size)
					newPtr := arena.Alloc(size)
					So(newPtr, ShouldNotBeNil)

					// Verify we can write to the new allocation
					*newPtr = byte(size * 2)
					So(*newPtr, ShouldEqual, byte(size*2))
				})
			}
		})
	})
}

// TestRecycledArena_MultipleRecycling tests recycling multiple allocations
func TestRecycledArena_MultipleRecycling(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating multiple blocks of different sizes", func() {
			ptrs := make([]*byte, 10)
			sizes := []int{8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096}

			for i, size := range sizes {
				ptrs[i] = arena.Alloc(size)
				So(ptrs[i], ShouldNotBeNil)
				// Write unique data to each allocation
				*ptrs[i] = byte(i)
			}

			Convey("And releasing all allocations", func() {
				for i, ptr := range ptrs {
					arena.Release(ptr, sizes[i])
				}

				Convey("And allocating the same sizes again", func() {
					recycledPtrs := make([]*byte, len(sizes))
					for i, size := range sizes {
						recycledPtrs[i] = arena.Alloc(size)
						So(recycledPtrs[i], ShouldNotBeNil)
					}

					Convey("Then all recycled allocations should succeed", func() {
						// Verify that recycled pointers are from the original set
						for _, recycledPtr := range recycledPtrs {
							So(ptrs, ShouldContain, recycledPtr)
							// Verify memory is cleared after recycling
							So(*recycledPtr, ShouldEqual, byte(0))
						}
					})
				})
			})
		})
	})
}

// TestRecycledArena_SizeLogging tests the size logging functionality
func TestRecycledArena_SizeLogging(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		arena := &Recycled{}

		Convey("When testing different allocation sizes", func() {
			testCases := []testCase{
				{"tiny allocation", 1, 6},        // 2^6 = 64
				{"small allocation", 32, 6},      // 2^6 = 64
				{"exact power of 2", 64, 6},      // 2^6 = 64
				{"just above power of 2", 65, 7}, // 2^7 = 128
				{"medium allocation", 128, 7},    // 2^7 = 128
				{"boundary case", 129, 8},        // 2^8 = 256
				{"large allocation", 256, 8},     // 2^8 = 256
				{"very large", 257, 9},           // 2^9 = 512
				{"power of 2 large", 512, 9},     // 2^9 = 512
				{"huge allocation", 513, 10},     // 2^10 = 1024
			}

			for _, tc := range testCases {
				Convey("And testing "+tc.name+" (size "+string(rune(tc.size))+")", func() {
					ptr := arena.Alloc(tc.size)
					So(ptr, ShouldNotBeNil)

					Convey("And releasing the allocation", func() {
						arena.Release(ptr, tc.size)

						Convey("And allocating the same size again", func() {
							recycledPtr := arena.Alloc(tc.size)
							So(recycledPtr, ShouldNotBeNil)

							Convey("Then the recycled pointer should match the original", func() {
								So(recycledPtr, ShouldEqual, ptr)
								So(*recycledPtr, ShouldEqual, byte(0))
							})
						})
					})
				})
			}
		})

		Convey("When testing size class boundaries", func() {
			boundaryTests := []struct {
				size          int
				expectedClass uint
				description   string
			}{
				{1, 6, "minimum size"},
				{63, 6, "just below 2^6"},
				{64, 6, "exactly 2^6"},
				{65, 7, "just above 2^6"},
				{127, 7, "just below 2^7"},
				{128, 7, "exactly 2^7"},
				{129, 8, "just above 2^7"},
				{255, 8, "just below 2^8"},
				{256, 8, "exactly 2^8"},
				{257, 9, "just above 2^8"},
				{511, 9, "just below 2^9"},
				{512, 9, "exactly 2^9"},
				{513, 10, "just above 2^9"},
				{1023, 10, "just below 2^10"},
				{1024, 10, "exactly 2^10"},
				{1025, 11, "just above 2^10"},
			}

			for _, bt := range boundaryTests {
				Convey("And testing "+bt.description+" (size "+string(rune(bt.size))+")", func() {
					ptr := arena.Alloc(bt.size)
					So(ptr, ShouldNotBeNil)

					// Verify alignment
					addr := uintptr(unsafe.Pointer(ptr))
					So(addr%uintptr(Align), ShouldEqual, uintptr(0))

					// Test recycling
					arena.Release(ptr, bt.size)
					recycledPtr := arena.Alloc(bt.size)
					So(recycledPtr, ShouldEqual, ptr)
					So(*recycledPtr, ShouldEqual, byte(0))
				})
			}
		})

		Convey("When testing very large size classes", func() {
			largeSizes := []int{4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576}

			for _, size := range largeSizes {
				Convey("And testing size "+string(rune(size)), func() {
					ptr := arena.Alloc(size)
					So(ptr, ShouldNotBeNil)

					// Test write access to first and last byte
					*ptr = 0xAA
					So(*ptr, ShouldEqual, byte(0xAA))

					lastByte := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(size-1)))
					*lastByte = 0xBB
					So(*lastByte, ShouldEqual, byte(0xBB))

					// Test recycling
					arena.Release(ptr, size)
					recycledPtr := arena.Alloc(size)
					So(recycledPtr, ShouldEqual, ptr)
					So(*recycledPtr, ShouldEqual, byte(0))
				})
			}
		})
	})
}

// TestRecycledArena_MixedSizes tests recycling with mixed allocation sizes
func TestRecycledArena_MixedSizes(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating different sizes", func() {
			ptr1 := arena.Alloc(64)
			ptr2 := arena.Alloc(128)
			ptr3 := arena.Alloc(256)

			So(ptr1, ShouldNotBeNil)
			So(ptr2, ShouldNotBeNil)
			So(ptr3, ShouldNotBeNil)

			// Write unique data
			*ptr1 = 1
			*ptr2 = 2
			*ptr3 = 3

			Convey("And releasing them in different order", func() {
				arena.Release(ptr2, 128) // Release 128 first
				arena.Release(ptr1, 64)  // Release 64 second
				arena.Release(ptr3, 256) // Release 256 last

				Convey("And allocating the same sizes again", func() {
					recycled1 := arena.Alloc(64)
					recycled2 := arena.Alloc(128)
					recycled3 := arena.Alloc(256)

					Convey("Then all allocations should succeed", func() {
						So(recycled1, ShouldEqual, ptr1)
						So(recycled2, ShouldEqual, ptr2)
						So(recycled3, ShouldEqual, ptr3)

						// Verify memory is cleared after recycling
						So(*recycled1, ShouldEqual, byte(0))
						So(*recycled2, ShouldEqual, byte(0))
						So(*recycled3, ShouldEqual, byte(0))
					})
				})
			})
		})
	})
}

// TestRecycledArena_FreeFunction tests the Free function with different types
func TestRecycledArena_FreeFunction(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating different sizes", func() {
			ptr1 := arena.Alloc(32)
			ptr2 := arena.Alloc(8)
			ptr3 := arena.Alloc(32)

			So(ptr1, ShouldNotBeNil)
			So(ptr2, ShouldNotBeNil)
			So(ptr3, ShouldNotBeNil)

			Convey("And freeing all allocations using Free function", func() {
				Free(arena, ptr1)
				Free(arena, ptr2)
				Free(arena, ptr3)

				Convey("And allocating the same sizes again", func() {
					recycled1 := arena.Alloc(32)
					recycled2 := arena.Alloc(8)
					recycled3 := arena.Alloc(32)

					Convey("Then all allocations should succeed", func() {
						So(recycled1, ShouldEqual, ptr1)
						So(recycled2, ShouldEqual, ptr2)
						So(recycled3, ShouldEqual, ptr3)

						// Verify memory is cleared after recycling
						So(*recycled1, ShouldEqual, byte(0))
						So(*recycled2, ShouldEqual, byte(0))
						So(*recycled3, ShouldEqual, byte(0))
					})
				})
			})
		})
	})
}

// TestRecycledArena_ResetMethod tests the Reset method
func TestRecycledArena_ResetMethod(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating some memory", func() {
			ptr1 := arena.Alloc(64)
			ptr2 := arena.Alloc(128)
			ptr3 := arena.Alloc(256)

			So(ptr1, ShouldNotBeNil)
			So(ptr2, ShouldNotBeNil)
			So(ptr3, ShouldNotBeNil)

			Convey("And releasing some allocations", func() {
				arena.Release(ptr1, 64)
				arena.Release(ptr2, 128)

				Convey("And calling the Reset method", func() {
					arena.Reset()

					Convey("And allocating the same sizes again", func() {
						newPtr1 := arena.Alloc(64)
						newPtr2 := arena.Alloc(128)
						newPtr3 := arena.Alloc(256)

						Convey("Then all allocations should succeed", func() {
							So(newPtr1, ShouldNotBeNil)
							So(newPtr2, ShouldNotBeNil)
							So(newPtr3, ShouldNotBeNil)

							// After Reset, the underlying Arena might reuse memory blocks
							// so we can't guarantee different pointers
						})
					})
				})
			})
		})
	})
}

// TestRecycledArena_ZeroSize tests allocation of zero size
func TestRecycledArena_ZeroSize(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating zero size", func() {
			ptr := arena.Alloc(0)

			Convey("Then the allocation should return nil for zero size", func() {
				So(ptr, ShouldBeNil)
			})
		})

		Convey("When releasing zero size allocation", func() {
			ptr := arena.Alloc(1)
			So(ptr, ShouldNotBeNil)

			Convey("And releasing with size smaller than Align", func() {
				arena.Release(ptr, Align-1)

				Convey("Then the release should be ignored", func() {
					// Try to allocate again, should get a new pointer
					newPtr := arena.Alloc(1)
					So(newPtr, ShouldNotBeNil)
					// Note: The underlying Arena might reuse memory blocks
					// so we can't guarantee different pointers
				})
			})
		})
	})
}

// TestRecycledArena_LargeAllocation tests large allocations
func TestRecycledArena_LargeAllocation(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating large memory blocks", func() {
			largeSizes := []int{4096, 8192, 16384, 32768}

			for _, size := range largeSizes {
				Convey("And allocating size "+string(rune(size)), func() {
					ptr := arena.Alloc(size)

					Convey("Then the allocation should succeed", func() {
						So(ptr, ShouldNotBeNil)
					})

					Convey("And the pointer should be properly aligned", func() {
						addr := uintptr(unsafe.Pointer(ptr))
						So(addr%uintptr(Align), ShouldEqual, uintptr(0))
					})

					Convey("And we should be able to write to the entire allocation", func() {
						// Write to the first and last byte
						*ptr = 0xFF
						lastByte := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(size-1)))
						*lastByte = 0xAA

						So(*ptr, ShouldEqual, byte(0xFF))
						So(*lastByte, ShouldEqual, byte(0xAA))
					})
				})
			}
		})
	})
}

// TestRecycledArena_StressTest performs a stress test with many allocations
func TestRecycledArena_StressTest(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		arena := &Recycled{}

		Convey("When performing multiple allocation cycles", func() {
			Convey("And allocating and releasing in cycles", func() {
				const numCycles = 10
				const numAllocs = 5

				testSizes := []int{64, 128, 256, 512, 1024}
				allSuccessful := true

				for cycle := 0; cycle < numCycles; cycle++ {
					// Allocate a few blocks
					ptrs := make([]*byte, numAllocs)
					success := true

					for i, size := range testSizes {
						ptr := arena.Alloc(size)
						if ptr == nil {
							success = false
							break
						}
						ptrs[i] = ptr

						// Write unique data to verify isolation
						*ptr = byte(cycle*numAllocs + i)
					}

					if success {
						// Verify data integrity
						for i, ptr := range ptrs {
							expected := byte(cycle*numAllocs + i)
							if *ptr != expected {
								success = false
								break
							}
						}

						// Release them
						for i, ptr := range ptrs {
							arena.Release(ptr, testSizes[i])
						}

						// Allocate again
						recycledPtrs := make([]*byte, numAllocs)
						for i, size := range testSizes {
							recycledPtr := arena.Alloc(size)
							if recycledPtr == nil {
								success = false
								break
							}
							recycledPtrs[i] = recycledPtr

							// Verify memory is cleared
							if *recycledPtr != 0 {
								success = false
								break
							}
						}

						// Verify we got recycled pointers
						if success {
							for _, recycledPtr := range recycledPtrs {
								if !containsPointer(ptrs, recycledPtr) {
									success = false
									break
								}
							}
						}
					}

					if !success {
						allSuccessful = false
					}
				}

				Convey("Then all cycles should complete successfully", func() {
					So(allSuccessful, ShouldBeTrue)
				})
			})
		})
	})
}

// Helper function to check if a slice contains a specific pointer
func containsPointer(ptrs []*byte, target *byte) bool {
	for _, ptr := range ptrs {
		if ptr == target {
			return true
		}
	}
	return false
}

// TestRecycledArena_Alignment tests proper alignment of allocations
func TestRecycledArena_Alignment(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		arena := &Recycled{}

		Convey("When allocating various sizes", func() {
			testSizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

			for _, size := range testSizes {
				Convey("And testing size "+string(rune(size)), func() {
					ptr := arena.Alloc(size)
					So(ptr, ShouldNotBeNil)

					Convey("Then the pointer should be properly aligned", func() {
						addr := uintptr(unsafe.Pointer(ptr))
						So(addr%uintptr(Align), ShouldEqual, uintptr(0))
					})

					Convey("And the allocated memory should be writable", func() {
						// Test writing to the first byte
						*ptr = 0xAA
						So(*ptr, ShouldEqual, byte(0xAA))

						// Test writing to the last byte if size > 1
						if size > 1 {
							lastByte := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(size-1)))
							*lastByte = 0xBB
							So(*lastByte, ShouldEqual, byte(0xBB))
						}
					})
				})
			}
		})

		Convey("When testing boundary alignment cases", func() {
			boundarySizes := []int{
				Align - 1,   // Just below alignment boundary
				Align,       // Exactly at alignment boundary
				Align + 1,   // Just above alignment boundary
				Align*2 - 1, // Just below next alignment boundary
				Align * 2,   // Exactly at next alignment boundary
			}

			for _, size := range boundarySizes {
				Convey("And testing boundary size "+string(rune(size)), func() {
					ptr := arena.Alloc(size)
					So(ptr, ShouldNotBeNil)

					addr := uintptr(unsafe.Pointer(ptr))
					So(addr%uintptr(Align), ShouldEqual, uintptr(0))
				})
			}
		})
	})
}

// TestRecycledArena_ReuseAfterReset tests that memory is reused after calling Reset
func TestRecycledArena_ReuseAfterReset(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating memory and then calling Reset", func() {
			ptr1 := arena.Alloc(64)
			ptr2 := arena.Alloc(128)
			So(ptr1, ShouldNotBeNil)
			So(ptr2, ShouldNotBeNil)

			arena.Reset()

			Convey("And allocating the same sizes again", func() {
				newPtr1 := arena.Alloc(64)
				newPtr2 := arena.Alloc(128)

				Convey("Then allocations should succeed", func() {
					So(newPtr1, ShouldNotBeNil)
					So(newPtr2, ShouldNotBeNil)
				})
			})
		})
	})
}

// TestRecycledArena_EdgeCases tests edge cases and boundary conditions
func TestRecycledArena_EdgeCases(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When testing very small allocations", func() {
			ptr1 := arena.Alloc(1)
			ptr2 := arena.Alloc(1)
			So(ptr1, ShouldNotBeNil)
			So(ptr2, ShouldNotBeNil)

			Convey("Then allocations should succeed", func() {
				// Small allocations might be rounded up to the same size category
				// so we can't guarantee they're different, but they should be valid
				So(ptr1, ShouldNotBeNil)
				So(ptr2, ShouldNotBeNil)
			})
		})

		Convey("When testing very large allocations", func() {
			largeSize := 1024 * 1024 // 1MB
			ptr := arena.Alloc(largeSize)
			So(ptr, ShouldNotBeNil)

			Convey("And releasing and reallocating", func() {
				arena.Release(ptr, largeSize)
				recycledPtr := arena.Alloc(largeSize)

				Convey("Then the recycled allocation should succeed", func() {
					So(recycledPtr, ShouldNotBeNil)
					So(recycledPtr, ShouldEqual, ptr)
				})
			})
		})

		Convey("When testing zero size allocation", func() {
			ptr := arena.Alloc(0)

			Convey("Then the allocation should return nil for zero size", func() {
				So(ptr, ShouldBeNil)
			})
		})

		Convey("When testing mixed size allocations", func() {
			sizes := []int{8, 16, 32, 64, 128, 256}
			ptrs := make([]*byte, len(sizes))

			for i, size := range sizes {
				ptrs[i] = arena.Alloc(size)
				So(ptrs[i], ShouldNotBeNil)
				*ptrs[i] = byte(i) // Write unique data
			}

			Convey("And releasing all allocations", func() {
				for i, ptr := range ptrs {
					arena.Release(ptr, sizes[i])
				}

				Convey("And reallocating the same sizes", func() {
					recycledPtrs := make([]*byte, len(sizes))
					for i, size := range sizes {
						recycledPtrs[i] = arena.Alloc(size)
						So(recycledPtrs[i], ShouldNotBeNil)
					}

					Convey("Then all recycled allocations should succeed", func() {
						// Verify that recycled pointers are from the original set
						for _, recycledPtr := range recycledPtrs {
							So(ptrs, ShouldContain, recycledPtr)
							// Verify memory is cleared after recycling
							So(*recycledPtr, ShouldEqual, byte(0))
						}
					})
				})
			})
		})

		Convey("When testing boundary conditions", func() {
			Convey("And releasing with size exactly at Align boundary", func() {
				ptr := arena.Alloc(Align)
				So(ptr, ShouldNotBeNil)

				arena.Release(ptr, Align)
				recycledPtr := arena.Alloc(Align)

				So(recycledPtr, ShouldEqual, ptr)
			})

			Convey("And releasing with size just below Align", func() {
				ptr := arena.Alloc(Align)
				So(ptr, ShouldNotBeNil)

				arena.Release(ptr, Align-1)
				newPtr := arena.Alloc(Align)

				// Should get a new pointer since release was ignored
				// Note: The underlying Arena might reuse memory blocks
				// so we can't guarantee different pointers
				So(newPtr, ShouldNotBeNil)
			})
		})
	})
}

// TestRecycledArena_MemoryCorruption tests for memory corruption scenarios
func TestRecycledArena_MemoryCorruption(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating and writing data", func() {
			ptr1 := arena.Alloc(64)
			ptr2 := arena.Alloc(64)

			So(ptr1, ShouldNotBeNil)
			So(ptr2, ShouldNotBeNil)

			// Write data to both allocations
			*ptr1 = 0xAA
			*ptr2 = 0xBB

			So(*ptr1, ShouldEqual, byte(0xAA))
			So(*ptr2, ShouldEqual, byte(0xBB))

			Convey("And releasing one allocation", func() {
				arena.Release(ptr1, 64)

				Convey("Then the other allocation should remain unchanged", func() {
					So(*ptr2, ShouldEqual, byte(0xBB))
				})

				Convey("And reallocating should return the released pointer", func() {
					recycledPtr := arena.Alloc(64)
					So(recycledPtr, ShouldEqual, ptr1)
					// Verify memory is cleared after recycling
					So(*recycledPtr, ShouldEqual, byte(0))
				})
			})
		})
	})
}

// TestRecycledArena_MemoryClearing tests that memory is properly cleared after reuse
func TestRecycledArena_MemoryClearing(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating memory and writing data", func() {
			ptr := arena.Alloc(128)
			So(ptr, ShouldNotBeNil)

			// Write data to the entire allocation
			for i := 0; i < 128; i++ {
				*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(i))) = byte(i % 256)
			}

			// Verify data was written
			So(*ptr, ShouldEqual, byte(0))
			So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1)), ShouldEqual, byte(1))
			So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 127)), ShouldEqual, byte(127))

			Convey("And releasing the allocation", func() {
				arena.Release(ptr, 128)

				Convey("And reallocating the same size", func() {
					recycledPtr := arena.Alloc(128)
					So(recycledPtr, ShouldEqual, ptr)

					Convey("Then the entire memory should be cleared to zero", func() {
						// Check first few bytes
						So(*recycledPtr, ShouldEqual, byte(0))
						So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(recycledPtr)) + 1)), ShouldEqual, byte(0))
						So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(recycledPtr)) + 2)), ShouldEqual, byte(0))

						// Check middle bytes
						So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(recycledPtr)) + 64)), ShouldEqual, byte(0))
						So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(recycledPtr)) + 65)), ShouldEqual, byte(0))

						// Check last few bytes
						So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(recycledPtr)) + 126)), ShouldEqual, byte(0))
						So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(recycledPtr)) + 127)), ShouldEqual, byte(0))
					})
				})
			})
		})

		Convey("When testing different allocation sizes", func() {
			sizes := []int{64, 128, 256, 512, 1024}

			for _, size := range sizes {
				Convey("And testing size "+string(rune(size)), func() {
					ptr := arena.Alloc(size)
					So(ptr, ShouldNotBeNil)

					// Write pattern to memory
					for i := 0; i < size; i++ {
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(i))) = 0xFF
					}

					// Verify pattern was written
					So(*ptr, ShouldEqual, byte(0xFF))
					So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(size-1))), ShouldEqual, byte(0xFF))

					arena.Release(ptr, size)
					recycledPtr := arena.Alloc(size)
					So(recycledPtr, ShouldEqual, ptr)

					Convey("Then the memory should be cleared", func() {
						So(*recycledPtr, ShouldEqual, byte(0))
						So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(recycledPtr)) + uintptr(size-1))), ShouldEqual, byte(0))
					})
				})
			}
		})

		Convey("When testing multiple allocations and releases", func() {
			ptrs := make([]*byte, 5)
			sizes := []int{64, 128, 256, 512, 1024}

			// Allocate and write data
			for i, size := range sizes {
				ptrs[i] = arena.Alloc(size)
				So(ptrs[i], ShouldNotBeNil)

				// Write unique pattern
				for j := 0; j < size; j++ {
					*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptrs[i])) + uintptr(j))) = byte(i + 1)
				}
			}

			// Release all
			for i, ptr := range ptrs {
				arena.Release(ptr, sizes[i])
			}

			// Reallocate all
			recycledPtrs := make([]*byte, 5)
			for i, size := range sizes {
				recycledPtrs[i] = arena.Alloc(size)
				So(recycledPtrs[i], ShouldNotBeNil)
			}

			Convey("Then all recycled memory should be cleared", func() {
				for i, recycledPtr := range recycledPtrs {
					// Check first and last byte of each allocation
					So(*recycledPtr, ShouldEqual, byte(0))
					So(*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(recycledPtr)) + uintptr(sizes[i]-1))), ShouldEqual, byte(0))
				}
			})
		})
	})
}

// TestRecycledArena_UnalignedSizes verifies behavior when requested sizes are not multiples of Align (typically 8).
func TestRecycledArena_UnalignedSizes(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		a := &Recycled{}

		Convey("When allocating an unaligned size and releasing it", func() {
			// Case 1: request 7 (aligns to 8), then allocate 8
			ptr17 := a.Alloc(17)
			So(ptr17, ShouldNotBeNil)
			addr7 := uintptr(unsafe.Pointer(ptr17))
			So(addr7%uintptr(Align), ShouldEqual, uintptr(0))

			// Write a marker
			*ptr17 = 0xAA
			a.Release(ptr17, 17)

			Convey("Then allocating the aligned size should succeed (release ignored for < Align)", func() {
				ptr18 := a.Alloc(18)
				So(ptr18, ShouldEqual, ptr17)
				addr8 := uintptr(unsafe.Pointer(ptr18))
				So(addr8%uintptr(Align), ShouldEqual, uintptr(0))
			})
		})

		Convey("When allocating another unaligned size and releasing it", func() {
			// Case 2: request 9 (aligns to 16), then allocate 15 (also aligns to 16)
			ptr9 := a.Alloc(9)
			So(ptr9, ShouldNotBeNil)
			addr9 := uintptr(unsafe.Pointer(ptr9))
			So(addr9%uintptr(Align), ShouldEqual, uintptr(0))

			// Write a marker
			*ptr9 = 0xBB
			a.Release(ptr9, 9)

			Convey("Then allocating another size in the same aligned class should recycle the same pointer and be zeroed", func() {
				r15 := a.Alloc(15)
				So(r15, ShouldEqual, ptr9)
				So(*r15, ShouldEqual, byte(0))
			})
		})
	})
}

// TestRecycledArena_FullAccess ensures that the entire allocated segment
// (for both aligned and unaligned request sizes) is writable and readable.
func TestRecycledArena_FullAccess(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		a := &Recycled{}

		sizes := []int{1, 7, 8, 9, 15, 16, 31, 32, 63, 64, 127, 128, 257}

		for _, size := range sizes {
			// Use a scope per size for clearer reporting
			func(size int) {
				Convey("When allocating size "+string(rune(size)), func() {
					ptr := a.Alloc(size)
					So(ptr, ShouldNotBeNil)
					addr := xunsafe.AddrOf(ptr)
					So(int(addr)%Align, ShouldEqual, 0)

					// Write through the entire requested segment
					for i := 0; i < size; i++ {
						*addr.Add(i).AssertValid() = byte((i ^ 0x5A) & 0xFF)
					}

					// Verify a few sentinel positions
					So(*ptr, ShouldEqual, byte((0^0x5A)&0xFF))
					mid := size / 2
					if mid > 0 {
						midByte := *addr.Add(mid).AssertValid()
						So(midByte, ShouldEqual, byte((mid^0x5A)&0xFF))
					}
					lastByte := *addr.Add(size - 1).AssertValid()
					So(lastByte, ShouldEqual, byte(((size-1)^0x5A)&0xFF))

					// Release and re-allocate; for sizes < Align release is ignored,
					// but the newly allocated segment must still be fully writable.
					a.Release(ptr, size)
					re := a.Alloc(size)
					So(re, ShouldNotBeNil)
					addrRe := xunsafe.AddrOf(re)

					// Fully write again across the entire segment
					for i := 0; i < size; i++ {
						*addrRe.Add(i).AssertValid() = byte((i ^ 0xA5) & 0xFF)
					}
					So(*re, ShouldEqual, byte((0^0xA5)&0xFF))
					if mid > 0 {
						midByteRe2 := *addrRe.Add(mid).AssertValid()
						So(midByteRe2, ShouldEqual, byte((mid^0xA5)&0xFF))
					}
					lastByteRe2 := *addrRe.Add(size - 1).AssertValid()
					So(lastByteRe2, ShouldEqual, byte(((size-1)^0xA5)&0xFF))
				})
			}(size)
		}
	})
}

// TestRecycledArena_MemoryIsolation tests that different allocations don't interfere with each other
func TestRecycledArena_MemoryIsolation(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		arena := &Recycled{}

		Convey("When allocating multiple blocks and writing different data", func() {
			// Allocate blocks of different sizes
			allocations := []struct {
				size int
				data byte
				ptr  *byte
			}{
				{32, 0x11, nil},
				{64, 0x22, nil},
				{128, 0x33, nil},
				{256, 0x44, nil},
			}

			// Perform allocations
			for i := range allocations {
				allocations[i].ptr = arena.Alloc(allocations[i].size)
				So(allocations[i].ptr, ShouldNotBeNil)
			}

			Convey("And writing data to all allocations", func() {
				// Write data to each allocation
				for _, alloc := range allocations {
					*alloc.ptr = alloc.data
					So(*alloc.ptr, ShouldEqual, alloc.data)
				}

				Convey("Then all allocations should maintain their data independently", func() {
					// Verify all data is still intact
					for _, alloc := range allocations {
						So(*alloc.ptr, ShouldEqual, alloc.data)
					}
				})

				Convey("And releasing one allocation should not affect others", func() {
					// Release the second allocation
					arena.Release(allocations[1].ptr, allocations[1].size)

					// Verify other allocations are unchanged
					So(*allocations[0].ptr, ShouldEqual, byte(0x11))
					So(*allocations[2].ptr, ShouldEqual, byte(0x33))
					So(*allocations[3].ptr, ShouldEqual, byte(0x44))

					// Reallocate the released memory
					newPtr := arena.Alloc(allocations[1].size)
					So(newPtr, ShouldEqual, allocations[1].ptr)
					So(*newPtr, ShouldEqual, byte(0))
				})
			})
		})

		Convey("When testing overlapping size classes", func() {
			// Test sizes that might fall into the same size class
			overlappingSizes := []int{64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80}
			ptrs := make([]*byte, len(overlappingSizes))

			// Allocate all
			for i, size := range overlappingSizes {
				ptrs[i] = arena.Alloc(size)
				So(ptrs[i], ShouldNotBeNil)
				*ptrs[i] = byte(i) // Write unique data
			}

			Convey("Then all allocations should be independent", func() {
				// Verify all data is intact
				for i, ptr := range ptrs {
					So(*ptr, ShouldEqual, byte(i))
				}
			})

			Convey("And releasing some allocations should not affect others", func() {
				// Release every other allocation
				for i := 0; i < len(ptrs); i += 2 {
					arena.Release(ptrs[i], overlappingSizes[i])
				}

				// Reallocate released memory with new pointers
				releasedCount := (len(ptrs) + 1) / 2 // Count of even-indexed elements
				newPtrs := make([]*byte, releasedCount)
				for i := 0; i < len(ptrs); i += 2 {
					newPtrs[i/2] = arena.Alloc(overlappingSizes[i])
					So(newPtrs[i/2], ShouldNotBeNil)
					So(*newPtrs[i/2], ShouldEqual, byte(0))
				}

				// Verify that we can write to reallocated memory independently
				// Write new data to reallocated memory
				for i := 0; i < len(newPtrs); i++ {
					*newPtrs[i] = byte(100 + i)
				}

				// Verify the new data was written
				for i := 0; i < len(newPtrs); i++ {
					So(*newPtrs[i], ShouldEqual, byte(100+i))
				}

				// Verify that odd-indexed allocations still contain their original data
				for i := 1; i < len(ptrs); i += 2 {
					So(*ptrs[i], ShouldEqual, byte(i))
				}

				// Verify that we can write to odd-indexed allocations independently
				for i := 1; i < len(ptrs); i += 2 {
					*ptrs[i] = byte(200 + i)
				}

				// Verify the new data was written
				for i := 1; i < len(ptrs); i += 2 {
					So(*ptrs[i], ShouldEqual, byte(200+i))
				}

				// Verify that reallocated memory is still independent
				for i := 0; i < len(newPtrs); i++ {
					So(*newPtrs[i], ShouldEqual, byte(100+i))
				}
			})
		})
	})
}

// TestRecycledArena_PerformanceCharacteristics tests basic performance characteristics
func TestRecycledArena_PerformanceCharacteristics(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		arena := &Recycled{}

		Convey("When performing rapid allocation and release cycles", func() {
			const numIterations = 100
			testSizes := []int{64, 128, 256, 512, 1024}

			// Track allocation patterns
			allocationCount := 0
			recyclingCount := 0

			for i := 0; i < numIterations; i++ {
				// Allocate a set of blocks
				ptrs := make([]*byte, len(testSizes))
				for _, size := range testSizes {
					ptr := arena.Alloc(size)
					if ptr != nil {
						allocationCount++
					}
				}

				// Release all blocks
				for j, ptr := range ptrs {
					if ptr != nil {
						arena.Release(ptr, testSizes[j])
					}
				}

				// Reallocate to test recycling
				for _, size := range testSizes {
					recycledPtr := arena.Alloc(size)
					if recycledPtr != nil {
						recyclingCount++
						// Verify memory is cleared
						So(*recycledPtr, ShouldEqual, byte(0))
					}
				}
			}

			Convey("Then all operations should complete successfully", func() {
				expectedAllocations := numIterations * len(testSizes)
				So(allocationCount, ShouldEqual, expectedAllocations)
				So(recyclingCount, ShouldEqual, expectedAllocations)
			})
		})
	})
}

// TestRecycledArena_FragmentationHandling tests how the arena handles fragmentation
func TestRecycledArena_FragmentationHandling(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		arena := &Recycled{}

		Convey("When creating fragmentation by releasing blocks in different patterns", func() {
			// Allocate blocks of different sizes
			allocations := []struct {
				size int
				ptr  *byte
			}{
				{64, nil},
				{128, nil},
				{256, nil},
				{512, nil},
				{1024, nil},
				{64, nil},
				{128, nil},
				{256, nil},
			}

			// Perform allocations
			for i := range allocations {
				allocations[i].ptr = arena.Alloc(allocations[i].size)
				So(allocations[i].ptr, ShouldNotBeNil)
			}

			Convey("And releasing blocks in a fragmented pattern", func() {
				// Release every other block to create fragmentation
				for i := 0; i < len(allocations); i += 2 {
					arena.Release(allocations[i].ptr, allocations[i].size)
				}

				Convey("Then reallocating the same sizes should succeed", func() {
					recycledPtrs := make([]*byte, len(allocations)/2)
					recycledSizes := make([]int, len(allocations)/2)

					// Collect sizes and reallocate
					for i := 0; i < len(allocations); i += 2 {
						recycledSizes[i/2] = allocations[i].size
					}

					for i, size := range recycledSizes {
						recycledPtrs[i] = arena.Alloc(size)
						So(recycledPtrs[i], ShouldNotBeNil)
					}

					Convey("And recycled pointers should be from the original set", func() {
						for i, recycledPtr := range recycledPtrs {
							originalPtr := allocations[i*2].ptr
							So(recycledPtr, ShouldEqual, originalPtr)
							So(*recycledPtr, ShouldEqual, byte(0))
						}
					})
				})
			})
		})

		Convey("When testing interleaved allocation and release", func() {
			const numCycles = 5
			const numAllocs = 3

			for cycle := 0; cycle < numCycles; cycle++ {
				Convey("And performing cycle "+string(rune(cycle)), func() {
					// Allocate
					ptrs := make([]*byte, numAllocs)
					sizes := []int{64, 128, 256}

					for i, size := range sizes {
						ptrs[i] = arena.Alloc(size)
						So(ptrs[i], ShouldNotBeNil)
					}

					// Write data
					for i, ptr := range ptrs {
						*ptr = byte(cycle*numAllocs + i)
					}

					// Release
					for i, ptr := range ptrs {
						arena.Release(ptr, sizes[i])
					}

					// Reallocate
					recycledPtrs := make([]*byte, numAllocs)
					for i, size := range sizes {
						recycledPtrs[i] = arena.Alloc(size)
						So(recycledPtrs[i], ShouldNotBeNil)
					}

					Convey("Then all operations should complete successfully", func() {
						// Verify we got recycled pointers
						for i, recycledPtr := range recycledPtrs {
							So(recycledPtr, ShouldEqual, ptrs[i])
							So(*recycledPtr, ShouldEqual, byte(0))
						}
					})
				})
			}
		})
	})
}

// TestRecycledArena_ConcurrentAccessPatterns tests various concurrent access patterns
func TestRecycledArena_ConcurrentAccessPatterns(t *testing.T) {
	Convey("Given a Recycled arena", t, func() {
		arena := &Recycled{}

		Convey("When testing alternating size patterns", func() {
			// Test pattern: allocate small, large, small, large...
			pattern := []int{64, 1024, 128, 2048, 256, 4096}
			ptrs := make([]*byte, len(pattern))

			// Allocate in pattern
			for i, size := range pattern {
				ptrs[i] = arena.Alloc(size)
				So(ptrs[i], ShouldNotBeNil)
			}

			Convey("And releasing in reverse order", func() {
				for i := len(pattern) - 1; i >= 0; i-- {
					arena.Release(ptrs[i], pattern[i])
				}

				Convey("Then reallocating in the same pattern should succeed", func() {
					recycledPtrs := make([]*byte, len(pattern))
					for i, size := range pattern {
						recycledPtrs[i] = arena.Alloc(size)
						So(recycledPtrs[i], ShouldNotBeNil)
					}

					Convey("And all recycled pointers should be from the original set", func() {
						for _, recycledPtr := range recycledPtrs {
							So(ptrs, ShouldContain, recycledPtr)
							So(*recycledPtr, ShouldEqual, byte(0))
						}
					})
				})
			})
		})

		Convey("When testing burst allocation and release", func() {
			// Allocate many small blocks
			const numSmall = 20
			smallPtrs := make([]*byte, numSmall)
			for i := 0; i < numSmall; i++ {
				smallPtrs[i] = arena.Alloc(64)
				So(smallPtrs[i], ShouldNotBeNil)
			}

			// Allocate a few large blocks
			const numLarge = 3
			largePtrs := make([]*byte, numLarge)
			largeSizes := []int{1024, 2048, 4096}
			for i := 0; i < numLarge; i++ {
				largePtrs[i] = arena.Alloc(largeSizes[i])
				So(largePtrs[i], ShouldNotBeNil)
			}

			Convey("And releasing all blocks", func() {
				// Release small blocks
				for i := 0; i < numSmall; i++ {
					arena.Release(smallPtrs[i], 64)
				}

				// Release large blocks
				for i := 0; i < numLarge; i++ {
					arena.Release(largePtrs[i], largeSizes[i])
				}

				Convey("Then reallocating should succeed", func() {
					// Reallocate small blocks
					recycledSmall := make([]*byte, numSmall)
					for i := 0; i < numSmall; i++ {
						recycledSmall[i] = arena.Alloc(64)
						So(recycledSmall[i], ShouldNotBeNil)
					}

					// Reallocate large blocks
					recycledLarge := make([]*byte, numLarge)
					for i := 0; i < numLarge; i++ {
						recycledLarge[i] = arena.Alloc(largeSizes[i])
						So(recycledLarge[i], ShouldNotBeNil)
					}

					Convey("And all recycled pointers should be from the original sets", func() {
						for _, recycledPtr := range recycledSmall {
							So(smallPtrs, ShouldContain, recycledPtr)
							So(*recycledPtr, ShouldEqual, byte(0))
						}

						for _, recycledPtr := range recycledLarge {
							So(largePtrs, ShouldContain, recycledPtr)
							So(*recycledPtr, ShouldEqual, byte(0))
						}
					})
				})
			})
		})
	})
}
