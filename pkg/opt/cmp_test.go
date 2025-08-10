//go:build go1.21

package opt_test

import (
	"slices"
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/opt"
)

func TestCmp(t *testing.T) {
	Convey("Given some new options", t, func() {
		some := Some(123)
		some2 := Some(456)
		none := None[int]()

		Convey("When sort slice of options", func() {
			s := []Option[int]{some2, none, some}

			Convey("Then sort the slices with sort.Sort", func() {
				So(IsSorted(s), ShouldBeFalse)

				sort.Sort(OptionSlice[int](s))

				So(IsSorted(s), ShouldBeTrue)
				So(s, ShouldResemble, []Option[int]{none, some, some2})
			})

			Convey("Then sort the slices with Sort", func() {
				So(IsSorted(s), ShouldBeFalse)

				Sort(OptionSlice[int](s))

				So(IsSorted(s), ShouldBeTrue)
				So(s, ShouldResemble, []Option[int]{none, some, some2})
			})

			Convey("Then sort the slices with OptionSlice", func() {
				So(IsSorted(s), ShouldBeFalse)

				os := OptionSlice[int](s)
				os.Sort()

				So(os.IsSorted(), ShouldBeTrue)
				So(s, ShouldResemble, []Option[int]{none, some, some2})
			})

			Convey("Then sort the slices with slices.SortFunc", func() {
				slices.SortFunc(s, Compare)

				So(slices.IsSortedFunc(s, Compare), ShouldBeTrue)
				So(s, ShouldResemble, []Option[int]{none, some, some2})
			})
		})

		Convey("When compare two options", func() {
			So(Compare(some, some2), ShouldBeLessThan, 0)
			So(Compare(none, some), ShouldBeLessThan, 0)
			So(Compare(some2, none), ShouldBeGreaterThan, 0)
			So(Compare(some, some), ShouldEqual, 0)
			So(Compare(none, none), ShouldEqual, 0)
		})

		Convey("When less two options", func() {
			So(Less(some, some2), ShouldBeTrue)
			So(Less(none, some), ShouldBeTrue)
			So(Less(some2, none), ShouldBeFalse)
			So(Less(some, some), ShouldBeFalse)
			So(Less(none, none), ShouldBeFalse)
		})

		Convey("When equal two options", func() {
			So(Equal(some, some), ShouldBeTrue)
			So(Equal(none, none), ShouldBeTrue)
			So(Equal(some, some2), ShouldBeFalse)
			So(Equal(none, some), ShouldBeFalse)
			So(Equal(some2, none), ShouldBeFalse)
		})
	})
}
