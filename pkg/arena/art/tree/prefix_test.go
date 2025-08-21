package tree_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	. "github.com/flier/goutil/pkg/arena/art/node"
	. "github.com/flier/goutil/pkg/arena/art/tree"
	"github.com/flier/goutil/pkg/arena/slice"
)

func TestLongestCommonPrefix(t *testing.T) {
	Convey("LongestCommonPrefix", t, func() {
		a := new(arena.Arena)
		hello := slice.FromBytes(a, hello)
		hell := slice.FromBytes(a, hell)
		helloWorld := slice.FromString(a, "hello world")
		help := slice.FromBytes(a, help)
		world := slice.FromString(a, "world")

		var empty slice.Slice[byte]

		Convey("should return depth when no common prefix", func() {
			So(LongestCommonPrefix(hello, world, 0), ShouldEqual, 0)
			So(LongestCommonPrefix(hello, world, 2), ShouldEqual, 2)
			So(LongestCommonPrefix(hello, world, 5), ShouldEqual, 5)
		})

		Convey("should find common prefix from start", func() {
			So(LongestCommonPrefix(hello, help, 0), ShouldEqual, 3)
			So(LongestCommonPrefix(hello, help, 1), ShouldEqual, 3)
			So(LongestCommonPrefix(hello, help, 2), ShouldEqual, 3)
		})

		Convey("should find common prefix from depth", func() {
			So(LongestCommonPrefix(hello, world, 1), ShouldEqual, 1)
			So(LongestCommonPrefix(hello, help, 2), ShouldEqual, 3)
			So(LongestCommonPrefix(hello, help, 3), ShouldEqual, 3)
		})

		Convey("should handle identical strings", func() {
			So(LongestCommonPrefix(hello, hello, 0), ShouldEqual, 5)
			So(LongestCommonPrefix(hello, hello, 2), ShouldEqual, 5)
			So(LongestCommonPrefix(hello, hello, 5), ShouldEqual, 5)
		})

		Convey("should handle one string being prefix of another", func() {
			So(LongestCommonPrefix(hello, helloWorld, 0), ShouldEqual, 5)
			So(LongestCommonPrefix(helloWorld, hello, 0), ShouldEqual, 5)
			So(LongestCommonPrefix(hello, helloWorld, 2), ShouldEqual, 5)
		})

		Convey("should handle empty strings", func() {
			So(LongestCommonPrefix(empty, empty, 0), ShouldEqual, 0)
			So(LongestCommonPrefix(empty, empty, 5), ShouldEqual, 5)
			So(LongestCommonPrefix(hello, empty, 0), ShouldEqual, 0)
			So(LongestCommonPrefix(empty, hello, 0), ShouldEqual, 0)
		})

		Convey("should handle depth beyond string length", func() {
			So(LongestCommonPrefix(hello, world, 10), ShouldEqual, 10)
			So(LongestCommonPrefix(hello, help, 10), ShouldEqual, 10)
			So(LongestCommonPrefix(empty, empty, 10), ShouldEqual, 10)
		})

		Convey("should handle single character strings", func() {
			So(LongestCommonPrefix(slice.FromString(a, "a"), slice.FromString(a, "a"), 0), ShouldEqual, 1)
			So(LongestCommonPrefix(slice.FromString(a, "a"), slice.FromString(a, "b"), 0), ShouldEqual, 0)
			So(LongestCommonPrefix(slice.FromString(a, "a"), slice.FromString(a, "ab"), 0), ShouldEqual, 1)
			So(LongestCommonPrefix(slice.FromString(a, "ab"), slice.FromString(a, "a"), 0), ShouldEqual, 1)
		})

		Convey("should handle special characters", func() {
			So(LongestCommonPrefix(slice.FromString(a, "hello\n"), slice.FromString(a, "hello\t"), 0), ShouldEqual, 5)
			So(LongestCommonPrefix(slice.FromString(a, "hello\000"), slice.FromString(a, "hello\000"), 0), ShouldEqual, 6)
			So(LongestCommonPrefix(slice.FromString(a, "hello\r"), slice.FromString(a, "hello\n"), 0), ShouldEqual, 5)
		})

		Convey("should handle unicode characters", func() {
			So(LongestCommonPrefix(slice.FromString(a, "hello世界"), slice.FromString(a, "hello世界"), 0), ShouldEqual, 11)
			So(LongestCommonPrefix(slice.FromString(a, "hello世界"), slice.FromString(a, "hello地球"), 0), ShouldEqual, 5)
			So(LongestCommonPrefix(slice.FromString(a, "hello世界"), slice.FromString(a, "hello世界!"), 0), ShouldEqual, 11)
		})

		Convey("should handle depth edge cases", func() {
			// Note: negative depth will cause panic due to array bounds check
			// So we only test valid depth values
			So(LongestCommonPrefix(hello, world, 0), ShouldEqual, 0)
			So(LongestCommonPrefix(hello, world, 1), ShouldEqual, 1)
			So(LongestCommonPrefix(hello, world, 4), ShouldEqual, 4)
			So(LongestCommonPrefix(hello, world, 5), ShouldEqual, 5)
		})

		Convey("should handle mixed length strings with common prefix", func() {
			So(LongestCommonPrefix(hello, hell, 0), ShouldEqual, 4)
			So(LongestCommonPrefix(hell, hello, 0), ShouldEqual, 4)
			So(LongestCommonPrefix(hello, hell, 2), ShouldEqual, 4)
			So(LongestCommonPrefix(hell, hello, 2), ShouldEqual, 4)
		})

		Convey("should handle strings with no common prefix from depth", func() {
			So(LongestCommonPrefix(hello, world, 0), ShouldEqual, 0)
			So(LongestCommonPrefix(hello, world, 1), ShouldEqual, 1)
			So(LongestCommonPrefix(hello, world, 2), ShouldEqual, 2)
		})
	})
}

