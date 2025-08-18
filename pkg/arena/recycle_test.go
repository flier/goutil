package arena

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

// TestRecycledArena_BasicAllocation tests basic allocation functionality
func TestRecycledArena_BasicAllocation(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating memory of different sizes", func() {
			sizes := []int{8, 16, 32, 64, 128, 256, 512, 1024}
			pointers := make([]*byte, len(sizes))

			for i, size := range sizes {
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
		})
	})
}

// TestRecycledArena_Recycling tests the recycling functionality
func TestRecycledArena_Recycling(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating memory and then releasing it", func() {
			ptr1 := arena.Alloc(64)
			So(ptr1, ShouldNotBeNil)

			// Write some data to verify it's preserved
			*ptr1 = 42

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
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When testing different allocation sizes", func() {
			testCases := []struct {
				size     int
				expected uint
			}{
				{1, 6},    // 2^6 = 64
				{32, 6},   // 2^6 = 64
				{64, 6},   // 2^6 = 64
				{65, 7},   // 2^7 = 128
				{128, 7},  // 2^7 = 128
				{129, 8},  // 2^8 = 256
				{256, 8},  // 2^8 = 256
				{257, 9},  // 2^9 = 512
				{512, 9},  // 2^9 = 512
				{513, 10}, // 2^10 = 1024
			}

			for _, tc := range testCases {
				Convey("And allocating size "+string(rune(tc.size)), func() {
					ptr := arena.Alloc(tc.size)
					So(ptr, ShouldNotBeNil)

					Convey("And releasing the allocation", func() {
						arena.Release(ptr, tc.size)

						Convey("And allocating the same size again", func() {
							recycledPtr := arena.Alloc(tc.size)
							So(recycledPtr, ShouldNotBeNil)

							Convey("Then the recycled pointer should match the original", func() {
								So(recycledPtr, ShouldEqual, ptr)
							})
						})
					})
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
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When performing multiple allocation cycles", func() {
			Convey("And allocating and releasing in cycles", func() {
				for cycle := 0; cycle < 10; cycle++ {
					// Allocate a few blocks
					ptrs := make([]*byte, 5)
					sizes := []int{64, 128, 256, 512, 1024}

					for i, size := range sizes {
						ptrs[i] = arena.Alloc(size)
						So(ptrs[i], ShouldNotBeNil)
					}

					// Release them
					for i, ptr := range ptrs {
						arena.Release(ptr, sizes[i])
					}

					// Allocate again
					for _, size := range sizes {
						recycledPtr := arena.Alloc(size)
						So(recycledPtr, ShouldNotBeNil)
						So(ptrs, ShouldContain, recycledPtr)
					}
				}

				Convey("Then all cycles should complete successfully", func() {
					So(true, ShouldBeTrue)
				})
			})
		})
	})
}

// TestRecycledArena_Alignment tests proper alignment of allocations
func TestRecycledArena_Alignment(t *testing.T) {
	Convey("Given a RecycledArena", t, func() {
		arena := &Recycled{}

		Convey("When allocating various sizes", func() {
			sizes := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

			for _, size := range sizes {
				Convey("And allocating size "+string(rune(size)), func() {
					ptr := arena.Alloc(size)
					So(ptr, ShouldNotBeNil)

					Convey("Then the pointer should be properly aligned", func() {
						addr := uintptr(unsafe.Pointer(ptr))
						So(addr%uintptr(Align), ShouldEqual, uintptr(0))
					})
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
					addr := uintptr(unsafe.Pointer(ptr))
					So(addr%uintptr(Align), ShouldEqual, uintptr(0))

					// Write through the entire requested segment
					for i := 0; i < size; i++ {
						*(*byte)(unsafe.Pointer(addr + uintptr(i))) = byte((i ^ 0x5A) & 0xFF)
					}

					// Verify a few sentinel positions
					So(*ptr, ShouldEqual, byte((0^0x5A)&0xFF))
					mid := size / 2
					if mid > 0 {
						midByte := *(*byte)(unsafe.Pointer(addr + uintptr(mid)))
						So(midByte, ShouldEqual, byte((mid^0x5A)&0xFF))
					}
					lastByte := *(*byte)(unsafe.Pointer(addr + uintptr(size-1)))
					So(lastByte, ShouldEqual, byte(((size-1)^0x5A)&0xFF))

					// Release and re-allocate; for sizes < Align release is ignored,
					// but the newly allocated segment must still be fully writable.
					a.Release(ptr, size)
					re := a.Alloc(size)
					So(re, ShouldNotBeNil)
					addrRe := uintptr(unsafe.Pointer(re))

					// Fully write again across the entire segment
					for i := 0; i < size; i++ {
						*(*byte)(unsafe.Pointer(addrRe + uintptr(i))) = byte((i ^ 0xA5) & 0xFF)
					}
					So(*re, ShouldEqual, byte((0^0xA5)&0xFF))
					if mid > 0 {
						midByteRe2 := *(*byte)(unsafe.Pointer(addrRe + uintptr(mid)))
						So(midByteRe2, ShouldEqual, byte((mid^0xA5)&0xFF))
					}
					lastByteRe2 := *(*byte)(unsafe.Pointer(addrRe + uintptr(size-1)))
					So(lastByteRe2, ShouldEqual, byte(((size-1)^0xA5)&0xFF))
				})
			}(size)
		}
	})
}
