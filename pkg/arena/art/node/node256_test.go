package node

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

func TestNode256(t *testing.T) {
	Convey("Given a Node256", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node256{})

		Convey("When checking basic properties", func() {
			So(node.Type(), ShouldEqual, TypeNode256)
			So(node.Full(), ShouldBeFalse)
			So(node.NumChildren, ShouldEqual, 0)
			So(node.Ref().Type(), ShouldEqual, TypeNode256)
		})

		Convey("When adding children", func() {
			// Create mock children
			children := make([]*Leaf, 10)
			for i := 0; i < 10; i++ {
				children[i] = arena.New(a, Leaf{Key: slice.Of(a, byte(i*25))})
			}

			Convey("Adding first child", func() {
				node.AddChild(42, children[0])
				So(node.NumChildren, ShouldEqual, 1)
				So(node.Children[42], ShouldEqual, children[0].Ref())
			})

			Convey("Adding multiple children", func() {
				node.AddChild(10, children[0])
				node.AddChild(20, children[1])
				node.AddChild(30, children[2])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Children[10], ShouldEqual, children[0].Ref())
				So(node.Children[20], ShouldEqual, children[1].Ref())
				So(node.Children[30], ShouldEqual, children[2].Ref())
			})

			Convey("Adding children with sparse keys", func() {
				// Add children with widely spaced keys
				node.AddChild(0, children[0])
				node.AddChild(128, children[1])
				node.AddChild(255, children[2])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Children[0], ShouldEqual, children[0].Ref())
				So(node.Children[128], ShouldEqual, children[1].Ref())
				So(node.Children[255], ShouldEqual, children[2].Ref())
			})

			Convey("Adding children at boundaries", func() {
				// Add children at the beginning and end of the byte range
				node.AddChild(0, children[0])
				node.AddChild(255, children[1])

				So(node.NumChildren, ShouldEqual, 2)
				So(node.Children[0], ShouldEqual, children[0].Ref())
				So(node.Children[255], ShouldEqual, children[1].Ref())
			})
		})

		Convey("When finding children", func() {
			// Setup children
			children := make([]*Leaf, 10)
			for i := 0; i < 10; i++ {
				children[i] = arena.New(a, Leaf{Key: slice.Of(a, byte(i*25))})
				node.AddChild(byte(i*25), children[i])
			}

			Convey("Finding existing children", func() {
				for i := 0; i < 10; i++ {
					found := node.FindChild(byte(i * 25))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, children[i].Ref())
				}
			})

			Convey("Finding non-existent children", func() {
				// Test keys that are not in the sparse distribution
				found := node.FindChild(1)
				So(found, ShouldNotBeNil)
				// Note: The actual value depends on arena initialization
				// We just check that it's not the expected child
				So(*found, ShouldNotEqual, children[0].Ref())

				found = node.FindChild(26) // Between 0 and 25
				So(found, ShouldNotBeNil)
				So(*found, ShouldNotEqual, children[1].Ref())
			})

			Convey("Finding children at boundaries", func() {
				// Test finding children at byte boundaries
				So(node.FindChild(0), ShouldNotBeNil)   // First child
				So(node.FindChild(225), ShouldNotBeNil) // Last child
			})
		})

		Convey("When checking capacity", func() {
			Convey("Empty node is not full", func() {
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 10 children is not full", func() {
				for i := 0; i < 10; i++ {
					child := arena.New(a, Leaf{Key: slice.Of(a, byte(i*25))})
					node.AddChild(byte(i*25), child)
				}
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with many children is not full", func() {
				// Add children to many different positions
				for i := 0; i < 100; i++ {
					child := arena.New(a, Leaf{Key: slice.Of(a, byte(i*2))})
					node.AddChild(byte(i*2), child)
				}
				So(node.Full(), ShouldBeFalse)
				So(node.NumChildren, ShouldEqual, 100)
			})

			Convey("Node with 256 children is full", func() {
				// Add children to all possible byte positions
				for i := 0; i < 256; i++ {
					child := arena.New(a, Leaf{Key: slice.Of(a, byte(i))})
					node.AddChild(byte(i), child)
				}
				So(node.Full(), ShouldBeTrue)
				So(node.NumChildren, ShouldEqual, 256)
			})
		})

		Convey("When growing (no-op)", func() {
			// Setup some children
			for i := 0; i < 10; i++ {
				child := arena.New(a, Leaf{Key: slice.Of(a, byte(i*25))})
				node.AddChild(byte(i*25), child)
			}

			Convey("Growing should return the same node", func() {
				newNode := node.Grow(a)
				So(newNode, ShouldEqual, node)
				So(newNode.Type(), ShouldEqual, TypeNode256)
			})

			Convey("Growing should not affect children", func() {
				originalChildren := node.NumChildren
				newNode := node.Grow(a)
				So(newNode.(*Node256).NumChildren, ShouldEqual, originalChildren)
			})
		})

		Convey("When getting minimum and maximum", func() {
			Convey("Empty node should return nil", func() {
				So(node.Minimum(), ShouldBeNil)
				So(node.Maximum(), ShouldBeNil)
			})

			Convey("Node with children should return correct min/max", func() {
				child1 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})
				child2 := arena.New(a, Leaf{Key: slice.FromString(a, "b")})
				child3 := arena.New(a, Leaf{Key: slice.FromString(a, "c")})

				node.AddChild('c', child3)
				node.AddChild('a', child1)
				node.AddChild('b', child2)

				So(node.Minimum(), ShouldEqual, child1)
				So(node.Maximum(), ShouldEqual, child3)
			})
		})
	})
}

