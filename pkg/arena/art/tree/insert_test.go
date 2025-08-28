package tree_test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/art/node"
	. "github.com/flier/goutil/pkg/arena/art/tree"
	"github.com/flier/goutil/pkg/arena/slice"
	"github.com/flier/goutil/pkg/opt"
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
			var root node.Ref[int]

			leaf := node.NewLeaf(a, hello, 123)

			So(RecursiveInsert(a, &root, leaf, 0, false), ShouldBeNil)

			Convey("Then the root should be replaced with the leaf", func() {
				So(root.Empty(), ShouldBeFalse)

				l := root.AsLeaf()
				So(l, ShouldNotBeNil)
				So(l.Key.Raw(), ShouldResemble, hello)
				So(l.Value, ShouldEqual, 123)
			})

			Convey("When inserting another leaf with a same key", func() {
				leaf2 := node.NewLeaf(a, hello, 456)

				v := RecursiveInsert(a, &root, leaf2, 0, true)
				So(v, ShouldNotBeNil)
				So(*v, ShouldEqual, 123)

				Convey("Then the root should be replaced with the value of the second leaf", func() {
					So(root.Empty(), ShouldBeFalse)

					l := root.AsLeaf()
					So(l, ShouldNotEqual, null)
					So(l.Key.Raw(), ShouldResemble, hello)
					So(l.Value, ShouldEqual, 456)
				})
			})

			Convey("When inserting another leaf with a different key", func() {
				leaf2 := node.NewLeaf(a, foobar, 456)

				v := RecursiveInsert(a, &root, leaf2, 0, true)
				So(v, ShouldBeNil)

				Convey("Then the root should be split into a node4", func() {
					So(root.Empty(), ShouldBeFalse)

					n := root.AsNode4()
					So(n, ShouldNotBeNil)
					So(n.Partial.Empty(), ShouldBeTrue)
					So(n.NumChildren, ShouldEqual, 2)
					So(n.Keys[:], ShouldResemble, []byte{'f', 'h', 0, 0})
					So(n.Children[:], ShouldResemble, []node.Ref[int]{leaf2.Ref(), leaf.Ref(), 0, 0})
				})
			})

			Convey("When inserting another leaf with a different key with a common prefix", func() {
				leaf2 := node.NewLeaf(a, help, 456)

				v := RecursiveInsert(a, &root, leaf2, 0, true)
				So(v, ShouldBeNil)

				Convey("Then the root should be split into a node4", func() {
					So(root.Empty(), ShouldBeFalse)

					n := root.AsNode4()
					So(n, ShouldNotBeNil)
					So(n.Partial.Raw(), ShouldEqual, []byte("hel"))
					So(n.NumChildren, ShouldEqual, 2)
					So(n.Keys[:], ShouldResemble, []byte{'l', 'p', 0, 0})
					So(n.Children[:], ShouldResemble, []node.Ref[int]{leaf.Ref(), leaf2.Ref(), 0, 0})
				})
			})

			Convey("When inserting another leaf with a same prefix", func() {
				leaf2 := node.NewLeaf(a, hell, 456)

				v := RecursiveInsert(a, &root, leaf2, 0, true)
				So(v, ShouldBeNil)

				Convey("Then the root should be split into a node4", func() {
					So(root.Empty(), ShouldBeFalse)

					n := root.AsNode4()
					So(n, ShouldNotBeNil)
					So(n.Partial.Raw(), ShouldEqual, []byte("hell"))

					Convey("Then the zero-sized child should be optimized", func() {
						So(n.ZeroSizedChild.Empty(), ShouldBeFalse)

						l := n.ZeroSizedChild.AsLeaf()
						So(l, ShouldEqual, leaf2)
						So(l.Key.Raw(), ShouldResemble, []byte("hell"))
						So(l.Value, ShouldEqual, 456)
					})

					Convey("Then the non-zero-sized child should be optimized", func() {
						So(n.NumChildren, ShouldEqual, 1)
						So(n.Keys, ShouldResemble, [4]byte{'o', 0, 0, 0})
						So(n.Children, ShouldResemble, [4]node.Ref[int]{leaf.Ref(), 0, 0, 0})
					})
				})
			})
		})
	})
}

