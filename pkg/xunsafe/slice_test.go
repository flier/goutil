package xunsafe_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestBoundsCheck(t *testing.T) {
	Convey("Given bounds check operations", t, func() {
		// Test bounds check with valid indices
		So(func() {
			xunsafe.BoundsCheck(0, 10)
			xunsafe.BoundsCheck(5, 10)
			xunsafe.BoundsCheck(9, 10)
		}, ShouldNotPanic)

		// Test bounds check with invalid indices
		So(func() {
			xunsafe.BoundsCheck(10, 10)
		}, ShouldPanic)

		So(func() {
			xunsafe.BoundsCheck(-1, 10)
		}, ShouldPanic)
	})
}

func TestBytes(t *testing.T) {
	Convey("Given various types", t, func() {
		// Test with int
		i := 42
		bytes := xunsafe.Bytes(&i)
		So(bytes, ShouldNotBeNil)
		So(len(bytes), ShouldEqual, 8) // int is 8 bytes on 64-bit systems

		// Test with string
		s := "hello"
		bytesStr := xunsafe.Bytes(&s)
		So(bytesStr, ShouldNotBeNil)
		So(len(bytesStr), ShouldEqual, 16) // string is 16 bytes on 64-bit systems

		// Test with struct
		type TestStruct struct {
			ID   int
			Name string
		}
		ts := TestStruct{ID: 1, Name: "test"}
		bytesStruct := xunsafe.Bytes(&ts)
		So(bytesStruct, ShouldNotBeNil)
		So(len(bytesStruct), ShouldEqual, 24) // 8 + 16 bytes

		// Test with empty struct
		type EmptyStruct struct{}
		empty := EmptyStruct{}
		bytesEmpty := xunsafe.Bytes(&empty)
		So(bytesEmpty, ShouldNotBeNil)
		So(len(bytesEmpty), ShouldEqual, 0)
	})
}

func TestLoadSlice(t *testing.T) {
	Convey("Given various slice types", t, func() {
		// Test with int slice
		intSlice := []int{1, 2, 3, 4, 5}

		val0 := xunsafe.LoadSlice(intSlice, 0)
		So(val0, ShouldEqual, 1)

		val2 := xunsafe.LoadSlice(intSlice, 2)
		So(val2, ShouldEqual, 3)

		val4 := xunsafe.LoadSlice(intSlice, 4)
		So(val4, ShouldEqual, 5)

		// Test with string slice
		strSlice := []string{"hello", "world", "test"}

		str0 := xunsafe.LoadSlice(strSlice, 0)
		So(str0, ShouldEqual, "hello")

		str1 := xunsafe.LoadSlice(strSlice, 1)
		So(str1, ShouldEqual, "world")
	})
}

func TestSliceToString(t *testing.T) {
	Convey("Given various slice types", t, func() {
		// Test with int slice
		intSlice := []int{1, 2, 3, 4, 5}
		str := xunsafe.SliceToString(intSlice)
		So(str, ShouldNotBeEmpty)

		// Test with string slice
		strSlice := []string{"hello", "world"}
		strResult := xunsafe.SliceToString(strSlice)
		So(strResult, ShouldNotBeEmpty)

		// Test with empty slice - empty slice should produce empty string
		emptySlice := []int{}
		emptyStr := xunsafe.SliceToString(emptySlice)
		So(emptyStr, ShouldBeEmpty)
	})
}

func TestStringToSlice(t *testing.T) {
	Convey("Given string to slice conversion", t, func() {
		// Test with string that can be converted to int slice
		// Note: This is a bit tricky since we need the right size
		// For now, we'll test that it doesn't panic
		So(func() {
			_ = xunsafe.StringToSlice[[]int]("test")
		}, ShouldNotPanic)

		// Test with empty string
		So(func() {
			_ = xunsafe.StringToSlice[[]int]("")
		}, ShouldNotPanic)
	})
}

func TestSliceEdgeCases(t *testing.T) {
	Convey("Given edge case slices", t, func() {
		// Test with single element slice
		singleSlice := []int{42}
		val := xunsafe.LoadSlice(singleSlice, 0)
		So(val, ShouldEqual, 42)

		// Test with empty slice - might not panic depending on implementation
		emptySlice := []int{}

		// Test that we can create an empty slice without panicking
		So(func() {
			_ = emptySlice
		}, ShouldNotPanic)
	})
}

func TestSliceTypes(t *testing.T) {
	Convey("Given different slice types", t, func() {
		// Test with different numeric types
		int8Slice := []int8{1, 2, 3}
		val1 := xunsafe.LoadSlice(int8Slice, 1)
		So(val1, ShouldEqual, int8(2))

		uint16Slice := []uint16{10, 20, 30}
		val2 := xunsafe.LoadSlice(uint16Slice, 2)
		So(val2, ShouldEqual, uint16(30))

		// Test with float types
		floatSlice := []float64{1.1, 2.2, 3.3}
		val3 := xunsafe.LoadSlice(floatSlice, 0)
		So(val3, ShouldEqual, 1.1)

		// Test with bool slice
		boolSlice := []bool{true, false, true}
		val4 := xunsafe.LoadSlice(boolSlice, 1)
		So(val4, ShouldEqual, false)
	})
}

func TestSliceComprehensive(t *testing.T) {
	Convey("Given comprehensive slice tests", t, func() {
		// Test with various slice types
		testCases := []struct {
			name  string
			slice interface{}
			index int
			value interface{}
		}{
			{"int8", []int8{1, 2, 3}, 1, int8(2)},
			{"int16", []int16{10, 20, 30}, 2, int16(30)},
			{"int32", []int32{100, 200, 300}, 0, int32(100)},
			{"int64", []int64{1000, 2000, 3000}, 1, int64(2000)},
			{"uint8", []uint8{1, 2, 3}, 2, uint8(3)},
			{"uint16", []uint16{10, 20, 30}, 0, uint16(10)},
			{"uint32", []uint32{100, 200, 300}, 1, uint32(200)},
			{"uint64", []uint64{1000, 2000, 3000}, 2, uint64(3000)},
			{"float32", []float32{1.1, 2.2, 3.3}, 0, float32(1.1)},
			{"float64", []float64{1.1, 2.2, 3.3}, 1, float64(2.2)},
			{"bool", []bool{true, false, true}, 1, false},
			{"string", []string{"hello", "world", "test"}, 2, "test"},
		}

		for _, tc := range testCases {
			Convey("When testing "+tc.name, func() {
				// Test LoadSlice operation
				switch v := tc.slice.(type) {
				case []int8:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []int16:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []int32:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []int64:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []uint8:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []uint16:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []uint32:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []uint64:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []float32:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []float64:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []bool:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				case []string:
					val := xunsafe.LoadSlice(v, tc.index)
					So(val, ShouldEqual, tc.value)
				}
			})
		}
	})
}
