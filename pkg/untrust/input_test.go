//go:build go1.20

package untrust_test

import (
	"io"
	"math"
	"strconv"
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/opt"
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
		Convey("When input is empty", func() {
			empty := untrust.Input([]byte(""))
			So(empty.Empty(), ShouldBeTrue)
			So(empty.Len(), ShouldEqual, 0)
			So(empty.Clone(), ShouldResemble, empty)
			So(empty.AsSliceLessSafe(), ShouldResemble, []byte(""))
		})

		Convey("When input has content", func() {
			foo := untrust.Input([]byte("foo"))
			So(foo.Empty(), ShouldBeFalse)
			So(foo.Len(), ShouldEqual, 3)
			So(foo.Clone(), ShouldResemble, foo)
			So(foo.AsSliceLessSafe(), ShouldResemble, []byte("foo"))
		})

		Convey("When input is nil", func() {
			var nilInput untrust.Input
			So(nilInput.Empty(), ShouldBeTrue)
			So(nilInput.Len(), ShouldEqual, 0)
			So(nilInput.Clone(), ShouldResemble, untrust.Input([]byte{}))
			So(nilInput.AsSliceLessSafe(), ShouldBeNil)
		})

		Convey("When input has special characters", func() {
			special := untrust.Input([]byte("hello\n\t\r\000world"))
			So(special.Empty(), ShouldBeFalse)
			So(special.Len(), ShouldEqual, 14) // \n, \t, \r, \000 are single bytes
			So(special.Clone(), ShouldResemble, special)
		})

		Convey("When input is very long", func() {
			longData := make([]byte, 10000)
			for i := range longData {
				longData[i] = byte(i % 256)
			}
			longInput := untrust.Input(longData)
			So(longInput.Empty(), ShouldBeFalse)
			So(longInput.Len(), ShouldEqual, 10000)
			So(longInput.Clone(), ShouldResemble, longInput)
		})
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

		Convey("Then read all with empty input", func() {
			emptyInput := untrust.Input([]byte{})
			_, err := untrust.ReadAll(emptyInput, untrust.ErrEndOfInput, func(r *untrust.Reader) (s any, err error) {
				So(r.AtEnd(), ShouldBeTrue)
				return
			})

			So(err, ShouldBeNil)
		})

		Convey("Then read all with nil input", func() {
			var nilInput untrust.Input
			_, err := untrust.ReadAll(nilInput, untrust.ErrEndOfInput, func(r *untrust.Reader) (s any, err error) {
				So(r.AtEnd(), ShouldBeTrue)
				return
			})

			So(err, ShouldBeNil)
		})
	})
}

