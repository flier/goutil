//go:build go1.22

package arena_test

import (
	"encoding/binary"
	"fmt"
	"math"
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

				So(a.Empty(), ShouldBeFalse)
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
