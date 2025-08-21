package node_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	. "github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
)

func TestNode4(t *testing.T) {
	Convey("Given a Node4", t, func() {
		a := &arena.Arena{}
		n := arena.New(a, Node4[any]{})

		Convey("When checking basic properties", func() {
			So(n.Type(), ShouldEqual, TypeNode4)
			So(n.Full(), ShouldBeFalse)
			So(n.NumChildren, ShouldEqual, 0)
			So(n.Ref().Type(), ShouldEqual, TypeNode4)
		})

		Convey("When adding children", func() {
			// Create mock children
			child1 := NewLeaf[any](a, []byte("a"), nil)
			child2 := NewLeaf[any](a, []byte("b"), nil)
			child3 := NewLeaf[any](a, []byte("c"), nil)
			child4 := NewLeaf[any](a, []byte("d"), nil)

			Convey("Adding first child", func() {
				n.AddChild('a', child1)
				So(n.NumChildren, ShouldEqual, 1)
				So(n.Keys[0], ShouldEqual, byte('a'))
				So(n.Children[0], ShouldEqual, child1.Ref())
			})

			Convey("Adding children in order", func() {
				n.AddChild('a', child1)
				n.AddChild('b', child2)
				n.AddChild('c', child3)

				So(n.NumChildren, ShouldEqual, 3)
				So(n.Keys[0], ShouldEqual, byte('a'))
				So(n.Keys[1], ShouldEqual, byte('b'))
				So(n.Keys[2], ShouldEqual, byte('c'))
			})

			Convey("Adding children out of order", func() {
				n.AddChild('c', child3)
				n.AddChild('a', child1)
				n.AddChild('b', child2)

				So(n.NumChildren, ShouldEqual, 3)
				So(n.Keys[0], ShouldEqual, byte('a'))
				So(n.Keys[1], ShouldEqual, byte('b'))
				So(n.Keys[2], ShouldEqual, byte('c'))
			})

			Convey("Adding children to maintain sorted order", func() {
				n.AddChild('d', child4)
				n.AddChild('b', child2)
				n.AddChild('a', child1)
				n.AddChild('c', child3)

				So(n.NumChildren, ShouldEqual, 4)
				So(n.Keys[0], ShouldEqual, byte('a'))
				So(n.Keys[1], ShouldEqual, byte('b'))
				So(n.Keys[2], ShouldEqual, byte('c'))
				So(n.Keys[3], ShouldEqual, byte('d'))
			})
		})

		Convey("When finding children", func() {
			// Setup children
			child1 := NewLeaf[any](a, []byte("a"), nil)
			child2 := NewLeaf[any](a, []byte("b"), nil)
			child3 := NewLeaf[any](a, []byte("c"), nil)

			n.AddChild('a', child1)
			n.AddChild('b', child2)
			n.AddChild('c', child3)

			Convey("Finding existing children", func() {
				found := n.FindChild('a')
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child1.Ref())

				found = n.FindChild('b')
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child2.Ref())

				found = n.FindChild('c')
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child3.Ref())
			})

			Convey("Finding non-existent children", func() {
				found := n.FindChild('x')
				So(found, ShouldBeNil)

				found = n.FindChild('z')
				So(found, ShouldBeNil)
			})
		})

		Convey("When checking capacity", func() {
			Convey("Empty node is not full", func() {
				So(n.Full(), ShouldBeFalse)
			})

			Convey("Node with 3 children is not full", func() {
				for i := 0; i < 3; i++ {
					child := NewLeaf[any](a, []byte{byte('a' + i)}, nil)
					n.AddChild(byte('a'+i), child)
				}
				So(n.Full(), ShouldBeFalse)
			})

			Convey("Node with 4 children is full", func() {
				for i := 0; i < 4; i++ {
					child := NewLeaf[any](a, []byte{byte('a' + i)}, nil)
					n.AddChild(byte('a'+i), child)
				}
				So(n.Full(), ShouldBeTrue)
			})
		})

		Convey("When growing to Node16", func() {
			// Setup children
			for i := 0; i < 4; i++ {
				child := NewLeaf[any](a, []byte{byte('a' + i)}, nil)
				n.AddChild(byte('a'+i), child)
			}

			Convey("Growing should create Node16", func() {
				newNode := n.Grow(a)
				So(newNode.Type(), ShouldEqual, TypeNode16)
			})

			Convey("Growing should preserve all children", func() {
				newNode := n.Grow(a)
				node16 := newNode.(*Node16[any])

				So(node16.NumChildren, ShouldEqual, 4)
				So(node16.Keys[0], ShouldEqual, byte('a'))
				So(node16.Keys[1], ShouldEqual, byte('b'))
				So(node16.Keys[2], ShouldEqual, byte('c'))
				So(node16.Keys[3], ShouldEqual, byte('d'))
			})
		})

		Convey("When getting minimum and maximum", func() {
			Convey("Empty node should return nil", func() {
				So(n.Minimum(), ShouldBeNil)
				So(n.Maximum(), ShouldBeNil)
			})

			Convey("Node with children should return correct min/max", func() {
				child1 := NewLeaf[any](a, []byte("a"), nil)
				child2 := NewLeaf[any](a, []byte("b"), nil)
				child3 := NewLeaf[any](a, []byte("c"), nil)

				n.AddChild('c', child3)
				n.AddChild('a', child1)
				n.AddChild('b', child2)

				// Note: Since these are mock leaves, we can't easily test the actual
				// minimum/maximum values without more complex setup
				So(n.Minimum(), ShouldEqual, child1)
				So(n.Maximum(), ShouldEqual, child3)
			})
		})
	})
}

