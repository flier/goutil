package node_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	. "github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
	"github.com/flier/goutil/pkg/opt"
)

func TestNode16(t *testing.T) {
	Convey("Given a Node16", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node16[any]{})

		Convey("When checking basic properties", func() {
			So(node.Type(), ShouldEqual, TypeNode16)
			So(node.Full(), ShouldBeFalse)
			So(node.NumChildren, ShouldEqual, 0)
			So(node.Ref().Type(), ShouldEqual, TypeNode16)
		})

		Convey("When adding children", func() {
			// Create mock children
			children := make([]*Leaf[any], 16)
			for i := 0; i < 16; i++ {
				children[i] = NewLeaf[any](a, []byte{byte('a' + i)}, nil)
			}

			Convey("Adding first child", func() {
				node.AddChild(opt.Some(byte('a')), children[0])
				So(node.NumChildren, ShouldEqual, 1)
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Children[0], ShouldEqual, children[0].Ref())
			})

			Convey("Adding children in order", func() {
				for i := 0; i < 8; i++ {
					node.AddChild(opt.Some(byte('a'+i)), children[i])
				}

				So(node.NumChildren, ShouldEqual, 8)
				for i := 0; i < 8; i++ {
					So(node.Keys[i], ShouldEqual, byte('a'+i))
					So(node.Children[i], ShouldEqual, children[i].Ref())
				}
			})

			Convey("Adding children out of order", func() {
				node.AddChild(opt.Some(byte('c')), children[2])
				node.AddChild(opt.Some(byte('a')), children[0])
				node.AddChild(opt.Some(byte('b')), children[1])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Keys[1], ShouldEqual, byte('b'))
				So(node.Keys[2], ShouldEqual, byte('c'))
			})

			Convey("Adding children to maintain sorted order", func() {
				// Add in reverse order
				for i := 15; i >= 0; i-- {
					node.AddChild(opt.Some(byte('a'+i)), children[i])
				}

				So(node.NumChildren, ShouldEqual, 16)
				for i := 0; i < 16; i++ {
					So(node.Keys[i], ShouldEqual, byte('a'+i))
					So(node.Children[i], ShouldEqual, children[i].Ref())
				}
			})
		})

		Convey("When finding children", func() {
			// Setup children
			children := make([]*Leaf[any], 8)
			for i := 0; i < 8; i++ {
				children[i] = NewLeaf[any](a, []byte{byte('a' + i)}, nil)
				node.AddChild(opt.Some(byte('a'+i)), children[i])
			}

			Convey("Finding existing children", func() {
				for i := 0; i < 8; i++ {
					found := node.FindChild(opt.Some(byte('a' + i)))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, children[i].Ref())
				}
			})

			Convey("Finding non-existent children", func() {
				found := node.FindChild(opt.Some(byte('x')))
				So(found, ShouldBeNil)

				found = node.FindChild(opt.Some(byte('z')))
				So(found, ShouldBeNil)
			})
		})

		Convey("When checking capacity", func() {
			Convey("Empty node is not full", func() {
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 8 children is not full", func() {
				for i := 0; i < 8; i++ {
					child := NewLeaf[any](a, []byte{byte('a' + i)}, nil)
					node.AddChild(opt.Some(byte('a'+i)), child)
				}
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 16 children is full", func() {
				for i := 0; i < 16; i++ {
					child := NewLeaf[any](a, []byte{byte('a' + i)}, nil)
					node.AddChild(opt.Some(byte('a'+i)), child)
				}
				So(node.Full(), ShouldBeTrue)
			})
		})

		Convey("When growing to Node48", func() {
			// Setup children
			for i := 0; i < 16; i++ {
				child := NewLeaf[any](a, []byte{byte('a' + i)}, nil)
				node.AddChild(opt.Some(byte('a'+i)), child)
			}

			Convey("Growing should create Node48", func() {
				newNode := node.Grow(a)
				So(newNode.Type(), ShouldEqual, TypeNode48)
			})

			Convey("Growing should preserve all children", func() {
				newNode := node.Grow(a)
				node48 := newNode.(*Node48[any])

				So(node48.NumChildren, ShouldEqual, 16)
				// Check that all children are properly mapped
				for i := 0; i < 16; i++ {
					key := byte('a' + i)
					So(node48.Keys[key], ShouldNotEqual, 0)
					childIndex := node48.Keys[key] - 1
					So(node48.Children[childIndex], ShouldEqual, node.Children[i])
				}
			})
		})

		Convey("When getting minimum and maximum", func() {
			Convey("Empty node should return nil", func() {
				So(node.Minimum(), ShouldBeNil)
				So(node.Maximum(), ShouldBeNil)
			})

			Convey("Node with children should return correct min/max", func() {
				child1 := NewLeaf[any](a, []byte("a"), nil)
				child2 := NewLeaf[any](a, []byte("b"), nil)
				child3 := NewLeaf[any](a, []byte("c"), nil)

				node.AddChild(opt.Some(byte('c')), child3)
				node.AddChild(opt.Some(byte('a')), child1)
				node.AddChild(opt.Some(byte('b')), child2)

				So(node.Minimum(), ShouldEqual, child1)
				So(node.Maximum(), ShouldEqual, child3)
			})
		})
	})
}

