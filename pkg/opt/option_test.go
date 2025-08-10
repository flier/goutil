package opt_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/opt"
)

func TestOption(t *testing.T) {
	Convey("Given a new option", t, func() {
		some := Some(123)

		Convey("It should have some value", func() {
			So(some.IsSome(), ShouldBeTrue)
			So(some.IsSomeAnd(func(v int) bool { return v < 0 }), ShouldBeFalse)
			So(some.String(), ShouldEqual, "Some(123)")

			So(some.IsNone(), ShouldBeFalse)
			So(some.IsNoneOr(func(v int) bool { return v > 0 }), ShouldBeTrue)

			So(some.Expect("some value"), ShouldEqual, 123)
			So(some.Unwrap(), ShouldEqual, 123)
			So(some.UnwrapOr(456), ShouldEqual, 123)
			So(some.UnwrapOrElse(func() int { return 456 }), ShouldEqual, 123)
			So(some.UnwrapOrDefault(), ShouldEqual, 123)

			n := 123
			So(Wrap(&n), ShouldEqual, some)
		})

		none := None[int]()

		Convey("It should have no value", func() {
			So(none.IsSome(), ShouldBeFalse)
			So(none.IsSomeAnd(func(v int) bool { return v > 0 }), ShouldBeFalse)
			So(none.String(), ShouldEqual, "None")

			So(none.IsNone(), ShouldBeTrue)
			So(none.IsNoneOr(func(v int) bool { return false }), ShouldBeTrue)

			So(func() { none.Unwrap() }, ShouldPanic)
			So(func() { none.Expect("no value") }, ShouldPanicWith, "no value")
			So(none.UnwrapOr(456), ShouldEqual, 456)
			So(none.UnwrapOrElse(func() int { return 456 }), ShouldEqual, 456)
			So(none.UnwrapOrDefault(), ShouldEqual, 0)

			So(Wrap[int](nil), ShouldEqual, none)
		})
	})
}
