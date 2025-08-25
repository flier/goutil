package xunsafe_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestVLA_Beyond(t *testing.T) {
	Convey("Given a VLA beyond operation", t, func() {
		// Test with simple header
		type Header struct {
			ID   int
			Name string
		}

		header := &Header{ID: 1, Name: "test"}
		vla := xunsafe.Beyond[int](header)

		So(vla, ShouldNotBeNil)
	})
}

func TestVLA_Get(t *testing.T) {
	Convey("Given a VLA Get operation", t, func() {
		// Test with simple header
		type Header struct {
			ID   int
			Name string
		}

		header := &Header{ID: 1, Name: "test"}
		vla := xunsafe.Beyond[int](header)

		// Test Get operation
		ptr := vla.Get(0)
		So(ptr, ShouldNotBeNil)

		// Test Get with different indices
		ptr1 := vla.Get(1)
		So(ptr1, ShouldNotBeNil)
	})
}

func TestVLA_ByteGet(t *testing.T) {
	Convey("Given a VLA ByteGet operation", t, func() {
		// Test with simple header
		type Header struct {
			ID   int
			Name string
		}

		header := &Header{ID: 1, Name: "test"}
		vla := xunsafe.Beyond[int](header)

		// Test ByteGet operation
		ptr := vla.ByteGet(0)
		So(ptr, ShouldNotBeNil)

		// Test ByteGet with different byte offsets
		ptr1 := vla.ByteGet(8)
		So(ptr1, ShouldNotBeNil)
	})
}

func TestVLA_Slice(t *testing.T) {
	Convey("Given a VLA Slice operation", t, func() {
		// Test with simple header
		type Header struct {
			ID   int
			Name string
		}

		header := &Header{ID: 1, Name: "test"}
		vla := xunsafe.Beyond[int](header)

		// Test Slice operation
		slice := vla.Slice(5)
		So(slice, ShouldNotBeNil)
		So(len(slice), ShouldEqual, 5)
	})
}

func TestVLA_Len(t *testing.T) {
	Convey("Given a VLA Len operation", t, func() {
		// Test with simple header
		type Header struct {
			ID   int
			Name string
		}

		header := &Header{ID: 1, Name: "test"}
		vla := xunsafe.Beyond[int](header)

		// Test Len operation
		length := vla.Len()
		So(length, ShouldEqual, 0)
	})
}

func TestVLA_Comprehensive(t *testing.T) {
	Convey("Given comprehensive VLA tests", t, func() {
		// Test with simple header
		header1 := &struct{ ID int }{ID: 1}
		vla1 := xunsafe.Beyond[int](header1)
		So(vla1, ShouldNotBeNil)

		// Test with complex header
		header2 := &struct {
			ID   int
			Name string
			Data []byte
		}{ID: 1, Name: "test", Data: []byte{1, 2, 3}}
		vla2 := xunsafe.Beyond[int](header2)
		So(vla2, ShouldNotBeNil)

		// Test with empty header
		header3 := &struct{}{}
		vla3 := xunsafe.Beyond[int](header3)
		So(vla3, ShouldNotBeNil)

		// Test operations on each VLA
		for i, vla := range []*xunsafe.VLA[int]{vla1, vla2, vla3} {
			Convey(fmt.Sprintf("When testing vla_%d", i), func() {
				// Test Get operation
				ptr := vla.Get(0)
				So(ptr, ShouldNotBeNil)

				// Test ByteGet operation
				bytePtr := vla.ByteGet(0)
				So(bytePtr, ShouldNotBeNil)

				// Test Slice operation
				slice := vla.Slice(3)
				So(slice, ShouldNotBeNil)
				So(len(slice), ShouldEqual, 3)

				// Test Len operation
				length := vla.Len()
				So(length, ShouldEqual, 0)
			})
		}
	})
}
