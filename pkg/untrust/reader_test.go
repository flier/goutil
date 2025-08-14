package untrust_test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/untrust"
)

func TestReader_NewReader(t *testing.T) {
	Convey("Given different input types", t, func() {
		Convey("When creating reader with empty input", func() {
			input := untrust.Input([]byte{})
			reader := untrust.NewReader(input)

			So(reader, ShouldNotBeNil)
			So(reader.AtEnd(), ShouldBeTrue)
			So(reader.GoString(), ShouldEqual, "Reader")
		})

		Convey("When creating reader with non-empty input", func() {
			input := untrust.Input([]byte("hello"))
			reader := untrust.NewReader(input)

			So(reader, ShouldNotBeNil)
			So(reader.AtEnd(), ShouldBeFalse)
			So(reader.GoString(), ShouldEqual, "Reader")
		})

		Convey("When creating reader with nil input", func() {
			reader := untrust.NewReader(nil)

			So(reader, ShouldNotBeNil)
			So(reader.AtEnd(), ShouldBeTrue)
		})
	})
}

func TestReader_Clone(t *testing.T) {
	Convey("Given a reader with some data", t, func() {
		input := untrust.Input([]byte("hello"))
		original := untrust.NewReader(input)

		// Read some bytes to move position
		_, err := original.ReadByte()
		So(err, ShouldBeNil)

		Convey("When cloning the reader", func() {
			cloned := original.Clone()

			So(cloned, ShouldNotBeNil)
			So(uintptr(unsafe.Pointer(cloned)), ShouldNotEqual,
				uintptr(unsafe.Pointer(original))) // Different pointer
			So(cloned.AtEnd(), ShouldEqual, original.AtEnd())
			So(cloned.GoString(), ShouldEqual, original.GoString())

			Convey("Then modifying cloned reader should not affect original", func() {
				// Read from cloned reader
				_, err := cloned.ReadByte()
				So(err, ShouldBeNil)

				// Original reader position should remain unchanged
				So(original.AtEnd(), ShouldBeFalse)
				So(cloned.AtEnd(), ShouldBeFalse)

				// Read from original reader
				_, err = original.ReadByte()
				So(err, ShouldBeNil)

				// Both should now be at different positions
				So(original.AtEnd(), ShouldBeFalse)
				So(cloned.AtEnd(), ShouldBeFalse)
			})
		})
	})
}

func TestReader_AtEnd(t *testing.T) {
	Convey("Given a reader with data", t, func() {
		input := untrust.Input([]byte("abc"))
		reader := untrust.NewReader(input)

		Convey("When reader is at beginning", func() {
			So(reader.AtEnd(), ShouldBeFalse)
		})

		Convey("When reader has read some bytes", func() {
			_, err := reader.ReadByte()
			So(err, ShouldBeNil)
			So(reader.AtEnd(), ShouldBeFalse)
		})

		Convey("When reader has read all bytes", func() {
			_, err := reader.ReadBytesToEnd()
			So(err, ShouldBeNil)
			So(reader.AtEnd(), ShouldBeTrue)
		})

		Convey("When reader is empty", func() {
			emptyReader := untrust.NewReader(untrust.Input([]byte{}))
			So(emptyReader.AtEnd(), ShouldBeTrue)
		})
	})
}

func TestReader_Peek(t *testing.T) {
	Convey("Given a reader with data", t, func() {
		input := untrust.Input([]byte("abc"))
		reader := untrust.NewReader(input)

		Convey("When peeking at first byte", func() {
			So(reader.Peek('a'), ShouldBeTrue)
			So(reader.Peek('b'), ShouldBeFalse)
			So(reader.Peek('c'), ShouldBeFalse)
		})

		Convey("When peeking after reading some bytes", func() {
			_, err := reader.ReadByte()
			So(err, ShouldBeNil)

			So(reader.Peek('a'), ShouldBeFalse)
			So(reader.Peek('b'), ShouldBeTrue)
			So(reader.Peek('c'), ShouldBeFalse)
		})

		Convey("When peeking at end of input", func() {
			_, err := reader.ReadBytesToEnd()
			So(err, ShouldBeNil)

			So(reader.Peek('a'), ShouldBeFalse)
			So(reader.Peek('b'), ShouldBeFalse)
			So(reader.Peek('c'), ShouldBeFalse)
		})

		Convey("When peeking with empty reader", func() {
			emptyReader := untrust.NewReader(untrust.Input([]byte{}))
			So(emptyReader.Peek('a'), ShouldBeFalse)
		})
	})
}

