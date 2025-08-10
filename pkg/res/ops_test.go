package res_test

import (
	"fmt"
	"io"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/res"
)

func val[T any](v T) func() T { return func() T { return v } }

func ExampleMap() {
	sizeof := func(s string) int { return len(s) }

	ok := Ok("hello world!")
	fmt.Println(Map(ok, sizeof))

	err := Err[string](io.EOF)
	fmt.Println(Map(err, sizeof))

	// Output:
	// Ok(12)
	// Err(EOF)
}

func ExampleMapOr() {
	sizeof := func(s string) int { return len(s) }

	ok := Ok("hello world!")
	fmt.Println(MapOr(ok, -1, sizeof))

	err := Err[string](io.EOF)
	fmt.Println(MapOr(err, -1, sizeof))

	// Output:
	// 12
	// -1
}

func ExampleMapOrElse() {
	sizeof := func(s string) int { return len(s) }

	ok := Ok("hello world!")
	fmt.Println(MapOrElse(ok, func() int { return 0 }, sizeof))

	err := Err[string](io.EOF)
	fmt.Println(MapOrElse(err, func() int { return 0 }, sizeof))

	// Output:
	// 12
	// 0
}

func ExampleFlatten() {
	x := Ok(Ok(6))
	fmt.Println(Flatten(x))

	y := Ok(Err[int](io.EOF))
	fmt.Println(Flatten(y))

	// Output:
	// Ok(6)
	// Err(EOF)
}

func TestOps(t *testing.T) {
	Convey("Given ok results", t, func() {
		ok := Ok(123)
		ok2 := Ok(456)
		okStr := Ok("foobar")
		err := Err[int](io.EOF)

		double := func(v int) int { return v * 2 }
		foobar := func() string { return "foobar" }

		Convey("Then map the value", func() {
			So(Map(ok, strconv.Itoa), ShouldEqual, Ok("123"))
			So(Map(err, strconv.Itoa), ShouldEqual, Err[string](io.EOF))

			So(ok.Map(double), ShouldEqual, Ok(246))
			So(err.Map(double), ShouldEqual, err)

			So(MapOr(ok, "456", strconv.Itoa), ShouldEqual, "123")
			So(MapOr(err, "456", strconv.Itoa), ShouldEqual, "456")

			So(ok.MapOr(456, double), ShouldEqual, 246)
			So(err.MapOr(456, double), ShouldEqual, 456)

			So(MapOrElse(ok, foobar, strconv.Itoa), ShouldEqual, "123")
			So(MapOrElse(err, foobar, strconv.Itoa), ShouldEqual, "foobar")

			So(ok.MapOrElse(val(456), double), ShouldEqual, 246)
			So(err.MapOrElse(val(456), double), ShouldEqual, 456)

			So(ok.MapErr(func(err error) error { return fmt.Errorf("err: %w", err) }), ShouldEqual, ok)
			So(err.MapErr(func(err error) error { return fmt.Errorf("err: %w", err) }).UnwrapErr().Error(), ShouldEqual, "err: EOF")
		})

		Convey("Then inspect the value", func() {
			So(ok.Inspect(func(i int) { So(i, ShouldEqual, 123) }), ShouldEqual, ok)
			So(err.Inspect(func(i int) { t.FailNow() }), ShouldEqual, err)

			So(ok.InspectErr(func(e error) { t.FailNow() }), ShouldEqual, ok)
			So(err.InspectErr(func(e error) { So(e, ShouldEqual, io.EOF) }), ShouldEqual, err)
		})

		Convey("Then and two results", func() {
			So(And(ok, okStr), ShouldEqual, okStr)
			So(And(ok, err), ShouldEqual, err)
			So(And(err, ok), ShouldEqual, err)
			So(And(err, err), ShouldEqual, err)

			So(ok.And(ok2), ShouldEqual, ok2)
			So(ok.And(err), ShouldEqual, err)
			So(err.And(ok), ShouldEqual, err)
			So(err.And(err), ShouldEqual, err)
		})

		Convey("Then call a function on the result value", func() {
			So(AndThen(ok, func(v int) Result[string] { return Ok(strconv.Itoa(v)) }), ShouldEqual, Ok("123"))
			So(AndThen(err, func(v int) Result[string] { return Ok(strconv.Itoa(v)) }).IsErr(), ShouldBeTrue)

			So(ok.AndThen(func(v int) Result[int] { return ok2 }), ShouldEqual, ok2)
			So(err.AndThen(func(v int) Result[int] { return ok2 }).IsErr(), ShouldBeTrue)
		})

		Convey("Then or two results", func() {
			So(ok.Or(ok2), ShouldEqual, ok)
			So(ok.Or(err), ShouldEqual, ok)
			So(err.Or(ok), ShouldEqual, ok)
			So(err.Or(err), ShouldEqual, err)
		})

		Convey("Then call a function if the result is err", func() {
			So(ok.OrElse(func(error) Result[int] { return ok2 }), ShouldEqual, ok)
			So(err.OrElse(func(error) Result[int] { return ok2 }), ShouldEqual, ok2)
		})

		Convey("Then flatten the result", func() {
			So(Flatten(Ok(ok)), ShouldEqual, ok)
			So(Flatten(Ok(err)), ShouldEqual, err)
			So(Flatten(Err[Result[int]](io.EOF)), ShouldEqual, err)
		})
	})
}
