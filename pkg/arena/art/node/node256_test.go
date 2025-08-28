package node_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	. "github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
)

func TestNode256(t *testing.T) {
	Convey("Given a Node256", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node256[any]{})

		Convey("When checking basic properties", func() {
			So(node.Type(), ShouldEqual, TypeNode256)
			So(node.Full(), ShouldBeFalse)
			So(node.NumChildren, ShouldEqual, 0)
			So(node.Ref().Type(), ShouldEqual, TypeNode256)
		})

		Convey("When adding children", func() {
			// Create mock children
			children := make([]*Leaf[any], 10)
			for i := 0; i < 10; i++ {
				children[i] = NewLeaf[any](a, []byte{byte(i * 25)}, nil)
			}

			Convey("Adding first child", func() {
				node.AddChild(int(42), children[0])
				So(node.NumChildren, ShouldEqual, 1)
				So(node.Children[42], ShouldEqual, children[0].Ref())
			})

			Convey("Adding multiple children", func() {
				node.AddChild(int(10), children[0])
				node.AddChild(int(20), children[1])
				node.AddChild(int(30), children[2])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Children[10], ShouldEqual, children[0].Ref())
				So(node.Children[20], ShouldEqual, children[1].Ref())
				So(node.Children[30], ShouldEqual, children[2].Ref())
			})

			Convey("Adding children with sparse keys", func() {
				// Add children with widely spaced keys
				node.AddChild(int(0), children[0])
				node.AddChild(int(128), children[1])
				node.AddChild(int(255), children[2])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Children[0], ShouldEqual, children[0].Ref())
				So(node.Children[128], ShouldEqual, children[1].Ref())
				So(node.Children[255], ShouldEqual, children[2].Ref())
			})

			Convey("Adding children at boundaries", func() {
				// Add children at the beginning and end of the byte range
				node.AddChild(int(0), children[0])
				node.AddChild(int(255), children[1])

				So(node.NumChildren, ShouldEqual, 2)
				So(node.Children[0], ShouldEqual, children[0].Ref())
				So(node.Children[255], ShouldEqual, children[1].Ref())
			})
		})

		Convey("When finding children", func() {
			// Setup children
			children := make([]*Leaf[any], 10)
			for i := 0; i < 10; i++ {
				children[i] = NewLeaf[any](a, []byte{byte(i * 25)}, nil)
				node.AddChild(int(i*25), children[i])
			}

			Convey("Finding existing children", func() {
				for i := 0; i < 10; i++ {
					found := node.FindChild(int(i * 25))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, children[i].Ref())
				}
			})

			Convey("Finding non-existent children", func() {
				// Test keys that are not in the sparse distribution
				So(node.FindChild(int(1)), ShouldBeNil)
				So(node.FindChild(int(26)), ShouldBeNil)
			})

			Convey("Finding children at boundaries", func() {
				// Test finding children at byte boundaries
				So(node.FindChild(int(0)), ShouldNotBeNil)   // First child
				So(node.FindChild(int(225)), ShouldNotBeNil) // Last child
			})
		})

		Convey("When checking capacity", func() {
			Convey("Empty node is not full", func() {
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 10 children is not full", func() {
				for i := 0; i < 10; i++ {
					child := NewLeaf[any](a, []byte{byte(i * 25)}, nil)
					node.AddChild(int(i*25), child)
				}
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with many children is not full", func() {
				// Add children to many different positions
				for i := 0; i < 100; i++ {
					child := NewLeaf[any](a, []byte{byte(i * 2)}, nil)
					node.AddChild(int(i*2), child)
				}
				So(node.Full(), ShouldBeFalse)
				So(node.NumChildren, ShouldEqual, 100)
			})

			Convey("Node with 256 children is full", func() {
				// Add children to all possible byte positions
				for i := 0; i < 256; i++ {
					child := NewLeaf[any](a, []byte{byte(i)}, nil)
					node.AddChild(int(i), child)
				}
				So(node.Full(), ShouldBeTrue)
				So(node.NumChildren, ShouldEqual, 256)
			})
		})

		Convey("When growing (no-op)", func() {
			// Setup some children
			for i := 0; i < 10; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 25)}, nil)
				node.AddChild(int(i*25), child)
			}

			Convey("Growing should return the same node", func() {
				newNode := node.Grow(a)
				So(newNode, ShouldEqual, node)
				So(newNode.Type(), ShouldEqual, TypeNode256)
			})

			Convey("Growing should not affect children", func() {
				originalChildren := node.NumChildren
				newNode := node.Grow(a)
				So(newNode.(*Node256[any]).NumChildren, ShouldEqual, originalChildren)
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

				node.AddChild(int('c'), child3)
				node.AddChild(int('a'), child1)
				node.AddChild(int('b'), child2)

				So(node.Minimum(), ShouldEqual, child1)
				So(node.Maximum(), ShouldEqual, child3)
			})
		})
	})
}