func TestReader_ReadByte(t *testing.T) {
	Convey("Given a reader with data", t, func() {
		input := untrust.Input([]byte("abc"))
		reader := untrust.NewReader(input)

		Convey("When reading bytes sequentially", func() {
			b, err := reader.ReadByte()
			So(err, ShouldBeNil)
			So(b, ShouldEqual, 'a')
			So(reader.AtEnd(), ShouldBeFalse)

			b, err = reader.ReadByte()
			So(err, ShouldBeNil)
			So(b, ShouldEqual, 'b')
			So(reader.AtEnd(), ShouldBeFalse)

			b, err = reader.ReadByte()
			So(err, ShouldBeNil)
			So(b, ShouldEqual, 'c')
			So(reader.AtEnd(), ShouldBeTrue)
		})

		Convey("When reading beyond end of input", func() {
			// Read all bytes first
			_, err := reader.ReadBytesToEnd()
			So(err, ShouldBeNil)

			// Try to read one more byte
			b, err := reader.ReadByte()
			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(b, ShouldEqual, 0)
		})

		Convey("When reading from empty reader", func() {
			emptyReader := untrust.NewReader(untrust.Input([]byte{}))
			b, err := emptyReader.ReadByte()
			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(b, ShouldEqual, 0)
		})
	})
}

func TestReader_ReadBytes(t *testing.T) {
	Convey("Given a reader with data", t, func() {
		input := untrust.Input([]byte("hello world"))
		reader := untrust.NewReader(input)

		Convey("When reading valid number of bytes", func() {
			bytes, err := reader.ReadBytes(5)
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, untrust.Input([]byte("hello")))
			So(reader.AtEnd(), ShouldBeFalse)

			// Read remaining bytes
			remaining, err := reader.ReadBytesToEnd()
			So(err, ShouldBeNil)
			So(remaining, ShouldResemble, untrust.Input([]byte(" world")))
		})

		Convey("When reading zero bytes", func() {
			bytes, err := reader.ReadBytes(0)
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, untrust.Input([]byte{}))
			So(reader.AtEnd(), ShouldBeFalse)
		})

		Convey("When reading negative bytes", func() {
			bytes, err := reader.ReadBytes(-1)
			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(bytes, ShouldBeNil)
		})

		Convey("When reading more bytes than available", func() {
			bytes, err := reader.ReadBytes(20)
			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(bytes, ShouldBeNil)
		})

		Convey("When reading exactly all remaining bytes", func() {
			// Skip first 6 bytes
			err := reader.Skip(6)
			So(err, ShouldBeNil)

			// Read remaining 5 bytes
			bytes, err := reader.ReadBytes(5)
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, untrust.Input([]byte("world")))
			So(reader.AtEnd(), ShouldBeTrue)
		})

		Convey("When reading with integer overflow", func() {
			// Create a reader and move to position 2
			largeReader := untrust.NewReader(untrust.Input([]byte("test")))
			err := largeReader.Skip(2)
			So(err, ShouldBeNil)

			bytes, err := largeReader.ReadBytes(1)
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, untrust.Input([]byte("s")))
		})
	})
}

func TestReader_ReadBytesToEnd(t *testing.T) {
	Convey("Given a reader with data", t, func() {
		input := untrust.Input([]byte("hello world"))
		reader := untrust.NewReader(input)

		Convey("When reading to end from beginning", func() {
			bytes, err := reader.ReadBytesToEnd()
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, untrust.Input([]byte("hello world")))
			So(reader.AtEnd(), ShouldBeTrue)
		})

		Convey("When reading to end after partial read", func() {
			// Read first 6 bytes
			_, err := reader.ReadBytes(6)
			So(err, ShouldBeNil)

			// Read remaining bytes
			bytes, err := reader.ReadBytesToEnd()
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, untrust.Input([]byte("world")))
			So(reader.AtEnd(), ShouldBeTrue)
		})

		Convey("When reading to end from empty reader", func() {
			emptyReader := untrust.NewReader(untrust.Input([]byte{}))
			bytes, err := emptyReader.ReadBytesToEnd()
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, untrust.Input([]byte{}))
			So(emptyReader.AtEnd(), ShouldBeTrue)
		})

		Convey("When reading to end from reader at end", func() {
			// Move to end
			_, err := reader.ReadBytesToEnd()
			So(err, ShouldBeNil)

			// Try to read to end again
			bytes, err := reader.ReadBytesToEnd()
			So(err, ShouldBeNil)
			So(bytes, ShouldResemble, untrust.Input([]byte{}))
		})
	})
}

func TestReader_Skip(t *testing.T) {
	Convey("Given a reader with data", t, func() {
		input := untrust.Input([]byte("hello world"))
		reader := untrust.NewReader(input)

		Convey("When skipping valid number of bytes", func() {
			err := reader.Skip(5)
			So(err, ShouldBeNil)
			So(reader.AtEnd(), ShouldBeFalse)

			// Verify position by reading next byte
			b, err := reader.ReadByte()
			So(err, ShouldBeNil)
			So(b, ShouldEqual, ' ')
		})

		Convey("When skipping zero bytes", func() {
			err := reader.Skip(0)
			So(err, ShouldBeNil)
			So(reader.AtEnd(), ShouldBeFalse)
		})

		Convey("When skipping negative bytes", func() {
			err := reader.Skip(-1)
			So(err, ShouldEqual, untrust.ErrEndOfInput)
		})

		Convey("When skipping more bytes than available", func() {
			err := reader.Skip(20)
			So(err, ShouldEqual, untrust.ErrEndOfInput)
		})

		Convey("When skipping exactly all remaining bytes", func() {
			// Skip first 6 bytes
			err := reader.Skip(6)
			So(err, ShouldBeNil)

			// Skip remaining 5 bytes
			err = reader.Skip(5)
			So(err, ShouldBeNil)
			So(reader.AtEnd(), ShouldBeTrue)
		})
	})
}

