package xunsafe_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestByteAdd(t *testing.T) {
	Convey("Given byte addition operations", t, func() {
		// Test adding byte offset to pointer
		arr := [5]int{1, 2, 3, 4, 5}
		basePtr := &arr[0]

		// Add byte offset to get pointer to arr[1]
		ptr1 := xunsafe.ByteAdd[int](basePtr, 8) // Assuming int is 8 bytes
		So(*ptr1, ShouldEqual, 2)

		// Add byte offset to get pointer to arr[2]
		ptr2 := xunsafe.ByteAdd[int](basePtr, 16) // 2 * 8 bytes
		So(*ptr2, ShouldEqual, 3)

		// Test with zero offset
		ptr0 := xunsafe.ByteAdd[int](basePtr, 0)
		So(*ptr0, ShouldEqual, 1)
	})
}

func TestByteSub(t *testing.T) {
	Convey("Given byte subtraction operations", t, func() {
		// Test subtracting pointers
		arr := [5]int{1, 2, 3, 4, 5}
		basePtr := &arr[0]
		ptr2 := &arr[2]
		ptr4 := &arr[4]

		// Calculate byte difference
		diff := xunsafe.ByteSub(ptr4, ptr2)
		So(diff, ShouldEqual, 16) // 2 * 8 bytes

		// Test with same pointer
		sameDiff := xunsafe.ByteSub(ptr2, ptr2)
		So(sameDiff, ShouldEqual, 0)

		// Test with base pointer
		baseDiff := xunsafe.ByteSub(ptr2, basePtr)
		So(baseDiff, ShouldEqual, 16) // 2 * 8 bytes
	})
}

func TestByteLoad(t *testing.T) {
	Convey("Given byte load operations", t, func() {
		// Test loading values at different byte offsets
		arr := [5]int{1, 2, 3, 4, 5}
		basePtr := &arr[0]

		// Load at byte offset 0
		val0 := xunsafe.ByteLoad[int](basePtr, 0)
		So(val0, ShouldEqual, 1)

		// Load at byte offset 8
		val1 := xunsafe.ByteLoad[int](basePtr, 8)
		So(val1, ShouldEqual, 2)

		// Load at byte offset 16
		val2 := xunsafe.ByteLoad[int](basePtr, 16)
		So(val2, ShouldEqual, 3)
	})
}

func TestByteStore(t *testing.T) {
	Convey("Given byte store operations", t, func() {
		// Test storing values at different byte offsets
		arr := [5]int{1, 2, 3, 4, 5}
		basePtr := &arr[0]

		// Store at byte offset 0
		xunsafe.ByteStore(basePtr, 0, 100)
		So(arr[0], ShouldEqual, 100)

		// Store at byte offset 8
		xunsafe.ByteStore(basePtr, 8, 200)
		So(arr[1], ShouldEqual, 200)

		// Store at byte offset 16
		xunsafe.ByteStore(basePtr, 16, 300)
		So(arr[2], ShouldEqual, 300)

		// Verify other elements unchanged
		So(arr[3], ShouldEqual, 4)
		So(arr[4], ShouldEqual, 5)
	})
}

func TestUntypedEdgeCases(t *testing.T) {
	Convey("Given edge cases", t, func() {
		// Test with nil pointers
		var nilPtr *int
		So(func() {
			xunsafe.ByteLoad[int](nilPtr, 0)
		}, ShouldPanic)

		So(func() {
			xunsafe.ByteStore(nilPtr, 0, 42)
		}, ShouldPanic)

		// Test with zero offset
		arr := [1]int{42}
		ptr := &arr[0]
		val := xunsafe.ByteLoad[int](ptr, 0)
		So(val, ShouldEqual, 42)
	})
}

func TestUntypedTypes(t *testing.T) {
	Convey("Given different types", t, func() {
		// Test with different types
		arr := [3]string{"hello", "world", "test"}
		basePtr := &arr[0]

		// Test string operations
		val0 := xunsafe.ByteLoad[string](basePtr, 0)
		So(val0, ShouldEqual, "hello")

		val1 := xunsafe.ByteLoad[string](basePtr, 16) // Assuming string is 16 bytes
		So(val1, ShouldEqual, "world")

		// Store new values
		xunsafe.ByteStore(basePtr, 0, "hi")
		xunsafe.ByteStore(basePtr, 16, "there")
		So(arr[0], ShouldEqual, "hi")
		So(arr[1], ShouldEqual, "there")
	})
}

func TestUntypedAlignment(t *testing.T) {
	Convey("Given alignment requirements", t, func() {
		// Test with different alignment requirements
		arr := [3]int32{1, 2, 3} // int32 is 4 bytes
		basePtr := &arr[0]

		// Test with 4-byte alignment
		val0 := xunsafe.ByteLoad[int32](basePtr, 0)
		So(val0, ShouldEqual, int32(1))

		val1 := xunsafe.ByteLoad[int32](basePtr, 4)
		So(val1, ShouldEqual, int32(2))

		val2 := xunsafe.ByteLoad[int32](basePtr, 8)
		So(val2, ShouldEqual, int32(3))
	})
}

func TestUntypedComprehensive(t *testing.T) {
	Convey("Given comprehensive untyped tests", t, func() {
		// Test with various types and byte offsets
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
			{"bool", [3]bool{true, false, true}, 1, false},
			{"string", [3]string{"hello", "world", "test"}, 2, "test"},
		}

		for _, tc := range testCases {
			Convey("When testing "+tc.name, func() {
				// Test ByteLoad operation
				switch v := tc.array.(type) {
				case [3]int8:
					val := xunsafe.ByteLoad[int8](&v[0], tc.index)
					So(val, ShouldEqual, tc.value)
				case [3]int16:
					val := xunsafe.ByteLoad[int16](&v[0], tc.index*2)
					So(val, ShouldEqual, tc.value)
				case [3]int32:
					val := xunsafe.ByteLoad[int32](&v[0], tc.index*4)
					So(val, ShouldEqual, tc.value)
				case [3]int64:
					val := xunsafe.ByteLoad[int64](&v[0], tc.index*8)
					So(val, ShouldEqual, tc.value)
				case [3]uint8:
					val := xunsafe.ByteLoad[uint8](&v[0], tc.index)
					So(val, ShouldEqual, tc.value)
				case [3]uint16:
					val := xunsafe.ByteLoad[uint16](&v[0], tc.index*2)
					So(val, ShouldEqual, tc.value)
				case [3]uint32:
					val := xunsafe.ByteLoad[uint32](&v[0], tc.index*4)
					So(val, ShouldEqual, tc.value)
				case [3]uint64:
					val := xunsafe.ByteLoad[uint64](&v[0], tc.index*8)
					So(val, ShouldEqual, tc.value)
				case [3]float32:
					val := xunsafe.ByteLoad[float32](&v[0], tc.index*4)
					So(val, ShouldEqual, tc.value)
				case [3]float64:
					val := xunsafe.ByteLoad[float64](&v[0], tc.index*8)
					So(val, ShouldEqual, tc.value)
				case [3]bool:
					val := xunsafe.ByteLoad[bool](&v[0], tc.index)
					So(val, ShouldEqual, tc.value)
				case [3]string:
					val := xunsafe.ByteLoad[string](&v[0], tc.index*16)
					So(val, ShouldEqual, tc.value)
				}
			})
		}
	})
}
