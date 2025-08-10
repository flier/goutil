//go:build go1.23

package res_test

import (
	"fmt"
	"io"
	"slices"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/res"
)

func ExampleCollect2() {
	// Some I/O results
	s := slices.All([]error{nil, io.EOF, nil})

	res, err := Collect2(s)
	fmt.Println(err, res)
	// Output:
	// EOF []
}

func ExampleCollect_ok() {
	// Some I/O results
	s := slices.Values([]Result[string]{
		Ok("hello"),
		Ok("world"),
	})

	fmt.Println(Collect(s))
	// Output: [hello world] <nil>
}

func ExampleCollect_err() {
	// Some I/O results
	s := slices.Values([]Result[string]{
		Ok("hello"),
		Err[string](io.EOF),
		Ok("world"),
	})

	fmt.Println(Collect(s))
	// Output: [] EOF
}

func TestIter(t *testing.T) {
	Convey("Given ok results", t, func() {
		ok := Ok(123)
		err := Err[int](io.EOF)

		Convey("Then iterate the result", func() {
			So(slices.Collect(ok.Iter()), ShouldResemble, []int{123})
			So(slices.Collect(err.Iter()), ShouldBeEmpty)
		})
	})
}

func TestCollect(t *testing.T) {
	Convey("Given some results", t, func() {
		results := []Result[int]{
			Ok(1),
			Ok(2),
			Ok(3),
			Ok(4),
		}

		Convey("It should collect ok values", func() {
			values, e := Collect(slices.Values(results))
			So(e, ShouldBeNil)
			So(values, ShouldResemble, []int{1, 2, 3, 4})

			_, e = Collect(slices.Values(append(results, Err[int](io.EOF))))
			So(e, ShouldWrap, io.EOF)
		})

		Convey("It should collect err values", func() {
			errs := []error{nil, nil, nil, nil}
			values, e := Collect2(slices.All(errs))
			So(e, ShouldBeNil)
			So(values, ShouldResemble, []int{0, 1, 2, 3})

			_, e = Collect2(slices.All(append(errs, io.EOF)))
			So(e, ShouldWrap, io.EOF)
		})
	})
}
