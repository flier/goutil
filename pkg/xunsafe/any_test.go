package xunsafe_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestAny(t *testing.T) {
	Convey("Given any type operations", t, func() {
		Convey("When extracting data from various types", func() {
			Convey("And extracting data from int", func() {
				i := 42
				data := xunsafe.AnyData(i)
				So(data, ShouldNotBeNil)
			})

			Convey("And extracting data from string", func() {
				s := "hello"
				dataStr := xunsafe.AnyData(s)
				So(dataStr, ShouldNotBeNil)
			})

			Convey("And extracting data from struct", func() {
				type TestStruct struct {
					ID   int
					Name string
				}
				ts := TestStruct{ID: 1, Name: "test"}
				dataStruct := xunsafe.AnyData(ts)
				So(dataStruct, ShouldNotBeNil)
			})

			Convey("And extracting data from pointer", func() {
				i := 42
				p := &i
				dataPtr := xunsafe.AnyData(p)
				So(dataPtr, ShouldNotBeNil)
			})
		})

		Convey("When extracting type information", func() {
			Convey("And extracting type from int", func() {
				i := 42
				typ := xunsafe.AnyType(i)
				So(typ, ShouldNotEqual, uintptr(0))
			})

			Convey("And extracting type from string", func() {
				s := "hello"
				typStr := xunsafe.AnyType(s)
				So(typStr, ShouldNotEqual, uintptr(0))
			})

			Convey("And extracting type from struct", func() {
				type TestStruct struct {
					ID   int
					Name string
				}
				ts := TestStruct{ID: 1, Name: "test"}
				typStruct := xunsafe.AnyType(ts)
				So(typStruct, ShouldNotEqual, uintptr(0))
			})
		})

		Convey("When working with any bytes", func() {
			Convey("And getting bytes from int", func() {
				i := 42
				bytes := xunsafe.AnyBytes(i)
				So(bytes, ShouldNotBeNil)
				So(len(bytes), ShouldBeGreaterThan, 0)
			})

			Convey("And getting bytes from string", func() {
				s := "hello"
				bytesStr := xunsafe.AnyBytes(s)
				So(bytesStr, ShouldNotBeNil)
				So(len(bytesStr), ShouldBeGreaterThan, 0)
			})

			Convey("And getting bytes from struct", func() {
				type TestStruct struct {
					ID   int
					Name string
				}
				ts := TestStruct{ID: 1, Name: "test"}
				bytesStruct := xunsafe.AnyBytes(ts)
				So(bytesStruct, ShouldNotBeNil)
				So(len(bytesStruct), ShouldBeGreaterThan, 0)
			})

			Convey("And getting bytes from nil", func() {
				bytesNil := xunsafe.AnyBytes(nil)
				So(bytesNil, ShouldBeNil)
			})
		})

		Convey("When creating any values", func() {
			Convey("And creating any from int data", func() {
				i := 42
				typ := xunsafe.AnyType(i)
				data := xunsafe.AnyData(i)

				anyVal := xunsafe.MakeAny(typ, data)
				So(anyVal, ShouldNotBeNil)

				// Test that we can extract the data back
				extractedData := xunsafe.AnyData(anyVal)
				So(extractedData, ShouldNotBeNil)
			})
		})

		Convey("When checking direct types", func() {
			Convey("And checking direct types", func() {
				p := &struct{}{}
				So(xunsafe.IsDirectAny(p), ShouldBeTrue)

				m := make(map[int]int)
				So(xunsafe.IsDirectAny(m), ShouldBeTrue)

				c := make(chan int)
				So(xunsafe.IsDirectAny(c), ShouldBeTrue)
			})

			Convey("And checking indirect types", func() {
				i := 42
				So(xunsafe.IsDirectAny(i), ShouldBeFalse)

				s := "hello"
				So(xunsafe.IsDirectAny(s), ShouldBeFalse)

				type TestStruct struct {
					ID   int
					Name string
				}
				ts := TestStruct{ID: 1, Name: "test"}
				So(xunsafe.IsDirectAny(ts), ShouldBeFalse)
			})
		})

		Convey("When checking generic direct types", func() {
			Convey("And checking indirect types", func() {
				So(xunsafe.IsDirect[int](), ShouldBeFalse)
				So(xunsafe.IsDirect[string](), ShouldBeFalse)
				So(xunsafe.IsDirect[[]byte](), ShouldBeFalse)
			})

			Convey("And checking direct types", func() {
				So(xunsafe.IsDirect[*int](), ShouldBeTrue)
				So(xunsafe.IsDirect[any](), ShouldBeTrue)
				So(xunsafe.IsDirect[map[int]int](), ShouldBeTrue)
				So(xunsafe.IsDirect[chan int](), ShouldBeTrue)
			})

			Convey("And checking single-field structs", func() {
				type SingleField struct {
					Value *int
				}
				So(xunsafe.IsDirect[SingleField](), ShouldBeTrue)
			})

			Convey("And checking multi-field structs", func() {
				type MultiField struct {
					ID   int
					Name string
				}
				So(xunsafe.IsDirect[MultiField](), ShouldBeFalse)
			})
		})

		Convey("When asserting inlined any", func() {
			Convey("And asserting direct types", func() {
				xunsafe.AssertInlinedAny[*int](t)
				xunsafe.AssertInlinedAny[any](t)
				xunsafe.AssertInlinedAny[map[int]int](t)
				xunsafe.AssertInlinedAny[chan int](t)
			})

			Convey("And asserting single-field structs", func() {
				type SingleField struct {
					Value *int
				}
				xunsafe.AssertInlinedAny[SingleField](t)
			})
		})

		Convey("When handling edge cases", func() {
			Convey("And working with empty struct", func() {
				type EmptyStruct struct{}
				empty := EmptyStruct{}
				bytesEmpty := xunsafe.AnyBytes(empty)
				So(bytesEmpty, ShouldNotBeNil)
			})

			Convey("And working with zero values", func() {
				var zeroInt int
				bytesZero := xunsafe.AnyBytes(zeroInt)
				So(bytesZero, ShouldNotBeNil)
			})

			Convey("And working with interface", func() {
				var iface any = "test"
				bytesIface := xunsafe.AnyBytes(iface)
				So(bytesIface, ShouldNotBeNil)
			})
		})

		Convey("When performing comprehensive tests", func() {
			Convey("Given various types", func() {
				testCases := []struct {
					name     string
					value    any
					expected bool
				}{
					{"int", 42, false},
					{"string", "hello", false},
					{"pointer", &struct{}{}, true},
					{"map", make(map[int]int), true},
					{"channel", make(chan int), true},
					{"func", func() {}, true},
					{"interface", any(42), false},
					{"nil", nil, false},
				}

				for _, tc := range testCases {
					Convey(fmt.Sprintf("When testing %s", tc.name), func() {
						if tc.value != nil {
							Convey("And testing AnyData", func() {
								data := xunsafe.AnyData(tc.value)
								So(data, ShouldNotBeNil)
							})

							Convey("And testing AnyType", func() {
								typ := xunsafe.AnyType(tc.value)
								So(typ, ShouldNotEqual, uintptr(0))
							})

							Convey("And testing AnyBytes", func() {
								bytes := xunsafe.AnyBytes(tc.value)
								if tc.value != nil {
									So(bytes, ShouldNotBeNil)
								}
							})

							Convey("And testing IsDirectAny", func() {
								isDirect := xunsafe.IsDirectAny(tc.value)
								So(isDirect, ShouldEqual, tc.expected)
							})
						}
					})
				}
			})
		})
	})
}
