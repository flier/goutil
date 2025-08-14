//go:build go1.22

package arena_test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
)

func TestArena_New(t *testing.T) {
	Convey("Given an arena", t, func() {
		a := &arena.Arena{}

		Convey("When creating a new value", func() {
			value := 42
			ptr := arena.New(a, value)

			So(ptr, ShouldNotBeNil)
			So(*ptr, ShouldEqual, 42)
		})

		Convey("When creating multiple values", func() {
			values := []int{1, 2, 3, 4, 5}
			ptrs := make([]*int, len(values))

			for i, v := range values {
				ptrs[i] = arena.New(a, v)
			}

			So(len(ptrs), ShouldEqual, 5)
			for i, ptr := range ptrs {
				So(ptr, ShouldNotBeNil)
				So(*ptr, ShouldEqual, values[i])
			}
		})

		Convey("When creating string values", func() {
			str := "hello world"
			ptr := arena.New(a, str)

			So(ptr, ShouldNotBeNil)
			So(*ptr, ShouldEqual, "hello world")
		})

		Convey("When creating struct values", func() {
			type testStruct struct {
				ID   int
				Name string
				Data []byte
			}

			ts := testStruct{
				ID:   123,
				Name: "test",
				Data: []byte("data"),
			}

			ptr := arena.New(a, ts)

			So(ptr, ShouldNotBeNil)
			So(ptr.ID, ShouldEqual, 123)
			So(ptr.Name, ShouldEqual, "test")
			So(ptr.Data, ShouldResemble, []byte("data"))
		})
	})
}

func TestArena_Alloc(t *testing.T) {
	Convey("Given an arena", t, func() {
		a := &arena.Arena{}

		Convey("When allocating small amounts of memory", func() {
			ptr1 := a.Alloc(8)
			So(ptr1, ShouldNotBeNil)

			ptr2 := a.Alloc(16)
			So(ptr2, ShouldNotBeNil)
			So(uintptr(unsafe.Pointer(ptr2)), ShouldNotEqual, uintptr(unsafe.Pointer(ptr1)))

			// Verify alignment
			So(uintptr(unsafe.Pointer(ptr1))%uintptr(arena.Align), ShouldEqual, uintptr(0))
			So(uintptr(unsafe.Pointer(ptr2))%uintptr(arena.Align), ShouldEqual, uintptr(0))
		})

		Convey("When allocating large amounts of memory", func() {
			largeSize := 1024 * 1024 // 1MB
			ptr := a.Alloc(largeSize)

			So(ptr, ShouldNotBeNil)
			So(uintptr(unsafe.Pointer(ptr))%uintptr(arena.Align), ShouldEqual, uintptr(0))
		})

		Convey("When allocating zero bytes", func() {
			ptr := a.Alloc(0)
			// Zero byte allocation might return nil or a valid pointer
			if ptr != nil {
				So(uintptr(unsafe.Pointer(ptr))%uintptr(arena.Align), ShouldEqual, uintptr(0))
			}
		})

		Convey("When allocating with alignment requirements", func() {
			// Test that allocations are properly aligned
			for i := 0; i < 10; i++ {
				size := 8 + i*8
				ptr := a.Alloc(size)
				So(uintptr(unsafe.Pointer(ptr))%uintptr(arena.Align), ShouldEqual, uintptr(0))
			}
		})
	})
}

func TestArena_Reserve(t *testing.T) {
	Convey("Given an arena", t, func() {
		a := &arena.Arena{}

		Convey("When reserving memory", func() {
			// Reserve a large amount
			a.Reserve(1024 * 1024) // 1MB

			// Should be able to allocate without growing
			ptr := a.Alloc(1024)
			So(ptr, ShouldNotBeNil)
		})

		Convey("When reserving multiple times", func() {
			a.Reserve(1000)
			a.Reserve(2000)
			a.Reserve(500)

			// Should be able to allocate the total reserved amount
			ptr := a.Alloc(3500)
			So(ptr, ShouldNotBeNil)
		})
	})
}

func TestArena_Grow(t *testing.T) {
	Convey("Given an arena", t, func() {
		a := &arena.Arena{}

		Convey("When growing automatically", func() {
			// Allocate more than the initial capacity
			largeSize := 1024 * 1024 // 1MB
			ptr := a.Alloc(largeSize)

			So(ptr, ShouldNotBeNil)
			So(a.Cap, ShouldBeGreaterThanOrEqualTo, largeSize)
		})

		Convey("When growing multiple times", func() {
			initialCap := a.Cap

			// Force multiple grows
			for i := 0; i < 5; i++ {
				size := 1024 * (i + 1)
				ptr := a.Alloc(size)
				So(ptr, ShouldNotBeNil)
			}

			So(a.Cap, ShouldBeGreaterThan, initialCap)
		})
	})
}

func TestArena_Free(t *testing.T) {
	Convey("Given an arena with allocated memory", t, func() {
		a := &arena.Arena{}

		// Allocate some memory
		ptr1 := a.Alloc(1000)
		ptr2 := a.Alloc(2000)
		So(ptr1, ShouldNotBeNil)
		So(ptr2, ShouldNotBeNil)

		Convey("When freeing the arena", func() {
			a.Free()

			So(a.Cap, ShouldBeGreaterThan, 0)
			So(a.Next, ShouldNotBeNil)
			So(a.End, ShouldNotBeNil)
		})

		Convey("When reusing after free", func() {
			a.Free()

			// Should be able to allocate again
			ptr3 := a.Alloc(500)
			So(ptr3, ShouldNotBeNil)
		})
	})
}

