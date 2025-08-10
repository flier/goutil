package opt_test

import (
	"fmt"
	"io"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/opt"
	"github.com/flier/goutil/pkg/res"
	"github.com/flier/goutil/pkg/tuple"
)

func ExampleMap() {
	sizeof := func(s string) int { return len(s) }

	some := Some("hello world!")
	fmt.Println(Map(some, sizeof))

	none := None[string]()
	fmt.Println(Map(none, sizeof))

	// Output:
	// Some(12)
	// None
}

func ExampleMapOr() {
	sizeof := func(s string) int { return len(s) }

	some := Some("hello world!")
	fmt.Println(MapOr(some, 42, sizeof))

	none := None[string]()
	fmt.Println(MapOr(none, 42, sizeof))

	// Output:
	// 12
	// 42
}

func ExampleMapOrElse() {
	sizeof := func(s string) int { return len(s) }

	some := Some("hello world!")
	fmt.Println(MapOrElse(some, func() int { return 42 }, sizeof))

	none := None[string]()
	fmt.Println(MapOrElse(none, func() int { return 42 }, sizeof))

	// Output:
	// 12
	// 42
}

func ExampleFlatten() {
	x := Some(Some(6))
	fmt.Println(Flatten(x))

	y := Some(None[int]())
	fmt.Println(Flatten(y))

	// Output:
	// Some(6)
	// None
}

func ExampleZip() {
	x := Some(42)
	y := Some("hello")
	z := None[rune]()

	fmt.Println(Zip(x, y))
	fmt.Println(Zip(x, z))

	// Output:
	// Some((42, hello))
	// None
}

func ExampleZipWith() {
	type Point struct {
		x, y float64
	}

	NewPoint := func(x, y float64) Point { return Point{x, y} }

	x := Some(17.5)
	y := Some(42.7)
	z := None[float64]()

	fmt.Println(ZipWith(x, y, NewPoint))
	fmt.Println(ZipWith(x, z, NewPoint))

	// Output:
	// Some({17.5 42.7})
	// None
}

func ExampleUnzip() {
	x := Some(tuple.New2("hello", 42))
	y := None[tuple.Tuple2[string, int]]()

	fmt.Println(Unzip(x))
	fmt.Println(Unzip(y))

	// Output:
	// Some(hello) Some(42)
	// None None
}

