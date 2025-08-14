//go:build go1.22

package slice_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

func TestSlice_Of(t *testing.T) {
	Convey("Given an arena", t, func() {
		a := &arena.Arena{}

		Convey("When creating slice from values", func() {
			values := []int{1, 2, 3, 4, 5}
			s := slice.Of(a, values...)

			So(s.Len(), ShouldEqual, 5)
			So(s.Cap(), ShouldBeGreaterThanOrEqualTo, 5)
			So(s.Ptr(), ShouldNotBeNil)

			// Verify values
			for i, expected := range values {
				So(s.Load(i), ShouldEqual, expected)
			}
		})

		Convey("When creating slice from empty values", func() {
			s := slice.Of[int](a)

			So(s.Len(), ShouldEqual, 0)
			So(s.Cap(), ShouldBeGreaterThanOrEqualTo, 0)
		})

		Convey("When creating slice from string values", func() {
			values := []string{"hello", "world", "test"}
			s := slice.Of(a, values...)

			So(s.Len(), ShouldEqual, 3)
			So(s.Load(0), ShouldEqual, "hello")
			So(s.Load(1), ShouldEqual, "world")
			So(s.Load(2), ShouldEqual, "test")
		})
	})
}

func TestSlice_Make(t *testing.T) {
	Convey("Given an arena", t, func() {
		a := &arena.Arena{}

		Convey("When making slice with specific length", func() {
			s := slice.Make[int](a, 10)

			So(s.Len(), ShouldEqual, 10)
			So(s.Cap(), ShouldBeGreaterThanOrEqualTo, 10)
			So(s.Ptr(), ShouldNotBeNil)
		})

		Convey("When making slice with zero length", func() {
			s := slice.Make[int](a, 0)

			So(s.Len(), ShouldEqual, 0)
			So(s.Cap(), ShouldBeGreaterThanOrEqualTo, 0)
		})

		Convey("When making slice with large length", func() {
			s := slice.Make[byte](a, 1000)

			So(s.Len(), ShouldEqual, 1000)
			So(s.Cap(), ShouldBeGreaterThanOrEqualTo, 1000)
		})
	})
}

func TestSlice_FromParts(t *testing.T) {
	Convey("Given slice parts", t, func() {
		Convey("When creating slice from parts", func() {
			var ptr *int
			len := uint32(5)
			cap := uint32(10)

			s := slice.FromParts(ptr, len, cap)

			So(s.Ptr(), ShouldEqual, ptr)
			So(s.Len(), ShouldEqual, 5)
			So(s.Cap(), ShouldEqual, 10)
		})
	})
}

func TestSlice_LoadStore(t *testing.T) {
	Convey("Given a slice with data", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 5)

		Convey("When storing and loading values", func() {
			// Store values
			s.Store(0, 100)
			s.Store(1, 200)
			s.Store(2, 300)

			// Load values
			So(s.Load(0), ShouldEqual, 100)
			So(s.Load(1), ShouldEqual, 200)
			So(s.Load(2), ShouldEqual, 300)
		})

		Convey("When storing and loading at boundaries", func() {
			// Store at first position
			s.Store(0, 999)
			So(s.Load(0), ShouldEqual, 999)

			// Store at last position
			s.Store(4, 888)
			So(s.Load(4), ShouldEqual, 888)
		})
	})
}

func TestSlice_Append(t *testing.T) {
	Convey("Given a slice", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 3)

		// Initialize with some values
		s.Store(0, 1)
		s.Store(1, 2)
		s.Store(2, 3)

		Convey("When appending elements", func() {
			s = s.Append(a, 4, 5, 6)

			So(s.Len(), ShouldEqual, 6)
			So(s.Load(3), ShouldEqual, 4)
			So(s.Load(4), ShouldEqual, 5)
			So(s.Load(5), ShouldEqual, 6)
		})

		Convey("When appending to empty slice", func() {
			empty := slice.Make[int](a, 0)
			empty = empty.Append(a, 1, 2, 3)

			So(empty.Len(), ShouldEqual, 3)
			So(empty.Load(0), ShouldEqual, 1)
			So(empty.Load(1), ShouldEqual, 2)
			So(empty.Load(2), ShouldEqual, 3)
		})

		Convey("When appending single element", func() {
			s = s.AppendOne(a, 4)

			So(s.Len(), ShouldEqual, 4)
			So(s.Load(3), ShouldEqual, 4)
		})
	})
}

func TestSlice_Grow(t *testing.T) {
	Convey("Given a slice", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 2)

		Convey("When growing capacity", func() {
			initialCap := s.Cap()
			s = s.Grow(a, 5)

			So(s.Cap(), ShouldBeGreaterThan, initialCap)
			So(s.Len(), ShouldEqual, 2) // Length should remain unchanged
		})

		Convey("When growing nil slice", func() {
			var s slice.Slice[int]
			s = s.Grow(a, 10)

			So(s.Cap(), ShouldBeGreaterThanOrEqualTo, 10)
			So(s.Len(), ShouldEqual, 0)
		})

		Convey("When growing multiple times", func() {
			s = s.Grow(a, 5)
			s = s.Grow(a, 10)
			s = s.Grow(a, 20)

			So(s.Cap(), ShouldBeGreaterThanOrEqualTo, 37) // 2 + 5 + 10 + 20
		})
	})
}

