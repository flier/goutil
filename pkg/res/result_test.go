package res_test

import (
	"errors"
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

func TestResult(t *testing.T) {
	Convey("Given a new result", t, func() {
		ok := Ok(123)

		isNeg := func(v int) bool { return v < 0 }
		isEof := func(v error) bool { return errors.Is(v, io.EOF) }

		Convey("It should be ok", func() {
			So(ok.IsOk(), ShouldBeTrue)
			So(ok.IsOkAnd(isNeg), ShouldBeFalse)
			So(ok.IsErr(), ShouldBeFalse)
			So(ok.IsErrAnd(isEof), ShouldBeFalse)

			So(ok.String(), ShouldEqual, "Ok(123)")

			So(ok.Expect("value"), ShouldEqual, 123)
			So(func() { _ = ok.ExpectErr("err") }, ShouldPanicWith, "err: 123")
			So(ok.Unwrap(), ShouldEqual, 123)
			So(func() { _ = ok.UnwrapErr() }, ShouldPanic)

			So(ok.UnwrapOr(456), ShouldEqual, 123)
			So(ok.UnwrapOrElse(func() int { return 456 }), ShouldEqual, 123)
			So(ok.UnwrapOrDefault(), ShouldEqual, 123)
		})

		err := Err[int](io.EOF)

		Convey("It should be err", func() {
			So(err.IsOk(), ShouldBeFalse)
			So(err.IsOkAnd(isNeg), ShouldBeFalse)
			So(err.IsErr(), ShouldBeTrue)
			So(err.IsErrAnd(isEof), ShouldBeTrue)

			So(err.String(), ShouldEqual, "Err(EOF)")

			So(func() { err.Expect("value") }, ShouldPanicWith, "value: EOF")
			So(err.ExpectErr("err"), ShouldEqual, io.EOF)
			So(func() { err.Unwrap() }, ShouldPanic)
			So(err.UnwrapErr(), ShouldEqual, io.EOF)

			So(err.UnwrapOr(456), ShouldEqual, 456)
			So(err.UnwrapOrElse(func() int { return 456 }), ShouldEqual, 456)
			So(err.UnwrapOrDefault(), ShouldEqual, 0)
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