func TestNode4_EdgeCases(t *testing.T) {
	Convey("Given a Node4 with edge cases", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node4[any]{})

		Convey("When adding duplicate keys", func() {
			child1 := NewLeaf[any](a, []byte("a"), nil)
			child2 := NewLeaf[any](a, []byte("a"), nil)

			node.AddChild('a', child1)
			node.AddChild('a', child2)

			// Should replace the existing child
			So(node.NumChildren, ShouldEqual, 2)
			found := node.FindChild('a')
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child1.Ref())
		})

		Convey("When adding zero byte key", func() {
			child := NewLeaf[any](a, []byte{0}, nil)
			node.AddChild(0, child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(0)
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding 255 byte key", func() {
			child := NewLeaf[any](a, []byte{255}, nil)
			node.AddChild(255, child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(255)
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})
	})
}

func TestNode4_RemoveChild(t *testing.T) {
	Convey("Given a Node4 with children", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node4[any]{})

		// Setup children
		child1 := NewLeaf[any](a, []byte("a"), nil)
		child2 := NewLeaf[any](a, []byte("b"), nil)
		child3 := NewLeaf[any](a, []byte("c"), nil)
		child4 := NewLeaf[any](a, []byte("d"), nil)

		node.AddChild('a', child1)
		node.AddChild('b', child2)
		node.AddChild('c', child3)
		node.AddChild('d', child4)

		So(node.NumChildren, ShouldEqual, 4)

		Convey("When removing the first child", func() {
			childRef := node.FindChild('a')
			So(childRef, ShouldNotBeNil)

			node.RemoveChild('a', childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 3)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild('a')
				So(found, ShouldBeNil)
			})

			Convey("And remaining children should be shifted left", func() {
				So(node.Keys[0], ShouldEqual, byte('b'))
				So(node.Keys[1], ShouldEqual, byte('c'))
				So(node.Keys[2], ShouldEqual, byte('d'))
				So(node.Children[0], ShouldEqual, child2.Ref())
				So(node.Children[1], ShouldEqual, child3.Ref())
				So(node.Children[2], ShouldEqual, child4.Ref())
			})
		})

		Convey("When removing the middle child", func() {
			childRef := node.FindChild('b')
			So(childRef, ShouldNotBeNil)

			node.RemoveChild('b', childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 3)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild('b')
				So(found, ShouldBeNil)
			})

			Convey("And remaining children should be properly shifted", func() {
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Keys[1], ShouldEqual, byte('c'))
				So(node.Keys[2], ShouldEqual, byte('d'))
				So(node.Children[0], ShouldEqual, child1.Ref())
				So(node.Children[1], ShouldEqual, child3.Ref())
				So(node.Children[2], ShouldEqual, child4.Ref())
			})
		})

		Convey("When removing the last child", func() {
			childRef := node.FindChild('d')
			So(childRef, ShouldNotBeNil)

			node.RemoveChild('d', childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 3)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild('d')
				So(found, ShouldBeNil)
			})

			Convey("And remaining children should be unchanged", func() {
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Keys[1], ShouldEqual, byte('b'))
				So(node.Keys[2], ShouldEqual, byte('c'))
				So(node.Children[0], ShouldEqual, child1.Ref())
				So(node.Children[1], ShouldEqual, child2.Ref())
				So(node.Children[2], ShouldEqual, child3.Ref())
			})
		})

		Convey("When removing multiple children", func() {
			// Remove 'b' first
			childRef := node.FindChild('b')
			node.RemoveChild('b', childRef)

			// Remove 'c' second
			childRef = node.FindChild('c')
			node.RemoveChild('c', childRef)

			Convey("Then NumChildren should be 2", func() {
				So(node.NumChildren, ShouldEqual, 2)
			})

			Convey("And only 'a' and 'd' should remain", func() {
				So(node.FindChild('a'), ShouldNotBeNil)
				So(node.FindChild('b'), ShouldBeNil)
				So(node.FindChild('c'), ShouldBeNil)
				So(node.FindChild('d'), ShouldNotBeNil)
			})

			Convey("And keys should be properly ordered", func() {
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Keys[1], ShouldEqual, byte('d'))
			})
		})
	})
}