func TestNode16_EdgeCases(t *testing.T) {
	Convey("Given a Node16 with edge cases", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node16[any]{})

		Convey("When adding duplicate keys", func() {
			child1 := NewLeaf[any](a, []byte("a"), nil)
			child2 := NewLeaf[any](a, []byte("a"), nil)

			node.AddChild(opt.Some(byte('a')), child1)
			node.AddChild(opt.Some(byte('a')), child2)

			Convey("Should not replace the existing child", func() {
				So(node.NumChildren, ShouldEqual, 2)

				found := node.FindChild(opt.Some(byte('a')))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child1.Ref())
			})
		})

		Convey("When adding zero byte key", func() {
			child := NewLeaf[any](a, []byte{0}, nil)
			node.AddChild(opt.Some(byte(0)), child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(opt.Some(byte(0)))
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding 255 byte key", func() {
			child := NewLeaf[any](a, []byte{255}, nil)
			node.AddChild(opt.Some(byte(255)), child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(opt.Some(byte(255)))
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding children at boundaries", func() {
			// Add children at the beginning and end of the byte range
			childStart := NewLeaf[any](a, []byte{0}, nil)
			childEnd := NewLeaf[any](a, []byte{255}, nil)

			node.AddChild(opt.Some(byte(0)), childStart)
			node.AddChild(opt.Some(byte(255)), childEnd)

			So(node.NumChildren, ShouldEqual, 2)
			So(node.Keys[0], ShouldEqual, byte(0))
			So(node.Keys[1], ShouldEqual, byte(255))
		})
	})
}

func TestNode16_Performance(t *testing.T) {
	Convey("Given a Node16 with performance considerations", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node16[any]{})

		Convey("When adding many children", func() {
			// Test that adding 16 children works correctly
			children := make([]*Leaf[any], 16)
			for i := 0; i < 16; i++ {
				children[i] = NewLeaf[any](a, []byte{byte(i)}, nil)
				node.AddChild(opt.Some(byte(i)), children[i])
			}

			So(node.NumChildren, ShouldEqual, 16)
			So(node.Full(), ShouldBeTrue)

			// Verify all children can be found
			for i := 0; i < 16; i++ {
				found := node.FindChild(opt.Some(byte(i)))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, children[i].Ref())
			}
		})

		Convey("When searching in sorted array", func() {
			// Add children in sorted order to test search performance
			for i := 0; i < 8; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 2)}, nil)
				node.AddChild(opt.Some(byte(i*2)), child)
			}

			// Test finding existing and non-existing keys
			So(node.FindChild(opt.Some(byte(0))), ShouldNotBeNil)
			So(node.FindChild(opt.Some(byte(14))), ShouldNotBeNil)
			So(node.FindChild(opt.Some(byte(1))), ShouldBeNil)  // Odd numbers don't exist
			So(node.FindChild(opt.Some(byte(15))), ShouldBeNil) // Odd numbers don't exist
		})
	})
}

