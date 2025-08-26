package xunsafe_test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestPointer(t *testing.T) {
	Convey("Given pointer operations", t, func() {
		Convey("When working with pointer casting", func() {
			Convey("And casting between different pointer types", func() {
				i := 42
				ptr := &i

				// Cast to uintptr
				uintptrPtr := xunsafe.Cast[uintptr, int](ptr)
				So(uintptrPtr, ShouldNotBeNil)

				// Cast to byte pointer
				bytePtr := xunsafe.Cast[byte, int](ptr)
				So(bytePtr, ShouldNotBeNil)

				// Cast back to int pointer
				intPtr := xunsafe.Cast[int, byte](bytePtr)
				So(intPtr, ShouldNotBeNil)
			})
		})

		Convey("When working with pointer arithmetic", func() {
			Convey("And adding offset to pointer", func() {
				arr := [5]int{1, 2, 3, 4, 5}
				basePtr := &arr[0]

				// Add offset to get pointer to arr[2]
				ptr2 := xunsafe.Add(basePtr, 2)
				So(*ptr2, ShouldEqual, 3)

				// Add offset to get pointer to arr[4]
				ptr4 := xunsafe.Add(basePtr, 4)
				So(*ptr4, ShouldEqual, 5)

				// Test with zero offset
				ptr0 := xunsafe.Add(basePtr, 0)
				So(*ptr0, ShouldEqual, 1)
			})

			Convey("And subtracting pointers", func() {
				arr := [5]int{1, 2, 3, 4, 5}
				basePtr := &arr[0]
				ptr2 := &arr[2]
				ptr4 := &arr[4]

				// Calculate difference
				diff := xunsafe.Sub(ptr4, ptr2)
				So(diff, ShouldEqual, 2)

				// Test with same pointer
				sameDiff := xunsafe.Sub(ptr2, ptr2)
				So(sameDiff, ShouldEqual, 0)

				// Test with base pointer
				baseDiff := xunsafe.Sub(ptr2, basePtr)
				So(baseDiff, ShouldEqual, 2)
			})
		})

		Convey("When working with pointer loading", func() {
			Convey("And loading values at different offsets", func() {
				arr := [5]int{1, 2, 3, 4, 5}
				basePtr := &arr[0]

				// Load at offset 0
				val0 := xunsafe.Load(basePtr, 0)
				So(val0, ShouldEqual, 1)

				// Load at offset 2
				val2 := xunsafe.Load(basePtr, 2)
				So(val2, ShouldEqual, 3)

				// Load at offset 4
				val4 := xunsafe.Load(basePtr, 4)
				So(val4, ShouldEqual, 5)
			})
		})

		Convey("When working with pointer storing", func() {
			Convey("And storing values at different offsets", func() {
				arr := [5]int{1, 2, 3, 4, 5}
				basePtr := &arr[0]

				// Store at offset 0
				xunsafe.Store(basePtr, 0, 100)
				So(arr[0], ShouldEqual, 100)

				// Store at offset 2
				xunsafe.Store(basePtr, 2, 300)
				So(arr[2], ShouldEqual, 300)

				// Store at offset 4
				xunsafe.Store(basePtr, 4, 500)
				So(arr[4], ShouldEqual, 500)

				// Verify other elements unchanged
				So(arr[1], ShouldEqual, 2)
				So(arr[3], ShouldEqual, 4)
			})
		})

		Convey("When working with write barrier operations", func() {
			Convey("And storing without write barriers", func() {
				var ptr *int
				var newPtr = new(int)
				*newPtr = 42

				// Store the pointer
				xunsafe.StoreNoWB(&ptr, newPtr)
				So(ptr, ShouldEqual, newPtr)
				So(*ptr, ShouldEqual, 42)
			})

			Convey("And storing untyped pointer without write barriers", func() {
				var ptr unsafe.Pointer
				var newPtr = unsafe.Pointer(new(int))

				// Store the pointer
				xunsafe.StoreNoWBUntyped(&ptr, newPtr)
				So(ptr, ShouldEqual, newPtr)
			})
		})

		Convey("When working with memory operations", func() {
			Convey("And copying elements between arrays", func() {
				src := [5]int{1, 2, 3, 4, 5}
				dst := [5]int{0, 0, 0, 0, 0}

				// Copy all elements
				xunsafe.Copy(&dst[0], &src[0], 5)
				So(dst, ShouldEqual, src)

				// Copy partial elements
				dst2 := [5]int{0, 0, 0, 0, 0}
				xunsafe.Copy(&dst2[0], &src[0], 3)
				So(dst2[0], ShouldEqual, 1)
				So(dst2[1], ShouldEqual, 2)
				So(dst2[2], ShouldEqual, 3)
				So(dst2[3], ShouldEqual, 0)
				So(dst2[4], ShouldEqual, 0)
			})

			Convey("And clearing elements", func() {
				arr := [5]int{1, 2, 3, 4, 5}

				// Clear first 3 elements
				xunsafe.Clear(&arr[0], 3)
				So(arr[0], ShouldEqual, 0)
				So(arr[1], ShouldEqual, 0)
				So(arr[2], ShouldEqual, 0)
				So(arr[3], ShouldEqual, 4)
				So(arr[4], ShouldEqual, 5)

				// Clear all elements
				xunsafe.Clear(&arr[0], 5)
				So(arr[0], ShouldEqual, 0)
				So(arr[1], ShouldEqual, 0)
				So(arr[2], ShouldEqual, 0)
				So(arr[3], ShouldEqual, 0)
				So(arr[4], ShouldEqual, 0)
			})
		})

		Convey("When working with edge cases", func() {
			Convey("And working with zero offset", func() {
				arr := [1]int{42}
				ptr := &arr[0]
				val := xunsafe.Load(ptr, 0)
				So(val, ShouldEqual, 42)
			})
		})

		Convey("When working with different types", func() {
			Convey("And working with string arrays", func() {
				arr := [3]string{"hello", "world", "test"}
				basePtr := &arr[0]

				// Test string operations
				val0 := xunsafe.Load(basePtr, 0)
				So(val0, ShouldEqual, "hello")

				val1 := xunsafe.Load(basePtr, 1)
				So(val1, ShouldEqual, "world")

				// Store new values
				xunsafe.Store(basePtr, 0, "hi")
				xunsafe.Store(basePtr, 1, "there")
				So(arr[0], ShouldEqual, "hi")
				So(arr[1], ShouldEqual, "there")
			})
		})

		Convey("When working with comprehensive type tests", func() {
			Convey("And testing various numeric types", func() {
				testCases := []struct {
					name  string
					array interface{}
					index int
					value interface{}
				}{
					{"int8", [3]int8{1, 2, 3}, 1, int8(2)},
					{"int16", [3]int16{10, 20, 30}, 2, int16(30)},
					{"int32", [3]int32{100, 200, 300}, 0, int32(100)},
					{"int64", [3]int64{1000, 2000, 3000}, 1, int64(2000)},
					{"uint8", [3]uint8{1, 2, 3}, 2, uint8(3)},
					{"uint16", [3]uint16{10, 20, 30}, 0, uint16(10)},
					{"uint32", [3]uint32{100, 200, 300}, 1, uint32(200)},
					{"uint64", [3]uint64{1000, 2000, 3000}, 2, uint64(3000)},
					{"float32", [3]float32{1.1, 2.2, 3.3}, 0, float32(1.1)},
					{"float64", [3]float64{1.1, 2.2, 3.3}, 1, float64(2.2)},
				}

				for _, tc := range testCases {
					Convey(tc.name, func() {
						// Test Load operation
						switch v := tc.array.(type) {
						case [3]int8:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						case [3]int16:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						case [3]int32:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						case [3]int64:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						case [3]uint8:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						case [3]uint16:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						case [3]uint32:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						case [3]uint64:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						case [3]float32:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						case [3]float64:
							val := xunsafe.Load(&v[0], tc.index)
							So(val, ShouldEqual, tc.value)
						}
					})
				}
			})
		})
	})
}