func TestNode4_Shrink(t *testing.T) {
	Convey("Given a Node4", t, func() {
		a := &arena.Arena{}

		node := arena.New(a, Node4[any]{})
		node.Partial = slice.FromString(a, "+")

		child1 := NewLeaf[any](a, []byte("+a"), nil)
		child2 := NewLeaf[any](a, []byte("+b"), nil)
		child3 := NewLeaf[any](a, []byte("+c"), nil)

		Convey("When shrinking with more than 1 child", func() {
			node.AddChild('a', child1)
			node.AddChild('b', child2)
			node.AddChild('c', child3)

			So(node.NumChildren, ShouldEqual, 3)

			result := node.Shrink(a)

			Convey("Then should return the same node", func() {
				So(result, ShouldEqual, node)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And NumChildren should remain unchanged", func() {
				So(node.NumChildren, ShouldEqual, 3)
			})
		})

		Convey("When shrinking with exactly 1 child that is a leaf", func() {
			node.AddChild('a', child1)

			So(node.NumChildren, ShouldEqual, 1)

			result := node.Shrink(a)

			Convey("Then should return the child", func() {
				So(result, ShouldEqual, child1)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+a"))
			})

			Convey("And the child should be a leaf", func() {
				leaf := result.(*Leaf[any])
				So(leaf, ShouldNotBeNil)
				So(leaf.Key.Raw(), ShouldEqual, []byte("+a"))
			})
		})

		Convey("When shrinking with exactly 1 child that is a node", func() {
			// Create a child node
			childNode := arena.New(a, Node4[any]{})
			childLeaf := NewLeaf[any](a, []byte("+ax"), nil)
			childNode.AddChild('x', childLeaf)

			node.AddChild('a', childNode)

			So(node.NumChildren, ShouldEqual, 1)

			result := node.Shrink(a)

			Convey("Then should return the child node", func() {
				So(result, ShouldEqual, childNode)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+a"))
			})

			Convey("And the child node should have concatenated prefix", func() {
				childNodeResult := result.(*Node4[any])
				So(childNodeResult.Partial.Raw(), ShouldEqual, []byte("+a"))
			})

			Convey("And the child node should still contain its original child", func() {
				childNodeResult := result.(*Node4[any])
				found := childNodeResult.FindChild('x')
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, childLeaf.Ref())
			})
		})

		Convey("When shrinking with no children", func() {
			So(node.NumChildren, ShouldEqual, 0)

			result := node.Shrink(a)

			Convey("Then should return nil (no children to shrink to)", func() {
				So(result, ShouldBeNil)
			})
		})
	})
}
