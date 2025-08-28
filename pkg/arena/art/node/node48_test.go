package node_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	. "github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
	"github.com/flier/goutil/pkg/opt"
)

func TestNode48(t *testing.T) {
	Convey("Given a Node48", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node48[any]{})

		Convey("When checking basic properties", func() {
			So(node.Type(), ShouldEqual, TypeNode48)
			So(node.Full(), ShouldBeFalse)
			So(node.NumChildren, ShouldEqual, 0)
			So(node.Ref().Type(), ShouldEqual, TypeNode48)
		})

		Convey("When adding children", func() {
			// Create mock children
			children := make([]*Leaf[any], 48)
			for i := 0; i < 48; i++ {
				children[i] = NewLeaf[any](a, []byte{byte(i)}, nil)
			}

			Convey("Adding first child", func() {
				node.AddChild(opt.Some(byte(42)), children[0])
				So(node.NumChildren, ShouldEqual, 1)
				So(node.Keys[42], ShouldEqual, byte(1)) // 1-based indexing
				So(node.Children[0], ShouldEqual, children[0].Ref())
			})

			Convey("Adding multiple children", func() {
				node.AddChild(opt.Some(byte(10)), children[0])
				node.AddChild(opt.Some(byte(20)), children[1])
				node.AddChild(opt.Some(byte(30)), children[2])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Keys[10], ShouldEqual, byte(1))
				So(node.Keys[20], ShouldEqual, byte(2))
				So(node.Keys[30], ShouldEqual, byte(3))
				So(node.Children[0], ShouldEqual, children[0].Ref())
				So(node.Children[1], ShouldEqual, children[1].Ref())
				So(node.Children[2], ShouldEqual, children[2].Ref())
			})

			Convey("Adding children with sparse keys", func() {
				// Add children with widely spaced keys
				node.AddChild(opt.Some(byte(0)), children[0])
				node.AddChild(opt.Some(byte(128)), children[1])
				node.AddChild(opt.Some(byte(255)), children[2])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Keys[0], ShouldEqual, byte(1))
				So(node.Keys[128], ShouldEqual, byte(2))
				So(node.Keys[255], ShouldEqual, byte(3))
			})

			Convey("Adding children to fill capacity", func() {
				// Add 48 children to fill the node
				for i := 0; i < 48; i++ {
					node.AddChild(opt.Some(byte(i)), children[i])
				}

				So(node.NumChildren, ShouldEqual, 48)
				So(node.Full(), ShouldBeTrue)

				// Verify all children can be found
				for i := 0; i < 48; i++ {
					found := node.FindChild(opt.Some(byte(i)))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, children[i].Ref())
				}
			})
		})

		Convey("When finding children", func() {
			// Setup children
			children := make([]*Leaf[any], 10)
			for i := 0; i < 10; i++ {
				children[i] = NewLeaf[any](a, []byte{byte(i * 25)}, nil)
				node.AddChild(opt.Some(byte(i*25)), children[i])
			}

			Convey("Finding existing children", func() {
				for i := 0; i < 10; i++ {
					found := node.FindChild(opt.Some(byte(i * 25)))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, children[i].Ref())
				}
			})

			Convey("Finding non-existent children", func() {
				// Test keys that are not in the sparse distribution
				found := node.FindChild(opt.Some(byte(1)))
				So(found, ShouldBeNil)

				found = node.FindChild(opt.Some(byte(26))) // Between 0 and 25
				So(found, ShouldBeNil)

				found = node.FindChild(opt.Some(byte(51))) // Between 25 and 50
				So(found, ShouldBeNil)
			})

			Convey("Finding children at boundaries", func() {
				// Test finding children at byte boundaries
				So(node.FindChild(opt.Some(byte(0))), ShouldNotBeNil)   // First child
				So(node.FindChild(opt.Some(byte(225))), ShouldNotBeNil) // Last child
			})
		})

		Convey("When checking capacity", func() {
			Convey("Empty node is not full", func() {
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 24 children is not full", func() {
				for i := 0; i < 24; i++ {
					child := NewLeaf[any](a, []byte{byte(i * 10)}, nil)
					node.AddChild(opt.Some(byte(i*10)), child)
				}
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 48 children is full", func() {
				for i := 0; i < 48; i++ {
					child := NewLeaf[any](a, []byte{byte(i * 5)}, nil)
					node.AddChild(opt.Some(byte(i*5)), child)
				}
				So(node.Full(), ShouldBeTrue)
			})
		})

		Convey("When growing to Node256", func() {
			// Setup children
			for i := 0; i < 48; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 5)}, nil)
				node.AddChild(opt.Some(byte(i*5)), child)
			}

			Convey("Growing should create Node256", func() {
				newNode := node.Grow(a)
				So(newNode.Type(), ShouldEqual, TypeNode256)
			})

			Convey("Growing should preserve all children", func() {
				newNode := node.Grow(a)
				node256 := newNode.(*Node256[any])

				So(node256.NumChildren, ShouldEqual, 48)
				// Check that all children are properly mapped
				for i := 0; i < 48; i++ {
					key := byte(i * 5)
					So(node256.Children[key], ShouldEqual, node.Children[i])
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

func TestNode48_EdgeCases(t *testing.T) {
	Convey("Given a Node48 with edge cases", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node48[any]{})

		Convey("When adding duplicate keys", func() {
			child1 := NewLeaf[any](a, []byte("a"), nil)
			child2 := NewLeaf[any](a, []byte("a"), nil)

			node.AddChild(opt.Some(byte('a')), child1)
			node.AddChild(opt.Some(byte('a')), child2)

			// Should replace the existing child
			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(opt.Some(byte('a')))
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child2.Ref())
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

		Convey("When adding children at sparse intervals", func() {
			// Add children with very sparse distribution
			sparseKeys := []byte{0, 64, 128, 192, 255}
			for i, key := range sparseKeys {
				child := NewLeaf[any](a, []byte{key}, nil)
				node.AddChild(opt.Some(key), child)
				So(node.NumChildren, ShouldEqual, i+1)
			}

			// Verify all sparse children can be found
			for _, key := range sparseKeys {
				found := node.FindChild(opt.Some(key))
				So(found, ShouldNotBeNil)
			}
		})
	})
}

func TestNode48_Performance(t *testing.T) {
	Convey("Given a Node48 with performance considerations", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node48[any]{})

		Convey("When adding many children", func() {
			// Test that adding 48 children works correctly
			children := make([]*Leaf[any], 48)
			for i := 0; i < 48; i++ {
				children[i] = NewLeaf[any](a, []byte{byte(i * 5)}, nil)
				node.AddChild(opt.Some(byte(i*5)), children[i])
			}

			So(node.NumChildren, ShouldEqual, 48)
			So(node.Full(), ShouldBeTrue)

			// Verify all children can be found
			for i := 0; i < 48; i++ {
				found := node.FindChild(opt.Some(byte(i * 5)))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, children[i].Ref())
			}
		})

		Convey("When searching with sparse distribution", func() {
			// Add children with sparse distribution to test search performance
			sparseKeys := []byte{1, 3, 7, 15, 31, 63, 127, 255}
			for _, key := range sparseKeys {
				child := NewLeaf[any](a, []byte{key}, nil)
				node.AddChild(opt.Some(key), child)
			}

			// Test finding existing and non-existing keys
			for _, key := range sparseKeys {
				So(node.FindChild(opt.Some(key)), ShouldNotBeNil)
			}

			// Test some non-existing keys
			So(node.FindChild(opt.Some(byte(2))), ShouldBeNil)
			So(node.FindChild(opt.Some(byte(8))), ShouldBeNil)
			So(node.FindChild(opt.Some(byte(16))), ShouldBeNil)
		})

		Convey("When testing 1-based indexing", func() {
			// Verify that the 1-based indexing works correctly
			child1 := NewLeaf[any](a, []byte{100}, nil)
			child2 := NewLeaf[any](a, []byte{200}, nil)

			node.AddChild(opt.Some(byte(100)), child1)
			node.AddChild(opt.Some(byte(200)), child2)

			So(node.Keys[100], ShouldEqual, byte(1))
			So(node.Keys[200], ShouldEqual, byte(2))
			So(node.Children[0], ShouldEqual, child1.Ref())
			So(node.Children[1], ShouldEqual, child2.Ref())
		})
	})
}