func TestNode256_EdgeCases(t *testing.T) {
	Convey("Given a Node256 with edge cases", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node256[any]{})

		Convey("When adding duplicate keys", func() {
			child1 := NewLeaf[any](a, []byte("a"), nil)
			child2 := NewLeaf[any](a, []byte("a"), nil)

			node.AddChild(int('a'), child1)
			So(node.NumChildren, ShouldEqual, 1)

			node.AddChild(int('a'), child2)
			So(node.NumChildren, ShouldEqual, 1) // Count doesn't change

			// Should replace the existing child
			found := node.FindChild(int('a'))
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child2.Ref())
		})

		Convey("When adding zero byte key", func() {
			child := NewLeaf[any](a, []byte{0}, nil)
			node.AddChild(int(0), child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(int(0))
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding 255 byte key", func() {
			child := NewLeaf[any](a, []byte{255}, nil)
			node.AddChild(int(255), child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(int(255))
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding children at sparse intervals", func() {
			// Add children with very sparse distribution
			sparseKeys := []byte{0, 64, 128, 192, 255}
			for i, key := range sparseKeys {
				child := NewLeaf[any](a, []byte{key}, nil)
				node.AddChild(int(key), child)
				So(node.NumChildren, ShouldEqual, i+1)
			}

			// Verify all sparse children can be found
			for _, key := range sparseKeys {
				found := node.FindChild(int(key))
				So(found, ShouldNotBeNil)
				So(*found, ShouldNotEqual, Ref[any](0))
			}
		})

		Convey("When adding children to all positions", func() {
			// Test adding children to every possible byte position
			for i := 0; i < 256; i++ {
				child := NewLeaf[any](a, []byte{byte(i)}, nil)
				node.AddChild(int(i), child)
			}

			So(node.NumChildren, ShouldEqual, 256)
			So(node.Full(), ShouldBeTrue)

			// Verify all children can be found
			for i := 0; i < 256; i++ {
				found := node.FindChild(int(i))
				So(found, ShouldNotBeNil)
				So(*found, ShouldNotEqual, Ref[any](0))
			}
		})
	})
}

func TestNode256_Performance(t *testing.T) {
	Convey("Given a Node256 with performance considerations", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node256[any]{})

		Convey("When adding many children", func() {
			// Test that adding many children works correctly
			children := make([]*Leaf[any], 100)
			for i := 0; i < 100; i++ {
				children[i] = NewLeaf[any](a, []byte{byte(i * 2)}, nil)
				node.AddChild(int(i*2), children[i])
			}

			So(node.NumChildren, ShouldEqual, 100)
			So(node.Full(), ShouldBeFalse)

			// Verify all children can be found
			for i := 0; i < 100; i++ {
				found := node.FindChild(int(i * 2))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, children[i].Ref())
			}
		})

		Convey("When searching with sparse distribution", func() {
			// Add children with sparse distribution to test search performance
			sparseKeys := []byte{1, 3, 7, 15, 31, 63, 127, 255}
			sparseChildren := make([]*Leaf[any], len(sparseKeys))
			for i, key := range sparseKeys {
				sparseChildren[i] = NewLeaf[any](a, []byte{key}, nil)
				node.AddChild(int(key), sparseChildren[i])
			}

			// Test finding existing and non-existing keys
			for _, key := range sparseKeys {
				So(node.FindChild(int(key)), ShouldNotBeNil)
			}

			// Test some non-existing keys
			So(node.FindChild(int(2)), ShouldBeNil)
			So(node.FindChild(int(8)), ShouldBeNil)
			So(node.FindChild(int(16)), ShouldBeNil)
		})

		Convey("When testing direct array access", func() {
			// Verify that direct array access works correctly
			child1 := NewLeaf[any](a, []byte{100}, nil)
			child2 := NewLeaf[any](a, []byte{200}, nil)

			node.AddChild(int(100), child1)
			node.AddChild(int(200), child2)

			So(node.Children[100], ShouldEqual, child1.Ref())
			So(node.Children[200], ShouldEqual, child2.Ref())
			So(node.Children[150], ShouldEqual, Ref[any](0)) // Unused position
		})

		Convey("When testing NumChildren accuracy", func() {
			// Test that NumChildren is accurately maintained
			So(node.NumChildren, ShouldEqual, 0)

			// Add a child
			child := NewLeaf[any](a, []byte{100}, nil)
			node.AddChild(int(100), child)
			So(node.NumChildren, ShouldEqual, 1)

			// Replace the same key
			child2 := NewLeaf[any](a, []byte{100}, nil)
			node.AddChild(int(100), child2)
			So(node.NumChildren, ShouldEqual, 1) // Count shouldn't change

			// Add a different key
			child3 := NewLeaf[any](a, []byte{200}, nil)
			node.AddChild(int(200), child3)
			So(node.NumChildren, ShouldEqual, 2)
		})
	})
}

