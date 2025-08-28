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
			children := make([]*Leaf[any], 4)
			for i := 0; i < 4; i++ {
				children[i] = NewLeaf[any](a, []byte{byte('a' + i)}, nil)
			}

			Convey("Adding first child", func() {
				n.AddChild(int('a'), children[0])
				So(n.NumChildren, ShouldEqual, 1)
				So(n.Keys[0], ShouldEqual, byte('a'))
				So(n.Children[0], ShouldEqual, children[0].Ref())
			})

			Convey("Adding multiple children", func() {
				n.AddChild(int('a'), children[0])
				n.AddChild(int('b'), children[1])
				n.AddChild(int('c'), children[2])

				So(n.NumChildren, ShouldEqual, 3)
				So(n.Keys[0], ShouldEqual, byte('a'))
				So(n.Keys[1], ShouldEqual, byte('b'))
				So(n.Keys[2], ShouldEqual, byte('c'))
				So(n.Children[0], ShouldEqual, children[0].Ref())
				So(n.Children[1], ShouldEqual, children[1].Ref())
				So(n.Children[2], ShouldEqual, children[2].Ref())
			})

			Convey("Adding children to maintain sorted order", func() {
				// Add in reverse order
				n.AddChild(int('d'), children[3])
				n.AddChild(int('b'), children[1])
				n.AddChild(int('a'), children[0])

				So(n.NumChildren, ShouldEqual, 3)
				So(n.Keys[0], ShouldEqual, byte('a'))
				So(n.Keys[1], ShouldEqual, byte('b'))
				So(n.Keys[2], ShouldEqual, byte('d'))
				So(n.Children[0], ShouldEqual, children[0].Ref())
				So(n.Children[1], ShouldEqual, children[1].Ref())
				So(n.Children[2], ShouldEqual, children[3].Ref())
			})
		})

		Convey("When finding children", func() {
			// Setup children
			children := make([]*Leaf[any], 4)
			for i := 0; i < 4; i++ {
				children[i] = NewLeaf[any](a, []byte{byte('a' + i)}, nil)
				n.AddChild(int('a'+i), children[i])
			}

			Convey("Finding existing children", func() {
				for i := 0; i < 4; i++ {
					found := n.FindChild(int('a' + i))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, children[i].Ref())
				}
			})

			Convey("Finding non-existent children", func() {
				found := n.FindChild(int('e'))
				So(found, ShouldBeNil)

				found = n.FindChild(int('z'))
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
					n.AddChild(int('a'+i), child)
				}
				So(n.Full(), ShouldBeFalse)
			})

			Convey("Node with 4 children is full", func() {
				for i := 0; i < 4; i++ {
					child := NewLeaf[any](a, []byte{byte('a' + i)}, nil)
					n.AddChild(int('a'+i), child)
				}
				So(n.Full(), ShouldBeTrue)
			})
		})

		Convey("When growing to Node16", func() {
			// Setup children
			for i := 0; i < 4; i++ {
				child := NewLeaf[any](a, []byte{byte('a' + i)}, nil)
				n.AddChild(int('a'+i), child)
			}

			Convey("Growing should create Node16", func() {
				newNode := n.Grow(a)
				So(newNode.Type(), ShouldEqual, TypeNode16)
			})

			Convey("Growing should preserve all children", func() {
				newNode := n.Grow(a)
				node16 := newNode.(*Node16[any])

				So(node16.NumChildren, ShouldEqual, 4)
				// Check that all children are properly mapped
				for i := 0; i < 4; i++ {
					key := byte('a' + i)
					found := node16.FindChild(int(key))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, n.Children[i])
				}
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

				n.AddChild(int('c'), child3)
				n.AddChild(int('a'), child1)
				n.AddChild(int('b'), child2)

				So(n.Minimum(), ShouldEqual, child1)
				So(n.Maximum(), ShouldEqual, child3)
			})
		})
	})
}