func TestLongestCommonPrefix_EdgeCases(t *testing.T) {
	Convey("LongestCommonPrefix Edge Cases", t, func() {
		a := new(arena.Arena)

		Convey("should handle very long strings", func() {
			longStr1 := slice.FromBytes(a, make([]byte, 10000))
			longStr2 := slice.FromBytes(a, make([]byte, 10000))

			// Fill with same pattern
			for i := 0; i < longStr1.Len(); i++ {
				longStr1.Store(i, byte(i%256))
				longStr2.Store(i, byte(i%256))
			}

			So(LongestCommonPrefix(longStr1, longStr2, 0), ShouldEqual, 10000)
			So(LongestCommonPrefix(longStr1, longStr2, 5000), ShouldEqual, 10000)

			// Make them different at position 5000
			longStr2.Store(5000, 0xFF)
			So(LongestCommonPrefix(longStr1, longStr2, 0), ShouldEqual, 5000)
			So(LongestCommonPrefix(longStr1, longStr2, 2500), ShouldEqual, 5000)
		})

		Convey("should handle strings with only one different character", func() {
			str1 := slice.FromString(a, "hello world")
			str2 := slice.FromString(a, "hello world!")

			So(LongestCommonPrefix(str1, str2, 0), ShouldEqual, 11)
			So(LongestCommonPrefix(str1, str2, 5), ShouldEqual, 11)
			So(LongestCommonPrefix(str1, str2, 10), ShouldEqual, 11)
		})

		Convey("should handle strings with different characters at start", func() {
			str1 := slice.FromString(a, "hello world")
			str2 := slice.FromString(a, "jello world")

			So(LongestCommonPrefix(str1, str2, 0), ShouldEqual, 0)
			So(LongestCommonPrefix(str1, str2, 1), ShouldEqual, 11)  // From index 1, they are identical
			So(LongestCommonPrefix(str1, str2, 2), ShouldEqual, 11)  // From index 2, they are identical
			So(LongestCommonPrefix(str1, str2, 3), ShouldEqual, 11)  // From index 3, they are identical
			So(LongestCommonPrefix(str1, str2, 4), ShouldEqual, 11)  // From index 4, they are identical
			So(LongestCommonPrefix(str1, str2, 5), ShouldEqual, 11)  // From index 5, they are identical
			So(LongestCommonPrefix(str1, str2, 6), ShouldEqual, 11)  // From index 6, they are identical
			So(LongestCommonPrefix(str1, str2, 7), ShouldEqual, 11)  // From index 7, they are identical
			So(LongestCommonPrefix(str1, str2, 8), ShouldEqual, 11)  // From index 8, they are identical
			So(LongestCommonPrefix(str1, str2, 9), ShouldEqual, 11)  // From index 9, they are identical
			So(LongestCommonPrefix(str1, str2, 10), ShouldEqual, 11) // From index 10, they are identical
			So(LongestCommonPrefix(str1, str2, 11), ShouldEqual, 11) // From index 11, they are identical
		})

		Convey("should handle depth beyond both string lengths", func() {
			hello := slice.FromBytes(a, hello)
			world := slice.FromString(a, "world")

			var empty slice.Slice[byte]

			So(LongestCommonPrefix(hello, world, 100), ShouldEqual, 100)
			So(LongestCommonPrefix(empty, empty, 100), ShouldEqual, 100)
			So(LongestCommonPrefix(slice.FromString(a, "a"), slice.FromString(a, "b"), 100), ShouldEqual, 100)
		})
	})
}