func TestSlice_SetLen(t *testing.T) {
	Convey("Given a slice", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 5)

		Convey("When setting length", func() {
			s = s.SetLen(3)
			So(s.Len(), ShouldEqual, 3)
		})

		Convey("When setting length to zero", func() {
			s = s.SetLen(0)
			So(s.Len(), ShouldEqual, 0)
		})

		Convey("When setting length to capacity", func() {
			s = s.SetLen(s.Cap())
			So(s.Len(), ShouldEqual, s.Cap())
		})
	})
}

func TestSlice_Raw(t *testing.T) {
	Convey("Given a slice", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 3)

		// Store some values
		s.Store(0, 100)
		s.Store(1, 200)
		s.Store(2, 300)

		Convey("When getting raw slice", func() {
			raw := s.Raw()

			So(len(raw), ShouldEqual, 3)
			So(raw[0], ShouldEqual, 100)
			So(raw[1], ShouldEqual, 200)
			So(raw[2], ShouldEqual, 300)
		})
	})
}

func TestSlice_Rest(t *testing.T) {
	Convey("Given a slice with capacity greater than length", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 3)
		s = s.Grow(a, 5) // Increase capacity

		Convey("When getting rest slice", func() {
			rest := s.Rest()

			So(len(rest), ShouldEqual, s.Cap()-s.Len())
		})
	})
}

func TestSlice_Addr(t *testing.T) {
	Convey("Given a slice", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 3)

		Convey("When converting to address slice", func() {
			addr := s.Addr()

			So(addr.Len, ShouldEqual, uint32(3))
			So(addr.Cap, ShouldEqual, uint32(s.Cap()))
		})
	})
}

func TestAddr_AssertValid(t *testing.T) {
	Convey("Given an address slice", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 3)
		addr := s.Addr()

		Convey("When asserting valid", func() {
			valid := addr.AssertValid()

			So(valid.Len(), ShouldEqual, 3)
			So(valid.Cap(), ShouldEqual, s.Cap())
		})
	})
}

func TestAddr_Untyped(t *testing.T) {
	Convey("Given an address slice", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 3)
		addr := s.Addr()

		Convey("When converting to untyped", func() {
			untyped := addr.Untyped()

			So(untyped.Len, ShouldEqual, uint32(3))
			So(untyped.Cap, ShouldEqual, uint32(s.Cap()))
		})
	})
}

func TestUntyped_OffArena(t *testing.T) {
	Convey("Given off-arena data", t, func() {
		data := []int{1, 2, 3, 4, 5}

		Convey("When creating off-arena slice", func() {
			untyped := slice.OffArena(&data[0], len(data))

			So(untyped.Len, ShouldEqual, uint32(5))
			So(untyped.Cap, ShouldEqual, uint32(5))
			So(untyped.OffArena(), ShouldBeTrue)
		})
	})
}

func TestCastUntyped(t *testing.T) {
	Convey("Given an untyped slice", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 3)
		untyped := s.Addr().Untyped()

		Convey("When casting to typed slice", func() {
			typed := slice.CastUntyped[int](untyped)

			So(typed.Len(), ShouldEqual, 3)
			So(typed.Cap(), ShouldEqual, s.Cap())
		})
	})
}

func TestSlice_Format(t *testing.T) {
	Convey("Given a slice", t, func() {
		a := &arena.Arena{}
		s := slice.Make[int](a, 3)

		// Store some values
		s.Store(0, 100)
		s.Store(1, 200)
		s.Store(2, 300)

		Convey("When formatting slice", func() {
			// This should not panic
			So(func() {
				_ = s.Raw()
			}, ShouldNotPanic)
		})
	})
}

func TestSlice_EdgeCases(t *testing.T) {
	Convey("Given edge cases", t, func() {
		Convey("When working with nil slice", func() {
			var s slice.Slice[int]

			So(s.Len(), ShouldEqual, 0)
			So(s.Cap(), ShouldEqual, 0)
			So(s.Ptr(), ShouldBeNil)
		})

		Convey("When working with zero-length slice", func() {
			a := &arena.Arena{}
			s := slice.Make[int](a, 0)

			So(s.Len(), ShouldEqual, 0)
			So(s.Cap(), ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}

func TestSlice_Performance(t *testing.T) {
	Convey("Given performance scenarios", t, func() {
		Convey("When creating many slices", func() {
			a := &arena.Arena{}

			So(func() {
				for i := 0; i < 1000; i++ {
					s := slice.Make[int](a, 100)
					So(s.Len(), ShouldEqual, 100)
				}
			}, ShouldNotPanic)
		})

		Convey("When appending many elements", func() {
			a := &arena.Arena{}
			s := slice.Make[int](a, 0)

			So(func() {
				for i := 0; i < 10000; i++ {
					s = s.AppendOne(a, i)
				}
				So(s.Len(), ShouldEqual, 10000)
			}, ShouldNotPanic)
		})
	})
}

func TestSlice_Concurrency(t *testing.T) {
	Convey("Given concurrent access", t, func() {
		Convey("When multiple goroutines create slices", func() {
			a := &arena.Arena{}

			So(func() {
				done := make(chan bool, 10)
				for i := 0; i < 10; i++ {
					go func() {
						defer func() { done <- true }()

						for j := 0; j < 100; j++ {
							s := slice.Make[int](a, 10)
							if s.Len() != 10 {
								panic("slice length mismatch")
							}
						}
					}()
				}

				// Wait for all goroutines to complete
				for i := 0; i < 10; i++ {
					<-done
				}
			}, ShouldNotPanic)
		})
	})
}
