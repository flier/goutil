package tree_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/node"
	. "github.com/flier/goutil/pkg/arena/art/tree"
	"github.com/flier/goutil/pkg/arena/slice"
)

// TestSearch tests the Search function
func TestSearch(t *testing.T) {
	Convey("Given a Search function", t, func() {
		a := new(arena.Arena)

		Convey("When searching in an empty tree", func() {
			var root node.Ref[int]

			result := Search(root, []byte("hello"))

			Convey("Then should return nil", func() {
				So(result, ShouldBeNil)
			})
		})

		Convey("When searching in a tree with single leaf", func() {
			// Create a single leaf
			leaf := node.NewLeaf(a, hello, 123)
			root := leaf.Ref()

			Convey("And searching for matching key", func() {
				result := Search(root, hello)

				Convey("Then should return the value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 123)
				})
			})

			Convey("And searching for non-matching key", func() {
				result := Search(root, []byte("world"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})

			Convey("And searching for partial key", func() {
				result := Search(root, []byte("hel"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})

			Convey("And searching for longer key", func() {
				result := Search(root, []byte("hello world"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})
		})

		Convey("When searching in a tree with Node4", func() {
			// Create a Node4 with two children
			node4 := arena.New(a, node.Node4[int]{})
			leaf1 := node.NewLeaf(a, hello, 123)
			leaf2 := node.NewLeaf(a, foobar, 456)

			// Add children to Node4
			node4.AddChild(int('h'), leaf1)
			node4.AddChild(int('f'), leaf2)
			root := node4.Ref()

			Convey("And searching for first child key", func() {
				result := Search(root, hello)

				Convey("Then should return the first value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 123)
				})
			})

			Convey("And searching for second child key", func() {
				result := Search(root, foobar)

				Convey("Then should return the second value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 456)
				})
			})

			Convey("And searching for non-existent key", func() {
				result := Search(root, []byte("xyz"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})
		})

		Convey("When searching in a tree with Node4 and prefix", func() {
			// Create a Node4 with common prefix
			node4 := arena.New(a, node.Node4[int]{})
			node4.Partial = slice.FromBytes(a, []byte("hel"))

			leaf1 := node.NewLeaf(a, hello, 123)
			leaf2 := node.NewLeaf(a, help, 456)

			// Add children to Node4
			node4.AddChild(int('l'), leaf1)
			node4.AddChild(int('p'), leaf2)
			root := node4.Ref()

			Convey("And searching for first child key", func() {
				result := Search(root, hello)

				Convey("Then should return the first value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 123)
				})
			})

			Convey("And searching for second child key", func() {
				result := Search(root, help)

				Convey("Then should return the second value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 456)
				})
			})

			Convey("And searching for key with different prefix", func() {
				result := Search(root, []byte("world"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})

			Convey("And searching for key with partial prefix match", func() {
				result := Search(root, []byte("he"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})

			Convey("And searching for key with longer prefix", func() {
				result := Search(root, []byte("hello world"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})
		})

		Convey("When searching in a tree with Node16", func() {
			// Create a Node16 with multiple children
			node16 := arena.New(a, node.Node16[int]{})
			leaves := make([]*node.Leaf[int], 8)
			values := []int{10, 20, 30, 40, 50, 60, 70, 80}

			for i := 0; i < 8; i++ {
				leaves[i] = node.NewLeaf(a, []byte{byte('a' + i)}, values[i])
				node16.AddChild(int('a'+i), leaves[i])
			}

			root := node16.Ref()

			Convey("And searching for existing keys", func() {
				for i := 0; i < 8; i++ {
					result := Search(root, []byte{byte('a' + i)})
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, values[i])
				}
			})

			Convey("And searching for non-existent key", func() {
				result := Search(root, []byte("x"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})
		})

		Convey("When searching in a tree with Node48", func() {
			// Create a Node48 with sparse keys
			node48 := arena.New(a, node.Node48[int]{})
			leaf1 := node.NewLeaf(a, []byte{0}, 100)
			leaf2 := node.NewLeaf(a, []byte{128}, 200)
			leaf3 := node.NewLeaf(a, []byte{255}, 300)

			node48.AddChild(int(0), leaf1)
			node48.AddChild(int(128), leaf2)
			node48.AddChild(int(255), leaf3)
			root := node48.Ref()

			Convey("And searching for existing keys", func() {
				result1 := Search(root, []byte{0})
				So(result1, ShouldNotEqual, null)
				So(*(*int)(result1), ShouldEqual, 100)

				result2 := Search(root, []byte{128})
				So(result2, ShouldNotEqual, null)
				So(*(*int)(result2), ShouldEqual, 200)

				result3 := Search(root, []byte{255})
				So(result3, ShouldNotEqual, null)
				So(*(*int)(result3), ShouldEqual, 300)
			})

			Convey("And searching for non-existent key", func() {
				result := Search(root, []byte{64})

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})
		})

		Convey("When searching in a tree with Node256", func() {
			// Create a Node256 with some children
			node256 := arena.New(a, node.Node256[int]{})
			leaf1 := node.NewLeaf(a, []byte{42}, 420)
			leaf2 := node.NewLeaf(a, []byte{200}, 2000)

			node256.AddChild(int(42), leaf1)
			node256.AddChild(int(200), leaf2)
			root := node256.Ref()

			Convey("And searching for existing keys", func() {
				result1 := Search(root, []byte{42})
				So(result1, ShouldNotEqual, null)
				So(*(*int)(result1), ShouldEqual, 420)

				result2 := Search(root, []byte{200})
				So(result2, ShouldNotEqual, null)
				So(*(*int)(result2), ShouldEqual, 2000)
			})

			Convey("And searching for non-existent key", func() {
				result := Search(root, []byte{100})

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})
		})

		Convey("When searching in a tree with multi-level structure", func() {
			// Create a tree by using the RecursiveInsert function to ensure correct structure
			var root node.Ref[int]

			// Insert hello
			helloLeaf := node.NewLeaf(a, hello, 123)
			RecursiveInsert(a, &root, helloLeaf, 0, false)

			// Insert help
			helpLeaf := node.NewLeaf(a, help, 456)
			RecursiveInsert(a, &root, helpLeaf, 0, false)

			Convey("And searching for 'hello'", func() {
				result := Search(root, hello)

				Convey("Then should return the correct value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 123)
				})
			})

			Convey("And searching for 'help'", func() {
				result := Search(root, help)

				Convey("Then should return the correct value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 456)
				})
			})

			Convey("And searching for non-existent key", func() {
				result := Search(root, []byte("world"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})

			Convey("And searching for partial key", func() {
				result := Search(root, []byte("he"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})
		})

		Convey("When searching with edge cases", func() {
			Convey("And searching with empty key", func() {
				leaf := node.NewLeaf(a, []byte{}, 999)
				root := leaf.Ref()

				result := Search(root, []byte{})

				Convey("Then should return the value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 999)
				})
			})

			Convey("And searching with zero byte key", func() {
				leaf := node.NewLeaf(a, []byte{0}, 888)
				root := leaf.Ref()

				result := Search(root, []byte{0})

				Convey("Then should return the value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 888)
				})
			})

			Convey("And searching with maximum byte key", func() {
				leaf := node.NewLeaf(a, []byte{255}, 777)
				root := leaf.Ref()

				result := Search(root, []byte{255})

				Convey("Then should return the value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 777)
				})
			})

			Convey("And searching with very long key", func() {
				longKey := make([]byte, 1000)
				for i := range longKey {
					longKey[i] = byte(i % 256)
				}

				leaf := node.NewLeaf(a, longKey, 666)
				root := leaf.Ref()

				result := Search(root, longKey)

				Convey("Then should return the value", func() {
					So(result, ShouldNotEqual, null)
					So(*(*int)(result), ShouldEqual, 666)
				})
			})
		})

		Convey("When searching with prefix mismatch scenarios", func() {
			Convey("And searching with key that has shorter prefix", func() {
				node4 := arena.New(a, node.Node4[int]{})
				node4.Partial = slice.FromBytes(a, []byte("hello"))

				leaf := node.NewLeaf(a, hello, 555)

				node4.AddChild(int(0), leaf)
				root := node4.Ref()

				// Search with key "hel" which is shorter than prefix "hello"
				result := Search(root, []byte("hel"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})

			Convey("And searching with key that has different prefix", func() {
				node4 := arena.New(a, node.Node4[int]{})
				node4.Partial = slice.FromBytes(a, []byte("hello"))

				leaf := node.NewLeaf(a, hello, 444)

				node4.AddChild(int(0), leaf)
				root := node4.Ref()

				// Search with key "world" which has different prefix
				result := Search(root, []byte("world"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})

			Convey("And searching with key that has partial prefix match", func() {
				node4 := arena.New(a, node.Node4[int]{})
				node4.Partial = slice.FromBytes(a, []byte("hello"))

				leaf := node.NewLeaf(a, hello, 333)

				node4.AddChild(int(0), leaf)
				root := node4.Ref()

				// Search with key "hellx" which matches prefix up to "hell" but differs at position 4
				result := Search(root, []byte("hellx"))

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})
			})
		})
	})
}