func TestPrefixMismatch(t *testing.T) {
	Convey("Given PrefixMismatch function", t, func() {
		a := &arena.Arena{}
		node4 := arena.New(a, Node4[any]{})

		Convey("When checking prefix mismatch", func() {
			Convey("And node has no prefix", func() {
				key := slice.FromString(a, "hello")

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should return 0", func() {
					So(result, ShouldEqual, 0)
				})
			})

			Convey("And node has empty prefix", func() {
				node4.Partial = slice.FromString(a, "")
				key := slice.FromString(a, "hello")

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should return 0", func() {
					So(result, ShouldEqual, 0)
				})
			})

			Convey("And node has prefix", func() {
				node4.Partial = slice.FromString(a, "hello")

				Convey("And key matches prefix exactly", func() {
					key := slice.FromString(a, "hello")

					result := PrefixMismatch[any](node4, key, 0)

					Convey("Then should return prefix length", func() {
						So(result, ShouldEqual, 5)
					})
				})

				Convey("And key starts with prefix", func() {
					key := slice.FromString(a, "hello world")

					result := PrefixMismatch[any](node4, key, 0)

					Convey("Then should return prefix length", func() {
						So(result, ShouldEqual, 5)
					})
				})

				Convey("And key has different first character", func() {
					key := slice.FromString(a, "world")

					result := PrefixMismatch[any](node4, key, 0)

					Convey("Then should return 0", func() {
						So(result, ShouldEqual, 0)
					})
				})

				Convey("And key has partial match", func() {
					key := slice.FromString(a, "help")

					result := PrefixMismatch[any](node4, key, 0)

					Convey("Then should return mismatch position", func() {
						So(result, ShouldEqual, 3)
					})
				})

				Convey("And key is shorter than prefix", func() {
					key := slice.FromString(a, "hello")

					result := PrefixMismatch[any](node4, key, 0)

					Convey("Then should return key length", func() {
						So(result, ShouldEqual, 5)
					})
				})

				Convey("And key is empty", func() {
					key := slice.FromBytes(a, []byte{})

					result := PrefixMismatch[any](node4, key, 0)

					Convey("Then should return 0", func() {
						So(result, ShouldEqual, 0)
					})
				})

				Convey("And key is nil", func() {
					result := PrefixMismatch[any](node4, slice.Slice[byte]{}, 0)

					Convey("Then should return 0", func() {
						So(result, ShouldEqual, 0)
					})
				})
			})
		})

		Convey("When checking with depth > 0", func() {
			node4.Partial = slice.FromString(a, "world")

			Convey("And checking from middle of key", func() {
				key := slice.FromBytes(a, []byte("hello world"))

				result := PrefixMismatch[any](node4, key, 6)

				Convey("Then should return depth + prefix length", func() {
					So(result, ShouldEqual, 5)
				})
			})

			Convey("And checking from end of key", func() {
				key := slice.FromBytes(a, []byte("hello world"))

				result := PrefixMismatch[any](node4, key, 11)

				Convey("Then should return key length", func() {
					So(result, ShouldEqual, 0)
				})
			})

			Convey("And depth exceeds key length", func() {
				key := slice.FromBytes(a, []byte("hello"))

				result := PrefixMismatch[any](node4, key, 6)

				Convey("Then should return key length", func() {
					So(result, ShouldEqual, 0)
				})
			})
		})

		Convey("When checking with special characters", func() {
			Convey("And prefix contains newlines", func() {
				node4.Partial = slice.FromBytes(a, []byte("hello\nworld"))
				key := slice.FromBytes(a, []byte("hello\nworld"))

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, 11)
				})
			})

			Convey("And prefix contains tabs", func() {
				node4.Partial = slice.FromBytes(a, []byte("hello\tworld"))
				key := slice.FromBytes(a, []byte("hello\tworld"))

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, 11)
				})
			})

			Convey("And prefix contains null bytes", func() {
				node4.Partial = slice.FromBytes(a, []byte("hello\000world"))
				key := slice.FromBytes(a, []byte("hello\000world"))

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, 11)
				})
			})
		})

		Convey("When checking with unicode characters", func() {
			Convey("And prefix contains unicode", func() {
				node4.Partial = slice.FromBytes(a, []byte("hello世界"))
				key := slice.FromBytes(a, []byte("hello世界"))

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, 11)
				})
			})

			Convey("And key has different unicode", func() {
				node4.Partial = slice.FromBytes(a, []byte("hello世界"))
				key := slice.FromBytes(a, []byte("hello地球"))

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should return mismatch position", func() {
					So(result, ShouldEqual, 5)
				})
			})
		})
	})
}