func TestNode4_EdgeCases(t *testing.T) {
	Convey("Given a Node4 with edge cases", t, func() {
		a := &arena.Arena{}
		n := arena.New(a, Node4[any]{})

		Convey("When adding duplicate keys", func() {
			child1 := NewLeaf[any](a, []byte("a"), nil)
			child2 := NewLeaf[any](a, []byte("a"), nil)

			n.AddChild(int('a'), child1)
			n.AddChild(int('a'), child2)

			Convey("Should not replace the existing child", func() {
				So(n.NumChildren, ShouldEqual, 2)

				found := n.FindChild(int('a'))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child1.Ref())
			})
		})

		Convey("When adding zero byte key", func() {
			child := NewLeaf[any](a, []byte{0}, nil)
			n.AddChild(int(0), child)

			So(n.NumChildren, ShouldEqual, 1)
			found := n.FindChild(int(0))
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding 255 byte key", func() {
			child := NewLeaf[any](a, []byte{255}, nil)
			n.AddChild(int(255), child)

			So(n.NumChildren, ShouldEqual, 1)
			found := n.FindChild(int(255))
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding children at boundaries", func() {
			// Add children at the beginning and end of the byte range
			childStart := NewLeaf[any](a, []byte{0}, nil)
			childEnd := NewLeaf[any](a, []byte{255}, nil)

			n.AddChild(int(0), childStart)
			n.AddChild(int(255), childEnd)

			So(n.NumChildren, ShouldEqual, 2)
			So(n.Keys[0], ShouldEqual, byte(0))
			So(n.Keys[1], ShouldEqual, byte(255))
		})
	})
}

func TestNode4_Performance(t *testing.T) {
	Convey("Given a Node4 with performance considerations", t, func() {
		a := &arena.Arena{}
		n := arena.New(a, Node4[any]{})

		Convey("When adding many children", func() {
			// Test that adding 4 children works correctly
			children := make([]*Leaf[any], 4)
			for i := 0; i < 4; i++ {
				children[i] = NewLeaf[any](a, []byte{byte(i)}, nil)
				n.AddChild(int(i), children[i])
			}

			So(n.NumChildren, ShouldEqual, 4)
			So(n.Full(), ShouldBeTrue)

			// Verify all children can be found
			for i := 0; i < 4; i++ {
				found := n.FindChild(int(i))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, children[i].Ref())
			}
		})

		Convey("When searching in sorted array", func() {
			// Add children in sorted order to test search performance
			for i := 0; i < 4; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 2)}, nil)
				n.AddChild(int(i*2), child)
			}

			// Test finding existing and non-existing keys
			So(n.FindChild(int(0)), ShouldNotBeNil)
			So(n.FindChild(int(6)), ShouldNotBeNil)
			So(n.FindChild(int(1)), ShouldBeNil) // Odd numbers don't exist
			So(n.FindChild(int(7)), ShouldBeNil) // Odd numbers don't exist
		})
	})
}

