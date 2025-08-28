package node_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	. "github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
)

func TestBase(t *testing.T) {
	Convey("Given a Base struct", t, func() {
		a := &arena.Arena{}

		Convey("When creating a new base", func() {
			base := &Base[any]{}

			Convey("Then should have default values", func() {
				So(base.NumChildren, ShouldEqual, 0)
				So(base.Partial.Empty(), ShouldBeTrue)
			})
		})

		Convey("When setting prefix", func() {
			base := &Base[any]{}
			prefix := slice.FromString(a, "hello")

			base.SetPrefix(prefix)

			Convey("Then prefix should be set correctly", func() {
				So(base.Prefix().Raw(), ShouldResemble, []byte("hello"))
				So(base.Partial.Raw(), ShouldResemble, []byte("hello"))
			})

			Convey("And prefix should be accessible via Prefix()", func() {
				So(base.Prefix().Raw(), ShouldResemble, []byte("hello"))
			})
		})

		Convey("When updating prefix", func() {
			base := &Base[any]{}
			prefix1 := slice.FromString(a, "hello")
			prefix2 := slice.FromString(a, "world")

			base.SetPrefix(prefix1)
			base.SetPrefix(prefix2)

			Convey("Then prefix should be updated", func() {
				So(base.Prefix().Raw(), ShouldResemble, []byte("world"))
				So(base.Partial.Raw(), ShouldResemble, []byte("world"))
			})
		})

		Convey("When setting empty prefix", func() {
			base := &Base[any]{}
			emptyPrefix := slice.FromBytes(a, []byte{})

			base.SetPrefix(emptyPrefix)

			Convey("Then prefix should be empty", func() {
				So(base.Prefix().Empty(), ShouldBeTrue)
				So(base.Partial.Empty(), ShouldBeTrue)
			})
		})

		Convey("When setting nil prefix", func() {
			base := &Base[any]{}
			var nilPrefix slice.Slice[byte]

			base.SetPrefix(nilPrefix)

			Convey("Then prefix should be empty", func() {
				So(base.Prefix().Empty(), ShouldBeTrue)
				So(base.Partial.Empty(), ShouldBeTrue)
			})
		})

		Convey("When setting very long prefix", func() {
			base := &Base[any]{}
			longPrefix := make([]byte, 1000)
			for i := range longPrefix {
				longPrefix[i] = byte(i % 256)
			}
			prefixSlice := slice.FromBytes(a, longPrefix)

			base.SetPrefix(prefixSlice)

			Convey("Then prefix should be set correctly", func() {
				So(base.Prefix().Len(), ShouldEqual, 1000)
				So(base.Partial.Len(), ShouldEqual, 1000)
				So(base.Prefix().Raw(), ShouldResemble, longPrefix)
			})
		})

		Convey("When setting prefix with special characters", func() {
			base := &Base[any]{}
			specialPrefix := slice.FromBytes(a, []byte("hello\n\t\r\000world"))

			base.SetPrefix(specialPrefix)

			Convey("Then prefix should preserve special characters", func() {
				So(base.Prefix().Raw(), ShouldResemble, []byte("hello\n\t\r\000world"))
				So(base.Partial.Raw(), ShouldResemble, []byte("hello\n\t\r\000world"))
			})
		})

		Convey("When setting prefix with unicode characters", func() {
			base := &Base[any]{}
			unicodePrefix := slice.FromBytes(a, []byte("hello世界"))

			base.SetPrefix(unicodePrefix)

			Convey("Then prefix should preserve unicode characters", func() {
				So(base.Prefix().Raw(), ShouldResemble, []byte("hello世界"))
				So(base.Partial.Raw(), ShouldResemble, []byte("hello世界"))
			})
		})
	})
}

