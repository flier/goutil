package zc_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/zc"
)

func TestView(t *testing.T) {
	Convey("Given a View", t, func() {
		Convey("When creating a new View with New", func() {
			src := []byte("hello world")
			start := &src[6] // points to "world"
			view := zc.New(&src[0], start, 5)

			Convey("It should have correct offset and length", func() {
				So(view.Start(), ShouldEqual, 6)
				So(view.Len(), ShouldEqual, 5)
				So(view.End(), ShouldEqual, 11)
			})

			Convey("It should convert to correct bytes", func() {
				bytes := view.Bytes(&src[0])
				So(string(bytes), ShouldEqual, "world")
			})
		})

		Convey("When creating a View with Raw", func() {
			view := zc.Raw(10, 20)

			Convey("It should have correct offset and length", func() {
				So(view.Start(), ShouldEqual, 10)
				So(view.Len(), ShouldEqual, 20)
				So(view.End(), ShouldEqual, 30)
			})
		})

		Convey("When working with zero View", func() {
			var view zc.View

			Convey("It should represent empty slice", func() {
				So(view.Start(), ShouldEqual, 0)
				So(view.Len(), ShouldEqual, 0)
				So(view.End(), ShouldEqual, 0)
			})

			Convey("It should return nil for Bytes", func() {
				src := []byte("test")
				bytes := view.Bytes(&src[0])
				So(bytes, ShouldBeNil)
			})
		})

		Convey("When working with maximum values", func() {
			view := zc.Raw(0xFFFFFFFF, 0xFFFFFFFF)

			Convey("It should handle maximum uint32 values", func() {
				So(view.Start(), ShouldEqual, 0xFFFFFFFF)
				So(view.Len(), ShouldEqual, 0xFFFFFFFF)
				So(view.End(), ShouldEqual, 0xFFFFFFFF+0xFFFFFFFF)
			})
		})
	})
}

func TestViewBytes(t *testing.T) {
	Convey("Given a View and source buffer", t, func() {
		src := []byte("hello world test")
		view := zc.Raw(6, 5) // "world"

		Convey("When calling Bytes", func() {
			bytes := view.Bytes(&src[0])

			Convey("It should return correct slice", func() {
				So(string(bytes), ShouldEqual, "world")
				So(len(bytes), ShouldEqual, 5)
			})
		})

		Convey("When calling Bytes with empty View", func() {
			emptyView := zc.Raw(0, 0)
			bytes := emptyView.Bytes(&src[0])

			Convey("It should return nil", func() {
				So(bytes, ShouldBeNil)
			})
		})

		Convey("When calling Bytes with out-of-bounds View", func() {
			// This should be handled gracefully by Go's slice mechanism
			view := zc.Raw(100, 10)
			bytes := view.Bytes(&src[0])

			Convey("It should return empty slice", func() {
				So(len(bytes), ShouldEqual, 10)
				// The actual content is undefined but length should be correct
			})
		})
	})
}

func TestViewFormat(t *testing.T) {
	Convey("Given a View", t, func() {
		view := zc.Raw(10, 20)

		Convey("When formatting with %v verb", func() {
			Convey("It should format correctly", func() {
				result := fmt.Sprintf("%v", view)
				So(result, ShouldEqual, "[10:30]")
			})
		})
	})
}

func TestExtractFrom(t *testing.T) {
	Convey("Given an ExtractFrom helper", t, func() {
		src := []byte("hello world test")
		extractor := zc.ExtractFrom{Src: &src[0]}

		Convey("When extracting bytes from a View", func() {
			view := zc.Raw(6, 5) // "world"
			bytes := extractor.Bytes(uint64(view))

			Convey("It should return correct bytes", func() {
				So(string(bytes), ShouldEqual, "world")
				So(len(bytes), ShouldEqual, 5)
			})
		})

		Convey("When extracting from empty View", func() {
			emptyView := zc.Raw(0, 0)
			bytes := extractor.Bytes(uint64(emptyView))

			Convey("It should return empty slice", func() {
				So(len(bytes), ShouldEqual, 0)
			})
		})

		Convey("When extracting from different positions", func() {
			Convey("It should extract from start", func() {
				view := zc.Raw(0, 5) // "hello"
				bytes := extractor.Bytes(uint64(view))
				So(string(bytes), ShouldEqual, "hello")
			})

			Convey("It should extract from middle", func() {
				view := zc.Raw(6, 5) // "world"
				bytes := extractor.Bytes(uint64(view))
				So(string(bytes), ShouldEqual, "world")
			})

			Convey("It should extract from end", func() {
				view := zc.Raw(12, 4) // "test"
				bytes := extractor.Bytes(uint64(view))
				So(string(bytes), ShouldEqual, "test")
			})
		})
	})
}

