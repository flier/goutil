package xunsafe_test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestIndirect(t *testing.T) {
	Convey("Given indirect type checks", t, func() {
		So(xunsafe.IsDirect[int](), ShouldBeFalse)
		So(xunsafe.IsDirect[string](), ShouldBeFalse)
		So(xunsafe.IsDirect[[]byte](), ShouldBeFalse)

		So(xunsafe.IsDirect[*int](), ShouldBeTrue)
		So(xunsafe.IsDirect[[1]*int](), ShouldBeTrue)
		So(xunsafe.IsDirect[any](), ShouldBeTrue)
		So(xunsafe.IsDirect[map[int]int](), ShouldBeTrue)
		So(xunsafe.IsDirect[chan int](), ShouldBeTrue)
		So(xunsafe.IsDirect[unsafe.Pointer](), ShouldBeTrue)
		So(xunsafe.IsDirect[struct{ _ *int }](), ShouldBeTrue)
		So(xunsafe.IsDirect[*struct{ _ *int }](), ShouldBeTrue)
	})
}

func TestAnyBytes(t *testing.T) {
	Convey("Given any bytes operations", t, func() {
		i := 0xaaaa
		p := &i
		So(xunsafe.IsDirectAny(i), ShouldBeFalse)
		So(xunsafe.IsDirectAny(p), ShouldBeTrue)

		So(xunsafe.AnyBytes(i), ShouldEqual, xunsafe.Bytes(&i))
		So(xunsafe.AnyBytes(p), ShouldEqual, xunsafe.Bytes(&p))

		p2 := struct{ p *int }{p}
		So(xunsafe.AnyBytes(p2), ShouldEqual, xunsafe.Bytes(&p2))
	})
}

func TestXunsafePC(t *testing.T) {
	Convey("Given a PC operation", t, func() {
		f := func() int { return 42 }
		pc := xunsafe.NewPC(f)

		t.Logf("%#x\n", pc)
		So(pc.Get()(), ShouldEqual, 42)
	})
}

func TestXunsafeComprehensive(t *testing.T) {
	Convey("Given comprehensive xunsafe tests", t, func() {
		// Test BitCast
		i := 42
		casted := xunsafe.BitCast[uint64](i)
		So(casted, ShouldEqual, uint64(42))

		// Test Ping
		ptr := &i
		xunsafe.Ping(ptr) // Should not panic

		// Test NoCopy
		var noCopy xunsafe.NoCopy
		_ = noCopy // Use the variable to avoid unused variable warning
	})
}
