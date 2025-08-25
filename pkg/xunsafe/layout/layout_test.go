package layout_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe/layout"
)

func TestLayout(t *testing.T) {
	Convey("Given type layout operations", t, func() {
		Convey("When working with basic type sizes", func() {
			Convey("And getting size of int", func() {
				size := layout.Size[int]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of string", func() {
				size := layout.Size[string]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of bool", func() {
				size := layout.Size[bool]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of float64", func() {
				size := layout.Size[float64]()
				So(size, ShouldBeGreaterThan, 0)
			})
		})

		Convey("When working with struct layouts", func() {
			Convey("And getting size of simple struct", func() {
				type simpleStruct struct {
					Field1 int
					Field2 string
				}
				size := layout.Size[simpleStruct]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of nested struct", func() {
				type nestedStruct struct {
					Field1 int
					Field2 struct {
						SubField1 string
						SubField2 bool
					}
				}
				size := layout.Size[nestedStruct]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of struct with pointers", func() {
				type pointerStruct struct {
					Field1 *int
					Field2 *string
				}
				size := layout.Size[pointerStruct]()
				So(size, ShouldBeGreaterThan, 0)
			})
		})

		Convey("When working with array and slice layouts", func() {
			Convey("And getting size of array", func() {
				size := layout.Size[[5]int]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of slice", func() {
				size := layout.Size[[]int]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of string array", func() {
				size := layout.Size[[3]string]()
				So(size, ShouldBeGreaterThan, 0)
			})
		})

		Convey("When working with interface layouts", func() {
			Convey("And getting size of interface", func() {
				size := layout.Size[interface{}]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of specific interface", func() {
				type testInterface interface {
					Method() string
				}
				size := layout.Size[testInterface]()
				So(size, ShouldBeGreaterThan, 0)
			})
		})

		Convey("When working with function layouts", func() {
			Convey("And getting size of function", func() {
				size := layout.Size[func()]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of function with parameters", func() {
				size := layout.Size[func(int, string) bool]()
				So(size, ShouldBeGreaterThan, 0)
			})

			Convey("And getting size of function with return values", func() {
				size := layout.Size[func() (int, error)]()
				So(size, ShouldBeGreaterThan, 0)
			})
		})

		Convey("When working with alignment operations", func() {
			Convey("And rounding up to alignment", func() {
				aligned := layout.RoundUp(15, 8)
				So(aligned, ShouldEqual, 16)
			})

			Convey("And rounding down to alignment", func() {
				aligned := layout.RoundDown(15, 8)
				So(aligned, ShouldEqual, 8)
			})

			Convey("And checking alignment values", func() {
				align := layout.Align[int]()
				So(align, ShouldBeGreaterThan, 0)

				alignString := layout.Align[string]()
				So(alignString, ShouldBeGreaterThan, 0)
			})
		})

		Convey("When working with padding calculations", func() {
			Convey("And calculating padding for aligned value", func() {
				padding := layout.Padding(15, 8)
				So(padding, ShouldEqual, 1)
			})

			Convey("And calculating padding for already aligned value", func() {
				padding := layout.Padding(16, 8)
				So(padding, ShouldEqual, 0)
			})

			Convey("And calculating padding for zero alignment", func() {
				padding := layout.Padding(15, 0)
				So(padding, ShouldEqual, 0)
			})
		})

		Convey("When working with layout information", func() {
			Convey("And getting layout of int", func() {
				layoutInfo := layout.Of[int]()
				So(layoutInfo.Size, ShouldBeGreaterThan, 0)
				So(layoutInfo.Align, ShouldBeGreaterThan, 0)
			})

			Convey("And getting layout of string", func() {
				layoutInfo := layout.Of[string]()
				So(layoutInfo.Size, ShouldBeGreaterThan, 0)
				So(layoutInfo.Align, ShouldBeGreaterThan, 0)
			})

			Convey("And getting layout of struct", func() {
				type testStruct struct {
					Field1 int
					Field2 string
				}
				layoutInfo := layout.Of[testStruct]()
				So(layoutInfo.Size, ShouldBeGreaterThan, 0)
				So(layoutInfo.Align, ShouldBeGreaterThan, 0)
			})
		})

		Convey("When working with layout operations", func() {
			Convey("And finding maximum layout", func() {
				layout1 := layout.Layout{Size: 8, Align: 4}
				layout2 := layout.Layout{Size: 16, Align: 8}

				maxLayout := layout1.Max(layout2)
				So(maxLayout.Size, ShouldEqual, 16)
				So(maxLayout.Align, ShouldEqual, 8)
			})

			Convey("And finding maximum with same size but different alignment", func() {
				layout1 := layout.Layout{Size: 8, Align: 4}
				layout3 := layout.Layout{Size: 8, Align: 16}

				maxLayout := layout1.Max(layout3)
				So(maxLayout.Size, ShouldEqual, 8)
				So(maxLayout.Align, ShouldEqual, 16)
			})
		})

		Convey("When working with slice padding", func() {
			Convey("And padding slice to alignment", func() {
				buf := []byte{1, 2, 3, 4, 5, 6, 7}
				padded := layout.PadSlice(buf, 8)
				So(len(padded), ShouldEqual, 8)
				So(padded[:7], ShouldResemble, buf)
				So(padded[7], ShouldEqual, byte(0))
			})

			Convey("And padding already aligned slice", func() {
				buf := []byte{1, 2, 3, 4, 5, 6, 7, 8}
				padded := layout.PadSlice(buf, 8)
				So(len(padded), ShouldEqual, 8)
				So(padded, ShouldResemble, buf)
			})
		})

		Convey("When working with edge cases", func() {
			Convey("And working with zero alignment", func() {
				aligned := layout.RoundUp(15, 0)
				So(aligned, ShouldEqual, 15)

				alignedDown := layout.RoundDown(15, 0)
				So(alignedDown, ShouldEqual, 15)
			})

			Convey("And working with negative values", func() {
				aligned := layout.RoundUp(-15, 8)
				So(aligned, ShouldEqual, -8)

				alignedDown := layout.RoundDown(-15, 8)
				So(alignedDown, ShouldEqual, -16)
			})
		})

		Convey("When working with comprehensive type tests", func() {
			Convey("And testing various numeric types", func() {
				testCases := []struct {
					name string
					typ  interface{}
				}{
					{"int8", int8(0)},
					{"int16", int16(0)},
					{"int32", int32(0)},
					{"int64", int64(0)},
					{"uint8", uint8(0)},
					{"uint16", uint16(0)},
					{"uint32", uint32(0)},
					{"uint64", uint64(0)},
					{"float32", float32(0)},
					{"float64", float64(0)},
					{"complex64", complex64(0)},
					{"complex128", complex128(0)},
				}

				for _, tc := range testCases {
					Convey(tc.name, func() {
						switch tc.typ.(type) {
						case int8:
							size := layout.Size[int8]()
							So(size, ShouldBeGreaterThan, 0)
						case int16:
							size := layout.Size[int16]()
							So(size, ShouldBeGreaterThan, 0)
						case int32:
							size := layout.Size[int32]()
							So(size, ShouldBeGreaterThan, 0)
						case int64:
							size := layout.Size[int64]()
							So(size, ShouldBeGreaterThan, 0)
						case uint8:
							size := layout.Size[uint8]()
							So(size, ShouldBeGreaterThan, 0)
						case uint16:
							size := layout.Size[uint16]()
							So(size, ShouldBeGreaterThan, 0)
						case uint32:
							size := layout.Size[uint32]()
							So(size, ShouldBeGreaterThan, 0)
						case uint64:
							size := layout.Size[uint64]()
							So(size, ShouldBeGreaterThan, 0)
						case float32:
							size := layout.Size[float32]()
							So(size, ShouldBeGreaterThan, 0)
						case float64:
							size := layout.Size[float64]()
							So(size, ShouldBeGreaterThan, 0)
						case complex64:
							size := layout.Size[complex64]()
							So(size, ShouldBeGreaterThan, 0)
						case complex128:
							size := layout.Size[complex128]()
							So(size, ShouldBeGreaterThan, 0)
						}
					})
				}
			})
		})
	})
}
