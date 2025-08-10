//go:build go1.23

package opt_test

import (
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/opt"
)

func TestOptionIter(t *testing.T) {
	Convey("Given some new options", t, func() {
		some := Some(123)
		none := None[int]()

		Convey("Then iterate the option", func() {
			So(slices.Collect(some.Iter()), ShouldResemble, []int{123})
			So(slices.Collect(none.Iter()), ShouldBeEmpty)
		})
	})
}