func TestNode256_EdgeCases(t *testing.T) {
	Convey("Given a Node256 with edge cases", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node256{})

		Convey("When adding duplicate keys", func() {
			child1 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})
			child2 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})

			node.AddChild('a', child1)
			So(node.NumChildren, ShouldEqual, 1)

			node.AddChild('a', child2)
			So(node.NumChildren, ShouldEqual, 1) // Count doesn't change

			// Should replace the existing child
			found := node.FindChild('a')
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child2.Ref())
		})

		Convey("When adding zero byte key", func() {
			child := arena.New(a, Leaf{Key: slice.Of(a, byte(0))})
			node.AddChild(0, child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(0)
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding 255 byte key", func() {
			child := arena.New(a, Leaf{Key: slice.Of(a, byte(255))})
			node.AddChild(255, child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(255)
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding children at sparse intervals", func() {
			// Add children with very sparse distribution
			sparseKeys := []byte{0, 64, 128, 192, 255}
			for i, key := range sparseKeys {
				child := arena.New(a, Leaf{Key: slice.Of(a, key)})
				node.AddChild(key, child)
				So(node.NumChildren, ShouldEqual, i+1)
			}

			// Verify all sparse children can be found
			for _, key := range sparseKeys {
				found := node.FindChild(key)
				So(found, ShouldNotBeNil)
				So(*found, ShouldNotEqual, Ref(0))
			}
		})

		Convey("When adding children to all positions", func() {
			// Test adding children to every possible byte position
			for i := 0; i < 256; i++ {
				child := arena.New(a, Leaf{Key: slice.Of(a, byte(i))})
				node.AddChild(byte(i), child)
			}

			So(node.NumChildren, ShouldEqual, 256)
			So(node.Full(), ShouldBeTrue)

			// Verify all children can be found
			for i := 0; i < 256; i++ {
				found := node.FindChild(byte(i))
				So(found, ShouldNotBeNil)
				So(*found, ShouldNotEqual, Ref(0))
			}
		})
	})
}

func TestNode256_Performance(t *testing.T) {
	Convey("Given a Node256 with performance considerations", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node256{})

		Convey("When adding many children", func() {
			// Test that adding many children works correctly
			children := make([]*Leaf, 100)
			for i := 0; i < 100; i++ {
				children[i] = arena.New(a, Leaf{Key: slice.Of(a, byte(i*2))})
				node.AddChild(byte(i*2), children[i])
			}

			So(node.NumChildren, ShouldEqual, 100)
			So(node.Full(), ShouldBeFalse)

			// Verify all children can be found
			for i := 0; i < 100; i++ {
				found := node.FindChild(byte(i * 2))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, children[i].Ref())
			}
		})

		Convey("When searching with sparse distribution", func() {
			// Add children with sparse distribution to test search performance
			sparseKeys := []byte{1, 3, 7, 15, 31, 63, 127, 255}
			sparseChildren := make([]*Leaf, len(sparseKeys))
			for i, key := range sparseKeys {
				sparseChildren[i] = arena.New(a, Leaf{Key: slice.Of(a, key)})
				node.AddChild(key, sparseChildren[i])
			}

			// Test finding existing and non-existing keys
			for _, key := range sparseKeys {
				So(node.FindChild(key), ShouldNotBeNil)
			}

			// Test some non-existing keys
			So(node.FindChild(2), ShouldNotBeNil)                            // Returns pointer to Ref
			So(*node.FindChild(2), ShouldNotEqual, sparseChildren[0].Ref())  // Check that it's not an expected child
			So(node.FindChild(8), ShouldNotBeNil)                            // Returns pointer to Ref
			So(*node.FindChild(8), ShouldNotEqual, sparseChildren[1].Ref())  // Check that it's not an expected child
			So(node.FindChild(16), ShouldNotBeNil)                           // Returns pointer to Ref
			So(*node.FindChild(16), ShouldNotEqual, sparseChildren[2].Ref()) // Check that it's not an expected child
		})

		Convey("When testing direct array access", func() {
			// Verify that direct array access works correctly
			child1 := arena.New(a, Leaf{Key: slice.Of(a, byte(100))})
			child2 := arena.New(a, Leaf{Key: slice.Of(a, byte(200))})

			node.AddChild(100, child1)
			node.AddChild(200, child2)

			So(node.Children[100], ShouldEqual, child1.Ref())
			So(node.Children[200], ShouldEqual, child2.Ref())
			So(node.Children[150], ShouldEqual, Ref(0)) // Unused position
		})

		Convey("When testing NumChildren accuracy", func() {
			// Test that NumChildren is accurately maintained
			So(node.NumChildren, ShouldEqual, 0)

			// Add a child
			child := arena.New(a, Leaf{Key: slice.Of(a, byte(100))})
			node.AddChild(100, child)
			So(node.NumChildren, ShouldEqual, 1)

			// Replace the same key
			child2 := arena.New(a, Leaf{Key: slice.Of(a, byte(100))})
			node.AddChild(100, child2)
			So(node.NumChildren, ShouldEqual, 1) // Count shouldn't change

			// Add a different key
			child3 := arena.New(a, Leaf{Key: slice.Of(a, byte(200))})
			node.AddChild(200, child3)
			So(node.NumChildren, ShouldEqual, 2)
		})
	})
}