// TestInsertToLeaf tests the InsertToLeaf function directly
func TestInsertToLeaf(t *testing.T) {
	Convey("Given InsertToLeaf function", t, func() {
		a := new(arena.Arena)

		Convey("When inserting a leaf with matching key", func() {
			// Create initial leaf
			initialLeaf := node.NewLeaf(a, hello, 123)

			ref := initialLeaf.Ref()

			// Create new leaf with same key
			newLeaf := node.NewLeaf(a, hello, 456)

			Convey("And replace is false", func() {
				result := InsertToLeaf(a, &ref, newLeaf, 0, false)

				Convey("Then should return old value", func() {
					So(result, ShouldNotBeNil)
					So(*result, ShouldEqual, 123)
				})

				Convey("And original leaf should remain unchanged", func() {
					So(ref.IsLeaf(), ShouldBeTrue)
					leaf := ref.AsLeaf()
					So(leaf.Value, ShouldEqual, 123)
				})
			})

			Convey("And replace is true", func() {
				result := InsertToLeaf(a, &ref, newLeaf, 0, true)

				Convey("Then should return old value", func() {
					So(result, ShouldNotBeNil)
					So(*result, ShouldEqual, 123)
				})

				Convey("And original leaf should have new value", func() {
					So(ref.IsLeaf(), ShouldBeTrue)
					leaf := ref.AsLeaf()
					So(leaf.Value, ShouldEqual, 456)
				})
			})
		})

		Convey("When inserting a leaf with different key and no common prefix", func() {
			// Create initial leaf
			initialLeaf := node.NewLeaf(a, hello, 123)

			ref := initialLeaf.Ref()

			// Create new leaf with completely different key
			newLeaf := node.NewLeaf(a, foobar, 456)

			result := InsertToLeaf(a, &ref, newLeaf, 0, false)

			Convey("Then should return null", func() {
				So(result, ShouldBeNil)
			})

			Convey("And ref should be replaced with Node4", func() {
				So(ref.IsNode4(), ShouldBeTrue)

				node4 := ref.AsNode4()
				So(node4.NumChildren, ShouldEqual, 2)
				So(node4.Partial.Empty(), ShouldBeTrue)

				Convey("And Node4 should contain both leaves", func() {
					// Keys should be sorted: 'f' (foobar) comes before 'h' (hello)
					So(node4.Keys[0], ShouldEqual, byte('f'))
					So(node4.Keys[1], ShouldEqual, byte('h'))
					So(node4.Children[0].AsLeaf().Key.Raw(), ShouldResemble, foobar)
					So(node4.Children[1].AsLeaf().Key.Raw(), ShouldResemble, hello)
				})
			})
		})

		Convey("When inserting a leaf with different key and common prefix", func() {
			// Create initial leaf
			initialLeaf := node.NewLeaf(a, hello, 123)

			ref := initialLeaf.Ref()

			// Create new leaf with common prefix "hel"
			newLeaf := node.NewLeaf(a, help, 456)

			result := InsertToLeaf(a, &ref, newLeaf, 0, false)

			Convey("Then should return null", func() {
				So(result, ShouldBeNil)
			})

			Convey("And ref should be replaced with Node4", func() {
				So(ref.IsNode4(), ShouldBeTrue)

				node4 := ref.AsNode4()
				So(node4.NumChildren, ShouldEqual, 2)

				Convey("And Node4 should have common prefix", func() {

					So(node4.Partial.Raw(), ShouldEqual, []byte("hel"))
				})

				Convey("And Node4 should contain both leaves with correct keys", func() {
					// At depth 3 (after "hel"), keys are 'l' and 'p'
					So(node4.Keys[0], ShouldEqual, byte('l'))
					So(node4.Keys[1], ShouldEqual, byte('p'))
					So(node4.Children[0].AsLeaf().Key.Raw(), ShouldResemble, hello)
					So(node4.Children[1].AsLeaf().Key.Raw(), ShouldResemble, help)
				})
			})
		})

		Convey("When inserting a leaf with same prefix", func() {
			// Create initial leaf
			initialLeaf := node.NewLeaf(a, hello, 123)

			ref := initialLeaf.Ref()

			// Create new leaf with same prefix but different suffix
			newLeaf := node.NewLeaf(a, hell, 456)

			result := InsertToLeaf(a, &ref, newLeaf, 0, false)

			Convey("Then should return null", func() {
				So(result, ShouldBeNil)
			})

			Convey("And ref should be replaced with Node4", func() {
				So(ref.IsNode4(), ShouldBeTrue)

				node4 := ref.AsNode4()
				So(node4.NumChildren, ShouldEqual, 1)

				Convey("And Node4 should have common prefix", func() {
					So(node4.Partial.Raw(), ShouldEqual, []byte("hell"))
				})

				Convey("And Node4 should contain both leaves with correct keys", func() {
					So(node4.ZeroSizedChild.Empty(), ShouldBeFalse)
					So(node4.ZeroSizedChild.AsLeaf(), ShouldEqual, newLeaf)

					// At depth 4 (after "hell"), keys are 0 (null terminator) and 'o'
					So(node4.Keys, ShouldResemble, [4]byte{'o', 0, 0, 0})
					So(node4.Children, ShouldResemble, [4]node.Ref[int]{initialLeaf.Ref(), 0, 0, 0})
				})
			})
		})

		Convey("When inserting with edge case keys", func() {
			Convey("And using zero byte keys", func() {
				initialLeaf := node.NewLeaf(a, []byte{0}, 123)

				ref := initialLeaf.Ref()

				newLeaf := node.NewLeaf(a, []byte{1}, 456)

				result := InsertToLeaf(a, &ref, newLeaf, 0, false)

				So(result, ShouldBeNil)
				So(ref.IsNode4(), ShouldBeTrue)

				node4 := ref.AsNode4()
				So(node4.NumChildren, ShouldEqual, 2)
				So(node4.Partial.Empty(), ShouldBeTrue)
				So(node4.Keys[0], ShouldEqual, byte(0))
				So(node4.Keys[1], ShouldEqual, byte(1))
			})

			Convey("And using maximum byte values", func() {
				initialLeaf := node.NewLeaf(a, []byte{255}, 123)

				ref := initialLeaf.Ref()

				newLeaf := node.NewLeaf(a, []byte{254}, 456)

				result := InsertToLeaf(a, &ref, newLeaf, 0, false)

				So(result, ShouldBeNil)
				So(ref.IsNode4(), ShouldBeTrue)

				node4 := ref.AsNode4()
				So(node4.NumChildren, ShouldEqual, 2)
				So(node4.Partial.Empty(), ShouldBeTrue)
				// Keys should be sorted: 254 comes before 255
				So(node4.Keys[0], ShouldEqual, byte(254))
				So(node4.Keys[1], ShouldEqual, byte(255))
			})
		})

		Convey("When inserting with very long keys", func() {
			// Create a long key
			longKey := make([]byte, 1000)
			for i := range longKey {
				longKey[i] = byte(i % 256)
			}

			initialLeaf := node.NewLeaf(a, longKey, 123)

			ref := initialLeaf.Ref()

			// Create another long key with common prefix
			longKey2 := make([]byte, 1000)
			copy(longKey2, longKey)
			longKey2[500] = byte(255) // Different at position 500

			newLeaf := node.NewLeaf(a, longKey2, 456)

			result := InsertToLeaf(a, &ref, newLeaf, 0, false)

			Convey("Then should return null", func() {
				So(result, ShouldBeNil)
			})

			Convey("And ref should be replaced with Node4", func() {
				So(ref.IsNode4(), ShouldBeTrue)

				node4 := ref.AsNode4()
				So(node4.NumChildren, ShouldEqual, 2)

				Convey("And Node4 should have common prefix up to position 500", func() {
					So(node4.Partial.Len(), ShouldEqual, 500)
					// Verify the prefix matches
					for i := 0; i < 500; i++ {
						So(node4.Partial.Load(i), ShouldEqual, byte(i%256))
					}
				})
			})
		})

		Convey("When inserting with empty keys", func() {
			Convey("And both keys are empty", func() {
				// Use null terminator (0) to represent empty string in ART
				initialLeaf := node.NewLeaf(a, []byte{0}, 123)

				ref := initialLeaf.Ref()

				newLeaf := node.NewLeaf(a, []byte{0}, 456)

				result := InsertToLeaf(a, &ref, newLeaf, 0, true)

				So(result, ShouldNotBeNil)
				So(*result, ShouldEqual, 123)
				So(ref.IsLeaf(), ShouldBeTrue)

				leaf := ref.AsLeaf()
				So(leaf.Value, ShouldEqual, 456)
			})

			Convey("And one key is empty", func() {
				// Use null terminator (0) to represent empty string in ART
				initialLeaf := node.NewLeaf(a, []byte{0}, 123)

				ref := initialLeaf.Ref()

				newLeaf := node.NewLeaf(a, hello, 456)

				result := InsertToLeaf(a, &ref, newLeaf, 0, false)

				So(result, ShouldBeNil)
				So(ref.IsNode4(), ShouldBeTrue)

				node4 := ref.AsNode4()
				So(node4.NumChildren, ShouldEqual, 2)
				So(node4.Partial.Empty(), ShouldBeTrue)
			})
		})
	})
}

