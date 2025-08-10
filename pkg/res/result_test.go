package res_test

import (
	"errors"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/res"
)

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

			So(Wrap(123, nil), ShouldEqual, Ok(123))
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

			So(Wrap(0, io.EOF), ShouldResemble, err)
		})
	})
}