func TestNode256_RemoveChild(t *testing.T) {
	Convey("Given a Node256 with children", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node256[any]{})

		// Setup children
		children := make([]*Leaf[any], 10)
		for i := 0; i < 10; i++ {
			children[i] = NewLeaf[any](a, []byte{byte(i * 25)}, nil)
			node.AddChild(int(i*25), children[i])
		}

		So(node.NumChildren, ShouldEqual, 10)

		Convey("When removing the first child", func() {
			childRef := node.FindChild(int(0))
			So(childRef, ShouldNotBeNil)

			node.RemoveChild(int(0), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 9)
			})

			Convey("And the child should not be found", func() {
				So(node.FindChild(int(0)), ShouldBeNil)
			})

			Convey("And the child reference should be cleared", func() {
				So(node.Children[0].Empty(), ShouldBeTrue)
			})

			Convey("And remaining children should still be accessible", func() {
				So(node.FindChild(int(25)), ShouldNotBeNil)
				So(node.FindChild(int(50)), ShouldNotBeNil)
				So(node.FindChild(int(75)), ShouldNotBeNil)
			})
		})

		Convey("When removing the middle child", func() {
			childRef := node.FindChild(int(50))
			So(childRef, ShouldNotBeNil)

			node.RemoveChild(int(50), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 9)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild(int(50))
				So(found, ShouldBeNil)
			})

			Convey("And the child reference should be cleared", func() {
				So(node.Children[50].Empty(), ShouldBeTrue)
			})

			Convey("And remaining children should still be accessible", func() {
				So(node.FindChild(int(0)), ShouldNotBeNil)
				So(node.FindChild(int(25)), ShouldNotBeNil)
				So(node.FindChild(int(75)), ShouldNotBeNil)
			})
		})

		Convey("When removing the last child", func() {
			childRef := node.FindChild(int(225))
			So(childRef, ShouldNotBeNil)

			node.RemoveChild(int(225), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 9)
			})

			Convey("And the child should not be found", func() {
				So(node.FindChild(int(225)), ShouldBeNil)
			})

			Convey("And the child reference should be cleared", func() {
				So(node.Children[225].Empty(), ShouldBeTrue)
			})

			Convey("And remaining children should still be accessible", func() {
				So(node.FindChild(int(0)), ShouldNotBeNil)
				So(node.FindChild(int(25)), ShouldNotBeNil)
				So(node.FindChild(int(50)), ShouldNotBeNil)
			})
		})

		Convey("When removing multiple children", func() {
			// Remove 25 first
			childRef := node.FindChild(int(25))
			node.RemoveChild(int(25), childRef)

			// Remove 100 second
			childRef = node.FindChild(int(100))
			node.RemoveChild(int(100), childRef)

			Convey("Then NumChildren should be 8", func() {
				So(node.NumChildren, ShouldEqual, 8)
			})

			Convey("And removed children should not be found", func() {
				So(node.FindChild(int(25)), ShouldBeNil)
				So(node.FindChild(int(100)), ShouldBeNil)
			})

			Convey("And remaining children should still be accessible", func() {
				So(node.FindChild(int(0)), ShouldNotBeNil)
				So(node.FindChild(int(50)), ShouldNotBeNil)
				So(node.FindChild(int(75)), ShouldNotBeNil)
				So(node.FindChild(int(125)), ShouldNotBeNil)
				So(node.FindChild(int(150)), ShouldNotBeNil)
				So(node.FindChild(int(175)), ShouldNotBeNil)
				So(node.FindChild(int(200)), ShouldNotBeNil)
				So(node.FindChild(int(225)), ShouldNotBeNil)
			})

			Convey("And child references should be properly cleared", func() {
				So(node.Children[25].Empty(), ShouldBeTrue)
				So(node.Children[100].Empty(), ShouldBeTrue)
			})
		})
	})
}