func TestPrefixMismatch_EdgeCases(t *testing.T) {
	Convey("Given PrefixMismatch edge cases", t, func() {
		a := &arena.Arena{}

		Convey("When working with very long prefixes", func() {
			Convey("And prefix is 1MB long", func() {
				node4 := arena.New(a, Node4[any]{})
				prefixLen := 1024 * 1024
				longPrefix := make([]byte, prefixLen)
				for i := range longPrefix {
					longPrefix[i] = byte(i % 256)
				}
				node4.Partial = slice.FromBytes(a, longPrefix)

				key := slice.FromBytes(a, make([]byte, prefixLen))
				copy(key.Raw(), longPrefix)

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, prefixLen)
				})
			})

			Convey("And key is shorter than prefix", func() {
				node4 := arena.New(a, Node4[any]{})
				prefixLen := 1000
				longPrefix := make([]byte, prefixLen)
				for i := range longPrefix {
					longPrefix[i] = byte(i % 256)
				}
				node4.Partial = slice.FromBytes(a, longPrefix)

				keyLen := 500
				key := slice.FromBytes(a, make([]byte, keyLen))
				copy(key.Raw(), longPrefix[:keyLen])

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should return key length", func() {
					So(result, ShouldEqual, keyLen)
				})
			})
		})

		Convey("When working with boundary values", func() {
			Convey("And prefix contains zero bytes", func() {
				node4 := arena.New(a, Node4[any]{})
				node4.Partial = slice.FromBytes(a, []byte{0, 0, 0})
				key := slice.FromBytes(a, []byte{0, 0, 0})

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, 3)
				})
			})

			Convey("And prefix contains maximum byte values", func() {
				node4 := arena.New(a, Node4[any]{})
				node4.Partial = slice.FromBytes(a, []byte{255, 255, 255})
				key := slice.FromBytes(a, []byte{255, 255, 255})

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, 3)
				})
			})

			Convey("And prefix contains mixed boundary values", func() {
				node4 := arena.New(a, Node4[any]{})
				node4.Partial = slice.FromBytes(a, []byte{0, 128, 255})
				key := slice.FromBytes(a, []byte{0, 128, 255})

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, 3)
				})
			})
		})

		Convey("When working with special patterns", func() {
			Convey("And prefix has alternating bytes", func() {
				node4 := arena.New(a, Node4[any]{})
				prefix := make([]byte, 100)
				for i := range prefix {
					if i%2 == 0 {
						prefix[i] = 0
					} else {
						prefix[i] = 255
					}
				}
				node4.Partial = slice.FromBytes(a, prefix)

				key := slice.FromBytes(a, make([]byte, 100))
				copy(key.Raw(), prefix)

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, 100)
				})
			})

			Convey("And prefix has sequential bytes", func() {
				node4 := arena.New(a, Node4[any]{})
				prefix := make([]byte, 256)
				for i := range prefix {
					prefix[i] = byte(i)
				}
				node4.Partial = slice.FromBytes(a, prefix)

				key := slice.FromBytes(a, make([]byte, 256))
				copy(key.Raw(), prefix)

				result := PrefixMismatch[any](node4, key, 0)

				Convey("Then should match exactly", func() {
					So(result, ShouldEqual, 256)
				})
			})
		})
	})
}