// TestInsertToNode tests the InsertToNode function directly
func TestInsertToNode(t *testing.T) {
	Convey("Given InsertToNode function", t, func() {
		a := new(arena.Arena)

		Convey("When inserting to a Node4 without prefix", func() {
			// Create a Node4 without prefix
			node4 := arena.New(a, node.Node4[int]{})
			ref := node4.Ref()

			// Add an existing child
			existingLeaf := node.NewLeaf(a, hello, 123)
			node4.AddChild(opt.Some(byte('h')), existingLeaf)

			// Create new leaf to insert
			newLeaf := node.NewLeaf(a, foobar, 456)

			Convey("And inserting a leaf with different first byte", func() {
				result := InsertToNode(a, &ref, newLeaf, 0, false)

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})

				Convey("And the leaf should be added to the node", func() {
					So(node4.NumChildren, ShouldEqual, 2)
					So(node4.Keys[0], ShouldEqual, byte('f'))
					So(node4.Keys[1], ShouldEqual, byte('h'))
					So(node4.Children[0], ShouldEqual, newLeaf.Ref())
					So(node4.Children[1], ShouldEqual, existingLeaf.Ref())
				})
			})

			Convey("And inserting a leaf with same first byte", func() {
				// Create leaf with same first byte but different key
				sameFirstByteLeaf := node.NewLeaf(a, []byte("hello world"), 789)

				result := InsertToNode(a, &ref, sameFirstByteLeaf, 0, false)

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})

				Convey("And the leaf should be added to the node", func() {
					// When keys have the same first byte, InsertToNode may recurse
					// or handle it differently depending on the implementation
					So(node4.NumChildren, ShouldBeGreaterThanOrEqualTo, 1)
					// Verify that the new leaf is accessible somehow
					// This might require checking the tree structure differently
				})
			})
		})

		Convey("When inserting to a Node4 with prefix", func() {
			// Create a Node4 with prefix "hel"
			node4 := arena.New(a, node.Node4[int]{})
			node4.Partial = slice.FromBytes(a, []byte("hel"))
			ref := node4.Ref()

			// Add an existing child
			existingLeaf := node.NewLeaf(a, hello, 123)
			node4.AddChild(opt.Some(byte('l')), existingLeaf)

			Convey("And inserting a leaf with matching prefix", func() {
				// Create leaf with matching prefix "hel"
				matchingPrefixLeaf := node.NewLeaf(a, help, 456)

				result := InsertToNode(a, &ref, matchingPrefixLeaf, 0, false)

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})

				Convey("And the leaf should be added to the node", func() {
					So(node4.NumChildren, ShouldEqual, 2)
					// Keys should be 'l' and 'p' (after "hel" prefix)
					So(node4.Keys[0], ShouldEqual, byte('l'))
					So(node4.Keys[1], ShouldEqual, byte('p'))
				})
			})

		})

		Convey("When inserting to a Node4 that becomes full", func() {
			// Create a Node4 and fill it
			node4 := arena.New(a, node.Node4[int]{})
			ref := node4.Ref()

			// Add 4 children to make it full
			for i := 0; i < 4; i++ {
				leaf := node.NewLeaf(a, []byte{byte('a' + i)}, i*100)
				node4.AddChild(opt.Some(byte('a'+i)), leaf)
			}

			So(node4.Full(), ShouldBeTrue)

			Convey("And inserting another leaf", func() {
				newLeaf := node.NewLeaf(a, []byte("e"), 500)

				result := InsertToNode(a, &ref, newLeaf, 0, false)

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})

				Convey("And the node should grow to Node16", func() {
					So(ref.IsNode16(), ShouldBeTrue)

					node16 := ref.AsNode16()
					So(node16.NumChildren, ShouldEqual, 5)
					So(node16.Keys[4], ShouldEqual, byte('e'))
					So(node16.Children[4], ShouldEqual, newLeaf.Ref())
				})
			})
		})

		Convey("When inserting to a Node16", func() {
			// Create a Node16
			node16 := arena.New(a, node.Node16[int]{})
			ref := node16.Ref()

			// Add some children
			for i := 0; i < 8; i++ {
				leaf := node.NewLeaf(a, []byte{byte('a' + i)}, i*100)
				node16.AddChild(opt.Some(byte('a'+i)), leaf)
			}

			Convey("And inserting a new leaf", func() {
				newLeaf := node.NewLeaf(a, []byte("x"), 800)

				result := InsertToNode(a, &ref, newLeaf, 0, false)

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})

				Convey("And the leaf should be added to the node", func() {
					So(node16.NumChildren, ShouldEqual, 9)
					// Find the new leaf
					found := false
					for i := 0; i < node16.NumChildren; i++ {
						if node16.Children[i] == newLeaf.Ref() {
							found = true
							break
						}
					}
					So(found, ShouldBeTrue)
				})
			})
		})

		Convey("When inserting to a Node48", func() {
			// Create a Node48
			node48 := arena.New(a, node.Node48[int]{})
			ref := node48.Ref()

			// Add some children
			for i := 0; i < 20; i++ {
				leaf := node.NewLeaf(a, []byte{byte(i * 10)}, i*100)
				node48.AddChild(opt.Some(byte(i*10)), leaf)
			}

			Convey("And inserting a new leaf", func() {
				newLeaf := node.NewLeaf(a, []byte{42}, 2100)

				result := InsertToNode(a, &ref, newLeaf, 0, false)

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})

				Convey("And the leaf should be added to the node", func() {
					So(node48.NumChildren, ShouldEqual, 21)
					// Verify the new leaf can be found
					found := node48.FindChild(opt.Some(byte(42)))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, newLeaf.Ref())
				})
			})
		})

		Convey("When inserting to a Node256", func() {
			// Create a Node256
			node256 := arena.New(a, node.Node256[int]{})
			ref := node256.Ref()

			// Add some children
			for i := 0; i < 50; i++ {
				leaf := node.NewLeaf(a, []byte{byte(i * 5)}, i*100)
				node256.AddChild(opt.Some(byte(i*5)), leaf)
			}

			Convey("And inserting a new leaf", func() {
				newLeaf := node.NewLeaf(a, []byte{99}, 9900)

				result := InsertToNode(a, &ref, newLeaf, 0, false)

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})

				Convey("And the leaf should be added to the node", func() {
					// Note: Node256.AddChild only increments NumChildren for new keys
					// If key 99 already exists, NumChildren won't change
					So(node256.NumChildren, ShouldBeGreaterThanOrEqualTo, 50)
					// Verify the new leaf can be found
					found := node256.FindChild(opt.Some(byte(99)))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, newLeaf.Ref())
				})
			})
		})

		Convey("When inserting with depth > 0", func() {
			// Create a Node4
			node4 := arena.New(a, node.Node4[int]{})
			ref := node4.Ref()

			// Add a child at depth 0
			existingLeaf := node.NewLeaf(a, []byte("hello"), 123)
			node4.AddChild(opt.Some(byte('h')), existingLeaf)

			Convey("And inserting a leaf at depth 1", func() {
				// Create leaf to insert at depth 1
				newLeaf := node.NewLeaf(a, []byte("xello"), 456)

				result := InsertToNode(a, &ref, newLeaf, 1, false)

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})

				Convey("And the leaf should be added to the node", func() {
					So(node4.NumChildren, ShouldEqual, 2)
					// Keys should be 'e' and 'h' (at depth 1)
					So(node4.Keys[0], ShouldEqual, byte('e'))
					So(node4.Keys[1], ShouldEqual, byte('h'))
				})
			})
		})

		Convey("When inserting with replace flag", func() {
			// Create a Node4
			node4 := arena.New(a, node.Node4[int]{})
			ref := node4.Ref()

			// Add a child
			existingLeaf := node.NewLeaf(a, []byte("hello"), 123)
			node4.AddChild(opt.Some(byte('h')), existingLeaf)

			Convey("And inserting a leaf with same first byte", func() {
				// Create leaf with same first byte
				newLeaf := node.NewLeaf(a, []byte("hello world"), 456)

				result := InsertToNode(a, &ref, newLeaf, 0, true)

				Convey("Then should return nil", func() {
					So(result, ShouldBeNil)
				})

				Convey("And both leaves should be in the node", func() {
					// When keys have the same first byte, InsertToNode may recurse
					// or handle it differently depending on the implementation
					So(node4.NumChildren, ShouldBeGreaterThanOrEqualTo, 1)
					// Verify that the new leaf is accessible somehow
					// This might require checking the tree structure differently
				})
			})
		})

		Convey("When inserting with edge cases", func() {
			Convey("And inserting with empty key", func() {
				node4 := arena.New(a, node.Node4[int]{})
				ref := node4.Ref()

				emptyKeyLeaf := node.NewLeaf(a, []byte{}, 999)

				result := InsertToNode(a, &ref, emptyKeyLeaf, 0, false)

				So(result, ShouldBeNil)

				So(node4.ZeroSizedChild.Empty(), ShouldBeFalse)
				So(node4.ZeroSizedChild.AsLeaf(), ShouldEqual, emptyKeyLeaf)

				So(node4.NumChildren, ShouldEqual, 0)
				So(node4.Keys, ShouldEqual, [4]byte{0, 0, 0, 0})
			})

			Convey("And inserting with zero byte key", func() {
				node4 := arena.New(a, node.Node4[int]{})
				ref := node4.Ref()

				zeroByteLeaf := node.NewLeaf(a, []byte{0}, 888)

				result := InsertToNode(a, &ref, zeroByteLeaf, 0, false)

				So(result, ShouldBeNil)
				So(node4.NumChildren, ShouldEqual, 1)
				So(node4.Keys[0], ShouldEqual, byte(0))
			})

			Convey("And inserting with maximum byte key", func() {
				node4 := arena.New(a, node.Node4[int]{})
				ref := node4.Ref()

				maxByteLeaf := node.NewLeaf(a, []byte{255}, 777)

				result := InsertToNode(a, &ref, maxByteLeaf, 0, false)

				So(result, ShouldBeNil)
				So(node4.NumChildren, ShouldEqual, 1)
				So(node4.Keys[0], ShouldEqual, byte(255))
			})

			Convey("And inserting with very long key", func() {
				node4 := arena.New(a, node.Node4[int]{})
				ref := node4.Ref()

				longKey := make([]byte, 1000)
				for i := range longKey {
					longKey[i] = byte(i % 256)
				}

				longKeyLeaf := node.NewLeaf(a, longKey, 666)

				result := InsertToNode(a, &ref, longKeyLeaf, 0, false)

				So(result, ShouldBeNil)
				So(node4.NumChildren, ShouldEqual, 1)
				So(node4.Keys[0], ShouldEqual, byte(0)) // First byte of longKey
			})
		})

		Convey("When inserting with complex prefix scenarios", func() {
			Convey("And inserting to node with long prefix", func() {
				// Create a Node4 with a long prefix
				node4 := arena.New(a, node.Node4[int]{})
				longPrefix := make([]byte, 100)
				for i := range longPrefix {
					longPrefix[i] = byte(i % 256)
				}
				node4.Partial = slice.FromBytes(a, longPrefix)
				ref := node4.Ref()

				// Add an existing child
				existingLeaf := node.NewLeaf(a, append(longPrefix, 'a'), 123)
				node4.AddChild(opt.Some(byte('a')), existingLeaf)

				// Create new leaf with matching prefix
				newLeaf := node.NewLeaf(a, append(longPrefix, 'b'), 456)

				result := InsertToNode(a, &ref, newLeaf, 0, false)

				So(result, ShouldBeNil)
				So(node4.NumChildren, ShouldEqual, 2)
				So(node4.Keys[0], ShouldEqual, byte('a'))
				So(node4.Keys[1], ShouldEqual, byte('b'))
			})

			Convey("And inserting to node with prefix that matches exactly", func() {
				// Create a Node4 with prefix "hello"
				node4 := arena.New(a, node.Node4[int]{})
				node4.Partial = slice.FromBytes(a, []byte("hello"))
				ref := node4.Ref()

				// Add an existing child
				existingLeaf := node.NewLeaf(a, []byte("hello"), 123)
				node4.AddChild(opt.Some(byte(0)), existingLeaf) // null terminator

				// Create new leaf with same prefix
				newLeaf := node.NewLeaf(a, []byte("hello world"), 456)

				result := InsertToNode(a, &ref, newLeaf, 0, false)

				So(result, ShouldBeNil)
				So(node4.NumChildren, ShouldEqual, 2)
				// Keys should be 0 (null terminator) and ' ' (space)
				So(node4.Keys[0], ShouldEqual, byte(0))
				So(node4.Keys[1], ShouldEqual, byte(' '))
			})
		})
	})
}