func TestBase_EdgeCases(t *testing.T) {
	Convey("Given Base edge cases", t, func() {
		a := &arena.Arena{}

		Convey("When working with boundary values", func() {
			Convey("And setting prefix with zero bytes", func() {
				base := &Base[any]{}
				zeroPrefix := slice.FromBytes(a, []byte{0, 0, 0})

				base.SetPrefix(zeroPrefix)

				Convey("Then prefix should handle zero bytes correctly", func() {
					So(base.Prefix().Len(), ShouldEqual, 3)
					So(base.Prefix().Raw(), ShouldResemble, []byte{0, 0, 0})
				})
			})

			Convey("And setting prefix with maximum byte values", func() {
				base := &Base[any]{}
				maxPrefix := slice.FromBytes(a, []byte{255, 255, 255})

				base.SetPrefix(maxPrefix)

				Convey("Then prefix should handle maximum bytes correctly", func() {
					So(base.Prefix().Len(), ShouldEqual, 3)
					So(base.Prefix().Raw(), ShouldResemble, []byte{255, 255, 255})
				})
			})

			Convey("And setting prefix with mixed boundary values", func() {
				base := &Base[any]{}
				mixedPrefix := slice.FromBytes(a, []byte{0, 128, 255})

				base.SetPrefix(mixedPrefix)

				Convey("Then prefix should handle mixed values correctly", func() {
					So(base.Prefix().Len(), ShouldEqual, 3)
					So(base.Prefix().Raw(), ShouldResemble, []byte{0, 128, 255})
				})
			})
		})

		Convey("When working with very long prefixes", func() {
			Convey("And setting prefix with 1MB length", func() {
				base := &Base[any]{}
				prefixLen := 1024 * 1024
				longPrefix := make([]byte, prefixLen)
				for i := range longPrefix {
					longPrefix[i] = byte(i % 256)
				}
				prefixSlice := slice.FromBytes(a, longPrefix)

				base.SetPrefix(prefixSlice)

				Convey("Then prefix should handle very long data correctly", func() {
					So(base.Prefix().Len(), ShouldEqual, prefixLen)
					So(base.Partial.Len(), ShouldEqual, prefixLen)
					// Verify first and last few bytes
					So(base.Prefix().Load(0), ShouldEqual, byte(0))
					So(base.Prefix().Load(prefixLen-1), ShouldEqual, byte((prefixLen-1)%256))
				})
			})

			Convey("And setting prefix with exactly 256 bytes", func() {
				base := &Base[any]{}
				prefixLen := 256
				prefix := make([]byte, prefixLen)
				for i := range prefix {
					prefix[i] = byte(i)
				}
				prefixSlice := slice.FromBytes(a, prefix)

				base.SetPrefix(prefixSlice)

				Convey("Then prefix should handle 256 bytes correctly", func() {
					So(base.Prefix().Len(), ShouldEqual, 256)
					So(base.Prefix().Load(0), ShouldEqual, byte(0))
					So(base.Prefix().Load(255), ShouldEqual, byte(255))
				})
			})
		})

		Convey("When working with special patterns", func() {
			Convey("And setting prefix with alternating bytes", func() {
				base := &Base[any]{}
				prefix := make([]byte, 100)
				for i := range prefix {
					if i%2 == 0 {
						prefix[i] = 0
					} else {
						prefix[i] = 255
					}
				}
				prefixSlice := slice.FromBytes(a, prefix)

				base.SetPrefix(prefixSlice)

				Convey("Then prefix should preserve alternating pattern", func() {
					So(base.Prefix().Len(), ShouldEqual, 100)
					So(base.Prefix().Load(0), ShouldEqual, byte(0))
					So(base.Prefix().Load(1), ShouldEqual, byte(255))
					So(base.Prefix().Load(98), ShouldEqual, byte(0))
					So(base.Prefix().Load(99), ShouldEqual, byte(255))
				})
			})

			Convey("And setting prefix with sequential bytes", func() {
				base := &Base[any]{}
				prefix := make([]byte, 256)
				for i := range prefix {
					prefix[i] = byte(i)
				}
				prefixSlice := slice.FromBytes(a, prefix)

				base.SetPrefix(prefixSlice)

				Convey("Then prefix should preserve sequential pattern", func() {
					So(base.Prefix().Len(), ShouldEqual, 256)
					So(base.Prefix().Load(0), ShouldEqual, byte(0))
					So(base.Prefix().Load(128), ShouldEqual, byte(128))
					So(base.Prefix().Load(255), ShouldEqual, byte(255))
				})
			})
		})
	})
}

func TestBase_Performance(t *testing.T) {
	Convey("Given Base performance considerations", t, func() {
		a := &arena.Arena{}

		Convey("When performing many prefix operations", func() {
			Convey("And setting prefix 1000 times", func() {
				base := &Base[any]{}

				for i := 0; i < 1000; i++ {
					prefix := slice.FromBytes(a, []byte{byte(i % 256)})
					base.SetPrefix(prefix)
				}

				Convey("Then final prefix should be correct", func() {
					So(base.Prefix().Len(), ShouldEqual, 1)
					So(base.Prefix().Load(0), ShouldEqual, byte(999%256))
				})
			})

			Convey("And setting prefixes of different lengths", func() {
				base := &Base[any]{}

				lengths := []int{1, 10, 100, 1000}
				for _, length := range lengths {
					prefix := make([]byte, length)
					for i := range prefix {
						prefix[i] = byte(i % 256)
					}
					prefixSlice := slice.FromBytes(a, prefix)
					base.SetPrefix(prefixSlice)
				}

				Convey("Then final prefix should have correct length", func() {
					So(base.Prefix().Len(), ShouldEqual, 1000)
					So(base.Prefix().Load(0), ShouldEqual, byte(0))
					So(base.Prefix().Load(999), ShouldEqual, byte(999%256))
				})
			})
		})

		Convey("When working with large data sets", func() {
			Convey("And setting multiple large prefixes", func() {
				base := &Base[any]{}

				// Set multiple large prefixes
				for i := 0; i < 10; i++ {
					prefixLen := 1000 + i*100
					prefix := make([]byte, prefixLen)
					for j := range prefix {
						prefix[j] = byte((i + j) % 256)
					}
					prefixSlice := slice.FromBytes(a, prefix)
					base.SetPrefix(prefixSlice)
				}

				Convey("Then final prefix should be correct", func() {
					expectedLen := 1000 + 9*100 // Last prefix length
					So(base.Prefix().Len(), ShouldEqual, expectedLen)
					So(base.Prefix().Load(0), ShouldEqual, byte(9)) // First byte of last prefix
				})
			})
		})
	})
}