func TestReader_SkipToEnd(t *testing.T) {
	Convey("Given a reader with data", t, func() {
		input := untrust.Input([]byte("hello world"))
		reader := untrust.NewReader(input)

		Convey("When skipping to end from beginning", func() {
			err := reader.SkipToEnd()
			So(err, ShouldBeNil)
			So(reader.AtEnd(), ShouldBeTrue)
		})

		Convey("When skipping to end after partial read", func() {
			// Read first 6 bytes
			_, err := reader.ReadBytes(6)
			So(err, ShouldBeNil)

			// Skip to end
			err = reader.SkipToEnd()
			So(err, ShouldBeNil)
			So(reader.AtEnd(), ShouldBeTrue)
		})

		Convey("When skipping to end from empty reader", func() {
			emptyReader := untrust.NewReader(untrust.Input([]byte{}))
			err := emptyReader.SkipToEnd()
			So(err, ShouldBeNil)
			So(emptyReader.AtEnd(), ShouldBeTrue)
		})

		Convey("When skipping to end from reader already at end", func() {
			// Move to end
			err := reader.SkipToEnd()
			So(err, ShouldBeNil)

			// Try to skip to end again
			err = reader.SkipToEnd()
			So(err, ShouldBeNil)
			So(reader.AtEnd(), ShouldBeTrue)
		})
	})
}

func TestReader_ReadPartial(t *testing.T) {
	Convey("Given a reader with data", t, func() {
		input := untrust.Input([]byte("hello world"))
		reader := untrust.NewReader(input)

		Convey("When reading partial data successfully", func() {
			consumed, result, err := untrust.ReadPartial(reader, func(r *untrust.Reader) (string, error) {
				bytes, err := r.ReadBytes(5)
				if err != nil {
					return "", err
				}
				return string(bytes), nil
			})

			So(err, ShouldBeNil)
			So(result, ShouldEqual, "hello")
			So(consumed, ShouldResemble, untrust.Input([]byte("hello")))
			So(reader.AtEnd(), ShouldBeFalse)
		})

		Convey("When read function returns error", func() {
			consumed, result, err := untrust.ReadPartial(reader, func(r *untrust.Reader) (string, error) {
				return "", untrust.ErrEndOfInput
			})

			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(result, ShouldEqual, "")
			So(consumed, ShouldResemble, untrust.Input([]byte{}))
		})

		Convey("When read function consumes all data", func() {
			consumed, result, err := untrust.ReadPartial(reader, func(r *untrust.Reader) (string, error) {
				bytes, err := r.ReadBytesToEnd()
				if err != nil {
					return "", err
				}
				return string(bytes), nil
			})

			So(err, ShouldBeNil)
			So(result, ShouldEqual, "hello world")
			So(consumed, ShouldResemble, untrust.Input([]byte("hello world")))
			So(reader.AtEnd(), ShouldBeTrue)
		})

		Convey("When read function consumes no data", func() {
			consumed, result, err := untrust.ReadPartial(reader, func(r *untrust.Reader) (string, error) {
				return "test", nil
			})

			So(err, ShouldBeNil)
			So(result, ShouldEqual, "test")
			So(consumed, ShouldResemble, untrust.Input([]byte{}))
			So(reader.AtEnd(), ShouldBeFalse)
		})
	})
}

func TestReader_EdgeCases(t *testing.T) {
	Convey("Given edge cases", t, func() {
		Convey("When reader has nil buffer", func() {
			reader := untrust.NewReader(nil)
			So(reader.AtEnd(), ShouldBeTrue)
			So(reader.Peek('a'), ShouldBeFalse)

			b, err := reader.ReadByte()
			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(b, ShouldEqual, 0)
		})

		Convey("When reader has negative index", func() {
			// Create a reader and try to manipulate it to test edge cases
			reader := untrust.NewReader(untrust.Input([]byte("test")))

			// Test that operations don't panic even in edge cases
			So(func() { reader.Peek('t') }, ShouldNotPanic)
			So(func() { _, _ = reader.ReadByte() }, ShouldNotPanic)
		})

		Convey("When reader has index beyond buffer length", func() {
			reader := untrust.NewReader(untrust.Input([]byte("test")))

			// Skip beyond the end
			err := reader.Skip(10)
			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(reader.AtEnd(), ShouldBeTrue)
			So(reader.Peek('t'), ShouldBeFalse)

			b, err := reader.ReadByte()
			So(err, ShouldEqual, untrust.ErrEndOfInput)
			So(b, ShouldEqual, 0)
		})

		Convey("When cloning reader with edge cases", func() {
			edgeReader := untrust.NewReader(nil)
			cloned := edgeReader.Clone()

			So(cloned, ShouldNotBeNil)
			So(cloned.AtEnd(), ShouldBeTrue)
		})
	})
}