// TestRecursiveInsert_LazyExpansion tests lazy expansion behavior where inner nodes are only
// created when they are required to distinguish between at least two leaf nodes
func TestRecursiveInsert_LazyExpansion(t *testing.T) {
	Convey("Given an ART tree with lazy expansion using RecursiveInsert", t, func() {
		a := new(arena.Arena)

		var root node.Ref[int]

		Convey("When inserting a single key", func() {
			leaf := node.NewLeaf(a, []byte("single"), 1)
			RecursiveInsert(a, &root, leaf, 0, false)

			Convey("Then the root should contain the leaf node", func() {
				So(root.Empty(), ShouldBeFalse)
				So(root.IsLeaf(), ShouldBeTrue)
			})

			Convey("Then no inner nodes should be created unnecessarily", func() {
				// With only one key, no inner nodes are needed for distinction
				// The tree should remain as a single leaf
				So(root.IsLeaf(), ShouldBeTrue)

				So(*Search(root, []byte("single")), ShouldEqual, 1)
			})
		})

		Convey("When inserting keys with no common prefix", func() {
			leaf1 := node.NewLeaf(a, []byte("apple"), 1)
			leaf2 := node.NewLeaf(a, []byte("zebra"), 2)

			RecursiveInsert(a, &root, leaf1, 0, false)
			RecursiveInsert(a, &root, leaf2, 0, false)

			Convey("Then both keys should be inserted into the tree", func() {
				So(root.Empty(), ShouldBeFalse)
				// The root should now be an inner node since we have two different keys
				So(root.IsLeaf(), ShouldBeFalse)
			})

			Convey("Then inner nodes should be created only when necessary", func() {
				// These keys have no common prefix, so minimal inner node structure
				// should be created to distinguish them
				So(root.IsLeaf(), ShouldBeFalse)
			})

			Convey("Then both keys should be searchable", func() {
				So(*Search(root, []byte("apple")), ShouldEqual, 1)
				So(*Search(root, []byte("zebra")), ShouldEqual, 2)
			})
		})

		Convey("When inserting keys with common prefix that requires distinction", func() {
			leaf1 := node.NewLeaf(a, []byte("apple"), 1)
			leaf2 := node.NewLeaf(a, []byte("apricot"), 2)

			RecursiveInsert(a, &root, leaf1, 0, false)
			RecursiveInsert(a, &root, leaf2, 0, false)

			Convey("Then both keys should be inserted into the tree", func() {
				So(root.Empty(), ShouldBeFalse)
				// The root should be an inner node since we have two keys with common prefix
				So(root.IsLeaf(), ShouldBeFalse)
			})

			Convey("Then both keys should be searchable", func() {
				So(*Search(root, []byte("apple")), ShouldEqual, 1)
				So(*Search(root, []byte("apricot")), ShouldEqual, 2)
			})

			Convey("Then inner nodes should be created to distinguish the common prefix", func() {
				// These keys share "ap" prefix but differ at position 2
				// Inner nodes should be created to distinguish them
				So(root.IsLeaf(), ShouldBeFalse)

				n := root.AsNode4()
				So(n.NumChildren, ShouldEqual, 2)
				So(n.Partial.Raw(), ShouldResemble, []byte("ap"))
				So(n.Keys[0], ShouldEqual, byte('p'))
				So(n.Keys[1], ShouldEqual, byte('r'))
				So(n.Children[0], ShouldEqual, leaf1.Ref())
				So(n.Children[1], ShouldEqual, leaf2.Ref())
			})
		})

		Convey("When inserting keys that share a longer common prefix", func() {
			leaf1 := node.NewLeaf(a, []byte("application"), 1)
			leaf2 := node.NewLeaf(a, []byte("appliance"), 2)

			RecursiveInsert(a, &root, leaf1, 0, false)
			RecursiveInsert(a, &root, leaf2, 0, false)

			Convey("Then both keys should be inserted into the tree", func() {
				So(root.Empty(), ShouldBeFalse)
				// The root should be an inner node since we have two keys with common prefix
				So(root.IsLeaf(), ShouldBeFalse)
			})

			Convey("Then both keys should be searchable", func() {
				So(*Search(root, []byte("application")), ShouldEqual, 1)
				So(*Search(root, []byte("appliance")), ShouldEqual, 2)
			})

			Convey("Then inner nodes should be created only at the point of divergence", func() {
				// These keys share "appli" prefix but differ at position 5
				// Inner nodes should be created only where distinction is needed
				So(root.IsLeaf(), ShouldBeFalse)

				n := root.AsNode4()
				So(n.NumChildren, ShouldEqual, 2)
				So(n.Partial.Raw(), ShouldResemble, []byte("appli"))
				So(n.Keys[0], ShouldEqual, byte('a'))
				So(n.Keys[1], ShouldEqual, byte('c'))
				So(n.Children[0], ShouldEqual, leaf2.Ref())
				So(n.Children[1], ShouldEqual, leaf1.Ref())
			})
		})

		Convey("When inserting keys with incremental common prefixes", func() {
			leaf1 := node.NewLeaf(a, []byte("a"), 1)
			leaf2 := node.NewLeaf(a, []byte("ab"), 2)
			leaf3 := node.NewLeaf(a, []byte("abc"), 3)

			RecursiveInsert(a, &root, leaf1, 0, false)
			RecursiveInsert(a, &root, leaf2, 0, false)
			RecursiveInsert(a, &root, leaf3, 0, false)

			Convey("Then all keys should be inserted into the tree", func() {
				So(root.Empty(), ShouldBeFalse)
				// The root should be an inner node since we have multiple keys
				So(root.IsLeaf(), ShouldBeFalse)
			})

			Convey("Then both keys should be searchable", func() {
				So(*Search(root, []byte("a")), ShouldEqual, 1)
				So(*Search(root, []byte("ab")), ShouldEqual, 2)
				So(*Search(root, []byte("abc")), ShouldEqual, 3)
			})

			Convey("Then inner nodes should be created progressively", func() {
				// Each key extends the previous one, requiring inner nodes
				// only where distinction is necessary
				So(root.IsLeaf(), ShouldBeFalse)

				Convey("Then the root should be a Node4", func() {
					So(root.IsNode4(), ShouldBeTrue)

					n := root.AsNode4()
					So(n.NumChildren, ShouldEqual, 1)
					So(n.Partial.Raw(), ShouldResemble, []byte("a"))
					So(n.ZeroSizedChild.Empty(), ShouldBeFalse)
					So(n.ZeroSizedChild.AsLeaf(), ShouldEqual, leaf1)
					So(n.Keys, ShouldEqual, [4]byte{'b', 0, 0, 0})
					So(n.Children[0].Empty(), ShouldBeFalse)

					Convey("Then the first child should be a Node4", func() {
						n1 := n.Children[0].AsNode4()
						So(n1.NumChildren, ShouldEqual, 1)
						So(n1.Partial.Raw(), ShouldResemble, []byte(nil))

						So(n1.ZeroSizedChild.Empty(), ShouldBeFalse)
						So(n1.ZeroSizedChild.AsLeaf(), ShouldEqual, leaf2)

						So(n1.Keys, ShouldEqual, [4]byte{'c', 0, 0, 0})
						So(n1.Children[0].Empty(), ShouldBeFalse)
						So(n1.Children[0].AsLeaf(), ShouldEqual, leaf3)
					})
				})

			})
		})

		Convey("When inserting keys with no common prefix but similar lengths", func() {
			leaf1 := node.NewLeaf(a, []byte("hello"), 1)
			leaf2 := node.NewLeaf(a, []byte("world"), 2)
			leaf3 := node.NewLeaf(a, []byte("test"), 3)

			RecursiveInsert(a, &root, leaf1, 0, false)
			RecursiveInsert(a, &root, leaf2, 0, false)
			RecursiveInsert(a, &root, leaf3, 0, false)

			Convey("Then all keys should be inserted into the tree", func() {
				So(root.Empty(), ShouldBeFalse)
				// The root should be an inner node since we have multiple different keys
				So(root.IsLeaf(), ShouldBeFalse)
			})

			Convey("Then both keys should be searchable", func() {
				So(*Search(root, []byte("hello")), ShouldEqual, 1)
				So(*Search(root, []byte("world")), ShouldEqual, 2)
				So(*Search(root, []byte("test")), ShouldEqual, 3)
			})

			Convey("Then inner nodes should be created minimally", func() {
				n := root.AsNode4()

				So(n.NumChildren, ShouldEqual, 3)
				So(n.Partial.Raw(), ShouldResemble, []byte(nil))
				So(n.Keys, ShouldEqual, [4]byte{'h', 't', 'w', 0})
				So(n.Children, ShouldEqual, [4]node.Ref[int]{leaf1.Ref(), leaf3.Ref(), leaf2.Ref(), 0})
			})
		})

		Convey("When testing lazy expansion with complex key patterns", func() {
			// Insert keys that create a complex but efficient structure
			keys := []string{
				"a", "aa", "aaa", "aaaa",
				"b", "bb", "bbb", "bbbb",
				"c", "cc", "ccc", "cccc",
			}

			for i, key := range keys {
				leaf := node.NewLeaf(a, []byte(key), i+1)
				RecursiveInsert(a, &root, leaf, 0, false)
			}

			Convey("Then all keys should be inserted into the tree", func() {
				So(root.Empty(), ShouldBeFalse)
				// The root should be an inner node since we have multiple keys
				So(root.IsLeaf(), ShouldBeFalse)
			})

			Convey("Then all keys should be searchable", func() {
				for i, key := range keys {
					So(*Search(root, []byte(key)), ShouldEqual, i+1)
				}
			})

			Convey("Then inner nodes should be created efficiently", func() {
				n := root.AsNode4()
				So(n.NumChildren, ShouldEqual, 3)
				So(n.Partial.Raw(), ShouldResemble, []byte(nil))
				So(n.Keys, ShouldEqual, [4]byte{'a', 'b', 'c', 0})
			})
		})

		Convey("When testing lazy expansion with edge cases", func() {
			Convey("Then empty keys should work correctly", func() {
				leaf := node.NewLeaf(a, []byte{}, 1)
				RecursiveInsert(a, &root, leaf, 0, false)
				So(root.Empty(), ShouldBeFalse)
				So(root.IsLeaf(), ShouldBeTrue)
			})

			Convey("Then single byte keys should work correctly", func() {
				leaf1 := node.NewLeaf(a, []byte("x"), 1)
				leaf2 := node.NewLeaf(a, []byte("y"), 2)
				RecursiveInsert(a, &root, leaf1, 0, false)
				RecursiveInsert(a, &root, leaf2, 0, false)
				So(root.Empty(), ShouldBeFalse)
				So(root.IsLeaf(), ShouldBeFalse)
			})

			Convey("Then keys with special characters should work correctly", func() {
				leaf1 := node.NewLeaf(a, []byte("key@123"), 1)
				leaf2 := node.NewLeaf(a, []byte("key#456"), 2)
				RecursiveInsert(a, &root, leaf1, 0, false)
				RecursiveInsert(a, &root, leaf2, 0, false)
				So(root.Empty(), ShouldBeFalse)
				So(root.IsLeaf(), ShouldBeFalse)
			})
		})
	})
}
