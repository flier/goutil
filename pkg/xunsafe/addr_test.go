//go:build go1.20

package xunsafe_test

import (
	"fmt"
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestAddr(t *testing.T) {
	Convey("Given address operations", t, func() {
		Convey("When getting address of various types", func() {
			Convey("And getting address of int", func() {
				i := 42
				addr := xunsafe.AddrOf(&i)
				So(uintptr(addr), ShouldEqual, uintptr(unsafe.Pointer(&i)))
			})

			Convey("And getting address of string", func() {
				s := "hello"
				addrStr := xunsafe.AddrOf(&s)
				So(uintptr(addrStr), ShouldEqual, uintptr(unsafe.Pointer(&s)))
			})

			Convey("And getting address of struct", func() {
				type TestStruct struct {
					ID   int
					Name string
				}
				ts := TestStruct{ID: 1, Name: "test"}
				addrStruct := xunsafe.AddrOf(&ts)
				So(uintptr(addrStruct), ShouldEqual, uintptr(unsafe.Pointer(&ts)))
			})
		})

		Convey("When getting end address of slices", func() {
			Convey("And getting end address of int slice", func() {
				s := []int{1, 2, 3, 4, 5}
				end := xunsafe.EndOf(s)
				So(uintptr(end), ShouldEqual,
					uintptr(unsafe.Add(unsafe.Pointer(unsafe.SliceData(s)), unsafe.Sizeof(*new(int))*uintptr(len(s)))))
			})

			Convey("And getting end address of string slice", func() {
				s := []string{"a", "b", "c"}
				end := xunsafe.EndOf(s)
				So(uintptr(end), ShouldEqual,
					uintptr(unsafe.Add(unsafe.Pointer(unsafe.SliceData(s)), unsafe.Sizeof(*new(string))*uintptr(len(s)))))
			})

			Convey("And getting end address of empty slice", func() {
				s := []int{}
				end := xunsafe.EndOf(s)
				So(uintptr(end), ShouldEqual, uintptr(unsafe.Pointer(unsafe.SliceData(s))))
			})
		})

		Convey("When asserting valid addresses", func() {
			Convey("And asserting address of int", func() {
				i := 42
				addr := xunsafe.AddrOf(&i)
				ptr := addr.AssertValid()
				So(ptr, ShouldEqual, &i)
				So(*ptr, ShouldEqual, 42)
			})

			Convey("And asserting address of string", func() {
				s := "hello"
				addr := xunsafe.AddrOf(&s)
				ptr := addr.AssertValid()
				So(ptr, ShouldEqual, &s)
				So(*ptr, ShouldEqual, "hello")
			})

			Convey("And asserting address of struct", func() {
				type TestStruct struct {
					ID   int
					Name string
				}
				ts := TestStruct{ID: 1, Name: "test"}
				addr := xunsafe.AddrOf(&ts)
				ptr := addr.AssertValid()
				So(ptr, ShouldEqual, &ts)
				So(ptr.ID, ShouldEqual, 1)
				So(ptr.Name, ShouldEqual, "test")
			})
		})

		Convey("When performing address arithmetic", func() {
			Convey("Given an array and base address", func() {
				arr := [5]int{1, 2, 3, 4, 5}
				baseAddr := xunsafe.AddrOf(unsafe.SliceData(arr[:]))

				Convey("And adding offset to get address of arr[2]", func() {
					addr2 := baseAddr.Add(2)
					ptr2 := addr2.AssertValid()
					So(*ptr2, ShouldEqual, 3)
				})

				Convey("And adding offset to get address of arr[4]", func() {
					addr4 := baseAddr.Add(4)
					ptr4 := addr4.AssertValid()
					So(*ptr4, ShouldEqual, 5)
				})

				Convey("And adding byte offset to get address of arr[1]", func() {
					addr := baseAddr.ByteAdd(int(unsafe.Sizeof(*new(int))))
					ptr := addr.AssertValid()
					So(*ptr, ShouldEqual, 2)
				})

				Convey("And subtracting addresses", func() {
					addr4 := baseAddr.Add(4)
					addr2 := baseAddr.Add(2)
					diff := addr4.Sub(addr2)
					So(diff, ShouldEqual, 2)
				})

				Convey("And subtracting same address", func() {
					addr2 := baseAddr.Add(2)
					sameDiff := addr2.Sub(addr2)
					So(sameDiff, ShouldEqual, 0)
				})
			})
		})

		Convey("When calculating padding", func() {
			Convey("Given an address", func() {
				addr := xunsafe.Addr[int](8)

				Convey("And calculating padding for 8-byte alignment", func() {
					padding := addr.Padding(8)
					So(padding, ShouldEqual, 0)
				})

				Convey("And calculating padding for 16-byte alignment", func() {
					padding16 := addr.Padding(16)
					So(padding16, ShouldEqual, 8)
				})

				Convey("And calculating padding for 4-byte alignment", func() {
					padding4 := addr.Padding(4)
					So(padding4, ShouldEqual, 0)
				})
			})
		})

		Convey("When rounding addresses", func() {
			Convey("Given an address", func() {
				addr := xunsafe.Addr[int](9)

				Convey("And rounding up to 8-byte alignment", func() {
					rounded8 := addr.RoundUpTo(8)
					So(rounded8, ShouldEqual, xunsafe.Addr[int](16))
				})

				Convey("And rounding up to 16-byte alignment", func() {
					rounded16 := addr.RoundUpTo(16)
					So(rounded16, ShouldEqual, xunsafe.Addr[int](16))
				})

				Convey("And rounding up to 4-byte alignment", func() {
					rounded4 := addr.RoundUpTo(4)
					So(rounded4, ShouldEqual, xunsafe.Addr[int](12))
				})
			})
		})

		Convey("When working with sign bits", func() {
			Convey("And detecting sign bit of positive address", func() {
				positiveAddr := xunsafe.Addr[int](0x7FFFFFFF)
				So(positiveAddr.SignBit(), ShouldBeFalse)
			})

			Convey("And detecting sign bit of negative address", func() {
				negativeAddr := xunsafe.Addr[int](-1)
				So(negativeAddr.SignBit(), ShouldBeTrue)
			})

			Convey("And detecting sign bit of zero address", func() {
				zeroAddr := xunsafe.Addr[int](0)
				So(zeroAddr.SignBit(), ShouldBeFalse)
			})

			Convey("And getting sign bit mask", func() {
				Convey("For positive address", func() {
					positiveAddr := xunsafe.Addr[int](0x7FFFFFFF)
					mask := positiveAddr.SignBitMask()
					So(mask, ShouldEqual, xunsafe.Addr[int](0))
				})

				Convey("For negative address", func() {
					negativeAddr := xunsafe.Addr[int](-1)
					maskNeg := negativeAddr.SignBitMask()
					So(maskNeg, ShouldEqual, xunsafe.Addr[int](-1))
				})
			})

			Convey("And clearing sign bit", func() {
				Convey("For negative address", func() {
					negativeAddr := xunsafe.Addr[int](-1)
					cleared := negativeAddr.ClearSignBit()
					So(cleared.SignBit(), ShouldBeFalse)
				})

				Convey("For positive address", func() {
					positiveAddr := xunsafe.Addr[int](0x7FFFFFFF)
					clearedPos := positiveAddr.ClearSignBit()
					So(clearedPos.SignBit(), ShouldBeFalse)
				})
			})
		})

		Convey("When formatting addresses", func() {
			Convey("And formatting with %v", func() {
				addr := xunsafe.Addr[int](0x12345678)
				result := fmt.Sprintf("%v", addr)
				So(result, ShouldContainSubstring, "0x12345678")
			})

			Convey("And formatting with %x", func() {
				addr := xunsafe.Addr[int](0x12345678)
				resultHex := fmt.Sprintf("%x", addr)
				So(resultHex, ShouldContainSubstring, "12345678")
			})

			Convey("And formatting zero address", func() {
				zeroAddr := xunsafe.Addr[int](0)
				zeroResult := fmt.Sprintf("%v", zeroAddr)
				So(zeroResult, ShouldContainSubstring, "0x0")
			})

			Convey("And formatting large address", func() {
				largeAddr := xunsafe.Addr[int](0x7FFFFFFFFFFFFFFF)
				largeResult := fmt.Sprintf("%v", largeAddr)
				So(largeResult, ShouldContainSubstring, "0x")
			})
		})

		Convey("When handling edge cases", func() {
			Convey("And working with very large addresses", func() {
				largeAddr := xunsafe.Addr[int](0x7FFFFFFF)
				So(largeAddr.SignBit(), ShouldBeFalse)
			})

			Convey("And working with very small addresses", func() {
				smallAddr := xunsafe.Addr[int](1)
				So(smallAddr.SignBit(), ShouldBeFalse)
			})

			Convey("And performing arithmetic with edge cases", func() {
				addr1 := xunsafe.Addr[int](0x7FFFFFFF)
				addr2 := addr1.Add(1)
				So(addr2, ShouldNotEqual, addr1)
			})

			Convey("And working with zero address", func() {
				zeroAddr := xunsafe.Addr[int](0)
				So(zeroAddr.SignBit(), ShouldBeFalse)
				So(zeroAddr.SignBitMask(), ShouldEqual, xunsafe.Addr[int](0))
				So(zeroAddr.ClearSignBit(), ShouldEqual, xunsafe.Addr[int](0))
			})

			Convey("And working with negative address", func() {
				negAddr := xunsafe.Addr[int](-1)
				So(negAddr.SignBit(), ShouldBeTrue)
				So(negAddr.SignBitMask(), ShouldEqual, xunsafe.Addr[int](-1))
				cleared := negAddr.ClearSignBit()
				So(cleared.SignBit(), ShouldBeFalse)
			})
		})

		Convey("When performing comprehensive operations", func() {
			Convey("Given an array of integers", func() {
				arr := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
				baseAddr := xunsafe.AddrOf(&arr[0])

				Convey("And accessing elements sequentially", func() {
					for i := 0; i < 10; i++ {
						addr := baseAddr.Add(i)
						ptr := addr.AssertValid()
						So(*ptr, ShouldEqual, i)
					}
				})

				Convey("And performing byte operations", func() {
					byteAddr := baseAddr.ByteAdd(8) // Assuming int is 8 bytes
					ptr := byteAddr.AssertValid()
					So(*ptr, ShouldEqual, 1)
				})

				Convey("And performing address arithmetic", func() {
					addr1 := baseAddr.Add(2)
					addr2 := baseAddr.Add(5)
					diff := addr2.Sub(addr1)
					So(diff, ShouldEqual, 3)
				})

				Convey("And calculating padding", func() {
					padding := baseAddr.Padding(16)
					So(padding, ShouldBeGreaterThanOrEqualTo, 0)
				})

				Convey("And rounding operations", func() {
					rounded := baseAddr.RoundUpTo(16)
					So(uintptr(rounded), ShouldBeGreaterThanOrEqualTo, uintptr(baseAddr))
				})
			})
		})
	})
}
