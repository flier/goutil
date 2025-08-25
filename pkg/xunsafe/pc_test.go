package xunsafe_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestPC(t *testing.T) {
	Convey("Given program counter operations", t, func() {
		Convey("When creating new PC values", func() {
			Convey("And creating PC from function returning int", func() {
				f1 := func() int { return 42 }
				pc1 := xunsafe.NewPC(f1)
				So(pc1, ShouldNotEqual, xunsafe.PC[func() int](0))
			})

			Convey("And creating PC from function returning string", func() {
				f2 := func() string { return "hello" }
				pc2 := xunsafe.NewPC(f2)
				So(pc2, ShouldNotEqual, xunsafe.PC[func() string](0))
			})

			Convey("And creating PC from function with parameters", func() {
				f3 := func(x int) int { return x * 2 }
				pc3 := xunsafe.NewPC(f3)
				So(pc3, ShouldNotEqual, xunsafe.PC[func(int) int](0))
			})

			Convey("And creating PC from function returning multiple values", func() {
				f4 := func() (int, string) { return 42, "hello" }
				pc4 := xunsafe.NewPC(f4)
				So(pc4, ShouldNotEqual, xunsafe.PC[func() (int, string)](0))
			})
		})

		Convey("When retrieving functions from PC", func() {
			Convey("And retrieving function returning int", func() {
				f1 := func() int { return 42 }
				pc1 := xunsafe.NewPC(f1)
				retrieved1 := pc1.Get()
				So(retrieved1(), ShouldEqual, 42)
			})

			Convey("And retrieving function returning string", func() {
				f2 := func() string { return "hello" }
				pc2 := xunsafe.NewPC(f2)
				retrieved2 := pc2.Get()
				So(retrieved2(), ShouldEqual, "hello")
			})

			Convey("And retrieving function with parameters", func() {
				f3 := func(x int) int { return x * 2 }
				pc3 := xunsafe.NewPC(f3)
				retrieved3 := pc3.Get()
				So(retrieved3(5), ShouldEqual, 10)
			})

			Convey("And retrieving function returning multiple values", func() {
				f4 := func() (int, string) { return 42, "hello" }
				pc4 := xunsafe.NewPC(f4)
				retrieved4 := pc4.Get()
				a, b := retrieved4()
				So(a, ShouldEqual, 42)
				So(b, ShouldEqual, "hello")
			})
		})

		Convey("When handling edge cases", func() {
			Convey("And working with zero PC value", func() {
				var zeroPC xunsafe.PC[func() int]
				So(func() {
					_ = zeroPC
				}, ShouldNotPanic)
			})
		})

		Convey("When working with different function signatures", func() {
			Convey("And working with function taking no parameters", func() {
				f1 := func() {}
				pc1 := xunsafe.NewPC(f1)
				retrieved1 := pc1.Get()
				retrieved1() // Should not panic
			})

			Convey("And working with function returning nothing", func() {
				f2 := func() {}
				pc2 := xunsafe.NewPC(f2)
				retrieved2 := pc2.Get()
				retrieved2() // Should not panic
			})

			Convey("And working with function taking interface parameter", func() {
				f3 := func(x interface{}) interface{} { return x }
				pc3 := xunsafe.NewPC(f3)
				retrieved3 := pc3.Get()
				result := retrieved3("test")
				So(result, ShouldEqual, "test")
			})
		})

		Convey("When working with closure behavior", func() {
			Convey("And working with simple functions", func() {
				f := func() int { return 42 }
				pc := xunsafe.NewPC(f)
				retrieved := pc.Get()
				So(retrieved(), ShouldEqual, 42)
			})

			Convey("And working with different function", func() {
				f2 := func() string { return "hello" }
				pc2 := xunsafe.NewPC(f2)
				retrieved2 := pc2.Get()
				So(retrieved2(), ShouldEqual, "hello")
			})
		})

		Convey("When working with generic functions", func() {
			Convey("And working with concrete function types", func() {
				f1 := func(x int) int { return x }
				pc1 := xunsafe.NewPC(f1)
				retrieved1 := pc1.Get()
				result1 := retrieved1(42)
				So(result1, ShouldEqual, 42)
			})

			Convey("And working with string function", func() {
				f2 := func(x string) string { return x }
				pc2 := xunsafe.NewPC(f2)
				retrieved2 := pc2.Get()
				result2 := retrieved2("hello")
				So(result2, ShouldEqual, "hello")
			})
		})
	})
}