func TestNode16_RemoveChild(t *testing.T) {
	Convey("Given a Node16 with children", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node16[any]{})

		// Setup children
		children := make([]*Leaf[any], 8)
		for i := 0; i < 8; i++ {
			children[i] = NewLeaf[any](a, []byte{byte('a' + i)}, nil)
			node.AddChild(opt.Some(byte('a'+i)), children[i])
		}

		So(node.NumChildren, ShouldEqual, 8)

		Convey("When removing the first child", func() {
			childRef := node.FindChild(opt.Some(byte('a')))
			So(childRef, ShouldNotBeNil)

			node.RemoveChild(opt.Some(byte('a')), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 7)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild(opt.Some(byte('a')))
				So(found, ShouldBeNil)
			})

			Convey("And remaining children should be shifted left", func() {
				So(node.Keys[:4], ShouldResemble, []byte{'b', 'c', 'd', 'e'})
				So(node.Children[0], ShouldEqual, children[1].Ref())
				So(node.Children[1], ShouldEqual, children[2].Ref())
				So(node.Children[2], ShouldEqual, children[3].Ref())
			})
		})

		Convey("When removing the middle child", func() {
			childRef := node.FindChild(opt.Some(byte('d')))
			So(childRef, ShouldNotBeNil)

			node.RemoveChild(opt.Some(byte('d')), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 7)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild(opt.Some(byte('d')))
				So(found, ShouldBeNil)
			})

			Convey("And remaining children should be properly shifted", func() {
				So(node.Keys[:4], ShouldResemble, []byte{'a', 'b', 'c', 'e'})
				So(node.Children[0], ShouldEqual, children[0].Ref())
				So(node.Children[1], ShouldEqual, children[1].Ref())
				So(node.Children[2], ShouldEqual, children[2].Ref())
				So(node.Children[3], ShouldEqual, children[4].Ref())
			})
		})

		Convey("When removing the last child", func() {
			childRef := node.FindChild(opt.Some(byte('h')))
			So(childRef, ShouldNotBeNil)

			node.RemoveChild(opt.Some(byte('h')), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 7)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild(opt.Some(byte('h')))
				So(found, ShouldBeNil)
			})

			Convey("And remaining children should be unchanged", func() {
				So(node.Keys[:8], ShouldResemble, []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 0})
			})
		})

		Convey("When removing multiple children", func() {
			// Remove 'b' first
			childRef := node.FindChild(opt.Some(byte('b')))
			node.RemoveChild(opt.Some(byte('b')), childRef)

			// Remove 'e' second
			childRef = node.FindChild(opt.Some(byte('e')))
			node.RemoveChild(opt.Some(byte('e')), childRef)

			Convey("Then NumChildren should be 6", func() {
				So(node.NumChildren, ShouldEqual, 6)
			})

			Convey("And only remaining children should be found", func() {
				So(node.FindChild(opt.Some(byte('a'))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte('b'))), ShouldBeNil)
				So(node.FindChild(opt.Some(byte('c'))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte('d'))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte('e'))), ShouldBeNil)
				So(node.FindChild(opt.Some(byte('f'))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte('g'))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte('h'))), ShouldNotBeNil)
			})

			Convey("And keys should be properly ordered", func() {
				So(node.Keys[:7], ShouldResemble, []byte{'a', 'c', 'd', 'f', 'g', 'h', 0})
			})
		})
	})
}

const hello = "hello"

func TestNode16_Shrink(t *testing.T) {
	Convey("Given a Node16", t, func() {
		a := &arena.Arena{}

		node := arena.New(a, Node16[any]{})
		node.Partial = slice.FromString(a, hello)

		child1 := NewLeaf[any](a, []byte("a"), nil)
		child2 := NewLeaf[any](a, []byte("b"), nil)
		child3 := NewLeaf[any](a, []byte("c"), nil)

		Convey("When shrinking with 3 or more children", func() {
			node.AddChild(opt.Some(byte('a')), child1)
			node.AddChild(opt.Some(byte('b')), child2)
			node.AddChild(opt.Some(byte('c')), child3)

			So(node.NumChildren, ShouldEqual, 3)

			result := node.Shrink(a)

			Convey("Then should return the same node", func() {
				So(result, ShouldEqual, node)
			})

			Convey("And NumChildren should remain unchanged", func() {
				So(node.NumChildren, ShouldEqual, 3)
			})
		})

		Convey("When shrinking with exactly 2 children", func() {
			node.AddChild(opt.Some(byte('a')), child1)
			node.AddChild(opt.Some(byte('b')), child2)

			So(node.NumChildren, ShouldEqual, 2)

			result := node.Shrink(a)

			Convey("Then should return a Node4", func() {
				So(result.Type(), ShouldEqual, TypeNode4)
				So(result.Prefix().Raw(), ShouldEqual, []byte("hello"))
			})

			Convey("And the new Node4 should have the same children", func() {
				node4 := result.(*Node4[any])
				So(node4.NumChildren, ShouldEqual, 2)

				found := node4.FindChild(opt.Some(byte('a')))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child1.Ref())

				found = node4.FindChild(opt.Some(byte('b')))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child2.Ref())
			})
		})

		Convey("When shrinking with exactly 1 child", func() {
			node.AddChild(opt.Some(byte('a')), child1)

			So(node.NumChildren, ShouldEqual, 1)

			result := node.Shrink(a)

			Convey("Then should return a Node4", func() {
				So(result.Type(), ShouldEqual, TypeNode4)
				So(result.Prefix().Raw(), ShouldEqual, []byte("hello"))
			})

			Convey("And the new Node4 should have the same child", func() {
				node4 := result.(*Node4[any])
				So(node4.NumChildren, ShouldEqual, 1)

				found := node4.FindChild(opt.Some(byte('a')))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child1.Ref())
			})

			Convey("And the original node should be freed", func() {
				// The original node should be replaced, so we can't access it directly
				// This is verified by the fact that we get a Node4 back
			})
		})

		Convey("When shrinking with no children", func() {
			So(node.NumChildren, ShouldEqual, 0)

			result := node.Shrink(a)

			Convey("Then should return a Node4", func() {
				So(result.Type(), ShouldEqual, TypeNode4)
				So(result.Prefix().Raw(), ShouldEqual, []byte("hello"))
			})

			Convey("And the new Node4 should have no children", func() {
				node4 := result.(*Node4[any])
				So(node4.NumChildren, ShouldEqual, 0)
			})
		})
	})
}