func TestOps(t *testing.T) {
	Convey("Given some new options", t, func() {
		some := Some(123)
		some2 := Some(456)
		none := None[int]()
		someStr := Some("foobar")

		ok := res.Ok(123)
		err := res.Err[int](io.EOF)

		double := func(v int) int { return v * 2 }
		mul := func(x, y int) int { return x * y }

		Convey("When map the value", func() {
			So(Map(some, strconv.Itoa), ShouldEqual, Some("123"))
			So(Map(none, strconv.Itoa).IsNone(), ShouldBeTrue)

			So(some.Map(double), ShouldEqual, Some(246))
			So(none.Map(double).IsNone(), ShouldBeTrue)

			So(MapOr(some, 456, double), ShouldEqual, 246)
			So(MapOr(none, 456, double), ShouldEqual, 456)

			So(some.MapOr(456, double), ShouldEqual, 246)
			So(none.MapOr(456, double), ShouldEqual, 456)

			So(MapOrElse(some, func() string { return "456" }, strconv.Itoa), ShouldEqual, "123")
			So(MapOrElse(none, func() string { return "456" }, strconv.Itoa), ShouldEqual, "456")

			So(some.MapOrElse(func() int { return 456 }, double), ShouldEqual, 246)
			So(none.MapOrElse(func() int { return 456 }, double), ShouldEqual, 456)
		})

		Convey("When inspect the value", func() {
			So(some.Inspect(func(v int) { So(v, ShouldEqual, 123) }), ShouldEqual, some)
			So(none.Inspect(func(v int) { t.FailNow() }), ShouldEqual, none)
		})

		Convey("When convert Result[T] to Option[T]", func() {
			So(Ok(ok), ShouldEqual, some)
			So(Err(ok).IsNone(), ShouldBeTrue)

			So(Ok(err).IsNone(), ShouldBeTrue)
			So(Err(err), ShouldEqual, Some(io.EOF))
		})

		Convey("When convert Option[T] to Result[T]", func() {
			So(some.OkOr(io.EOF), ShouldEqual, ok)
			So(none.OkOr(io.EOF), ShouldEqual, err)

			So(some.OkOrElse(func() error { return io.EOF }), ShouldEqual, ok)
			So(none.OkOrElse(func() error { return io.EOF }), ShouldEqual, err)
		})

		Convey("When and two options", func() {
			So(And(some, someStr), ShouldEqual, someStr)
			So(And(some, none), ShouldEqual, none)
			So(And(none, some), ShouldEqual, none)
			So(And(none, none), ShouldEqual, none)

			So(some.And(some2), ShouldEqual, some2)
			So(some.And(none), ShouldEqual, none)
			So(none.And(some), ShouldEqual, none)
			So(none.And(none), ShouldEqual, none)
		})

		Convey("When call a function on the option value", func() {
			So(AndThen(some, func(v int) Option[string] { return Some(strconv.Itoa(v)) }), ShouldEqual, Some("123"))
			So(AndThen(none, func(v int) Option[string] { return Some(strconv.Itoa(v)) }).IsNone(), ShouldBeTrue)

			So(some.AndThen(func(v int) Option[int] { return some2 }), ShouldEqual, some2)
			So(none.AndThen(func(v int) Option[int] { return some2 }), ShouldEqual, none)
		})

		Convey("When filter the option", func() {
			So(some.Filter(func(v int) bool { return v > 0 }), ShouldEqual, some)
			So(some.Filter(func(v int) bool { return v < 0 }), ShouldEqual, none)
			So(none.Filter(func(v int) bool { return v > 0 }), ShouldEqual, none)
		})

		Convey("When or two options", func() {
			So(some.Or(some2), ShouldEqual, some)
			So(some.Or(none), ShouldEqual, some)
			So(none.Or(some), ShouldEqual, some)
			So(none.Or(none), ShouldEqual, none)
		})

		Convey("When call a function if the option is none", func() {
			So(some.OrElse(func() Option[int] { return some2 }), ShouldEqual, some)
			So(none.OrElse(func() Option[int] { return some2 }), ShouldEqual, some2)
		})

		Convey("When xor two options", func() {
			So(some.Xor(some2), ShouldEqual, none)
			So(some.Xor(none), ShouldEqual, some)
			So(none.Xor(some), ShouldEqual, some)
			So(none.Xor(none), ShouldEqual, none)
		})

		Convey("When flatten the option", func() {
			So(Flatten(Some(some)), ShouldEqual, some)
			So(Flatten(Some(none)), ShouldEqual, none)
			So(Flatten(None[Option[int]]()), ShouldEqual, none)
		})

		Convey("When zip two optoins", func() {
			So(Zip(some, someStr), ShouldEqual, Some(tuple.New2(123, "foobar")))
			So(Zip(some, none), ShouldEqual, None[tuple.Tuple2[int, int]]())
			So(Zip(none, someStr), ShouldEqual, None[tuple.Tuple2[int, string]]())
			So(Zip(none, none), ShouldEqual, None[tuple.Tuple2[int, int]]())
		})

		Convey("When zip two optoins with function", func() {
			So(ZipWith(some, some2, mul), ShouldEqual, Some(123*456))
			So(ZipWith(some, none, mul), ShouldEqual, none)
			So(ZipWith(none, some2, mul), ShouldEqual, none)
			So(ZipWith(none, none, mul), ShouldEqual, none)
		})

		Convey("When unzip to two optoins", func() {
			x, y := Unzip(Some(tuple.New2(123, "foobar")))
			So(x, ShouldEqual, some)
			So(y, ShouldEqual, someStr)
		})

		Convey("When unzip none", func() {
			x, y := Unzip(None[tuple.Tuple2[int, int]]())
			So(x, ShouldEqual, none)
			So(y, ShouldEqual, none)
		})
	})
}