func TestNode4_RemoveChild(t *testing.T) {
	Convey("Given a Node4 with children", t, func() {
		a := &arena.Arena{}
		n := arena.New(a, Node4[any]{})

		// Setup children
		children := make([]*Leaf[any], 4)
		for i := 0; i < 4; i++ {
			children[i] = NewLeaf[any](a, []byte{byte('a' + i)}, nil)
			n.AddChild(int('a'+i), children[i])
		}

		So(n.NumChildren, ShouldEqual, 4)

		Convey("When removing the first child", func() {
			childRef := n.FindChild(int('a'))
			So(childRef, ShouldNotBeNil)

			n.RemoveChild(int('a'), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(n.NumChildren, ShouldEqual, 3)
			})

			Convey("And the child should not be found", func() {
				found := n.FindChild(int('a'))
				So(found, ShouldBeNil)
			})

			Convey("And remaining children should be shifted left", func() {
				So(n.Keys[:3], ShouldResemble, []byte{'b', 'c', 'd'})
				So(n.Children[0], ShouldEqual, children[1].Ref())
				So(n.Children[1], ShouldEqual, children[2].Ref())
				So(n.Children[2], ShouldEqual, children[3].Ref())
			})
		})

		Convey("When removing the middle child", func() {
			childRef := n.FindChild(int('b'))
			So(childRef, ShouldNotBeNil)

			n.RemoveChild(int('b'), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(n.NumChildren, ShouldEqual, 3)
			})

			Convey("And the child should not be found", func() {
				found := n.FindChild(int('b'))
				So(found, ShouldBeNil)
			})

			Convey("And remaining children should be properly shifted", func() {
				So(n.Keys[:3], ShouldResemble, []byte{'a', 'c', 'd'})
				So(n.Children[0], ShouldEqual, children[0].Ref())
				So(n.Children[1], ShouldEqual, children[2].Ref())
				So(n.Children[2], ShouldEqual, children[3].Ref())
			})
		})

		Convey("When removing the last child", func() {
			childRef := n.FindChild(int('d'))
			So(childRef, ShouldNotBeNil)

			n.RemoveChild(int('d'), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(n.NumChildren, ShouldEqual, 3)
			})

			Convey("And the child should not be found", func() {
				found := n.FindChild(int('d'))
				So(found, ShouldBeNil)
			})

			Convey("And remaining children should be unchanged", func() {
				So(n.Keys[:3], ShouldResemble, []byte{'a', 'b', 'c'})
			})
		})

		Convey("When removing multiple children", func() {
			// Remove 'b' first
			childRef := n.FindChild(int('b'))
			n.RemoveChild(int('b'), childRef)

			// Remove 'c' second
			childRef = n.FindChild(int('c'))
			n.RemoveChild(int('c'), childRef)

			Convey("Then NumChildren should be 2", func() {
				So(n.NumChildren, ShouldEqual, 2)
			})

			Convey("And only remaining children should be found", func() {
				So(n.FindChild(int('a')), ShouldNotBeNil)
				So(n.FindChild(int('b')), ShouldBeNil)
				So(n.FindChild(int('c')), ShouldBeNil)
				So(n.FindChild(int('d')), ShouldNotBeNil)
			})

			Convey("And keys should be properly ordered", func() {
				So(n.Keys[:2], ShouldResemble, []byte{'a', 'd'})
			})
		})
	})
}

func TestNode4_Shrink(t *testing.T) {
	Convey("Given a Node4", t, func() {
		a := &arena.Arena{}

		n := arena.New(a, Node4[any]{})
		n.Partial = slice.FromString(a, hello)

		child1 := NewLeaf[any](a, []byte("a"), nil)
		child2 := NewLeaf[any](a, []byte("b"), nil)
		child3 := NewLeaf[any](a, []byte("c"), nil)

		Convey("When shrinking with 3 or more children", func() {
			n.AddChild(int('a'), child1)
			n.AddChild(int('b'), child2)
			n.AddChild(int('c'), child3)

			So(n.NumChildren, ShouldEqual, 3)

			result := n.Shrink(a)

			Convey("Then should return the same node", func() {
				So(result, ShouldEqual, n)
			})

			Convey("And NumChildren should remain unchanged", func() {
				So(n.NumChildren, ShouldEqual, 3)
			})
		})

		Convey("When shrinking with exactly 2 children", func() {
			n.AddChild(int('a'), child1)
			n.AddChild(int('b'), child2)

			So(n.NumChildren, ShouldEqual, 2)

			result := n.Shrink(a)

			Convey("Then should return the same node", func() {
				So(result, ShouldEqual, n)
			})

			Convey("And NumChildren should remain unchanged", func() {
				So(n.NumChildren, ShouldEqual, 2)
			})
		})

		Convey("When shrinking with exactly 1 child", func() {
			n.AddChild(int('a'), child1)

			So(n.NumChildren, ShouldEqual, 1)

			result := n.Shrink(a)

			Convey("Then should return the child node", func() {
				So(result, ShouldEqual, child1)
			})

			Convey("And the original node should be freed", func() {
				// The original node is freed during shrinking
				So(result, ShouldNotEqual, n)
			})
		})

		Convey("When shrinking with no children", func() {
			So(n.NumChildren, ShouldEqual, 0)

			result := n.Shrink(a)

			Convey("Then should return nil", func() {
				So(result, ShouldBeNil)
			})

			Convey("And the original node should be freed", func() {
				// The original node is freed during shrinking
				So(result, ShouldNotEqual, n)
			})
		})
	})
}
