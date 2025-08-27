package xerrors_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/xerrors"
)

type CustomError struct {
	message string
}

func (e CustomError) Error() string {
	return e.message
}

type AnotherError struct {
	code int
	msg  string
}

func (e *AnotherError) Error() string {
	return e.msg
}

func TestAsA(t *testing.T) {
	Convey("Given a custom error type", t, func() {
		err := CustomError{message: "test error"}
		aerr := &AnotherError{code: 1, msg: "another error"}

		Convey("Should work with direct error (value)", func() {
			e, ok := AsA[CustomError](err)

			So(ok, ShouldBeTrue)
			So(e, ShouldEqual, err)
		})

		Convey("Should work with direct error (pointer)", func() {
			e, ok := AsA[*AnotherError](aerr)

			So(ok, ShouldBeTrue)
			So(e, ShouldEqual, aerr)
		})

		Convey("Should work with wrapped error (value)", func() {
			wrappedErr := fmt.Errorf("wrapped: %w", err)

			e, ok := AsA[CustomError](wrappedErr)

			So(ok, ShouldBeTrue)
			So(e, ShouldEqual, err)
		})

		Convey("Should work with wrapped error type (pointer)", func() {
			e, ok := AsA[*AnotherError](aerr)

			So(ok, ShouldBeTrue)
			So(e, ShouldEqual, aerr)
		})

		Convey("Should work with multiple layers of wrapping (value)", func() {
			err1 := fmt.Errorf("first: %w", err)
			err2 := fmt.Errorf("custom: %w", err1)

			e, ok := AsA[CustomError](err2)

			So(ok, ShouldBeTrue)
			So(e, ShouldEqual, err)
		})

		Convey("Should work with multiple layers of wrapping (pointer)", func() {
			err1 := fmt.Errorf("first: %w", aerr)
			err2 := fmt.Errorf("custom: %w", err1)

			e, ok := AsA[*AnotherError](err2)

			So(ok, ShouldBeTrue)
			So(e, ShouldEqual, aerr)
		})

		Convey("Should not work with non-matching type", func() {
			_, ok := AsA[CustomError](aerr)

			So(ok, ShouldBeFalse)
		})
	})
}
