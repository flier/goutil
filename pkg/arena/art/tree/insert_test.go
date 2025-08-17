//go:build go1.22

package tree

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
)

var (
	hell   = []byte("hell")
	hello  = []byte("hello")
	help   = []byte("help")
	foobar = []byte("foobar")
	null   = unsafe.Pointer(nil)
)

func TestInsert(t *testing.T) {
	Convey("Given a ART tree", t, func() {
		a := new(arena.Arena)

		Convey("When inserting a leaf to an empty tree", func() {
			var root node.Ref

			leaf := arena.New(a, node.Leaf{
				Key:   slice.FromBytes(a, hello),
				Value: unsafe.Pointer(arena.New(a, 123)),
			})

			So(RecursiveInsert(a, &root, leaf, 0, false), ShouldEqual, null)

			Convey("Then the root should be replaced with the leaf", func() {
				So(root.Empty(), ShouldBeFalse)

				l := root.AsLeaf()
				So(l, ShouldNotEqual, null)
				So(l.Key.Raw(), ShouldResemble, hello)
				So(*(*int)(l.Value), ShouldEqual, 123)
			})

			Convey("When inserting another leaf with a same key", func() {
				leaf2 := arena.New(a, node.Leaf{
					Key:   slice.FromBytes(a, hello),
					Value: unsafe.Pointer(arena.New(a, 456)),
				})

				v := RecursiveInsert(a, &root, leaf2, 0, true)
				So(v, ShouldNotEqual, null)
				So(*(*int)(v), ShouldEqual, 123)

				Convey("Then the root should be replaced with the value of the second leaf", func() {
					So(root.Empty(), ShouldBeFalse)

					l := root.AsLeaf()
					So(l, ShouldNotEqual, null)
					So(l.Key.Raw(), ShouldResemble, hello)
					So(*(*int)(l.Value), ShouldEqual, 456)
				})
			})

			Convey("When inserting another leaf with a different key", func() {
				leaf2 := arena.New(a, node.Leaf{
					Key:   slice.FromBytes(a, foobar),
					Value: unsafe.Pointer(arena.New(a, 456)),
				})

				v := RecursiveInsert(a, &root, leaf2, 0, true)
				So(v, ShouldEqual, null)

				Convey("Then the root should be split into a node4", func() {
					So(root.Empty(), ShouldBeFalse)

					n := root.AsNode4()
					So(n, ShouldNotEqual, null)
					So(n.Partial.Empty(), ShouldBeTrue)
					So(n.NumChildren, ShouldEqual, 2)
					So(n.Keys[:], ShouldResemble, []byte{'f', 'h', 0, 0})
					So(n.Children[:], ShouldResemble, []node.Ref{leaf2.Ref(), leaf.Ref(), 0, 0})
				})
			})

			Convey("When inserting another leaf with a different key with a common prefix", func() {
				leaf2 := arena.New(a, node.Leaf{
					Key:   slice.FromBytes(a, help),
					Value: unsafe.Pointer(arena.New(a, 456)),
				})

				v := RecursiveInsert(a, &root, leaf2, 0, true)
				So(v, ShouldEqual, null)

				Convey("Then the root should be split into a node4", func() {
					So(root.Empty(), ShouldBeFalse)

					n := root.AsNode4()
					So(n, ShouldNotEqual, null)
					So(n.Partial.Raw(), ShouldEqual, []byte("hel"))
					So(n.NumChildren, ShouldEqual, 2)
					So(n.Keys[:], ShouldResemble, []byte{'l', 'p', 0, 0})
					So(n.Children[:], ShouldResemble, []node.Ref{leaf.Ref(), leaf2.Ref(), 0, 0})
				})
			})

			Convey("When inserting another leaf with a same prefix", func() {
				leaf2 := arena.New(a, node.Leaf{
					Key:   slice.FromBytes(a, hell),
					Value: unsafe.Pointer(arena.New(a, 456)),
				})

				v := RecursiveInsert(a, &root, leaf2, 0, true)
				So(v, ShouldEqual, null)

				Convey("Then the root should be split into a node4", func() {
					So(root.Empty(), ShouldBeFalse)

					n := root.AsNode4()
					So(n, ShouldNotEqual, null)
					So(n.Partial.Raw(), ShouldEqual, []byte("hell"))
					So(n.NumChildren, ShouldEqual, 2)
					So(n.Keys[:], ShouldResemble, []byte{0, 'o', 0, 0})
					So(n.Children[:], ShouldResemble, []node.Ref{leaf2.Ref(), leaf.Ref(), 0, 0})
				})
			})
		})
	})
}

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
			for i := range longStr1.Len() {
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