func TestReadAllOptional(t *testing.T) {
	Convey("Given some input", t, func() {
		input := opt.Some(untrust.Input([]byte("foo")))

		Convey("Then read all input", func() {
			_, err := untrust.ReadAllOptional(input, untrust.ErrEndOfInput,
				func(r opt.Option[*untrust.Reader]) (s any, err error) {
					if r.IsSome() {
						r := r.Unwrap()

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
					}

					return
				})

			So(err, ShouldBeNil)
		})

		Convey("Then read all with unconsumed input", func() {
			_, err := untrust.ReadAllOptional(input, untrust.ErrEndOfInput,
				func(r opt.Option[*untrust.Reader]) (s any, err error) {
					if r.IsSome() {
						r := r.Unwrap()

						b, err := r.ReadByte()
						So(err, ShouldBeNil)
						So(b, ShouldEqual, 'f')

						So(r.AtEnd(), ShouldBeFalse)
					}

					return
				})

			So(err, ShouldWrap, untrust.ErrEndOfInput)
		})

		Convey("Then read all without input", func() {
			_, err := untrust.ReadAllOptional(opt.None[untrust.Input](), untrust.ErrEndOfInput,
				func(r opt.Option[*untrust.Reader]) (s any, err error) {
					So(r.IsNone(), ShouldBeTrue)

					return
				})

			So(err, ShouldBeNil)
		})

		Convey("Then read all with error returned", func() {
			_, err := untrust.ReadAllOptional(input, untrust.ErrEndOfInput,
				func(r opt.Option[*untrust.Reader]) (s any, err error) {
					err = io.ErrShortBuffer

					return
				})

			So(err, ShouldWrap, io.ErrShortBuffer)
		})

		Convey("Then read after skipping must not panic", func() {
			_, err := untrust.ReadAllOptional(input, untrust.ErrEndOfInput,
				func(r opt.Option[*untrust.Reader]) (s string, err error) {
					if r.IsSome() {
						r := r.Unwrap()

						_, err = r.ReadBytesToEnd()
						So(err, ShouldBeNil)

						_, err = r.ReadByte()
						So(err, ShouldWrap, untrust.ErrEndOfInput)

						_, err = r.ReadBytesToEnd()
					}

					return
				})

			So(err, ShouldBeNil)
		})

		Convey("Then read all with empty optional input", func() {
			emptyInput := opt.Some(untrust.Input([]byte{}))
			_, err := untrust.ReadAllOptional(emptyInput, untrust.ErrEndOfInput,
				func(r opt.Option[*untrust.Reader]) (s any, err error) {
					if r.IsSome() {
						reader := r.Unwrap()
						So(reader.AtEnd(), ShouldBeTrue)
					}
					return
				})

			So(err, ShouldBeNil)
		})

		Convey("Then read all with nil optional input", func() {
			var nilInput opt.Option[untrust.Input]
			_, err := untrust.ReadAllOptional(nilInput, untrust.ErrEndOfInput,
				func(r opt.Option[*untrust.Reader]) (s any, err error) {
					So(r.IsNone(), ShouldBeTrue)
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

		Convey("Then read partial with empty input", func() {
			emptyInput := untrust.Input([]byte{})
			reader := untrust.NewReader(emptyInput)

			consumed, result, err := untrust.ReadPartial(reader, func(r *untrust.Reader) (string, error) {
				return "empty", nil
			})

			So(err, ShouldBeNil)
			So(result, ShouldEqual, "empty")
			So(consumed, ShouldResemble, untrust.Input([]byte{}))
		})

		Convey("Then read partial with nil input", func() {
			var nilInput untrust.Input
			reader := untrust.NewReader(nilInput)

			consumed, result, err := untrust.ReadPartial(reader, func(r *untrust.Reader) (string, error) {
				return "nil", nil
			})

			So(err, ShouldBeNil)
			So(result, ShouldEqual, "nil")
			So(consumed, ShouldResemble, untrust.Input(nil))
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

		Convey("Then read a reader", func() {
			r := untrust.NewReader(input)

			b, err := r.ReadByte()
			So(err, ShouldBeNil)
			So(b, ShouldEqual, 'f')

			Convey("Then clone the reader", func() {
				cr := r.Clone()

				rest, err := cr.ReadBytesToEnd()
				So(err, ShouldBeNil)
				So(rest.AsSliceLessSafe(), ShouldEqual, []byte("oo"))
				So(cr.AtEnd(), ShouldBeTrue)

				Convey("Then read the original reader", func() {
					So(r.AtEnd(), ShouldBeFalse)

					rest, err := r.ReadBytesToEnd()
					So(err, ShouldBeNil)
					So(rest.AsSliceLessSafe(), ShouldEqual, []byte("oo"))
					So(r.AtEnd(), ShouldBeTrue)
				})
			})
		})

		Convey("Then read bytes with empty input", func() {
			emptyInput := untrust.Input([]byte{})
			reader := untrust.NewReader(emptyInput)

			bytes, err := reader.ReadBytes(1)
			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(bytes, ShouldBeNil)
		})

		Convey("Then read bytes with nil input", func() {
			var nilInput untrust.Input
			reader := untrust.NewReader(nilInput)

			bytes, err := reader.ReadBytes(1)
			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(bytes, ShouldBeNil)
		})
	})
}

func TestInput_Clone(t *testing.T) {
	Convey("Given different input types", t, func() {
		Convey("When cloning empty input", func() {
			empty := untrust.Input([]byte{})
			cloned := empty.Clone()

			So(cloned, ShouldResemble, empty)
			So(cloned, ShouldResemble, empty) // Same content, different slice
		})

		Convey("When cloning non-empty input", func() {
			original := untrust.Input([]byte("hello world"))
			cloned := original.Clone()

			So(cloned, ShouldResemble, original)
			So(uintptr(unsafe.Pointer(unsafe.SliceData(cloned))), ShouldNotEqual,
				uintptr(unsafe.Pointer(unsafe.SliceData(original)))) // Different slice
			So(string(cloned), ShouldEqual, "hello world")
		})

		Convey("When cloning nil input", func() {
			var nilInput untrust.Input
			cloned := nilInput.Clone()

			So(cloned, ShouldResemble, untrust.Input([]byte{}))
			So(cloned, ShouldNotBeNil)
		})

		Convey("When cloning large input", func() {
			largeData := make([]byte, 1000)
			for i := range largeData {
				largeData[i] = byte(i % 256)
			}
			original := untrust.Input(largeData)
			cloned := original.Clone()

			So(cloned, ShouldResemble, original)
			So(uintptr(unsafe.Pointer(unsafe.SliceData(cloned))), ShouldNotEqual,
				uintptr(unsafe.Pointer(unsafe.SliceData(original)))) // Different slice
			So(cloned.Len(), ShouldEqual, 1000)
		})
	})
}

func TestInput_AsSliceLessSafe(t *testing.T) {
	Convey("Given different input types", t, func() {
		Convey("When getting slice from empty input", func() {
			empty := untrust.Input([]byte{})
			slice := empty.AsSliceLessSafe()

			So(slice, ShouldResemble, []byte{})
			So(len(slice), ShouldEqual, 0)
		})

		Convey("When getting slice from non-empty input", func() {
			original := untrust.Input([]byte("test"))
			slice := original.AsSliceLessSafe()

			So(slice, ShouldResemble, []byte("test"))
			So(string(slice), ShouldEqual, "test")
		})

		Convey("When getting slice from nil input", func() {
			var nilInput untrust.Input
			slice := nilInput.AsSliceLessSafe()

			So(slice, ShouldBeNil)
		})

		Convey("When modifying returned slice", func() {
			original := untrust.Input([]byte("hello"))
			slice := original.AsSliceLessSafe()

			// Modify the slice
			slice[0] = 'H'

			// Original should be affected since it's the same underlying array
			So(string(original), ShouldEqual, "Hello")
		})
	})
}
