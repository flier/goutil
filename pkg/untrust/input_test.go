package untrust_test

import (
	"io"
	"math"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/untrust"
)

func TestDebug(t *testing.T) {
	Convey("Given some input", t, func() {
		for _, s := range []string{"", "foo"} {
			Convey("When read input: "+strconv.Quote(s), func() {
				input := untrust.Input([]byte(s))

				So(input.GoString(), ShouldEqual, "Input")

				_, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (s string, err error) {
					So(r.GoString(), ShouldEqual, "Reader")
					err = r.SkipToEnd()
					So(r.GoString(), ShouldEqual, "Reader")

					return
				})

				So(err, ShouldBeNil)
			})
		}
	})
}

func TestInput(t *testing.T) {
	Convey("Given some input", t, func() {
		empty := untrust.Input([]byte(""))
		So(empty.Empty(), ShouldBeTrue)
		So(empty.Len(), ShouldEqual, 0)
		So(empty.Clone(), ShouldResemble, empty)
		So(empty.AsSliceLessSafe(), ShouldResemble, []byte(""))

		foo := untrust.Input([]byte("foo"))
		So(foo.Empty(), ShouldBeFalse)
		So(foo.Len(), ShouldEqual, 3)
		So(foo.Clone(), ShouldResemble, foo)
		So(foo.AsSliceLessSafe(), ShouldResemble, []byte("foo"))
	})
}

func TestReadAll(t *testing.T) {
	Convey("Given some input", t, func() {
		input := untrust.Input([]byte("foo"))

		Convey("Then read all input", func() {
			_, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (s any, err error) {
				So(r.Peek('f'), ShouldBeTrue)

				b, err := r.ReadByte()
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 'f')

				So(r.Peek('f'), ShouldBeFalse)

				b, err = r.ReadByte()
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 'o')

				b, err = r.ReadByte()
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 'o')

				So(r.AtEnd(), ShouldBeTrue)

				So(r.Peek('f'), ShouldBeFalse)

				return
			})

			So(err, ShouldBeNil)
		})

		Convey("Then read all with unconsumed input", func() {
			_, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (s any, err error) {
				b, err := r.ReadByte()
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 'f')

				So(r.AtEnd(), ShouldBeFalse)

				return
			})

			So(err, ShouldWrap, untrust.ErrEndOfInput)
		})

		Convey("Then read all with error returned", func() {
			_, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (s any, err error) {
				err = io.ErrShortBuffer

				return
			})

			So(err, ShouldWrap, io.ErrShortBuffer)
		})

		Convey("Then read after skipping must not panic", func() {
			_, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (s string, err error) {
				_, err = r.ReadBytesToEnd()
				So(err, ShouldBeNil)

				_, err = r.ReadByte()
				So(err, ShouldWrap, untrust.ErrEndOfInput)

				_, err = r.ReadBytesToEnd()

				return
			})

			So(err, ShouldBeNil)
		})
	})
}

func TestReadPartial(t *testing.T) {
	Convey("Given some input", t, func() {
		input := untrust.Input([]byte("foobar"))

		Convey("Then read partial", func() {
			r, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (s string, err error) {
				i, res, err := untrust.ReadPartial(r, func(r *untrust.Reader) (s string, err error) {
					res, err := r.ReadBytes(3)

					return string(res), err
				})
				So(err, ShouldBeNil)
				So(res, ShouldEqual, "foo")
				So(i, ShouldResemble, untrust.Input([]byte("foo")))

				i, err = r.ReadBytesToEnd()
				So(err, ShouldBeNil)
				So(i, ShouldResemble, untrust.Input([]byte("bar")))

				s = res + string(i)

				return
			})

			So(err, ShouldBeNil)
			So(r, ShouldEqual, "foobar")
		})
	})
}

func TestReadBytes(t *testing.T) {
	Convey("Given some input", t, func() {
		input := untrust.Input([]byte("foo"))

		Convey("Then read bytes", func() {
			res, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (res string, err error) {
				var buf untrust.Input

				buf, err = r.ReadBytes(2)
				So(err, ShouldBeNil)

				res = string(buf.AsSliceLessSafe())

				err = r.Skip(1)

				return
			})

			So(err, ShouldBeNil)
			So(res, ShouldEqual, "fo")
		})

		Convey("Then read overflowed bytes", func() {
			_, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (res string, err error) {
				_, err = r.ReadBytes(math.MaxInt)

				return
			})

			So(err, ShouldWrap, untrust.ErrEndOfInput)
		})

		Convey("Then read negative bytes", func() {
			_, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (res string, err error) {
				_, err = r.ReadBytes(-1)

				return
			})

			So(err, ShouldWrap, untrust.ErrEndOfInput)
		})

		Convey("Then read too many bytes", func() {
			_, err := untrust.ReadAll(input, untrust.ErrEndOfInput, func(r *untrust.Reader) (res string, err error) {
				_, err = r.ReadBytes(12)

				return
			})

			So(err, ShouldWrap, untrust.ErrEndOfInput)
		})
	})
}