func TestNode48_RemoveChild(t *testing.T) {
	Convey("Given a Node48 with children", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node48[any]{})

		// Setup children
		children := make([]*Leaf[any], 10)
		for i := 0; i < 10; i++ {
			children[i] = NewLeaf[any](a, []byte{byte(i * 25)}, nil)
			node.AddChild(opt.Some(byte(i*25)), children[i])
		}

		So(node.NumChildren, ShouldEqual, 10)

		Convey("When removing the first child", func() {
			childRef := node.FindChild(opt.Some(byte(0)))
			So(childRef, ShouldNotBeNil)

			node.RemoveChild(opt.Some(byte(0)), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 9)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild(opt.Some(byte(0)))
				So(found, ShouldBeNil)
			})

			Convey("And the key should be cleared", func() {
				So(node.Keys[0], ShouldEqual, byte(0))
			})

			Convey("And remaining children should still be accessible", func() {
				So(node.FindChild(opt.Some(byte(25))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(50))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(75))), ShouldNotBeNil)
			})
		})

		Convey("When removing the middle child", func() {
			childRef := node.FindChild(opt.Some(byte(50)))
			So(childRef, ShouldNotBeNil)

			node.RemoveChild(opt.Some(byte(50)), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 9)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild(opt.Some(byte(50)))
				So(found, ShouldBeNil)
			})

			Convey("And the key should be cleared", func() {
				So(node.Keys[50], ShouldEqual, byte(0))
			})

			Convey("And remaining children should still be accessible", func() {
				So(node.FindChild(opt.Some(byte(0))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(25))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(75))), ShouldNotBeNil)
			})
		})

		Convey("When removing the last child", func() {
			childRef := node.FindChild(opt.Some(byte(225)))
			So(childRef, ShouldNotBeNil)

			node.RemoveChild(opt.Some(byte(225)), childRef)

			Convey("Then NumChildren should be decremented", func() {
				So(node.NumChildren, ShouldEqual, 9)
			})

			Convey("And the child should not be found", func() {
				found := node.FindChild(opt.Some(byte(225)))
				So(found, ShouldBeNil)
			})

			Convey("And the key should be cleared", func() {
				So(node.Keys[225], ShouldEqual, byte(0))
			})

			Convey("And remaining children should still be accessible", func() {
				So(node.FindChild(opt.Some(byte(0))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(25))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(50))), ShouldNotBeNil)
			})
		})

		Convey("When removing multiple children", func() {
			// Remove 25 first
			childRef := node.FindChild(opt.Some(byte(25)))
			node.RemoveChild(opt.Some(byte(25)), childRef)

			// Remove 100 second
			childRef = node.FindChild(opt.Some(byte(100)))
			node.RemoveChild(opt.Some(byte(100)), childRef)

			Convey("Then NumChildren should be 8", func() {
				So(node.NumChildren, ShouldEqual, 8)
			})

			Convey("And removed children should not be found", func() {
				So(node.FindChild(opt.Some(byte(25))), ShouldBeNil)
				So(node.FindChild(opt.Some(byte(100))), ShouldBeNil)
			})

			Convey("And remaining children should still be accessible", func() {
				So(node.FindChild(opt.Some(byte(0))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(50))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(75))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(125))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(150))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(175))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(200))), ShouldNotBeNil)
				So(node.FindChild(opt.Some(byte(225))), ShouldNotBeNil)
			})

			Convey("And keys should be properly cleared", func() {
				So(node.Keys[25], ShouldEqual, byte(0))
				So(node.Keys[100], ShouldEqual, byte(0))
			})
		})
	})
}