func TestArena_KeepAlive(t *testing.T) {
	Convey("Given an arena", t, func() {
		a := &arena.Arena{}

		Convey("When keeping values alive", func() {
			value := "keep me alive"
			// This should not panic
			So(func() {
				a.KeepAlive(value)
			}, ShouldNotPanic)
		})

		Convey("When keeping multiple values alive", func() {
			values := []string{"value1", "value2", "value3"}

			// This should not panic
			So(func() {
				for _, v := range values {
					a.KeepAlive(v)
				}
			}, ShouldNotPanic)
		})
	})
}

func TestArena_Log(t *testing.T) {
	Convey("Given an arena", t, func() {
		a := &arena.Arena{}

		Convey("When logging operations", func() {
			// This should not panic
			So(func() {
				a.Log("test", "test message")
			}, ShouldNotPanic)
		})
	})
}

func TestArena_EdgeCases(t *testing.T) {
	Convey("Given edge cases", t, func() {
		Convey("When using zero arena", func() {
			var a arena.Arena

			// Should be able to allocate
			ptr := a.Alloc(8)
			So(ptr, ShouldNotBeNil)
		})

		Convey("When allocating very large amounts", func() {
			a := &arena.Arena{}

			// Test with a reasonable large size
			largeSize := 10 * 1024 * 1024 // 10MB
			ptr := a.Alloc(largeSize)

			So(ptr, ShouldNotBeNil)
			So(a.Cap, ShouldBeGreaterThanOrEqualTo, largeSize)
		})

		Convey("When allocating with odd sizes", func() {
			a := &arena.Arena{}

			// Test various odd sizes
			sizes := []int{1, 3, 7, 15, 31, 63, 127, 255, 511, 1023}

			for _, size := range sizes {
				ptr := a.Alloc(size)
				So(ptr, ShouldNotBeNil)
				So(uintptr(unsafe.Pointer(ptr))%uintptr(arena.Align), ShouldEqual, uintptr(0))
			}
		})
	})
}

func TestArena_Performance(t *testing.T) {
	Convey("Given performance scenarios", t, func() {
		Convey("When allocating many small objects", func() {
			a := &arena.Arena{}

			So(func() {
				for i := 0; i < 10000; i++ {
					ptr := a.Alloc(8)
					So(ptr, ShouldNotBeNil)
				}
			}, ShouldNotPanic)
		})

		Convey("When allocating mixed sizes", func() {
			a := &arena.Arena{}

			So(func() {
				for i := 0; i < 1000; i++ {
					size := 8 + (i%100)*8
					ptr := a.Alloc(size)
					So(ptr, ShouldNotBeNil)
				}
			}, ShouldNotPanic)
		})

		Convey("When freeing and reusing", func() {
			a := &arena.Arena{}

			So(func() {
				for i := 0; i < 10; i++ {
					// Allocate some memory
					for j := 0; j < 100; j++ {
						ptr := a.Alloc(8)
						So(ptr, ShouldNotBeNil)
					}

					// Free and reuse
					a.Free()
				}
			}, ShouldNotPanic)
		})
	})
}

func TestSuggestSize(t *testing.T) {
	Convey("Given size suggestions", t, func() {
		Convey("When suggesting sizes", func() {
			// Test that sizes are rounded up to powers of 2
			So(arena.SuggestSize(0), ShouldEqual, 0)
			So(arena.SuggestSize(1), ShouldEqual, 64)
			So(arena.SuggestSize(63), ShouldEqual, 64)
			So(arena.SuggestSize(64), ShouldEqual, 64)
			So(arena.SuggestSize(65), ShouldEqual, 128)
			So(arena.SuggestSize(127), ShouldEqual, 128)
			So(arena.SuggestSize(128), ShouldEqual, 128)
		})

		Convey("When suggesting large sizes", func() {
			So(arena.SuggestSize(1024), ShouldEqual, 1024)
			So(arena.SuggestSize(1025), ShouldEqual, 2048)
			So(arena.SuggestSize(2047), ShouldEqual, 2048)
			So(arena.SuggestSize(1024*1024), ShouldEqual, 1024*1024)
		})
	})
}

func TestAllocTraceable(t *testing.T) {
	Convey("Given traceable allocations", t, func() {
		Convey("When allocating traceable memory", func() {
			a := &arena.Arena{}
			ptr := arena.AllocTraceable(100, unsafe.Pointer(a))

			So(ptr, ShouldNotBeNil)
			So(uintptr(unsafe.Pointer(ptr))%uintptr(arena.Align), ShouldEqual, uintptr(0))
		})

		Convey("When allocating different sizes", func() {
			a := &arena.Arena{}
			sizes := []int{8, 16, 32, 64, 128, 256, 512, 1024}

			for _, size := range sizes {
				ptr := arena.AllocTraceable(size, unsafe.Pointer(a))
				So(ptr, ShouldNotBeNil)
				So(uintptr(unsafe.Pointer(ptr))%uintptr(arena.Align), ShouldEqual, uintptr(0))
			}
		})
	})
}
