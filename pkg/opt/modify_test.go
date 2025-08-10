package opt_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/opt"
)

func TestModify(t *testing.T) {
	Convey("Given some new option", t, func() {
		some := Some(123)
		none := None[int]()

		Convey("Then insert value", func() {
			So(*some.Insert(456), ShouldEqual, 456)
			So(*none.Insert(456), ShouldEqual, 456)
		})

		Convey("Then get or insert value", func() {
			So(*some.GetOrInsert(456), ShouldEqual, 123)
			So(*none.GetOrInsert(456), ShouldEqual, 456)
		})

		Convey("Then get or insert default value", func() {
			So(*some.GetOrInsertDefault(), ShouldEqual, 123)
			So(*none.GetOrInsertDefault(), ShouldEqual, 0)
		})

		Convey("Then get or insert value with function", func() {
			So(*some.GetOrInsertWith(func() int { return 456 }), ShouldEqual, 123)
			So(*none.GetOrInsertWith(func() int { return 456 }), ShouldEqual, 456)
		})

		Convey("Then take the value", func() {
			So(some.Take(), ShouldEqual, Some(123))
			So(some.IsNone(), ShouldBeTrue)
			So(none.Take().IsNone(), ShouldBeTrue)
		})

		Convey("Then take the value if predicate evaluates to true", func() {
			So(some.TakeIf(func(i int) bool { return i < 0 }), ShouldEqual, none)
			So(some.IsSome(), ShouldBeTrue)

			So(some.TakeIf(func(i int) bool { return i > 0 }), ShouldEqual, Some(123))
			So(some.IsNone(), ShouldBeTrue)

			So(none.TakeIf(func(i int) bool { return i > 0 }).IsNone(), ShouldBeTrue)
		})

		Convey("Then replace the value", func() {
			So(some.Replace(456), ShouldEqual, Some(123))
			So(some, ShouldEqual, Some(456))

			So(none.Replace(456).IsNone(), ShouldBeTrue)
			So(none, ShouldEqual, Some(456))
		})
	})
}