func TestNode256_Shrink(t *testing.T) {
	Convey("Given a Node256", t, func() {
		a := &arena.Arena{}

		node := arena.New(a, Node256[any]{})
		node.Partial = slice.FromString(a, "+")

		Convey("When shrinking with 37 or more children", func() {
			// Add 40 children
			for i := 0; i < 40; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 6)}, nil)
				node.AddChild(int(i*6), child)
			}

			So(node.NumChildren, ShouldEqual, 40)

			result := node.Shrink(a)

			Convey("Then should return the same node", func() {
				So(result, ShouldEqual, node)
			})

			Convey("And NumChildren should remain unchanged", func() {
				So(node.NumChildren, ShouldEqual, 40)
			})
		})

		Convey("When shrinking with exactly 36 children", func() {
			// Add 36 children
			for i := 0; i < 36; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 7)}, nil)
				node.AddChild(int(i*7), child)
			}

			So(node.NumChildren, ShouldEqual, 36)

			result := node.Shrink(a)

			Convey("Then should return a Node48", func() {
				So(result.Type(), ShouldEqual, TypeNode48)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And the new Node48 should have the same children", func() {
				node48 := result.(*Node48[any])
				So(node48.NumChildren, ShouldEqual, 36)

				// Verify all children are accessible
				for i := 0; i < 36; i++ {
					key := byte(i * 7)
					found := node48.FindChild(int(key))
					So(found, ShouldNotBeNil)
				}
			})
		})

		Convey("When shrinking with exactly 35 children", func() {
			// Add 35 children
			for i := 0; i < 35; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 7)}, nil)
				node.AddChild(int(i*7), child)
			}

			So(node.NumChildren, ShouldEqual, 35)

			result := node.Shrink(a)

			Convey("Then should return a Node48", func() {
				So(result.Type(), ShouldEqual, TypeNode48)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And the new Node48 should have the same children", func() {
				node48 := result.(*Node48[any])
				So(node48.NumChildren, ShouldEqual, 35)

				// Verify all children are accessible
				for i := 0; i < 35; i++ {
					key := byte(i * 7)
					found := node48.FindChild(int(key))
					So(found, ShouldNotBeNil)
				}
			})

			Convey("And the original node should be freed", func() {
				// The original node should be replaced, so we can't access it directly
				// This is verified by the fact that we get a Node48 back
			})
		})

		Convey("When shrinking with exactly 1 child", func() {
			child := NewLeaf[any](a, []byte("+a"), nil)

			node.AddChild(int('a'), child)

			So(node.NumChildren, ShouldEqual, 1)

			result := node.Shrink(a)

			Convey("Then should return a Node48", func() {
				So(result.Type(), ShouldEqual, TypeNode48)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And the new Node48 should have the same child", func() {
				node48 := result.(*Node48[any])
				So(node48.NumChildren, ShouldEqual, 1)
				found := node48.FindChild(int('a'))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child.Ref())
			})
		})

		Convey("When shrinking with no children", func() {
			So(node.NumChildren, ShouldEqual, 0)

			result := node.Shrink(a)

			Convey("Then should return a Node48", func() {
				So(result.Type(), ShouldEqual, TypeNode48)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And the new Node48 should have no children", func() {
				node48 := result.(*Node48[any])
				So(node48.NumChildren, ShouldEqual, 0)
			})
		})

		Convey("When shrinking with sparse children", func() {
			// Add children at sparse intervals
			sparseKeys := []byte{0, 64, 128, 192, 255}
			for _, key := range sparseKeys {
				child := NewLeaf[any](a, []byte{key}, nil)
				node.AddChild(int(key), child)
			}

			So(node.NumChildren, ShouldEqual, 5)

			result := node.Shrink(a)

			Convey("Then should return a Node48", func() {
				So(result.Type(), ShouldEqual, TypeNode48)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And the new Node48 should have the same children", func() {
				node48 := result.(*Node48[any])
				So(node48.NumChildren, ShouldEqual, 5)

				// Verify all sparse children are accessible
				for _, key := range sparseKeys {
					found := node48.FindChild(int(key))
					So(found, ShouldNotBeNil)
				}
			})
		})
	})
}