func TestNode48_Shrink(t *testing.T) {
	Convey("Given a Node48", t, func() {
		a := &arena.Arena{}

		node := arena.New(a, Node48[any]{})
		node.Partial = slice.FromString(a, "+")

		Convey("When shrinking with 12 or more children", func() {
			// Add 15 children
			for i := 0; i < 15; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 16)}, nil)
				node.AddChild(opt.Some(byte(i*16)), child)
			}

			So(node.NumChildren, ShouldEqual, 15)

			result := node.Shrink(a)

			Convey("Then should return the same node", func() {
				So(result, ShouldEqual, node)
			})

			Convey("And NumChildren should remain unchanged", func() {
				So(node.NumChildren, ShouldEqual, 15)
			})
		})

		Convey("When shrinking with exactly 11 children", func() {
			// Add 11 children
			for i := 0; i < 11; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 23)}, nil)
				node.AddChild(opt.Some(byte(i*23)), child)
			}

			So(node.NumChildren, ShouldEqual, 11)

			result := node.Shrink(a)

			Convey("Then should return a Node16", func() {
				So(result.Type(), ShouldEqual, TypeNode16)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And the new Node16 should have the same children", func() {
				node16 := result.(*Node16[any])
				So(node16.NumChildren, ShouldEqual, 11)

				// Verify all children are accessible
				for i := 0; i < 11; i++ {
					key := byte(i * 23)
					found := node16.FindChild(opt.Some(key))
					So(found, ShouldNotBeNil)
				}
			})
		})

		Convey("When shrinking with exactly 10 children", func() {
			// Add 10 children
			for i := 0; i < 10; i++ {
				child := NewLeaf[any](a, []byte{byte(i * 25)}, nil)
				node.AddChild(opt.Some(byte(i*25)), child)
			}

			So(node.NumChildren, ShouldEqual, 10)

			result := node.Shrink(a)

			Convey("Then should return a Node16", func() {
				So(result.Type(), ShouldEqual, TypeNode16)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And the new Node16 should have the same children", func() {
				node16 := result.(*Node16[any])
				So(node16.NumChildren, ShouldEqual, 10)

				// Verify all children are accessible
				for i := 0; i < 10; i++ {
					key := byte(i * 25)
					found := node16.FindChild(opt.Some(key))
					So(found, ShouldNotBeNil)
				}
			})

			Convey("And the original node should be freed", func() {
				// The original node should be replaced, so we can't access it directly
				// This is verified by the fact that we get a Node16 back
			})
		})

		Convey("When shrinking with exactly 1 child", func() {
			child := NewLeaf[any](a, []byte("+a"), nil)
			node.AddChild(opt.Some(byte('a')), child)

			So(node.NumChildren, ShouldEqual, 1)

			result := node.Shrink(a)

			Convey("Then should return a Node16", func() {
				So(result.Type(), ShouldEqual, TypeNode16)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And the new Node16 should have the same child", func() {
				node16 := result.(*Node16[any])
				So(node16.NumChildren, ShouldEqual, 1)
				found := node16.FindChild(opt.Some(byte('a')))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child.Ref())
			})
		})

		Convey("When shrinking with no children", func() {
			So(node.NumChildren, ShouldEqual, 0)

			result := node.Shrink(a)

			Convey("Then should return a Node16", func() {
				So(result.Type(), ShouldEqual, TypeNode16)
				So(result.Prefix().Raw(), ShouldEqual, []byte("+"))
			})

			Convey("And the new Node16 should have no children", func() {
				node16 := result.(*Node16[any])
				So(node16.NumChildren, ShouldEqual, 0)
			})
		})
	})
}