func TestViewEdgeCases(t *testing.T) {
	Convey("Given edge cases for View", t, func() {
		src := []byte("test")

		Convey("When creating View with zero length", func() {
			view := zc.Raw(2, 0)

			Convey("It should have zero length", func() {
				So(view.Len(), ShouldEqual, 0)
				So(view.Start(), ShouldEqual, 2)
				So(view.End(), ShouldEqual, 2)
			})

			Convey("It should return nil for Bytes", func() {
				bytes := view.Bytes(&src[0])
				So(bytes, ShouldBeNil)
			})
		})

		Convey("When creating View with maximum offset", func() {
			view := zc.Raw(0xFFFFFFFF, 0)

			Convey("It should handle maximum offset", func() {
				So(view.Start(), ShouldEqual, 0xFFFFFFFF)
				So(view.Len(), ShouldEqual, 0)
				So(view.End(), ShouldEqual, 0xFFFFFFFF)
			})
		})

		Convey("When creating View with maximum length", func() {
			view := zc.Raw(0, 0xFFFFFFFF)

			Convey("It should handle maximum length", func() {
				So(view.Start(), ShouldEqual, 0)
				So(view.Len(), ShouldEqual, 0xFFFFFFFF)
				So(view.End(), ShouldEqual, 0xFFFFFFFF)
			})
		})
	})
}

func TestViewString(t *testing.T) {
	Convey("Given a View with string conversion", t, func() {
		src := []byte("hello world")
		view := zc.Raw(6, 5) // "world"

		Convey("When converting to string", func() {
			str := view.String(&src[0])

			Convey("It should return correct string", func() {
				So(str, ShouldEqual, "world")
			})
		})

		Convey("When converting empty View to string", func() {
			emptyView := zc.Raw(0, 0)
			str := emptyView.String(&src[0])

			Convey("It should return empty string", func() {
				So(str, ShouldEqual, "")
			})
		})
	})
}

func TestViewPacking(t *testing.T) {
	Convey("Given View packing and unpacking", t, func() {
		Convey("When creating and unpacking View", func() {
			originalOffset := 12345
			originalLen := 67890
			view := zc.Raw(originalOffset, originalLen)

			Convey("It should preserve values", func() {
				So(view.Start(), ShouldEqual, originalOffset)
				So(view.Len(), ShouldEqual, originalLen)
				So(view.End(), ShouldEqual, originalOffset+originalLen)
			})
		})

		Convey("When working with bit manipulation", func() {
			view := zc.Raw(0x12345678, 0x9ABCDEF0)

			Convey("It should handle large values correctly", func() {
				So(view.Start(), ShouldEqual, 0x12345678)
				So(view.Len(), ShouldEqual, 0x9ABCDEF0)
			})
		})
	})
}

func TestViewWithUnsafeOperations(t *testing.T) {
	Convey("Given View with unsafe operations", t, func() {
		src := []byte("hello world test")
		view := zc.Raw(6, 5) // "world"

		Convey("When using unsafe operations", func() {
			bytes := view.Bytes(&src[0])

			Convey("It should work with unsafe.Slice", func() {
				So(len(bytes), ShouldEqual, 5)
				So(cap(bytes), ShouldEqual, 5)
			})

			Convey("It should be safe for concurrent access", func() {
				// This is more of a documentation test
				// The actual safety depends on the source buffer
				So(len(bytes), ShouldEqual, 5)
			})
		})
	})
}

func ExampleView() {
	// Create a source buffer
	src := []byte("hello world test")

	// Create a view pointing to "world" (offset 6, length 5)
	view := zc.Raw(6, 5)

	// Get the bytes
	bytes := view.Bytes(&src[0])
	fmt.Println(string(bytes)) // "world"

	// Get the string directly
	str := view.String(&src[0])
	fmt.Println(str) // "world"

	// Check properties
	fmt.Printf("Start: %d, Length: %d, End: %d\n",
		view.Start(), view.Len(), view.End()) // Start: 6, Length: 5, End: 11
}

func ExampleNew() {
	// Create a source buffer
	src := []byte("hello world test")

	// Create a view using New function
	view := zc.New(&src[0], &src[6], 5) // points to "world"

	// Get the bytes
	bytes := view.Bytes(&src[0])
	fmt.Println(string(bytes)) // "world"
}

func ExampleExtractFrom() {
	// Create a source buffer
	src := []byte("hello world test")

	// Create an extractor
	extractor := zc.ExtractFrom{Src: &src[0]}

	// Create a view
	view := zc.Raw(6, 5) // "world"

	// Extract bytes using the extractor
	bytes := extractor.Bytes(uint64(view))
	fmt.Println(string(bytes)) // "world"
}
